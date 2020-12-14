// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	"sync"
	"time"

	ds_models "github.com/edgexfoundry/device-sdk-go/pkg/models"
	sdk "github.com/edgexfoundry/device-sdk-go/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	rpio "github.com/stianeikeland/go-rpio"
)

var (
	// Use mcu pin 22, corresponds to GPIO3 on the pi
	pin = rpio.Pin(21)
)

type GPIODriver struct {
	lc           logger.LoggingClient
	asyncCh      chan<- *ds_models.AsyncValues
	switchButton bool
}

var once sync.Once
var driver *GPIODriver
var sdkService sdk.DeviceService

func NewGPIODriver() ds_models.ProtocolDriver {
	once.Do(func() {
		driver = new(GPIODriver)
	})
	return driver
}

// Initialize performs protocol-specific initialization for the device service.
func (s *GPIODriver) Initialize(lc logger.LoggingClient, asyncCh chan<- *ds_models.AsyncValues, deviceCh chan<- []ds_models.DiscoveredDevice) error {
	s.lc = lc
	s.asyncCh = asyncCh
	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *GPIODriver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []ds_models.CommandRequest) (res []*ds_models.CommandValue, err error) {

	for _, req := range reqs {
		//s.lc.Debug(fmt.Sprintf("GPIODriver.HandleReadCommand: device: %s operation: %v attributes: %v", addr.Name, req.RO.Operation, req.DeviceObject.Attributes))
		s.lc.Debug(fmt.Sprintf("device: %v resource: %v attributes: %v", deviceName, req.DeviceResourceName, req.Attributes))
	}

	now := time.Now().UnixNano() / int64(time.Millisecond)
	val, err := s.readPin()

	if err != nil {
		s.lc.Error(fmt.Sprintf("GPIODriver.HandleReadCommands; %v", err))
		return
	}

	s.lc.Debug(fmt.Sprintf("Moisture detected value: ", val))

	//cv, _ := ds_models.NewInt16Value(&reqs[0].RO, now, int16(val))
	cv, _ := ds_models.NewInt16Value(reqs[0].DeviceResourceName, now, int16(val))
	res = append(res, cv)

	return res, nil
}

func (s *GPIODriver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []ds_models.CommandRequest,
	params []*ds_models.CommandValue) error {

	s.lc.Debug(fmt.Sprintf("GPIODriver.HandleWriteCommands not supported"))
	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *GPIODriver) Stop(force bool) error {
	s.lc.Debug(fmt.Sprintf("GPIODriver.Stop called: force=%v", force))
	return nil
}

// RemoveDevice handles protocol-specific cleanup when a device is removed.
func (s *GPIODriver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	return nil
}

// AddDevice handles protocol specific init when a device is added
func (d *GPIODriver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	return nil
}

// UpdateDevice handles protocol specific actions when a device is updated
func (d *GPIODriver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	return nil
}

func (s *GPIODriver) readPin() (result int, err error) {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		s.lc.Error(fmt.Sprintf("GPIODriver.readPin error opening map"))
		return 0, err
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Pull up and read value
	pin.PullUp()
	st := pin.Read()
	if st == rpio.Low {
		s.lc.Debug(fmt.Sprintf("GPIODriver.readPin water detected"))
		return 1, nil
	}
	s.lc.Debug(fmt.Sprintf("GPIODriver.readPin water not detected"))
	return 0, nil
}

// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2019 Dell Technologies
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	ds_models "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/edgex-go/pkg/clients/logging"
	"github.com/edgexfoundry/edgex-go/pkg/models"
	"time"
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

// DisconnectDevice handles protocol-specific cleanup when a device
// is removed.
func (s *GPIODriver) DisconnectDevice(address *models.Addressable) error {
	return nil
}

// Initialize performs protocol-specific initialization for the device
// service.
func (s *GPIODriver) Initialize(lc logger.LoggingClient, asyncCh chan<- *ds_models.AsyncValues) error {
	s.lc = lc
	s.asyncCh = asyncCh
	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *GPIODriver) HandleReadCommands(addr *models.Addressable, reqs []ds_models.CommandRequest) (res []*ds_models.CommandValue, err error) {

	for _, req := range reqs {
		s.lc.Debug(fmt.Sprintf("GPIODriver.HandleReadCommand: device: %s operation: %v attributes: %v", addr.Name, req.RO.Operation, req.DeviceObject.Attributes))
	}

	now := time.Now().UnixNano() / int64(time.Millisecond)
	val, err := s.readPin()

	if err != nil {
		s.lc.Error(fmt.Sprintf("GPIODriver.HandleReadCommands; %v", err))
		return
	}
	
	s.lc.Debug(fmt.Sprintf("Moisture detected value: ", val))

	cv, _ := ds_models.NewInt16Value(&reqs[0].RO, now, int16(val))
	res = append(res, cv)

	return
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource (aka DeviceObject).
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *GPIODriver) HandleWriteCommands(addr *models.Addressable, reqs []ds_models.CommandRequest,
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
	st:= pin.Read()
	if st == rpio.Low {
		s.lc.Debug(fmt.Sprintf("GPIODriver.readPin water detected"))
		return 1, nil
	}
	s.lc.Debug(fmt.Sprintf("GPIODriver.readPin water not detected"))
	return 0, nil
}

// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
	"github.com/edgexfoundry/device-gpio-moist-go/internal/driver"
	"github.com/edgexfoundry/device-gpio-moist-go"
)

const (
	serviceName string = "device-gpio-moist-go"
)

func main() {
	sd := driver.NewGPIODriver()
	startup.Bootstrap(serviceName, device_gpio.Version, &sd)
}

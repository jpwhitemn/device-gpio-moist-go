# Jim's Device GPIO Moisture sensor in Go
=========================================

Simple EdgeX Foundry Device Service.  The service reads the GPIO pin #21 connected to a moisture sensor.  If it reads "high" on pin, then water is detected and a 1 is sent to core data.  On low read of pin, no water is detected and 0 is sent to core data.

See https://www.instructables.com/id/Soil-Moisture-Sensor-Raspberry-Pi/ for details on the moisture detector and how to connect the sensor to the Pi's GPIO pins.

See https://github.com/stianeikeland/go-rpio for the Go Lang library to read GPIO pins.



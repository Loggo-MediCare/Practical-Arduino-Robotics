// Port of Chapter07/plotter_demo/plotter_demo.ino.
package main

import (
	"machine"
	"time"

	"practicalarduinorobotics/tinygo/internal/numfmt"
	"practicalarduinorobotics/tinygo/internal/trig"
)

const sinFreq = 2.0
const cosFreq = 4.0
const sinAmp = 10
const cosAmp = 2

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})

	for {
		// Argument for 1Hz sin and cos waves.
		arg := 2 * trig.Pi * float64(time.Now().UnixMilli()) / 1000.0
		sinwave := sinAmp * trig.Sin(sinFreq*arg)
		coswave := cosAmp * trig.Cos(cosFreq*arg)

		machine.Serial.Write([]byte("sin:"))
		machine.Serial.Write([]byte(numfmt.Fixed(float32(sinwave), 2)))
		machine.Serial.Write([]byte("\tsin*cos:"))
		machine.Serial.Write([]byte(numfmt.Fixed(float32(sinwave*coswave), 2)))
		machine.Serial.Write([]byte("\tmin:"))
		machine.Serial.Write([]byte(numfmt.Fixed(-sinAmp*cosAmp, 2)))
		machine.Serial.Write([]byte("\tmax:"))
		machine.Serial.Write([]byte(numfmt.Fixed(sinAmp*cosAmp, 2)))
		machine.Serial.Write([]byte("\r\n"))
		time.Sleep(10 * time.Millisecond)
	}
}

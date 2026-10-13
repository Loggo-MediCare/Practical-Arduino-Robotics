// Port of Chapter02/Analog_potentiometer/Analog_potentiometer.ino.
package main

import (
	"machine"
	"time"

	"practicalarduinorobotics/tinygo/internal/numfmt"
)

// analogRead10 reads an ADC channel and scales it back down to the
// 0..1023 range that Arduino's analogRead() returns (TinyGo's ADC.Get()
// returns a value scaled to the full 16 bit range).
func analogRead10(pin machine.Pin) int32 {
	adc := machine.ADC{Pin: pin}
	adc.Configure(machine.ADCConfig{})
	return int32(adc.Get() >> 6)
}

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	machine.InitADC()

	for {
		analogValue := analogRead10(machine.ADC0)
		if time.Now().UnixMilli()%50 == 0 {
			machine.Serial.Write([]byte(numfmt.Int(analogValue) + "\r\n"))
		}
	}
}

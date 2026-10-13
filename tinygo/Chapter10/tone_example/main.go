// Port of Chapter10/tone_example/tone_example.ino.
package main

import (
	"machine"
	"time"
)

const speakerPin = machine.D10 // Timer1 channel B.
const analogInputPin = machine.ADC0

var frequencies = [5]int{440, 523, 659, 523, 440}
var durations = [5]int{200, 200, 200, 200, 450}

var pwm = machine.Timer1
var toneChannel uint8

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})

	machine.InitADC()
	adc := machine.ADC{Pin: analogInputPin}
	adc.Configure(machine.ADCConfig{})

	pwm.Configure(machine.PWMConfig{})
	toneChannel, _ = pwm.Channel(speakerPin)

	for {
		// Create tone from analog input.
		freq := 100 * (int(adc.Get()>>12) + 1)
		tone(freq)
		time.Sleep(5 * time.Millisecond)

		// Play custom jingle whenever Serial data arrives.
		if machine.Serial.Buffered() > 0 {
			jingle()
			// Clear Serial input buffer.
			for machine.Serial.Buffered() > 0 {
				machine.Serial.ReadByte()
			}
		}
	}
}

func jingle() {
	for i := 0; i < 5; i++ {
		tone(frequencies[i])
		time.Sleep(time.Duration(durations[i]) * time.Millisecond)
	}
}

// tone plays a continuous square wave at frequency Hz on speakerPin via
// hardware PWM at 50% duty cycle, mirroring Arduino's non-blocking tone().
func tone(frequency int) {
	if frequency <= 0 {
		return
	}
	pwm.SetPeriod(uint64(1e9 / frequency))
	pwm.Set(toneChannel, pwm.Top()/2)
}

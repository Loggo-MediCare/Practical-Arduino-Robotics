// Port of Chapter03/dc_motor/dc_motor.ino.
package main

import "machine"

const dirPin = machine.D12
const brakePin = machine.D9
const pwmPin = machine.D3
const analogPin = machine.ADC5

func mapRange(x, inMin, inMax, outMin, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func main() {
	dirPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	brakePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	brakePin.Low() // Disable brake functionality.

	machine.InitADC()
	adc := machine.ADC{Pin: analogPin}
	adc.Configure(machine.ADCConfig{})

	pwm := machine.Timer2 // Timer2 drives pin D3.
	pwm.Configure(machine.PWMConfig{})
	channel, _ := pwm.Channel(pwmPin)

	for {
		analogVal := int(adc.Get() >> 6) // scale back down to 0..1023
		var dutyPwm int
		if analogVal < 512 {
			// Turn forward.
			dirPin.High()
			dutyPwm = mapRange(analogVal, 0, 511, 255, 0)
		} else {
			// Turn backward.
			dirPin.Low()
			dutyPwm = mapRange(analogVal, 512, 1023, 0, 255)
		}
		pwm.Set(channel, uint32(dutyPwm)*pwm.Top()/255)
	}
}

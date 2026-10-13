// Port of Chapter03/stepper/stepper.ino.
package main

import (
	"machine"

	"practicalarduinorobotics/tinygo/internal/stepper"
)

// Phase A pins.
const dirPinA = machine.D12
const brakePinA = machine.D9
const pwmPinA = machine.D3

// Phase B pins.
const dirPinB = machine.D13
const brakePinB = machine.D8
const pwmPinB = machine.D11

// Analog input pins.
const analogPinPower = machine.ADC5
const analogPinSpeed = machine.ADC4

const stepsPerRev = 200

func mapRange(x, inMin, inMax, outMin, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func main() {
	brakePinA.Configure(machine.PinConfig{Mode: machine.PinOutput})
	brakePinA.Low()
	brakePinB.Configure(machine.PinConfig{Mode: machine.PinOutput})
	brakePinB.Low()

	bipolarStepper := stepper.New(stepsPerRev, dirPinA, dirPinB)

	machine.InitADC()
	powerAdc := machine.ADC{Pin: analogPinPower}
	powerAdc.Configure(machine.ADCConfig{})
	speedAdc := machine.ADC{Pin: analogPinSpeed}
	speedAdc.Configure(machine.ADCConfig{})

	// Timer2 drives both pin D3 (channel A) and pin D11 (channel B).
	pwm := machine.Timer2
	pwm.Configure(machine.PWMConfig{})
	channelA, _ := pwm.Channel(pwmPinA)
	channelB, _ := pwm.Channel(pwmPinB)

	for {
		// Map power potentiometer to a PWM output.
		powerInput := int(powerAdc.Get() >> 6)
		outputPwm := mapRange(powerInput, 0, 1023, 0, 255)
		// Map speed potentiometer to an rpm output.
		speedInput := int(speedAdc.Get() >> 6)
		rpm := mapRange(speedInput, 0, 1023, 0, 100)
		// Drive both phases with the same PWM.
		pwm.Set(channelA, uint32(outputPwm)*pwm.Top()/255)
		pwm.Set(channelB, uint32(outputPwm)*pwm.Top()/255)
		// Set the motor speed.
		bipolarStepper.SetSpeed(rpm)
		// Only step if rpm is greater than 0.
		if rpm > 0 {
			bipolarStepper.Step(1)
		}
	}
}

// Port of Chapter03/rc_servo/rc_servo.ino.
//
// Note: the original sketch drives the servo from pin 3, relying on
// Arduino's Servo library, which bit-bangs a 20ms/50Hz period via Timer1
// compare interrupts and so works on any pin. TinyGo's machine.PWM instead
// exposes each AVR hardware timer's real period limits directly: Timer2
// (which drives pin 3) is an 8 bit timer that tops out at ~16.4ms, too short
// for a 20ms servo period. Timer1 (16 bit, pins 9/10) can do it, so this port
// uses pin 9 instead.
package main

import "machine"

const inputPin = machine.ADC5
const servoPin = machine.D9 // Timer1 channel A; must be a 16 bit PWM pin.

func mapRange(x, inMin, inMax, outMin, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func main() {
	machine.InitADC()
	adc := machine.ADC{Pin: inputPin}
	adc.Configure(machine.ADCConfig{})

	// A hobby servo expects a 50Hz signal with a 1000-2000us pulse width.
	pwm := machine.Timer1
	pwm.Configure(machine.PWMConfig{Period: 20_000_000}) // 20ms, in nanoseconds.
	channel, _ := pwm.Channel(servoPin)

	for {
		// Read the potentiometer.
		input := int(adc.Get() >> 6) // scale back down to 0..1023
		// Map the analog value from 0..1023 to a pulse width between 1000 and 2000us.
		pulseWidthUs := mapRange(input, 0, 1023, 1000, 2000)
		pwm.Set(channel, uint32(pulseWidthUs)*pwm.Top()/20000)
	}
}

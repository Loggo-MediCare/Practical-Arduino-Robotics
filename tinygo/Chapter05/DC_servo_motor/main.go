// Port of Chapter05/DC_servo_motor/DC_servo_motor.ino.
package main

import (
	"machine"
	"time"

	"practicalarduinorobotics/tinygo/internal/numfmt"
)

// Pins for motor control.
const dirPin0 = machine.D8
const dirPin1 = machine.D9
const pwmPin = machine.D5

// Encoder pins (Arduino Uno's only interrupt-capable pins).
const pinA = machine.D2
const pinB = machine.D3

const controlInterval = 10 * time.Millisecond
const printInterval = 20 * time.Millisecond

var encoderPosition int32

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	pinA.Configure(machine.PinConfig{Mode: machine.PinInput})
	pinB.Configure(machine.PinConfig{Mode: machine.PinInput})
	dirPin0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dirPin1.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.InitADC()
	targetAdc := machine.ADC{Pin: machine.ADC0}
	targetAdc.Configure(machine.ADCConfig{})

	pwm := machine.Timer0 // Timer0 drives pin D5.
	pwm.Configure(machine.PWMConfig{})
	pwmChannel, _ := pwm.Channel(pwmPin)

	attachEncoderInterrupts()

	var targetPosition int32
	var iTerm float32
	var lastError int32
	var controllerOutput float32

	lastControlTime := time.Now()
	lastPrintTime := time.Now()

	for {
		now := time.Now()
		if now.Sub(lastControlTime) >= controlInterval {
			lastControlTime = lastControlTime.Add(controlInterval)
			targetPosition = int32(targetAdc.Get() >> 6)
			controllerOutput, lastError, iTerm = motorPidController(
				targetPosition, encoderPosition, lastError, iTerm)
			setMotorPwm(pwm, pwmChannel, int(controllerOutput))
		}
		if now.Sub(lastPrintTime) >= printInterval {
			lastPrintTime = lastPrintTime.Add(printInterval)
			machine.Serial.Write([]byte("Target_position:"))
			machine.Serial.Write([]byte(numfmt.Int(targetPosition)))
			machine.Serial.Write([]byte("\tActual_position:"))
			machine.Serial.Write([]byte(numfmt.Int(encoderPosition)))
			machine.Serial.Write([]byte("\tController_output:"))
			machine.Serial.Write([]byte(numfmt.Fixed(controllerOutput, 2)))
			machine.Serial.Write([]byte("\r\n"))
		}
	}
}

func motorPidController(targetPosition, encoderPos int32, lastError int32, iTerm float32) (output float32, newLastError int32, newITerm float32) {
	const kp = 1.0
	const ki = 0.0
	const kd = 60.0
	intervalMs := float32(controlInterval.Milliseconds())

	errorVal := targetPosition - encoderPos
	pTerm := kp * float32(errorVal)
	newITerm = iTerm + ki*float32(errorVal+lastError)/2*intervalMs
	newITerm = constrainF(newITerm, -255, 255)
	dTerm := kd * float32(errorVal-lastError) / intervalMs
	output = pTerm + newITerm + dTerm
	newLastError = errorVal
	return
}

func constrainF(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func setMotorPwm(pwm machine.PWM, channel uint8, motorPwm int) {
	if motorPwm > 0 {
		dirPin0.High()
		dirPin1.Low()
	} else {
		dirPin0.Low()
		dirPin1.High()
		motorPwm = -motorPwm
	}
	if motorPwm < 0 {
		motorPwm = 0
	}
	if motorPwm > 255 {
		motorPwm = 255
	}
	pwm.Set(channel, uint32(motorPwm)*pwm.Top()/255)
}

// attachEncoderInterrupts wires up quadrature decoding, mirroring the
// four dynamically-reattached RISING/FALLING handlers of the original
// sketch with a pair of PinToggle handlers that check the current pin level.
func attachEncoderInterrupts() {
	pinA.SetInterrupt(machine.PinToggle, func(p machine.Pin) {
		if p.Get() {
			if !pinB.Get() {
				encoderPosition++
			} else {
				encoderPosition--
			}
		} else {
			if !pinB.Get() {
				encoderPosition--
			} else {
				encoderPosition++
			}
		}
	})
	pinB.SetInterrupt(machine.PinToggle, func(p machine.Pin) {
		if p.Get() {
			if !pinA.Get() {
				encoderPosition--
			} else {
				encoderPosition++
			}
		} else {
			if !pinA.Get() {
				encoderPosition++
			} else {
				encoderPosition--
			}
		}
	})
}

// Package stepper ports the 2-pin mode of the Arduino Stepper library, which
// is the only mode Chapter03/stepper.ino exercises: two direction pins wired
// to a driver board are stepped through a 4-state full-step sequence.
package stepper

import (
	"machine"
	"time"
)

type Stepper struct {
	stepsPerRev int
	pin1, pin2  machine.Pin

	stepDelay    time.Duration
	lastStepTime time.Time
	direction    int
	stepNumber   int
}

func New(stepsPerRev int, pin1, pin2 machine.Pin) *Stepper {
	pin1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return &Stepper{
		stepsPerRev: stepsPerRev,
		pin1:        pin1,
		pin2:        pin2,
	}
}

// SetSpeed sets the motor speed in revolutions per minute.
func (s *Stepper) SetSpeed(rpm int) {
	if rpm <= 0 {
		return
	}
	s.stepDelay = time.Duration(60*1000*1000/(s.stepsPerRev*rpm)) * time.Microsecond
}

// Step moves the motor stepsToMove steps, blocking until the configured
// speed's timing allows each step, exactly like the original library.
func (s *Stepper) Step(stepsToMove int) {
	stepsLeft := stepsToMove
	if stepsLeft < 0 {
		stepsLeft = -stepsLeft
		s.direction = 0
	} else {
		s.direction = 1
	}

	for stepsLeft > 0 {
		now := time.Now()
		if now.Sub(s.lastStepTime) >= s.stepDelay {
			s.lastStepTime = now
			if s.direction == 1 {
				s.stepNumber++
				if s.stepNumber == s.stepsPerRev {
					s.stepNumber = 0
				}
			} else {
				if s.stepNumber == 0 {
					s.stepNumber = s.stepsPerRev
				}
				s.stepNumber--
			}
			stepsLeft--
			s.stepMotor(s.stepNumber % 4)
		}
	}
}

func (s *Stepper) stepMotor(step int) {
	switch step {
	case 0:
		s.pin1.Low()
		s.pin2.Low()
	case 1:
		s.pin1.Low()
		s.pin2.High()
	case 2:
		s.pin1.High()
		s.pin2.High()
	case 3:
		s.pin1.High()
		s.pin2.Low()
	}
}

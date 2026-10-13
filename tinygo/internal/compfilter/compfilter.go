// Package compfilter ports Chapter13/CompFilter: a complementary filter for
// tilt-angle estimation from gyroscope and accelerometer readings.
package compfilter

import "math"

type CompFilter struct {
	angleFiltered float32
	angleGyro     float32
	angleAccel    float32
	filterCoeff   float32
	gyroOffset    float32
	initialized   bool
}

func New(coeff float32) *CompFilter {
	return &CompFilter{filterCoeff: coeff}
}

func (c *CompFilter) AngleFiltered() float32 { return c.angleFiltered }
func (c *CompFilter) AngleAccel() float32    { return c.angleAccel }
func (c *CompFilter) AngleGyro() float32     { return c.angleGyro }

func (c *CompFilter) SetGyroOffset(offset float32) {
	c.gyroOffset = offset
}

// Update performs one incremental tilt estimation step.
// rateZ is the angular rate around one axis in rad/s, accelX/accelY are
// accelerations along the other two axes, and samplingInterval is in
// microseconds.
func (c *CompFilter) Update(rateZ, accelX, accelY, samplingInterval float32) {
	c.angleAccel = float32(math.Atan2(float64(accelX), float64(accelY)))
	if !c.initialized {
		c.angleGyro = c.angleAccel
		c.angleFiltered = c.angleAccel
		c.initialized = true
	}
	gyroDelta := 0.000001 * ((c.gyroOffset - rateZ) * samplingInterval)
	c.angleGyro += gyroDelta
	c.angleFiltered = c.filterCoeff*(c.angleFiltered+gyroDelta) + (1.0-c.filterCoeff)*c.angleAccel
}

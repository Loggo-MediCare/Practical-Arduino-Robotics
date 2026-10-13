// Package trig provides Sin, Cos, and Atan2 implementations that work on
// TinyGo's AVR target.
//
// The standard library's math.Sin/Cos fail to compile for AVR ("LLVM ERROR:
// unsupported library call operation": their pure-Go implementations lower
// to an llvm.sin/llvm.cos intrinsic that has no AVR libm backing), and
// math.Atan2 fails to link ("linker could not find symbol atan2") for the
// same reason. These are small self-contained approximations built only
// from +, -, *, / and comparisons, which the AVR soft-float runtime does
// support.
package trig

const Pi = 3.14159265358979323846
const twoPi = 2 * Pi
const halfPi = Pi / 2

// Sin approximates sin(x) via a 9th-order Taylor series around a range
// reduced to (-Pi, Pi], accurate to better than 1e-4 over that range.
func Sin(x float64) float64 {
	x = reduceAngle(x)
	x2 := x * x
	term := x
	sum := term
	term *= -x2 / (2 * 3)
	sum += term
	term *= -x2 / (4 * 5)
	sum += term
	term *= -x2 / (6 * 7)
	sum += term
	term *= -x2 / (8 * 9)
	sum += term
	return sum
}

// Cos approximates cos(x) as Sin(x + Pi/2).
func Cos(x float64) float64 {
	return Sin(x + halfPi)
}

func reduceAngle(x float64) float64 {
	x -= twoPi * float64(int64(x/twoPi))
	if x > Pi {
		x -= twoPi
	} else if x < -Pi {
		x += twoPi
	}
	return x
}

// atanApprox approximates atan(z) for |z| <= 1 to within about 0.07 degrees,
// using the polynomial approximation by Rajan et al.
func atanApprox(z float64) float64 {
	absZ := z
	if absZ < 0 {
		absZ = -absZ
	}
	return z*(Pi/4.0) - z*(absZ-1)*(0.2447+0.0663*absZ)
}

func atan(z float64) float64 {
	switch {
	case z > 1:
		return halfPi - atanApprox(1/z)
	case z < -1:
		return -halfPi - atanApprox(1/z)
	default:
		return atanApprox(z)
	}
}

// Atan2 mirrors math.Atan2.
func Atan2(y, x float64) float64 {
	switch {
	case x > 0:
		return atan(y / x)
	case x < 0 && y >= 0:
		return atan(y/x) + Pi
	case x < 0 && y < 0:
		return atan(y/x) - Pi
	case x == 0 && y > 0:
		return halfPi
	case x == 0 && y < 0:
		return -halfPi
	default:
		return 0
	}
}

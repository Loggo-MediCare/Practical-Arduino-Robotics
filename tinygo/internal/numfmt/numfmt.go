// Package numfmt provides small integer and fixed-point number-to-string
// formatters for Serial output.
//
// It exists instead of using strconv.FormatInt/FormatFloat directly because
// those general-purpose formatters (correctly-rounded arbitrary-precision
// float formatting in particular) pull in several KB of code and, in
// FormatFloat's case, enough transient state to blow Arduino Uno's 2KB of
// RAM on sketches that otherwise fit comfortably. A plain 32 bit
// divide-by-10 loop is all any sketch here actually needs.
package numfmt

// Int formats a 32 bit signed integer in base 10.
func Int(v int32) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf [11]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// Fixed formats v with a fixed number of decimal digits.
func Fixed(v float32, decimals int) string {
	neg := v < 0
	if neg {
		v = -v
	}
	scale := float32(1)
	for i := 0; i < decimals; i++ {
		scale *= 10
	}
	scaled := int32(v*scale + 0.5)
	intPart := scaled / int32(scale)
	fracPart := scaled % int32(scale)
	fracStr := Int(fracPart)
	for len(fracStr) < decimals {
		fracStr = "0" + fracStr
	}
	s := Int(intPart) + "." + fracStr
	if neg {
		s = "-" + s
	}
	return s
}

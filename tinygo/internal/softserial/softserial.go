// Package softserial is a minimal bit-banged 8N1 UART used in place of
// Arduino's SoftwareSerial library, which TinyGo has no equivalent for.
//
// Unlike SoftwareSerial (which uses a pin-change interrupt to catch the start
// bit asynchronously), this is a purely cooperative implementation: Available
// must be polled often enough from the main loop to catch the falling edge of
// a start bit, and ReadByte then blocks for the ~1 bit-period-per-bit
// duration of the whole byte. This matches the four call sites in this repo
// (Serial_GPS, BLE_chat, BLE_control_pad, sine_wave_generator), all of which
// have nothing else to do while receiving.
package softserial

import (
	"machine"
	"time"
)

type Port struct {
	rx, tx   machine.Pin
	bitDelay time.Duration
}

func New(rx, tx machine.Pin, baud int) *Port {
	return &Port{rx: rx, tx: tx, bitDelay: time.Second / time.Duration(baud)}
}

func (p *Port) Configure() {
	p.tx.Configure(machine.PinConfig{Mode: machine.PinOutput})
	p.tx.High() // idle high
	p.rx.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
}

func (p *Port) WriteByte(b byte) {
	p.tx.Low() // start bit
	time.Sleep(p.bitDelay)
	for i := 0; i < 8; i++ {
		p.tx.Set(b&(1<<uint(i)) != 0)
		time.Sleep(p.bitDelay)
	}
	p.tx.High() // stop bit
	time.Sleep(p.bitDelay)
}

func (p *Port) Write(data []byte) {
	for _, b := range data {
		p.WriteByte(b)
	}
}

func (p *Port) WriteString(s string) {
	for i := 0; i < len(s); i++ {
		p.WriteByte(s[i])
	}
}

// Available reports whether a start bit is currently on the RX line.
func (p *Port) Available() bool {
	return !p.rx.Get() // idle high, start bit pulls low
}

// ReadByte blocks for one byte period to receive a full byte. Only call it
// once Available() has reported true.
func (p *Port) ReadByte() byte {
	time.Sleep(p.bitDelay / 2) // sample in the middle of each bit
	var b byte
	for i := 0; i < 8; i++ {
		time.Sleep(p.bitDelay)
		if p.rx.Get() {
			b |= 1 << uint(i)
		}
	}
	time.Sleep(p.bitDelay) // stop bit
	return b
}

// ReadFrame blocks until n bytes have been received, one at a time.
func (p *Port) ReadFrame(n int) []byte {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		for !p.Available() {
		}
		buf[i] = p.ReadByte()
	}
	return buf
}

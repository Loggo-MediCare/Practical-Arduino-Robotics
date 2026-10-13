// Port of Chapter02/SPI_MCP3008_adc/SPI_MCP3008_adc.ino.
package main

import (
	"machine"
	"time"

	"practicalarduinorobotics/tinygo/internal/numfmt"
)

const csPin = machine.D4
const mcpSetSgl = 1

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 1e6,
		LSBFirst:  false,
		Mode:      machine.Mode0,
	})
	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	csPin.High()

	for {
		value, _ := mcpReadChannel(0)
		machine.Serial.Write([]byte(numfmt.Int(int32(value))))
		if time.Now().UnixMilli()%100 == 0 {
			machine.Serial.Write([]byte("\t0\t1023"))
		}
		machine.Serial.Write([]byte("\r\n"))
	}
}

// mcpReadChannel performs a single-ended analog read of the given channel on
// an MCP3008 over SPI.
func mcpReadChannel(channel uint8) (int, error) {
	if channel > 7 {
		return -1, nil
	}
	csPin.Low()
	machine.SPI0.Transfer(0x01)
	msb, _ := machine.SPI0.Transfer(mcpSetSgl<<7 | channel<<4)
	msb &= 0b11
	lsb, _ := machine.SPI0.Transfer(0)
	csPin.High()
	return int(msb)<<8 | int(lsb), nil
}

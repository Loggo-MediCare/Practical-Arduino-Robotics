// Port of Chapter02/I2C_bno055_temperature/I2C_bno055_temperature.ino.
package main

import (
	"machine"
	"time"

	"practicalarduinorobotics/tinygo/internal/numfmt"
)

const bnoAddress = 0x28
const bnoTempRegister = 0x34

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	machine.I2C0.Configure(machine.I2CConfig{})
	time.Sleep(1000 * time.Millisecond)

	for {
		if time.Now().UnixMilli()%200 == 0 {
			machine.Serial.Write([]byte("BNO055 temperature [C]: "))
			machine.Serial.Write([]byte(numfmt.Int(int32(bnoReadTemp())) + "\r\n"))
		}
	}
}

// bnoReadTemp reads the BNO055's temperature register over I2C.
func bnoReadTemp() uint8 {
	var temp [1]byte
	machine.I2C0.Tx(bnoAddress, []byte{bnoTempRegister}, temp[:])
	return temp[0]
}

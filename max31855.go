package max31855

import (
	"errors"
	"fmt"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
)

// ErrOpenCircuit - Thermocouple is not connected
var ErrOpenCircuit error = errors.New("Thermocouple is not connected")

// ErrShortToGround - Short Circuit to Ground
var ErrShortToGround error = errors.New("Short Circuit to Ground")

// ErrShortToVcc - Short Circuit to Power
var ErrShortToVcc error = errors.New("Short Circuit to Power")

// ErrReadingValue - Error Reading Value
var ErrReadingValue error = errors.New("Error Reading Value")

// Dev - A handle to contain the SPI connection
type Dev struct {
	c spi.Conn
}

// New - Connects to the MAX31855
func New(p spi.Port) (*Dev, error) {
	c, err := p.Connect(5*physic.MegaHertz, spi.Mode0, 8)

	if err != nil {
		return nil, err
	}

	d := &Dev{
		c: c,
	}

	return d, nil
}

// Temp - contains the temperature at both ends of the thermcouple
type Temp struct {
	Thermocouple physic.Temperature
	Internal     physic.Temperature
}

// GetTemp - Gets the current temperature in Celcius
func (d *Dev) GetTemp() (Temp, error) {
	raw := make([]byte, 4)

	if err := d.c.Tx(nil, raw); err != nil {
		return Temp{}, err
	}

	if raw[3]&0x01 != 0 {
		return Temp{}, ErrOpenCircuit
	}

	if raw[3]&0x02 != 0 {
		return Temp{}, ErrShortToGround
	}

	if raw[3]&0x04 != 0 {
		return Temp{}, ErrShortToVcc
	}

	var internal physic.Temperature
	var thermocouple physic.Temperature

	thermocoupleWord := ((uint16(raw[0]) << 8) | uint16(raw[1])) >> 2
	thermocouple.Set(fmt.Sprintf("%fC", float64(int16(thermocoupleWord))*0.25))
	internalWord := ((uint16(raw[2]) << 8) | uint16(raw[3])) >> 4
	internal.Set(fmt.Sprintf("%fC", float64(int16(internalWord))*0.0625))

	temp := Temp{
		Internal:     internal,
		Thermocouple: thermocouple,
	}

	return temp, nil
}

package max31855

import (
	"encoding/binary"
	"errors"

	"periph.io/x/periph/conn/spi"

	. "periph.io/x/periph/conn/physic"
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
	c, err := p.Connect(5*MegaHertz, spi.Mode0, 8)

	if err != nil {
		return nil, err
	}

	d := &Dev{
		c: c,
	}

	return d, nil
}

// GetTemp - Gets the current temperature in nano Kelvins
// Can easily be converted to more useful units with .Celcius(), .Farenheit(), etc.
func (d *Dev) GetTemp() (Temperature, error) {
	raw := make([]byte, 4)

	if err := d.c.Tx(nil, raw); err != nil {
		return 0, err
	}

	if raw[3]&0x01 != 0 {
		return 0, ErrOpenCircuit
	}

	if raw[3]&0x02 != 0 {
		return 0, ErrShortToGround
	}

	if raw[3]&0x04 != 0 {
		return 0, ErrShortToVcc
	}

	v := int32(binary.BigEndian.Uint32(raw))

	if v&0x07 != 0 {
		return 0, ErrReadingValue
	}

	v >>= 18

	t := ZeroCelsius + Temperature(v/4)*Celsius

	return t, nil
}

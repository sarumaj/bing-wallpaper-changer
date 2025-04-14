package types

import (
	"fmt"
	"strconv"

	"github.com/spf13/pflag"
)

type Percent float32

var _ pflag.Value = (*Percent)(nil)

func (p Percent) String() string {
	return fmt.Sprintf("%.2f", p)
}

func (p Percent) Float32() float32 {
	return float32(p)
}

func (p *Percent) Set(value string) error {
	f, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return err
	}

	if f < 0.0 || f > 100.0 {
		return fmt.Errorf("percent value must be between 0.0 and 100.0")
	}

	*p = Percent(f)
	return nil
}

func (p Percent) Type() string {
	return "float64"
}

package types

import (
	"fmt"

	"github.com/spf13/pflag"
)

var _ pflag.Value = &Enum[string, []string]{}

// Enum represents an enumeration flag.
type Enum[K comparable, L ~[]K] struct {
	value  K
	values L
	alias  func(K) string
}

// Set sets the enum value from the given string.
func (e *Enum[K, L]) Set(value string) error {
	for _, v := range e.values {
		if fmt.Sprint(v) == value || (e.alias != nil && e.alias(v) == value) {
			e.value = v
			return nil
		}
	}
	return fmt.Errorf("unknown value: %s", value)
}

// SetAlias sets the alias function of the enum.
func (e *Enum[K, L]) SetAlias(alias func(K) string) { e.alias = alias }

// SetDefault sets the default value of the enum.
func (e *Enum[K, L]) SetDefault(value K) { e.value = value }

// SetValues sets the values of the enum.
func (e *Enum[K, L]) SetValues(values ...K) { e.values = values }

// String returns the string representation of the enum value.
func (e Enum[K, L]) String() string {
	for _, v := range e.values {
		if v == e.value {
			if i, ok := any(v).(fmt.Stringer); ok {
				return i.String()
			}

			if e.alias != nil {
				return e.alias(v)
			}

			return fmt.Sprint(v)
		}
	}
	return "unknown"
}

// Type returns the type of the enum.
func (e Enum[K, L]) Type() string { return fmt.Sprintf("Enum[%T]", e.value) }

// Value returns the enum value.
func (e Enum[K, L]) Value() K { return e.value }

// Values returns the enum values.
func (e Enum[K, L]) Values() L { return e.values }

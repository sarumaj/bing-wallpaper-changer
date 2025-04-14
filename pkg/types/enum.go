package types

import (
	"encoding/json"
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

// MarshalJSON marshals the enum value to JSON.
func (e Enum[K, L]) MarshalJSON() ([]byte, error) {
	var aux = struct {
		Value  K      `json:"value"`
		Values L      `json:"values"`
		Alias  string `json:"alias"`
	}{
		Value:  e.value,
		Values: e.values,
	}

	if e.alias != nil {
		aux.Alias = e.alias(e.value)
	}

	return json.Marshal(aux)
}

// UnmarshalJSON unmarshals the enum value from JSON.
func (e *Enum[K, L]) UnmarshalJSON(data []byte) error {
	var aux struct {
		Value  K      `json:"value"`
		Values L      `json:"values"`
		Alias  string `json:"alias"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.value = aux.Value
	e.values = aux.Values

	if aux.Alias != "" {
		e.alias = func(v K) string {
			return aux.Alias
		}
	}

	return nil
}

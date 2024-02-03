package types

import "fmt"

const (
	Today Day = iota
	Yesterday
	TheDayBeforeYesterday
	ThreeDaysAgo
	FourDaysAgo
	FiveDaysAgo
	SixDaysAgo
	SevenDaysAgo
)

// Day is an enum type for relative days.
// 0 is today, 1 is yesterday, 2 is the day before yesterday, and so on.
// 7 is the highest value, which is seven days ago.
type Day int

// IsValid returns true if the Day is valid.
func (d Day) IsValid() error {
	if d < Today || d > SevenDaysAgo {
		return fmt.Errorf("invalid day: %d, allowed values are 0 to 7", d)
	}

	return nil
}

// String returns the string representation of the Day.
func (d Day) String() string {
	switch d {
	case Today:
		return "Today"

	case Yesterday:
		return "Yesterday"

	case TheDayBeforeYesterday:
		return "The day before yesterday"

	case ThreeDaysAgo:
		return "Three days ago"

	case FourDaysAgo:
		return "Four days ago"

	case FiveDaysAgo:
		return "Five days ago"

	case SixDaysAgo:
		return "Six days ago"

	case SevenDaysAgo:
		return "Seven days ago"

	default:
		return "Unknown"
	}
}

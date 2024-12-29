package types

import "fmt"

const (
	DayToday Day = iota
	DayYesterday
	DayTheDayBeforeYesterday
	DayThreeDaysAgo
	DayFourDaysAgo
	DayFiveDaysAgo
	DaySixDaysAgo
	DaySevenDaysAgo
)

// Day is an enum type for relative days.
// 0 is today, 1 is yesterday, 2 is the day before yesterday, and so on.
// 7 is the highest value, which is seven days ago.
type Day int

// IsValid returns true if the Day is valid.
func (d Day) IsValid() error {
	if d < DayToday || d > DaySevenDaysAgo {
		return fmt.Errorf("invalid day: %d, allowed values are 0 to 7", d)
	}

	return nil
}

// String returns the string representation of the Day.
func (d Day) String() string {
	s, ok := map[Day]string{
		DayToday:                 "Today",
		DayYesterday:             "Yesterday",
		DayTheDayBeforeYesterday: "The day before yesterday",
		DayThreeDaysAgo:          "Three days ago",
		DayFourDaysAgo:           "Four days ago",
		DayFiveDaysAgo:           "Five days ago",
		DaySixDaysAgo:            "Six days ago",
		DaySevenDaysAgo:          "Seven days ago",
	}[d]
	if !ok {
		return "Unknown"
	}

	return s
}

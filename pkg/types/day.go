package types

import (
	"fmt"
	"strings"
)

var AllowedDays = Days{DayToday, Day1Ago, Day2Ago, Day3Ago, Day4Ago, Day5Ago, Day6Ago, Day7Ago}

const (
	DayToday Day = iota
	Day1Ago
	Day2Ago
	Day3Ago
	Day4Ago
	Day5Ago
	Day6Ago
	Day7Ago
)

// Day is an enum type for relative days.
// 0 is today, 1 is yesterday, 2 is the day before yesterday, and so on.
// 7 is the highest value, which is seven days ago.
type Day int

// String returns the string representation of the Day.
func (d Day) String() string {
	if d > 7 || d < 0 {
		return "unknown"
	}

	if d == 0 {
		return "today"
	}

	return fmt.Sprintf("%d days ago", d)
}

// Days is a slice of Day.
type Days []Day

// Contains checks if the Days contains the given Day.
func (d Days) Contains(day Day) bool {
	for _, v := range d {
		if v == day {
			return true
		}
	}

	return false
}

// String returns the string representation of the Days.
func (d Days) String() string {
	var s []string
	for _, v := range d {
		s = append(s, v.String())
	}

	return strings.Join(s, ", ")
}

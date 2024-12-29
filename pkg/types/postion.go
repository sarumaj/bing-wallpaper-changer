package types

import "strings"

const (
	PositionTopLeft Position = iota
	PositionTopRight
	PositionBottomLeft
	PositionBottomRight
	PositionTopCenter
	PositionBottomCenter
	PositionCenterLeft
	PositionCenterRight
	PositionCenter
)

// Position is an enum type for relative positions.
type Position int

// String returns the string representation of the Position.
func (p Position) String() string {
	s, ok := map[Position]string{
		PositionTopLeft:      "TopLeft",
		PositionTopRight:     "TopRight",
		PositionBottomLeft:   "BottomLeft",
		PositionBottomRight:  "BottomRight",
		PositionTopCenter:    "TopCenter",
		PositionBottomCenter: "BottomCenter",
		PositionCenterLeft:   "CenterLeft",
		PositionCenterRight:  "CenterRight",
		PositionCenter:       "Center",
	}[p]
	if !ok {
		return "Unknown"
	}

	return s
}

// Positions is a slice of Position.
type Positions []Position

// String returns the string representation of the Positions.
func (p Positions) String() string {
	var s []string
	for _, v := range p {
		s = append(s, v.String())
	}

	return strings.Join(s, ", ")
}

package types

import "strings"

const (
	TopLeft Position = iota
	TopRight
	BottomLeft
	BottomRight
	TopCenter
	BottomCenter
	CenterLeft
	CenterRight
	Center
)

// Position is an enum type for relative positions.
type Position int

// String returns the string representation of the Position.
func (p Position) String() string {
	switch p {
	case TopLeft:
		return "TopLeft"
	case TopRight:
		return "TopRight"
	case BottomLeft:
		return "BottomLeft"
	case BottomRight:
		return "BottomRight"
	case TopCenter:
		return "TopCenter"
	case BottomCenter:
		return "BottomCenter"
	case CenterLeft:
		return "CenterLeft"
	case CenterRight:
		return "CenterRight"
	case Center:
		return "Center"
	default:
		return "Unknown"
	}
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

package mysqlroundrobinconnector

import (
	"time"
)

// UnknownLocationsErr indicate given addr (name) not found in register locations.
type UnknownLocationsErr struct {
	Name string
}

func (e *UnknownLocationsErr) Error() string {
	return "[UnknownLocationsErr: " + e.Name + "]"
}

// EmptyLocationsErr indicate given locations for register is empty.
type EmptyLocationsErr struct {
	Name string
}

func (e *EmptyLocationsErr) Error() string {
	return "[EmptyLocationsErr: " + e.Name + "]"
}

// LocationSyntaxErr indicate failed on parsing given location address.
type LocationSyntaxErr struct {
	LocationText string
}

func (e *LocationSyntaxErr) Error() string {
	return "[LocationSyntaxErr: " + e.LocationText + "]"
}

// TimeoutErr indicate operation is timeout.
type TimeoutErr struct {
	ReferenceTime time.Time
	DeadlineTime  time.Time
}

func (e *TimeoutErr) Error() string {
	return "[TimeoutErr: reference=" + e.ReferenceTime.String() + "; deadline=" + e.DeadlineTime.String() + "]"
}

// DialsErr indicate all attempt of dialing are failed.
type DialsErr struct {
	Errors []error
}

func (e *DialsErr) append(err error) {
	e.Errors = append(e.Errors, err)
}

func (e *DialsErr) Error() string {
	result := "[DialsErr: "
	for idx, err := range e.Errors {
		if idx != 0 {
			result += "; " + err.Error()
		} else {
			result += err.Error()
		}
	}
	result += "]"
	return result
}

package mysqlroundrobinconnector

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

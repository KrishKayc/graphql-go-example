package booker

//NoCabsError returned when there are no nearby cabs for the pick up location
type NoCabsError struct {
}

func (e *NoCabsError) Error() string {
	return "No nearby cabs currently. Please refresh"
}

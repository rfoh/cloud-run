package entity

type Location struct {
	City  string
	State string
}

func NewLocation(city, state string) *Location {
	return &Location{
		City:  city,
		State: state,
	}
}

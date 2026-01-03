package entity

import "fmt"

type CEPNotFoundError struct {
	CEP string
}

func (e *CEPNotFoundError) Error() string {
	return fmt.Sprintf("CEP %s not found", e.CEP)
}

type InvalidCEPError struct {
	Message string
}

func (e *InvalidCEPError) Error() string {
	return fmt.Sprintf("invalid CEP: %s", e.Message)
}

type TemperatureError struct {
	City    string
	Message string
}

func (e *TemperatureError) Error() string {
	return fmt.Sprintf("error fetching temperature for %s: %s", e.City, e.Message)
}

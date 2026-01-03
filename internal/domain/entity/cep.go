package entity

import (
	"regexp"
	"strings"
)

type CEP struct {
	value string
}

func NewCEP(cep string) (*CEP, error) {
	if err := validateCEP(cep); err != nil {
		return nil, err
	}

	cleanCEP := strings.ReplaceAll(cep, "-", "")
	return &CEP{value: cleanCEP}, nil
}

func (c *CEP) Value() string {
	return c.value
}

func (c *CEP) String() string {
	if len(c.value) == 8 {
		return c.value[:5] + "-" + c.value[5:]
	}
	return c.value
}

func validateCEP(cep string) error {
	cep = strings.ReplaceAll(cep, "-", "")

	if len(cep) != 8 {
		return &InvalidCEPError{Message: "CEP deve ter 8 dígitos"}
	}

	if !regexp.MustCompile(`^\d{8}$`).MatchString(cep) {
		return &InvalidCEPError{Message: "CEP deve conter apenas números"}
	}

	return nil
}

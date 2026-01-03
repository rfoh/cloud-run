package entity_test

import (
	"testing"

	"rfoh/cloud-run/internal/domain/entity"

	"github.com/stretchr/testify/assert"
)

func TestNewCEP_Valid(t *testing.T) {
	cep, err := entity.NewCEP("01310100")

	assert.Nil(t, err)
	assert.NotNil(t, cep)
	assert.Equal(t, "01310100", cep.Value())
	assert.Equal(t, "01310-100", cep.String())
}

func TestNewCEP_ValidWithHyphen(t *testing.T) {
	cep, err := entity.NewCEP("01310-100")

	assert.Nil(t, err)
	assert.NotNil(t, cep)
	assert.Equal(t, "01310100", cep.Value())
	assert.Equal(t, "01310-100", cep.String())
}

func TestNewCEP_InvalidLength(t *testing.T) {
	cep, err := entity.NewCEP("12345")

	assert.Error(t, err)
	assert.Nil(t, cep)
	assert.Contains(t, err.Error(), "8 dígitos")
}

func TestNewCEP_InvalidCharacters(t *testing.T) {
	cep, err := entity.NewCEP("0131010a")

	assert.Error(t, err)
	assert.Nil(t, cep)
	assert.Contains(t, err.Error(), "números")
}

func TestNewCEP_Empty(t *testing.T) {
	cep, err := entity.NewCEP("")

	assert.Error(t, err)
	assert.Nil(t, cep)
}

func TestNewTemperature(t *testing.T) {
	temp := entity.NewTemperature(25.0)

	assert.NotNil(t, temp)
	assert.Equal(t, 25.0, temp.Celsius)
	assert.Equal(t, 77.0, temp.Fahrenheit)
	assert.Equal(t, 298.0, temp.Kelvin)
}

func TestNewTemperature_Zero(t *testing.T) {
	temp := entity.NewTemperature(0.0)

	assert.NotNil(t, temp)
	assert.Equal(t, 0.0, temp.Celsius)
	assert.Equal(t, 32.0, temp.Fahrenheit)
	assert.Equal(t, 273.0, temp.Kelvin)
}

func TestNewTemperature_Negative(t *testing.T) {
	temp := entity.NewTemperature(-10.0)

	assert.NotNil(t, temp)
	assert.Equal(t, -10.0, temp.Celsius)
	assert.Equal(t, 14.0, temp.Fahrenheit)
	assert.Equal(t, 263.0, temp.Kelvin)
}

func TestNewTemperature_HighValue(t *testing.T) {
	temp := entity.NewTemperature(40.5)

	assert.NotNil(t, temp)
	assert.Equal(t, 40.5, temp.Celsius)
	assert.Equal(t, 104.9, temp.Fahrenheit)
	assert.Equal(t, 313.5, temp.Kelvin)
}

func TestNewLocation(t *testing.T) {
	location := entity.NewLocation("São Paulo", "SP")

	assert.NotNil(t, location)
	assert.Equal(t, "São Paulo", location.City)
	assert.Equal(t, "SP", location.State)
}

func TestNewLocation_Empty(t *testing.T) {
	location := entity.NewLocation("", "")

	assert.NotNil(t, location)
	assert.Equal(t, "", location.City)
	assert.Equal(t, "", location.State)
}

package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"rfoh/cloud-run/internal/application/dto"
	"rfoh/cloud-run/internal/domain/entity"

	httpAdapter "rfoh/cloud-run/internal/adapter/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGetTemperatureByCEPUseCase struct {
	mock.Mock
}

func (m *MockGetTemperatureByCEPUseCase) Execute(input *dto.GetTemperatureByCEPInput) (*dto.GetTemperatureByCEPOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GetTemperatureByCEPOutput), args.Error(1)
}

func TestWeatherHandler_HandleSuccess(t *testing.T) {
	mockUseCase := new(MockGetTemperatureByCEPUseCase)
	expected := &dto.GetTemperatureByCEPOutput{
		TempC: 28.5,
		TempF: 28.5,
		TempK: 28.5,
	}

	mockUseCase.On("Execute", mock.MatchedBy(func(input *dto.GetTemperatureByCEPInput) bool {
		return input.CEP == "01310100"
	})).Return(expected, nil).Once()

	handler := httpAdapter.NewWeatherHandler(mockUseCase)
	req := httptest.NewRequest("GET", "/?cep=01310100", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), `"temp_C":28.5`)
	assert.Contains(t, w.Body.String(), `"temp_F":28.5`)
	assert.Contains(t, w.Body.String(), `"temp_K":28.5`)
	mockUseCase.AssertExpectations(t)
}

func TestWeatherHandler_HandleInvalidCEP(t *testing.T) {
	mockUseCase := new(MockGetTemperatureByCEPUseCase)

	mockUseCase.On("Execute", mock.MatchedBy(func(input *dto.GetTemperatureByCEPInput) bool {
		return input.CEP == "12345"
	})).Return(nil, &entity.InvalidCEPError{Message: "CEP must have 8 digits"}).Once()

	handler := httpAdapter.NewWeatherHandler(mockUseCase)
	req := httptest.NewRequest("GET", "/?cep=12345", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "invalid zipcode")
	mockUseCase.AssertExpectations(t)
}

func TestWeatherHandler_HandleNotFoundCEP(t *testing.T) {
	mockUseCase := new(MockGetTemperatureByCEPUseCase)

	mockUseCase.On("Execute", mock.MatchedBy(func(input *dto.GetTemperatureByCEPInput) bool {
		return input.CEP == "00000000"
	})).Return(nil, &entity.CEPNotFoundError{CEP: "00000000"}).Once()

	handler := httpAdapter.NewWeatherHandler(mockUseCase)
	req := httptest.NewRequest("GET", "/?cep=00000000", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "can not find zipcode")
	mockUseCase.AssertExpectations(t)
}

func TestWeatherHandler_HandleMissingCEP(t *testing.T) {
	mockUseCase := new(MockGetTemperatureByCEPUseCase)

	handler := httpAdapter.NewWeatherHandler(mockUseCase)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "missing cep parameter")
	mockUseCase.AssertNotCalled(t, "Execute")
}

func TestWeatherHandler_HandleEmptyCEP(t *testing.T) {
	mockUseCase := new(MockGetTemperatureByCEPUseCase)

	handler := httpAdapter.NewWeatherHandler(mockUseCase)
	req := httptest.NewRequest("GET", "/?cep=", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing cep parameter")
	mockUseCase.AssertNotCalled(t, "Execute")
}

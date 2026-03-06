//go:build e2e
// +build e2e

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var appURL string

func init() {
	appURL = os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8080"
	}
}

func waitForApp(maxAttempts int) error {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	for i := 0; i < maxAttempts; i++ {
		resp, err := client.Get(appURL)
		if err == nil && resp.StatusCode != http.StatusNotFound {
			resp.Body.Close()
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("application not ready after %d attempts", maxAttempts)
}

func makeRequest(t *testing.T, cep string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("%s/?cep=%s", appURL, cep)
	return client.Get(url)
}

func TestE2E_HealthCheck(t *testing.T) {
	// Wait for app to be ready
	err := waitForApp(10)
	require.NoError(t, err, "Application should be ready")
}

func TestE2E_ValidCEP_Maceio(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/?cep=57052710", appURL)

	resp, err := client.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, _ := io.ReadAll(resp.Body)
	var result map[string]float64
	json.Unmarshal(body, &result)

	assert.Greater(t, result["temp_C"], 0.0)
	assert.Greater(t, result["temp_F"], 0.0)
	assert.Greater(t, result["temp_K"], 0.0)

	t.Logf("Maceió - Celsius: %.2f, Fahrenheit: %.2f, Kelvin: %.2f",
		result["temp_C"], result["temp_F"], result["temp_K"])
}

func TestE2E_InvalidCEP_Format(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/?cep=12345", appURL)

	resp, err := client.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	json.Unmarshal(body, &result)

	assert.Equal(t, "invalid zipcode", result["message"])
}

func TestE2E_NonexistentCEP(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/?cep=99999999", appURL)

	resp, err := client.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	json.Unmarshal(body, &result)

	assert.Equal(t, "can not find zipcode", result["message"])
}

func TestE2E_MissingCEPParameter(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := appURL

	resp, err := client.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	json.Unmarshal(body, &result)

	assert.Equal(t, "missing cep parameter", result["message"])
}

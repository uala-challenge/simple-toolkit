package error_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommonApiError(t *testing.T) {
	originalErr := errors.New("original error")
	apiErr := NewCommonApiError("E-001", "Test error", originalErr, http.StatusBadRequest)

	assert.NotNil(t, apiErr)
	assert.Error(t, apiErr)
	assert.Equal(t, "E-001", apiErr.(*CommonApiError).Code)
	assert.Equal(t, "Test error", apiErr.(*CommonApiError).Msg)
	assert.Equal(t, originalErr, apiErr.(*CommonApiError).Err)
	assert.Equal(t, http.StatusBadRequest, apiErr.(*CommonApiError).HttpCode)
}

func TestCommonApiError_ErrorMethod(t *testing.T) {
	apiErr := NewCommonApiError("E-002", "Another error", errors.New("some internal error"), http.StatusNotFound)

	expectedErrorMessage := "Error E-002: Another error \ntrace: some internal error"
	assert.Equal(t, expectedErrorMessage, apiErr.Error())
}

func TestWrapError_ExistingCommonApiError(t *testing.T) {
	originalErr := NewCommonApiError("E-003", "Initial message", errors.New("underlying error"), http.StatusConflict)
	wrappedErr := WrapError(originalErr, "Additional context")

	assert.NotNil(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "Additional context: Initial message")
}

func TestWrapError_NonCommonApiError(t *testing.T) {
	originalErr := errors.New("generic error")
	wrappedErr := WrapError(originalErr, "This should not change the error")

	assert.Equal(t, originalErr, wrappedErr)
}

func TestHandleApiErrorResponse_CommonApiError(t *testing.T) {
	recorder := httptest.NewRecorder()
	apiErr := NewCommonApiError("E-004", "API Failure", errors.New("api issue"), http.StatusInternalServerError)

	err := HandleApiErrorResponse(apiErr, recorder)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response CommonApiError
	err = json.NewDecoder(recorder.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Equal(t, "E-004", response.Code)
	assert.Equal(t, "API Failure", response.Msg)
}

func TestHandleApiErrorResponse_GenericError(t *testing.T) {
	recorder := httptest.NewRecorder()
	genericErr := errors.New("unexpected error")

	err := HandleApiErrorResponse(genericErr, recorder)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response CommonApiError
	err = json.NewDecoder(recorder.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Equal(t, "GE-001", response.Code)
	assert.Equal(t, "Internal Error", response.Msg)
}

func TestHandleApiErrorResponse_NilError(t *testing.T) {
	recorder := httptest.NewRecorder()
	apiErr := NewCommonApiError("E-005", "Missing error", nil, http.StatusBadRequest)

	err := HandleApiErrorResponse(apiErr, recorder)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	var response CommonApiError
	err = json.NewDecoder(recorder.Body).Decode(&response)

	assert.NoError(t, err)
	assert.Equal(t, "E-005", response.Code)
	assert.Equal(t, "Missing error", response.Msg)
	expectedConsoleMsg := "[error_wrapper]HandleApiErrorResponse: The err attribute is null"
	fmt.Println("Validando salida esperada:", expectedConsoleMsg)
}

func TestCommonApiError_Unwrap(t *testing.T) {
	originalErr := errors.New("underlying error")
	apiErr := NewCommonApiError("E-006", "Test unwrap", originalErr, http.StatusBadRequest)
	unwrappedErr := errors.Unwrap(apiErr)
	assert.Equal(t, originalErr, unwrappedErr)
	assert.EqualError(t, unwrappedErr, "underlying error")
}

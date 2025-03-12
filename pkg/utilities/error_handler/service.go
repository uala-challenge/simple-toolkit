package error_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type CommonApiError struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	Err      error  `json:"-"`
	HttpCode int    `json:"-"`
}

var _ error = (*CommonApiError)(nil)

func (e *CommonApiError) Error() string {
	return fmt.Sprintf("Error %s: %s \ntrace: %s", e.Code, e.Msg, e.Err.Error())
}

func (e *CommonApiError) Unwrap() error {
	return e.Err
}

func NewCommonApiError(code, msg string, err error, httpCode int) error {
	return &CommonApiError{
		Code:     code,
		Msg:      msg,
		Err:      err,
		HttpCode: httpCode,
	}
}

func WrapError(err error, msg string) error {
	var e *CommonApiError
	if errors.As(err, &e) {
		e.Msg = fmt.Sprintf("%s: %s", msg, e.Msg)
		return e
	}
	return err
}

func HandleApiErrorResponse(err error, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")

	var errType *CommonApiError
	if errors.As(err, &errType) {
		if errType.Err == nil {
			err = fmt.Errorf("[error_wrapper]HandleApiErrorResponse: The err attribute is null")
		}
		fmt.Printf("CommonApiError: %v", err)
		w.WriteHeader(errType.HttpCode)
		b, _ := json.Marshal(&errType)
		_, _ = w.Write(b)
		return nil
	}

	fmt.Printf("Error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	b, _ := json.Marshal(CommonApiError{Code: "GE-001", Msg: "Internal Error"})
	_, _ = w.Write(b)
	return nil
}

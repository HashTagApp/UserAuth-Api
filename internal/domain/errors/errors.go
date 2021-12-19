package errors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

const (
	applicationError = `APPLICATION`
	domainError      = `DOMAIN`
	validationError  = `VALIDATION`
)

type DomainError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Line    int    `json:"line"`
	File    string `json:"file"`
	Trace   string `json:"trace"`
}

func (e DomainError) Error() string {
	return fmt.Sprintf("Error Occurred %d", e.Code)
}

type ApplicationError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Line    int    `json:"line"`
	File    string `json:"file"`
	Trace   string `json:"trace"`
}

func (e ApplicationError) Error() string {
	return fmt.Sprintf("Error Occurred %d", e.Code)
}

type ValidationError struct {
	Message interface{} `json:"message"`
	Code    int         `json:"code"`
	Line    int         `json:"line"`
	File    string      `json:"file"`
	Trace   string      `json:"trace"`
}

// ErrorResponseValidation ...
type ErrorResponseValidation struct {
	Errors struct {
		Message string      `json:"message"`
		Code    int         `json:"code"`
		Fields  interface{} `json:"fields"`
	} `json:"errors"`
}

// ErrorResponseDomain ...
type ErrorResponseDomain struct {
	Errors struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Line    int    `json:"line,omitempty"`
		File    string `json:"file,omitempty"`
		Trace   string `json:"trace,omitempty"`
	} `json:"errors"`
}

// ErrorResponseApplication ...
type ErrorResponseApplication struct {
	Errors struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Line    int    `json:"line,omitempty"`
		File    string `json:"file,omitempty"`
		Trace   string `json:"trace,omitempty"`
	} `json:"errors"`
}

type ErrorResponseAuthentication struct {
	Errors struct {
		Message interface{} `json:"message"`
		Code    int         `json:"code"`
	} `json:"errors"`
}

type AppOutDateError struct {
	Message interface{} `json:"message"`
	Code    int         `json:"code"`
	Type    string      `json:"type"`
}

type GeneralError struct {
	Code             int         `json:"code"`
	Line             int         `json:"line"`
	CorrelationId    string      `json:"correlationId"`
	Message          string      `json:"message"`
	DeveloperMessage string      `json:"developerMessage"`
	File             string      `json:"file"`
	Trace            string      `json:"trace"`
	Type             string      `json:"type"`
	Fields           interface{} `json:"fields"`
}

func (e GeneralError) Error() string {
	return "General Error"
}

func (e AppOutDateError) Error() string {
	return "App has outdated"
}

func (e ValidationError) Error() string {
	return "Validation Error"
}

type AuthenticationError struct {
	Message interface{} `json:"message"`
	Code    int         `json:"code"`
	Type    string      `json:"type"`
}

func (e AuthenticationError) Error() string {
	return "Authentication Error"
}

func filePath() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return file, line
	}

	return "", 0
}

func NewDomainError(message string, code int, trace string) DomainError {
	path, line := filePath()

	return DomainError{
		Message: message,
		Code:    code,
		Line:    line,
		File:    path,
		Trace:   trace,
	}
}

func NewApplicationError(message string, code int, trace string) ApplicationError {
	path, line := filePath()

	return ApplicationError{
		Message: message,
		Code:    code,
		Line:    line,
		File:    path,
		Trace:   trace,
	}
}

func NewValidationError(messages map[string]string, code int, trace string) ValidationError {
	path, line := filePath()

	ve := ValidationError{
		Message: messages,
		Code:    code,
		Line:    line,
		File:    path,
		Trace:   trace,
	}
	return ve
}

func NewAuthenticationError(messages interface{}, code int) AuthenticationError {

	return AuthenticationError{
		Message: messages,
		Code:    code,
	}
}

func NewAppOutDatedError(messages interface{}, code int) AppOutDateError {
	ve := AppOutDateError{
		Message: messages,
		Code:    code,
	}
	return ve
}

func NewApplicationV2Error(correlationId string, message string, developerMessage string, code int, trace string) GeneralError {
	path, line := filePath()

	return GeneralError{
		Code:             code,
		Line:             line,
		CorrelationId:    correlationId,
		Message:          message,
		DeveloperMessage: developerMessage,
		File:             path,
		Trace:            trace,
		Type:             "APPLICATION",
	}
}
func NewValidationV2Error(correlationId string, message string, developerMessage string, code int, trace string) GeneralError {
	path, line := filePath()

	return GeneralError{
		Code:             code,
		Line:             line,
		CorrelationId:    correlationId,
		Message:          message,
		DeveloperMessage: developerMessage,
		File:             path,
		Trace:            trace,
		Type:             "VALIDATION",
	}
}

func IsDomain(err error) bool {
	if reflect.TypeOf(err) == reflect.TypeOf(GeneralError{}) {
		return err.(GeneralError).Type == domainError
	}
	return reflect.TypeOf(err) == reflect.TypeOf(DomainError{})
}

func IsApplication(err error) bool {
	if reflect.TypeOf(err) == reflect.TypeOf(GeneralError{}) {
		return err.(GeneralError).Type == applicationError
	}
	return reflect.TypeOf(err) == reflect.TypeOf(ApplicationError{})
}

func IsValidationError(err error) bool {
	if reflect.TypeOf(err) == reflect.TypeOf(GeneralError{}) {
		return err.(GeneralError).Type == validationError
	}
	return reflect.TypeOf(err) == reflect.TypeOf(ValidationError{})
}

func IsAuthenticationError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(AuthenticationError{})
}

func IsAppOutDatedError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(AppOutDateError{})
}

// ErrorEncoder ...
func ErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	switch reflect.TypeOf(err) {
	case reflect.TypeOf(ApplicationError{}):
		encodeApplicationErrorResponse(ctx, err, w)

	case reflect.TypeOf(AuthenticationError{}):
		encodeAuthenticationErrorResponse(ctx, err, w)

	case reflect.TypeOf(DomainError{}):
		encodeDomainErrorResponse(ctx, err, w)

	case reflect.TypeOf(ValidationError{}):
		encodeValidationErrorResponse(ctx, err, w)

	case reflect.TypeOf(AppOutDateError{}):
		encodeAppOutDateErrorResponse(ctx, err, w)

	}
}

func encodeValidationErrorResponse(ctx context.Context, err error, w http.ResponseWriter) error {
	derr := err.(ValidationError)
	res := ErrorResponseValidation{}
	res.Errors.Message = "Validation Error"
	res.Errors.Fields = derr.Message
	res.Errors.Code = derr.Code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	return json.NewEncoder(w).Encode(res)

}

func encodeDomainErrorResponse(ctx context.Context, err error, w http.ResponseWriter) error {
	derr := err.(DomainError)
	res := ErrorResponseDomain{}
	w.Header().Set("Content-Type", "application/json")

	res.Errors.Code = derr.Code // 400
	res.Errors.Message = derr.Message

	return json.NewEncoder(w).Encode(res)
}

func encodeApplicationErrorResponse(ctx context.Context, err error, w http.ResponseWriter) error {
	aerr := err.(ApplicationError)
	res := ErrorResponseApplication{}
	res.Errors.Code = aerr.Code // 500
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(res)
}

func encodeAuthenticationErrorResponse(ctx context.Context, err error, w http.ResponseWriter) error {
	aerr := err.(AuthenticationError)
	res := ErrorResponseAuthentication{}
	res.Errors.Message = aerr.Message
	res.Errors.Code = aerr.Code //500
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(res)
}

func encodeAppOutDateErrorResponse(ctx context.Context, err error, w http.ResponseWriter) error {
	derr := err.(AppOutDateError)
	res := AppOutDateError{}
	res.Message = derr.Message
	res.Code = derr.Code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUpgradeRequired)
	return json.NewEncoder(w).Encode(res)
}

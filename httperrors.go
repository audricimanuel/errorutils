package errorutils

import (
	"errors"
	"net/http"
)

// ERROR HTTP
var (
	// 5XX
	ErrorInternalServer      = NewHttpError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR)
	ErrorNotImplemented      = NewHttpError(http.StatusNotImplemented, NOT_IMPLEMENTED)
	ErrorBadGateway          = NewHttpError(http.StatusBadGateway, BAD_GATEWAY)
	ErrorServiceNotAvailable = NewHttpError(http.StatusServiceUnavailable, SERVICE_UNAVAILABLE)

	// 4XX
	ErrorBadRequest       = NewHttpError(http.StatusBadRequest, BAD_REQUEST)
	ErrorInvalidPayload   = NewHttpError(http.StatusBadRequest, INVALID_PAYLOAD)
	ErrorUnauthorized     = NewHttpError(http.StatusUnauthorized, UNAUTHORIZED)
	ErrorPaymentRequired  = NewHttpError(http.StatusPaymentRequired, PAYMENT_REQUIRED)
	ErrorForbidden        = NewHttpError(http.StatusForbidden, NO_PERMISSION)
	ErrorNotFound         = NewHttpError(http.StatusNotFound, DATA_NOT_FOUND)
	ErrorMethodNotAllowed = NewHttpError(http.StatusMethodNotAllowed, METHOD_NOT_ALLOWED)
	ErrorRequestTimeout   = NewHttpError(http.StatusRequestTimeout, REQUEST_TIMEOUT)
	ErrorDuplicateData    = NewHttpError(http.StatusConflict, ALREADY_EXISTS)
	ErrorLengthRequired   = NewHttpError(http.StatusLengthRequired, MINIMUM_LENGTH_REQUIRED)
	ErrorMaxSize          = NewHttpError(http.StatusRequestEntityTooLarge, MAX_SIZE_EXCEEDED)
	ErrorLoginRequired    = NewHttpError(http.StatusUnauthorized, LOGIN_REQUIRED)
	ErrorTokenRequired    = NewHttpError(http.StatusUnauthorized, TOKEN_REQUIRED)
	ErrorTokenExpired     = NewHttpError(http.StatusUnauthorized, TOKEN_EXPIRED)
	ErrorInvalidToken     = NewHttpError(http.StatusUnauthorized, INVALID_TOKEN)

	// 2XX
	ErrorNoContent = NewHttpError(http.StatusNoContent, NO_CONTENT)
)

type (
	HttpError interface {
		Error() string
		CustomMessage(message string) *HttpErrorImpl
	}
	HttpErrorImpl struct {
		Status  int
		Message string
	}
)

// NewHttpError creating new custom error that implements error object, with both customized status code (HTTP status) and message.
func NewHttpError(status int, message string) HttpError {
	return &HttpErrorImpl{Status: status, Message: message}
}

// Error returning a string of error message
func (e *HttpErrorImpl) Error() string {
	return e.Message
}

// CustomMessage to define custom message when getting any error that implements HttpError
// Example:
//
//	func Abc() error {
//		errorMessage := "this is an example of error message"
//		log.Error(errorMessage)
//		return ErrorInternalServer.CustomMessage(errorMessage)
//	}
//
//	func Abc2() error {
//		if _, err := time.Parse("2006-01-02", "23 Aug 2024"); err != nil {
//			log.Error("error when parsing string to time object: ", err)
//			return ErrorInternalServer.CustomMessage(err.Error())
//		}
//		return nil
//	}
func (e *HttpErrorImpl) CustomMessage(message string) *HttpErrorImpl {
	e.Message = message
	return e
}

// ToHttpError converting any errors to HttpError
func ToHttpError(err error, statusCode int) HttpError {
	var result HttpErrorImpl
	if err != nil {
		result.Message = err.Error()
		result.Status = statusCode
	}
	return &result
}

// GetStatusCode to get status code from any errors.
//
//	Regardless it's a kind of HttpErrorImpl, or it will return 500 and its error by default if not.
//	Will return 200 and "success" message if error is nil
func GetStatusCode(err error) (int, string) {
	if err == nil {
		return http.StatusOK, SUCCESS
	}

	var httpErr *HttpErrorImpl
	if errors.As(err, &httpErr) {
		return httpErr.Status, httpErr.Message
	}

	return http.StatusInternalServerError, err.Error()
}

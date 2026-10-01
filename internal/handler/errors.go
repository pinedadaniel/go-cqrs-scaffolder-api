package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrCodeUnauthorized        = errors.New("UNAUTHORIZED")
	ErrCodeNotFound            = errors.New("NOT_FOUND")
	ErrCodeBadRequest          = errors.New("BAD_REQUEST")
	ErrCodeForbidden           = errors.New("FORBIDDEN")
	ErrCodeConflict            = errors.New("CONFLICT")
	ErrCodeUnprocessableEntity = errors.New("UNPROCESSABLE_ENTITY")
	ErrCodeTooManyRequests     = errors.New("TOO_MANY_REQUESTS")
	ErrCodeInternalServerError = errors.New("INTERNAL_SERVER_ERROR")
	ErrCodeBadGateway          = errors.New("BAD_GATEWAY")
	ErrCodeServiceUnavailable  = errors.New("SERVICE_UNAVAILABLE")
	ErrCodeGatewayTimeout      = errors.New("GATEWAY_TIMEOUT")
	ErrCodeCircuitBreakerOpen  = errors.New("CIRCUIT_BREAKER_OPEN")
)

type CustomWebError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *CustomWebError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Status, e.Code, e.Message)
}

func CreateWebError(status int, code, message string) *CustomWebError {
	return &CustomWebError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func GlobalErrorHandler(c *gin.Context, err error) *CustomWebError {

	var webErr *CustomWebError

	switch {

	// 4xx Client Errors
	case errors.Is(err, ErrCodeNotFound):
		webErr = CreateWebError(http.StatusNotFound, ErrCodeNotFound.Error(), err.Error())
	case errors.Is(err, ErrCodeBadRequest):
		webErr = CreateWebError(http.StatusBadRequest, ErrCodeBadRequest.Error(), err.Error())
	case errors.Is(err, ErrCodeUnauthorized):
		webErr = CreateWebError(http.StatusUnauthorized, ErrCodeUnauthorized.Error(), err.Error())
	case errors.Is(err, ErrCodeForbidden):
		webErr = CreateWebError(http.StatusForbidden, ErrCodeForbidden.Error(), err.Error())
	case errors.Is(err, ErrCodeNotFound):
		webErr = CreateWebError(http.StatusNotFound, ErrCodeNotFound.Error(), err.Error())
	case errors.Is(err, ErrCodeConflict):
		webErr = CreateWebError(http.StatusConflict, ErrCodeConflict.Error(), err.Error())
	case errors.Is(err, ErrCodeUnprocessableEntity):
		webErr = CreateWebError(http.StatusUnprocessableEntity, ErrCodeUnprocessableEntity.Error(), err.Error())
	case errors.Is(err, ErrCodeTooManyRequests):
		addRetryAfterHeader(c, err)
		webErr = CreateWebError(http.StatusTooManyRequests, ErrCodeTooManyRequests.Error(), err.Error())

	// 5xx Server Errors
	case errors.Is(err, ErrCodeInternalServerError):
		webErr = CreateWebError(http.StatusBadGateway, ErrCodeInternalServerError.Error(), err.Error())
	case errors.Is(err, ErrCodeBadGateway):
		webErr = CreateWebError(http.StatusBadGateway, ErrCodeBadGateway.Error(), err.Error())
	case errors.Is(err, ErrCodeServiceUnavailable):
		addRetryAfterHeader(c, err)
		webErr = CreateWebError(http.StatusServiceUnavailable, ErrCodeServiceUnavailable.Error(), err.Error())
	case errors.Is(err, ErrCodeGatewayTimeout):
		webErr = CreateWebError(http.StatusGatewayTimeout, ErrCodeGatewayTimeout.Error(), err.Error())

	// Circuit Breaker Errors
	case errors.Is(err, ErrCodeCircuitBreakerOpen):
		addRetryAfterHeader(c, err)
		webErr = CreateWebError(http.StatusServiceUnavailable, ErrCodeCircuitBreakerOpen.Error(), err.Error())

	default:
		webErr = CreateWebError(http.StatusInternalServerError, ErrCodeInternalServerError.Error(), err.Error())
	}

	return webErr
}

func addRetryAfterHeader(c *gin.Context, err error) {
	c.Get("Retry-after")
	c.Header("Retry-After", err.Error())
}

func WriteJSONError(c *gin.Context, webErr *CustomWebError) {
	c.JSON(webErr.Status, Response[any]{
		Success: false,
		Error: &APIError{
			Code:    webErr.Code,
			Message: webErr.Message,
		},
	})
	//c.Abort()
}

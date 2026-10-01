package health

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pinedadaniel/scaffolder-api-go/internal/handler"
)

const (
	AppName                      = "Scaffolder"
	ErrCodeCallerIDRequired      = AppName + "CALLER-ID-REQUIRED"
	ErrCodeDecodeJSONRequestBody = AppName + "DECODE-JSON-REQUEST-BODY"
)

var (
	ErrCallerIDRequired      = errors.New("header x-caller-id is required")
	ErrDecodeJSONRequestBody = errors.New("error decoding request body")
)

func handleError(ctx *gin.Context, err error) {
	var webErr *handler.CustomWebError

	switch {
	case errors.Is(err, ErrCallerIDRequired):
		webErr = handler.CreateWebError(http.StatusBadRequest, ErrCodeCallerIDRequired, err.Error())
	case errors.Is(err, ErrDecodeJSONRequestBody):
		webErr = handler.CreateWebError(http.StatusBadRequest, ErrCodeDecodeJSONRequestBody, err.Error())
	default:
		webErr = handler.GlobalErrorHandler(ctx, err)
	}

	handler.WriteJSONError(ctx, webErr)
}

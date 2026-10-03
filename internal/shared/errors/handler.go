package apperrors

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

const errorKey = "error"

// HTTPErrorHandler is a custom error handler for Echo framework.
func HTTPErrorHandler(c *echo.Context, err error) {
	if handled := handleValidationError(c, err); handled {
		return
	}

	if handled := handleApplicationError(c, err); handled {
		return
	}

	if handled := handleEchoError(c, err); handled {
		return
	}

	handleInternalError(c, err)
}

func handleValidationError(c *echo.Context, err error) bool {
	validationErr, ok := errors.AsType[*ValidationError](err)
	if !ok {
		return false
	}

	writeErrorResponse(c, http.StatusBadRequest, Error{
		Code:    CodeValidation,
		Message: "Request validation failed",
		Fields:  validationErr.Fields,
	})

	return true
}

func handleApplicationError(c *echo.Context, err error) bool {
	appErr, ok := errors.AsType[*Error](err)
	if !ok {
		return false
	}

	writeErrorResponse(c, statusCodeFor(appErr.Code), *appErr)

	return true
}

func handleEchoError(c *echo.Context, err error) bool {
	echoErr, ok := errors.AsType[*echo.HTTPError](err)
	if !ok {
		return false
	}

	writeErrorResponse(c, echoErr.Code, Error{
		Code:    codeForHTTPStatus(echoErr.Code),
		Message: echoErrorMessage(echoErr.Code),
	})

	return true
}

func handleInternalError(c *echo.Context, err error) {
	c.Logger().Error(err.Error())

	writeErrorResponse(c, http.StatusInternalServerError, Error{
		Code:    CodeInternalError,
		Message: "Internal server error",
	})
}

func writeErrorResponse(
	c *echo.Context,
	status int,
	appError Error,
) {
	if err := c.JSON(status, map[string]any{
		errorKey: appError,
	}); err != nil {
		c.Logger().Error(err.Error())
	}
}

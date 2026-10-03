package apperrors

import "net/http"

func statusCodeFor(code string) int {
	switch code {
	case CodeNotFound:
		return http.StatusNotFound

	case CodeConflict:
		return http.StatusConflict

	case CodeBadRequest:
		return http.StatusBadRequest

	case CodeValidation:
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}

func codeForHTTPStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return CodeBadRequest

	case http.StatusUnauthorized:
		return CodeUnauthorized

	case http.StatusForbidden:
		return CodeForbidden

	case http.StatusNotFound:
		return CodeNotFound

	case http.StatusConflict:
		return CodeConflict

	default:
		return CodeInternalError
	}
}

func echoErrorMessage(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "Bad request"

	case http.StatusUnauthorized:
		return "Unauthorized"

	case http.StatusForbidden:
		return "Forbidden"

	case http.StatusNotFound:
		return "Not found"

	case http.StatusMethodNotAllowed:
		return "Method not allowed"

	default:
		return "Request failed"
	}
}

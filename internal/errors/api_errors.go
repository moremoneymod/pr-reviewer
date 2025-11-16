package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorCode string

const (
	ErrorCodeTeamExists     ErrorCode = "TEAM_EXISTS"
	ErrorCodePRExists       ErrorCode = "PR_EXISTS"
	ErrorCodePRMerged       ErrorCode = "PR_MERGED"
	ErrorCodeNotAssigned    ErrorCode = "NOT_ASSIGNED"
	ErrorCodeNoCandidate    ErrorCode = "NO_CANDIDATE"
	ErrorCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrorCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrorCodeBadRequest     ErrorCode = "BAD_REQUEST"
	ErrorCodeInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
)

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(code ErrorCode, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetails{
			Code:    string(code),
			Message: message,
		},
	}
}

func ValidationError(errs validator.ValidationErrors) ErrorResponse {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return NewErrorResponse(ErrorCodeBadRequest, strings.Join(errMsgs, ","))
}

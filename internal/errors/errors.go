package server_errors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SError struct {
	Message   string
	Code      int
	ErrorCode int
}

func (e *SError) Error() string {
	return e.Message
}

func (e *SError) Unwrap() (int, gin.H) {
	return e.Code, gin.H{
		"error": e.Message,
		"code":  strconv.Itoa(e.ErrorCode),
	}
}

var (
	CredentialError            = SError{Code: http.StatusBadRequest, Message: "Credentials are not correct!", ErrorCode: 100}
	InternalError              = SError{Code: http.StatusInternalServerError, Message: "Internal error!", ErrorCode: 101}
	UserAlreadyExistsError     = SError{Code: http.StatusBadRequest, Message: "User already exists!", ErrorCode: 102}
	InvalidInput               = SError{Code: http.StatusBadRequest, Message: "Invalid request data or format!", ErrorCode: 103}
	InvalidToken               = SError{Code: http.StatusBadRequest, Message: "Invalid token", ErrorCode: 104}
	TokenMalformed             = SError{Code: http.StatusBadRequest, Message: "Token is malformed or tempered!", ErrorCode: 105}
	TokenExpired               = SError{Code: http.StatusBadRequest, Message: "Token expired!", ErrorCode: 106}
	TokenSignatureInvalid      = SError{Code: http.StatusBadRequest, Message: "Token signature is not correct", ErrorCode: 107}
	InvalidAuthorizationHeader = SError{Code: http.StatusUnauthorized, Message: "Invalid authorization format", ErrorCode: 108}
	Unauthorized               = SError{Code: http.StatusUnauthorized, Message: "You are not authorized", ErrorCode: 109}
	ItemNotFound               = SError{Code: http.StatusNotFound, Message: "Requested item does not exists!", ErrorCode: 110}
	UserNotFound               = SError{Code: http.StatusNotFound, Message: "Requested user does not exists!", ErrorCode: 111}
	EmptyUpdate                = SError{Code: http.StatusNotFound, Message: "No fields to update", ErrorCode: 112}
	AccountIsNotActive         = SError{Code: http.StatusForbidden, Message: "Requested account is not active", ErrorCode: 113}
	InvalidVerificationCode    = SError{Code: http.StatusBadRequest, Message: "Invalid verification code!", ErrorCode: 114}
	FileRequired               = SError{Code: http.StatusBadRequest, Message: "File is required", ErrorCode: 115}
	InvalidFileFormat          = SError{Code: http.StatusBadRequest, Message: "Only .png, .jpg and .jpeg files are allowed", ErrorCode: 116}
	CategoryNotFound           = SError{Code: http.StatusBadRequest, Message: "Requested category was not found!", ErrorCode: 116}
)

func ValidationErrorBuilder(errList *[]string) *SError {
	message := strings.Join(*errList, "\n")
	return &SError{
		Code:      http.StatusBadRequest,
		Message:   message,
		ErrorCode: 117,
	}
}

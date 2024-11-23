package server_errors

import "net/http"

type SError struct {
	Message string
	Code    int
}

func (e *SError) Error() string {
	return e.Message
}

func (e *SError) Unwrap() (int, map[string]string) {
	return e.Code, map[string]string{"error": e.Message}
}

var (
	CredentialError            = SError{Code: http.StatusBadRequest, Message: "Credentials are not correct!"}
	InternalError              = SError{Code: http.StatusInternalServerError, Message: "Internal error!"}
	UserAlreadyExistsError     = SError{Code: http.StatusBadRequest, Message: "User already exists!"}
	InvalidInput               = SError{Code: http.StatusBadRequest, Message: "Invalid request data or format!"}
	InvalidToken               = SError{Code: http.StatusBadRequest, Message: "Invalid token"}
	TokenMalformed             = SError{Code: http.StatusBadRequest, Message: "Token is malformed or tempered!"}
	TokenExpired               = SError{Code: http.StatusBadRequest, Message: "Token expired!"}
	TokenSignatureInvalid      = SError{Code: http.StatusBadRequest, Message: "Token signature is not correct"}
	InvalidAuthorizationHeader = SError{Code: http.StatusUnauthorized, Message: "Invalid authorization format"}
	Unauthorized               = SError{Code: http.StatusUnauthorized, Message: "You are not authorized"}
)

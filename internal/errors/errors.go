package server_errors

import "net/http"

type SError struct {
    Message     string
    Code        int
}

func (e *SError) Error() string {
    return e.Message
}

var CredentialError = SError{Code: http.StatusBadRequest, Message: "Credentials are not correct!"}
var InternalError = SError{Code: http.StatusInternalServerError, Message: "Internal error!"}

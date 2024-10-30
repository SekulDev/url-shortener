package apperrors

import "errors"

var InvalidUrlError = errors.New("invalid url")
var TooManyRequests = errors.New("too many requests")

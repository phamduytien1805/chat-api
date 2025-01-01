package http_utils

type ResponseCode int32

const (
	OK ResponseCode = iota * -1
	ERROR
	ERROR_VALIDATION
	ERROR_UNIQUE
)

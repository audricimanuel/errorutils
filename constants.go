package errorutils

// STATUS
const (
	SUCCESS                 = "success"
	INTERNAL_SERVER_ERROR   = "internal server error"
	NOT_IMPLEMENTED         = "not implemented"
	BAD_GATEWAY             = "bad gateway"
	SERVICE_UNAVAILABLE     = "service unavailable"
	DATA_NOT_FOUND          = "your requested item is not found"
	METHOD_NOT_ALLOWED      = "method not allowed"
	REQUEST_TIMEOUT         = "request timeout"
	BAD_REQUEST             = "bad request"
	PAYMENT_REQUIRED        = "payment required"
	MINIMUM_LENGTH_REQUIRED = "minimum length not exceeded"
	UNAUTHORIZED            = "unauthorized"
	LOGIN_REQUIRED          = "login required"
	TOKEN_REQUIRED          = "token required"
	INVALID_TOKEN           = "invalid token"
	TOKEN_EXPIRED           = "token expired"
	NO_CONTENT              = "no content"
	MAX_SIZE_EXCEEDED       = "maximum size exceeded"
	ALREADY_EXISTS          = "your item already exist"
	INVALID_PAYLOAD         = "invalid payload"
	NO_PERMISSION           = "you have no permission to access this resource"
)

// FORMAT TIME
const (
	FORMAT_DATE_DEFAULT     = "2006-01-02"
	FORMAT_DATETIME_DEFAULT = "2006-01-02 15:04:05"
)

package errorutils

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

func DefineSQLError(err error) HttpError {
	// if data is not found
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorNotFound
	}

	// if the assertion of dbErr fails, "ok" will be false
	if dbErr, ok := err.(*pq.Error); ok {
		switch dbErr.Code {
		case "23505":
			// duplicate data
			return ErrorDuplicateData
		case "22001":
			// max length exceeded
			return ErrorMaxSize
		default:
			return ErrorInternalServer.CustomMessage(dbErr.Error())
		}
	}
	return ErrorInternalServer.CustomMessage(err.Error())
}

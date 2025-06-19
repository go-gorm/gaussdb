package gaussdb

import (
	"encoding/json"

	"github.com/HuaweiCloudDeveloper/gaussdb-go/gaussdbconn"

	"gorm.io/gorm"
)

// The error codes to map GaussDB errors to gorm errors, here is the GaussDB error codes reference https://support.huaweicloud.com/gaussdb/index.html.
var errCodes = map[string]error{
	"23505": gorm.ErrDuplicatedKey,
	"23503": gorm.ErrForeignKeyViolated,
	"42703": gorm.ErrInvalidField,
	"23514": gorm.ErrCheckConstraintViolated,
}

type ErrMessage struct {
	Code     string
	Severity string
	Message  string
}

// Translate it will translate the error to native gorm errors.
func (dialector Dialector) Translate(err error) error {
	if pgErr, ok := err.(*gaussdbconn.GaussdbError); ok {
		if translatedErr, found := errCodes[pgErr.Code]; found {
			return translatedErr
		}
		return err
	}

	parsedErr, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		return err
	}

	var errMsg ErrMessage
	unmarshalErr := json.Unmarshal(parsedErr, &errMsg)
	if unmarshalErr != nil {
		return err
	}

	if translatedErr, found := errCodes[errMsg.Code]; found {
		return translatedErr
	}
	return err
}

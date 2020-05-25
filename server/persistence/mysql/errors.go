package mysql

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

var (
	errDuplicateEntry = errors.New(`mysql duplicate entry`)
)

func convertMysqlError(err error) error {
	if err == nil {
		return nil
	}
	switch t := err.(type) {
	case *mysql.MySQLError:
		switch t.Number {
		case 1062:
			return errDuplicateEntry
		}
	}
	return err
}

func IsLockWaitTimeout(err error) bool {
	if err == nil {
		return false
	}
	if merr, ok := err.(*mysql.MySQLError); ok {
		// Error 1205: Lock wait timeout exceeded; try restarting transaction
		return merr.Number == 1205
	}
	return false
}

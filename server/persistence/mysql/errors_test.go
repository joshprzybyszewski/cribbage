package mysql

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestConvertMysqlError(t *testing.T) {
	testCases := []struct {
		input       error
		returnsSame bool
		output      error
	}{{
		input:  nil,
		output: nil,
	}, {
		input:       errors.New(`whatever`),
		returnsSame: true,
	}, {
		input: &mysql.MySQLError{
			Number: 1205,
		},
		returnsSame: true,
	}, {
		input: &mysql.MySQLError{
			Number: 1062,
		},
		output: errDuplicateEntry,
	}}

	for _, tc := range testCases {
		actOutput := convertMysqlError(tc.input)
		if tc.returnsSame {
			assert.Same(t, tc.input, actOutput)
		} else {
			assert.Equal(t, tc.output, actOutput)
		}
	}
}

func TestIsLockWaitTimeout(t *testing.T) {
	assert.False(t, IsLockWaitTimeout(nil))
	assert.False(t, IsLockWaitTimeout(errors.New(`whodathunkit`)))
	assert.False(t, IsLockWaitTimeout(&mysql.MySQLError{}))

	assert.True(t, IsLockWaitTimeout(&mysql.MySQLError{
		Number: 1205,
	}))
}

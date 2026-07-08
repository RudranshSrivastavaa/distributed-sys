package errors

import "errors"


var ErrOptimisticLockConflict =
	errors.New("inventory was modified by another transaction")
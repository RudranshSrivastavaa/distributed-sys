package errors

import "errors"

var ErrTemporaryFailure = errors.New(
    "temporary email provider failure",
)

var ErrInvalidRecipient = errors.New(
    "invalid recipient",
)
var ErrDuplicateEvent = errors.New("duplicate event");
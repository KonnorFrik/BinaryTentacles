package usecase

import "errors"

// This errors safe for show to user. '.Error()' can be use as error message.
var (
	// ErrNoMarkets - found 0 markets for any reason.
	ErrNoMarkets = errors.New("no available markets")
	// ErrInvalidInput - got bad/corrupted input for any reason.
	ErrInvalidInput = errors.New("invalid input")
	// ErrForbidden - access denien for any reason.
	ErrForbidden = errors.New("fordibbed")
)

// This errors unsafe for show to user. Error message must be substituted.
var (
	// ErrInternal - error from any other package for any reason.
	ErrInternal = errors.New("internal")
)

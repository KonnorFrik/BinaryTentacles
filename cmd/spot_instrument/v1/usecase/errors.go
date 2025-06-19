package usecase

import "errors"

var (
	// ErrNoMarkets - found 0 markets for any reason.
	ErrNoMarkets = errors.New("no available markets")
	// ErrInternal - error from any other package for any reason.
	ErrInternal = errors.New("internal")
)

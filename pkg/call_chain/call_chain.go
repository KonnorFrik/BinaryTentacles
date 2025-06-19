package callChain

import "context"

type ChainElement func(context.Context) error

// CallChain - chain for call anything.
// If any elements return non-nil error - next elements will not be called.
type CallChain struct {
	elements []ChainElement
}

// New - create a new call chain with elements 'elems'.
func New(elems ...ChainElement) *CallChain {
	var chain CallChain
	chain.elements = elems
	return &chain
}

// Call - start call a chain.
// Return number of element who return non-nil error.
// If chain complete successfully return zero-value.
func (cc *CallChain) Call(ctx context.Context) (int, error) {
	for i, el := range cc.elements {
		if err := el(ctx); err != nil {
			return i + 1, err
		}
	}

	return 0, nil
}

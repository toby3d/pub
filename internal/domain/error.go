package domain

import (
	"fmt"

	"golang.org/x/xerrors"
)

// Error represent a custom error implementation with HTTP status codes support.
//
//nolint:tagliatelle
type Error struct {
	Description string        `json:"error_description,omitempty"`
	Frame       xerrors.Frame `json:"-"`
	Code        int           `json:"error"`
}

func (e Error) Error() string {
	return fmt.Sprint(e)
}

func (e Error) Format(f fmt.State, r rune) {
	xerrors.FormatError(e, f, r)
}

func (e Error) FormatError(p xerrors.Printer) error {
	p.Printf("%d: %s", e.Code, e.Description)

	if !p.Detail() {
		return e
	}

	e.Frame.Format(p)

	return nil
}

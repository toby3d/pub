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

func (err Error) Error() string {
	return fmt.Sprint(err)
}

func (err Error) Format(f fmt.State, r rune) {
	xerrors.FormatError(err, f, r)
}

func (err Error) FormatError(p xerrors.Printer) error {
	p.Printf("%d: %s", err.Code, err.Description)

	if !p.Detail() {
		return err
	}

	err.Frame.Format(p)

	return nil
}

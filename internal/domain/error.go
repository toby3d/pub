package domain

import (
	"fmt"

	"golang.org/x/xerrors"
)

//nolint: tagliatelle
type Error struct {
	Code        string        `json:"error"`
	Description string        `json:"error_description,omitempty"`
	Frame       xerrors.Frame `json:"-"`
}

func (err Error) Error() string {
	return fmt.Sprint(err)
}

func (err Error) Format(f fmt.State, r rune) {
	xerrors.FormatError(err, f, r)
}

func (err Error) FormatError(p xerrors.Printer) error {
	p.Printf("%s: %s", err.Code, err.Description)

	if !p.Detail() {
		return err
	}

	err.Frame.Format(p)

	return nil
}

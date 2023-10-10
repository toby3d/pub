package domain

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/xerrors"
)

// RSVP represent a enum status relation to some event.
type RSVP struct {
	rsvp string
}

var (
	RSVPUnd        RSVP = RSVP{}             // "und"
	RSVPInterested RSVP = RSVP{"interested"} // "interested"
	RSVPMaybe      RSVP = RSVP{"maybe"}      // "maybe"
	RSVPNo         RSVP = RSVP{"no"}         // "no"
	RSVPYes        RSVP = RSVP{"yes"}        // "yes"
)

var ErrRSVPSyntax error = Error{
	Description: fmt.Sprintf("got unsupported RSVP enum value, expect '%s', '%s', '%s' or '%s'", RSVPInterested,
		RSVPMaybe, RSVPNo, RSVPYes),
	Frame: xerrors.Caller(1),
	Code:  http.StatusBadRequest,
}

var stringsRSVPs = map[string]RSVP{
	RSVPInterested.rsvp: RSVPInterested,
	RSVPMaybe.rsvp:      RSVPMaybe,
	RSVPNo.rsvp:         RSVPNo,
	RSVPYes.rsvp:        RSVPYes,
}

func ParseRSVP(v string) (RSVP, error) {
	// NOTE(toby3d): Case-insensitive values, normalized to lowercase.
	if out, ok := stringsRSVPs[strings.ToLower(v)]; ok {
		return out, nil
	}

	return RSVPUnd, fmt.Errorf("cannot parse '%s' as RSVP enum: %w", v, ErrRSVPSyntax)
}

func (rsvp RSVP) String() string {
	if rsvp.rsvp != "" {
		return rsvp.rsvp
	}

	return "und"
}

func (rsvp RSVP) GoString() string {
	return "domain.RSVP(" + rsvp.String() + ")"
}

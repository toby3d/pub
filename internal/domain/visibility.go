package domain

// Visibility describes visibility enum.
//
// See: https://indieweb.org/Micropub-extensions#Visibility
type Visibility struct {
	visibility string
}

var (
	VisibilityUnd      = Visibility{}
	VisibilityPublic   = Visibility{"public"}   // "public"
	VisibilityUnlisted = Visibility{"unlisted"} // "unlisted"
	VisibilityPrivate  = Visibility{"private"}  // "private"
)

func (v Visibility) String() string {
	if v.visibility == "" {
		return "und"
	}

	return v.visibility
}

func (v Visibility) GoString() string {
	return "domain.Visibility(" + v.String() + ")"
}

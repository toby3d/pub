package domain

// See: https://indieweb.org/Micropub-extensions#Post_Status
type PostStatus struct {
	postStatus string
}

var (
	PostStatusUnd       = PostStatus{""}          // "und"
	PostStatusDraft     = PostStatus{"draft"}     // "draft"
	PostStatusPublished = PostStatus{"published"} // "published"
)

func (ps PostStatus) String() string {
	if ps.postStatus == "" {
		return "und"
	}

	return ps.postStatus
}

func (ps PostStatus) GoString() string {
	return "domain.PostStatus(" + ps.String() + ")"
}

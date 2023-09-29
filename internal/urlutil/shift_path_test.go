package urlutil_test

import (
	"testing"

	"source.toby3d.me/toby3d/pub/internal/urlutil"
)

func TestShiftPath(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		input, expHead, expTail string
	}{
		"root":    {input: "/", expHead: "", expTail: "/"},
		"file":    {input: "/foo", expHead: "foo", expTail: "/"},
		"dir":     {input: "/foo/", expHead: "foo", expTail: "/"},
		"dirfile": {input: "/foo/bar", expHead: "foo", expTail: "/bar"},
		"subdir":  {input: "/foo/bar/", expHead: "foo", expTail: "/bar"},
	} {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			head, tail := urlutil.ShiftPath(tc.input)

			if head != tc.expHead {
				t.Errorf("ShiftPath(%s) = '%s', want '%s'", tc.input, head, tc.expHead)
			}

			if tail != tc.expTail {
				t.Errorf("ShiftPath(%s) = '%s', want '%s'", tc.input, tail, tc.expTail)
			}
		})
	}
}

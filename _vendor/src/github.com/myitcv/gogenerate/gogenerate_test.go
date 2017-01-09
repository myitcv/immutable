// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

// Copyright (c) 2012 - 2016 modelogiq GmbH <www.modelogiq.com>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package gogenerate

import "testing"

func TestCommentRegexp(t *testing.T) {
	comm := "banana"
	r, err := commentRegex(comm)
	if err != nil {
		t.Fatalf("Expected call to result in no error")
	}

	checks := []struct {
		s string
		r bool
	}{
		{GoGeneratePrefix + " " + comm, true},
		{GoGeneratePrefix + " " + comm + " ", true},
		{GoGeneratePrefix + " " + comm + "  some arguments\n", true},
		{GoGeneratePrefix + " " + comm + "\n", true},
		{GoGeneratePrefix + " " + comm + "\t", false},
		{GoGeneratePrefix + " " + comm + "\t", false},
	}

	for _, c := range checks {
		v := r.MatchString(c.s)
		if v != c.r {
			t.Errorf("commentRegex /%v/.MatchString(%q) does not equal %v", r, c.s, c.r)
		}
	}
}

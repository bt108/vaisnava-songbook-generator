package generator

import (
	"strings"
	"testing"
)

func TestResolveFragments(t *testing.T) {
	fragments := map[string][]Entry{
		"vandeguro": []Entry{{Word: "vande", Gloss: "I praise"}, {Word: "guroḥ", Gloss: "of guru"}},
	}
	entries, err := resolveFragments([]Entry{{Word: "prāptasya", Gloss: "attained"}, {Fragment: "vandeguro"}}, fragments)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 || entries[1].Word != "vande" || entries[2].Gloss != "of guru" {
		t.Fatalf("resolved entries = %#v", entries)
	}
}

func TestResolveFragmentsKeepsBreakInsideFragment(t *testing.T) {
	fragments := map[string][]Entry{
		"vanderupa": []Entry{
			{Word: "vande", Gloss: "I praise"},
			{Word: "rūpa-", Gloss: "Rūpa-"},
			{Word: "sanātanau", Gloss: "Sanātana"},
			{Word: "raghu-", Gloss: "Raghu-"},
			{Word: "yugau", Gloss: "the-pair"},
			{Break: true},
			{Word: "śrī-", Gloss: "Śrī-"},
			{Word: "jīva-", Gloss: "Jīva-"},
			{Word: "gopālakau", Gloss: "Gopālaka"},
		},
	}
	entries, err := resolveFragments([]Entry{{Fragment: "vanderupa"}}, fragments)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 9 || !entries[5].Break || entries[6].Word != "śrī-" {
		t.Fatalf("resolved entries = %#v, want break preserved inside fragment", entries)
	}
}

func TestGlossGapMultiEntry(t *testing.T) {
	entries := []Entry{
		{Word: "vande", Gloss: "I praise"},
		{Word: "svādv-", Gloss: "sweet-"},
		{Word: "anna-", Gloss: "food-"},
		{Word: "tṛptān", Gloss: "satisfied"},
		{Break: true},
		{Word: "hari-", Gloss: "Hari-"},
		{Word: "bhakta-", Gloss: "devotee-"},
		{Word: "saṅghān", Gloss: "multitudes"},
	}
	cases := []struct {
		name string
		i    int
		want string
	}{
		{"plain word gets a thin space", 0, `\,`},
		{"gloss ending hyphen joins tight", 1, ""},
		{"gloss ending hyphen joins tight", 2, ""},
		{"entry before a break is last of its line", 3, ""},
		{"gloss ending hyphen joins tight", 5, ""},
		{"last entry of pada gets no space", 7, ""},
	}
	for _, tc := range cases {
		if got := glossGap(entries, tc.i); got != tc.want {
			t.Errorf("glossGap(entries, %d) = %q, want %q", tc.i, got, tc.want)
		}
	}
}

func TestResolveFragmentsRejectsUnknownFragment(t *testing.T) {
	_, err := resolveFragments([]Entry{{Fragment: "missing"}}, nil)
	if err == nil || !strings.Contains(err.Error(), `unknown fragment "missing"`) {
		t.Fatalf("error = %v, want unknown fragment error", err)
	}
}

package generator

import "testing"

func TestGlossGap(t *testing.T) {
	tests := []struct {
		name    string
		entries []Entry
		i       int
		want    string
	}{
		{
			name:    "normal entry followed by another",
			entries: []Entry{{Gloss: "feet"}, {Gloss: "of"}},
			i:       0,
			want:    `\,`,
		},
		{
			name:    "compound join (gloss ends with hyphen)",
			entries: []Entry{{Gloss: "Viṣṇu-"}, {Gloss: "feet"}},
			i:       0,
			want:    "",
		},
		{
			name:    "last entry in pada",
			entries: []Entry{{Gloss: "feet"}},
			i:       0,
			want:    "",
		},
		{
			name:    "entry followed by break",
			entries: []Entry{{Gloss: "feet"}, {Break: true}},
			i:       0,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := glossGap(tt.entries, tt.i)
			if got != tt.want {
				t.Errorf("glossGap() = %q, want %q", got, tt.want)
			}
		})
	}
}

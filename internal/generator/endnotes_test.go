package generator

import (
	"reflect"
	"testing"
)

func TestActiveEndnotes(t *testing.T) {
	n1, n2 := "note one", "note two"
	chapter := Chapter{
		Content: []Verse{
			{
				Title: "Śrī Śrī Gurv-aṣṭaka",
				Slokas: []Sloka{
					{Number: "1", Footnotes: []Footnote{{Term: "saṁsāra", Note: &n1}}},
					{Number: "2", Footnotes: []Footnote{{Term: "bhāj", Note: &n2}, {Term: "skip", Note: nil}}},
				},
			},
			{
				Title: "Another Song",
				Slokas: []Sloka{
					{Number: "3", Footnotes: []Footnote{{Term: "hari", Note: &n1}}},
				},
			},
		},
	}
	want := []Endnote{
		{VerseLabel: "Śrī Śrī Gurv-aṣṭaka v.1", Term: "saṁsāra", Note: "note one"},
		{VerseLabel: "v.2", Term: "bhāj", Note: "note two"},
		{VerseLabel: "Another Song v.3", Term: "hari", Note: "note one"},
	}
	got := chapter.ActiveEndnotes()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ActiveEndnotes() = %#v, want %#v", got, want)
	}
}

func TestActiveEndnotesSongTitleOnlyOnNoteBearingSlokas(t *testing.T) {
	n := "a note"
	chapter := Chapter{
		Content: []Verse{
			{
				Title: "Song Without First-Verse Notes",
				Slokas: []Sloka{
					{Number: "1", Footnotes: []Footnote{{Term: "skip", Note: nil}}},
					{Number: "2", Footnotes: []Footnote{{Term: "bhāj", Note: &n}}},
				},
			},
			{
				Title: "Song With Notes On First Sloka",
				Slokas: []Sloka{
					{Number: "1", Footnotes: []Footnote{{Term: "hari", Note: &n}}},
					{Number: "2", Footnotes: []Footnote{{Term: "skip", Note: nil}}},
				},
			},
		},
	}
	want := []Endnote{
		{VerseLabel: "Song Without First-Verse Notes v.2", Term: "bhāj", Note: "a note"},
		{VerseLabel: "Song With Notes On First Sloka v.1", Term: "hari", Note: "a note"},
	}
	got := chapter.ActiveEndnotes()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ActiveEndnotes() = %#v, want %#v", got, want)
	}
}

package generator

import (
	"testing"
)

func TestNumberLabel(t *testing.T) {
	cases := []struct {
		name   string
		number VerseNumber
		want   string
	}{
		{"plain number suppressed", "1", ""},
		{"another plain number suppressed", "2", ""},
		{"empty number suppressed", "", ""},
		{"non-numeric label rendered", "First Refrain", "First Refrain"},
	}
	for _, tc := range cases {
		if got := (Sloka{Number: tc.number}).NumberLabel(); got != tc.want {
			t.Errorf("%s: NumberLabel() = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestAuthorLabel(t *testing.T) {
	cases := []struct {
		name   string
		author string
		want   string
	}{
		{"known author kept", "Viśvanātha Cakravartī Ṭhākura", "Viśvanātha Cakravartī Ṭhākura"},
		{"blank suppressed", "", ""},
		{"unknown placeholder suppressed", "unknown", ""},
		{"unknown case-insensitive", "Unknown", ""},
		{"unknown with padding suppressed", " unknown ", ""},
	}
	for _, tc := range cases {
		if got := (Verse{Author: tc.author}).AuthorLabel(); got != tc.want {
			t.Errorf("%s: AuthorLabel() = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestSlokaShouldClearPage(t *testing.T) {
	trueValue, falseValue := true, false
	cases := []struct {
		name string
		in   *bool
		want bool
	}{
		{"nil defaults to false", nil, false},
		{"explicit true", &trueValue, true},
		{"explicit false", &falseValue, false},
	}
	for _, tc := range cases {
		if got := (Sloka{ClearPage: tc.in}).ShouldClearPage(); got != tc.want {
			t.Errorf("%s: Sloka{ClearPage=%v}.ShouldClearPage() = %v, want %v", tc.name, tc.in, got, tc.want)
		}
	}
}

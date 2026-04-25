package generator

import (
	"testing"
)

func TestMergeEntriesRespectsBreak(t *testing.T) {
	entries := []Entry{
		{Word: "radha-", Gloss: "Rādhā-"},
		{Break: true},
		{Word: "kṛṣṇa", Gloss: "Kṛṣṇa"},
	}
	merged := mergeEntries(entries)
	if len(merged) != 3 {
		t.Fatalf("mergeEntries produced %d entries, want 3: %#v", len(merged), merged)
	}
	if merged[0].Word != "radha-" {
		t.Errorf("first entry word = %q, want %q", merged[0].Word, "radha-")
	}
	if !merged[1].Break {
		t.Errorf("second entry should be a break, got %#v", merged[1])
	}
	if merged[2].Word != "kṛṣṇa" {
		t.Errorf("third entry word = %q, want %q", merged[2].Word, "kṛṣṇa")
	}
}

func TestMergeEntriesJoinsAcrossNormalEntries(t *testing.T) {
	entries := []Entry{
		{Word: "viṣṇu-", Gloss: "Viṣṇu-"},
		{Word: "pādāya", Gloss: "feet"},
	}
	merged := mergeEntries(entries)
	if len(merged) != 1 {
		t.Fatalf("mergeEntries produced %d entries, want 1: %#v", len(merged), merged)
	}
	if merged[0].Word != "viṣṇu-pādāya" {
		t.Errorf("merged word = %q, want %q", merged[0].Word, "viṣṇu-pādāya")
	}
	if merged[0].Gloss != "Viṣṇu-feet" {
		t.Errorf("merged gloss = %q, want %q", merged[0].Gloss, "Viṣṇu-feet")
	}
}

func TestMergeEntriesPassesBreakThrough(t *testing.T) {
	entries := []Entry{
		{Break: true},
		{Word: "kṛṣṇa", Gloss: "Kṛṣṇa"},
	}
	merged := mergeEntries(entries)
	if len(merged) != 2 {
		t.Fatalf("mergeEntries produced %d entries, want 2: %#v", len(merged), merged)
	}
	if !merged[0].Break {
		t.Errorf("first entry should be a break")
	}
}

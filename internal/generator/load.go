package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func loadChapter(dataDir, chapterID string, fragments map[string][]Entry) (Chapter, error) {
	path := filepath.Join(dataDir, "chapters", chapterID+".yaml")
	chapterData, err := loadYAML[chapterFile](path)
	if err != nil {
		return Chapter{}, err
	}

	chapter := Chapter{Title: chapterData.Title, Intro: chapterData.Intro}
	for _, item := range chapterData.Content {
		if item.Prose != "" {
			if item.Include != "" {
				return Chapter{}, fmt.Errorf("chapter %s: content item has both prose and include", chapterID)
			}
			chapter.Prose = append(chapter.Prose, item.Prose)
			continue
		}
		if item.Include == "" {
			return Chapter{}, fmt.Errorf("chapter %s: content item missing include", chapterID)
		}
		verse, err := loadVerse(dataDir, item.Include, fragments)
		if err != nil {
			return Chapter{}, err
		}
		if len(verse.Slokas) > 0 {
			verse.Slokas[0].SongTitle = verse.Title
			verse.Slokas[0].AuthorLabel = verse.AuthorLabel()
			verse.Slokas[0].SongLanguage = verse.Language
		}
		chapter.Content = append(chapter.Content, verse)
	}
	return chapter, nil
}

func loadVerse(dataDir, name string, fragments map[string][]Entry) (Verse, error) {
	path := filepath.Join(dataDir, "verse", name+".yaml")
	verse, err := loadYAML[Verse](path)
	if err != nil {
		return Verse{}, err
	}
	for k, v := range verse.Fragments {
		if fragments == nil {
			fragments = make(map[string][]Entry)
		}
		fragments[k] = v
	}
	if err := resolveSlokaFragments(verse.Slokas, fragments); err != nil {
		return Verse{}, err
	}
	mergeCompounds(verse.Slokas)
	return verse, nil
}

func resolveSlokaFragments(slokas []Sloka, fragments map[string][]Entry) error {
	for slokaIndex := range slokas {
		for padaIndex := range slokas[slokaIndex].Padas {
			entries, err := resolveFragments(slokas[slokaIndex].Padas[padaIndex].Entries, fragments)
			if err != nil {
				return err
			}
			slokas[slokaIndex].Padas[padaIndex].Entries = entries
		}
	}
	return nil
}

// mergeCompounds joins compound entries that are split across hyphens.
// When a word ends with "-" (e.g. "viṣṇu-"), it is merged with subsequent
// entries until the compound word no longer ends with "-", producing a single
// Entry per compound (e.g. "viṣṇu-pādāya").
func mergeCompounds(slokas []Sloka) {
	for i := range slokas {
		for j := range slokas[i].Padas {
			slokas[i].Padas[j].Entries = mergeEntries(slokas[i].Padas[j].Entries)
		}
	}
}

func mergeEntries(entries []Entry) []Entry {
	merged := make([]Entry, 0, len(entries))
	i := 0
	for i < len(entries) {
		entry := entries[i]
		// Break entries and non-compound entries pass through unchanged.
		if entry.Break || !strings.HasSuffix(entry.Word, "-") {
			merged = append(merged, entry)
			i++
			continue
		}

		// Start building a compound from this entry.
		compoundWord := entry.Word
		compoundGloss := entry.Gloss
		if !strings.HasSuffix(compoundGloss, "-") {
			compoundGloss += "-"
		}

		// Consume subsequent entries as long as the compound still ends with "-".
		for i+1 < len(entries) {
			next := entries[i+1]
			// A break entry terminates the compound — do not merge across breaks.
			if next.Break {
				break
			}
			compoundWord += next.Word
			compoundGloss += next.Gloss
			i++
			if !strings.HasSuffix(compoundWord, "-") {
				break
			}
			// Ensure trailing "-" on gloss before consuming more parts.
			if !strings.HasSuffix(compoundGloss, "-") {
				compoundGloss += "-"
			}
		}

		merged = append(merged, Entry{Word: compoundWord, Gloss: compoundGloss})
		i++
	}
	return merged
}

func resolveFragments(entries []Entry, fragments map[string][]Entry) ([]Entry, error) {
	resolved := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if entry.Fragment == "" {
			resolved = append(resolved, entry)
			continue
		}
		fragment, ok := fragments[entry.Fragment]
		if !ok {
			return nil, fmt.Errorf("unknown fragment %q", entry.Fragment)
		}
		resolved = append(resolved, fragment...)
	}
	return resolved, nil
}

func loadYAML[T any](path string) (T, error) {
	var data T
	contents, err := os.ReadFile(path)
	if err != nil {
		return data, fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(contents, &data); err != nil {
		return data, fmt.Errorf("malformed YAML in %s: %w", path, err)
	}
	return data, nil
}

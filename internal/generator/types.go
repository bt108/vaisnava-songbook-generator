package generator

import (
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Book struct {
	Title      BookTitle          `yaml:"title"`
	Copyright  string             `yaml:"copyright"`
	Dedication string             `yaml:"dedication"`
	Intro      string             `yaml:"intro"`
	Chapters   []ChapterReference `yaml:"chapters"`
	Fragments  map[string][]Entry `yaml:"fragments"`
}

type BookTitle struct {
	Parts    []TitlePart `yaml:"parts"`
	Tagline  string      `yaml:"tagline"`
	Subtitle *Subtitle   `yaml:"subtitle"`
}

type Subtitle struct {
	Parts   []TitlePart `yaml:"parts"`
	Tagline string      `yaml:"tagline"`
}

type TitlePart struct {
	Sanskrit string `yaml:"sanskrit"`
	English  string `yaml:"english"`
}

type ChapterReference struct {
	ID      string `yaml:"id"`
	Enabled *bool  `yaml:"enabled"`
}

func (chapter ChapterReference) IsEnabled() bool {
	return chapter.Enabled == nil || *chapter.Enabled
}

type Chapter struct {
	Title string
	Intro string
	// Prose holds the free-form LaTeX body blocks of a prose (non-song)
	// chapter, in content order. They render between the intro and any
	// verses; a chapter is either prose-only or song-based.
	Prose   []string
	Content []Verse
}

type chapterFile struct {
	Title   string        `yaml:"title"`
	Intro   string        `yaml:"intro"`
	Content []ContentItem `yaml:"content"`
}

type ContentItem struct {
	Include string `yaml:"include"`
	Prose   string `yaml:"prose"`
}

type Verse struct {
	Title     string             `yaml:"title"`
	Author    string             `yaml:"author"`
	Language  string             `yaml:"language"`
	Slokas    []Sloka            `yaml:"slokas"`
	Fragments map[string][]Entry `yaml:"fragments"`
}

type Sloka struct {
	Number      VerseNumber `yaml:"number"`
	ClearPage   *bool       `yaml:"clearPage"`
	Padas       []Pada      `yaml:"padas"`
	Translation string      `yaml:"translation,omitempty"`
	Footnotes   []Footnote  `yaml:"footnotes"`
	// SongTitle is the title of the enclosing song, shown in place of the
	// sloka's number on the first sloka of each song. Set at load time; not
	// read from YAML.
	SongTitle string
	// AuthorLabel is the author shown in the song header on the first sloka
	// of each song. Set at load time; not read from YAML.
	AuthorLabel string
	// SongLanguage is the language code of the enclosing song (e.g. "bn" for
	// Bengali), used to append the \bengalimark superscript to Bengali song
	// titles. Set at load time; not read from YAML.
	SongLanguage string
}

type Pada struct {
	Entries []Entry `yaml:"entries"`
}

type VerseNumber string

func (number *VerseNumber) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag == "!!null" {
		*number = ""
		return nil
	}
	*number = VerseNumber(value.Value)
	return nil
}

// IsNumeric reports whether the verse number is a plain number (e.g. "1",
// "2"). Non-numeric labels such as "First Refrain" are not numeric.
func (number VerseNumber) IsNumeric() bool {
	_, err := strconv.Atoi(string(number))
	return err == nil
}

type Entry struct {
	Word     string `yaml:"word"`
	Gloss    string `yaml:"gloss"`
	Fragment string `yaml:"fragment"`
	Break    bool   `yaml:"break"`
}

// ShouldClearPage reports whether a page break should be emitted after this
// sloka. Absent means false (slokas flow together by default).
func (sloka Sloka) ShouldClearPage() bool {
	return sloka.ClearPage != nil && *sloka.ClearPage
}

// NumberLabel returns the sloka's header label: the number only when it is a
// non-numeric label such as "First Refrain". Plain numbers are suppressed.
func (sloka Sloka) NumberLabel() string {
	if sloka.Number.IsNumeric() {
		return ""
	}
	return string(sloka.Number)
}

// AuthorLabel returns the author shown in the song header, or "" when the
// verse has no known author. The "unknown" placeholder used by imports
// (any case, with surrounding whitespace) is suppressed.
func (verse Verse) AuthorLabel() string {
	if strings.EqualFold(strings.TrimSpace(verse.Author), "unknown") {
		return ""
	}
	return verse.Author
}

type Footnote struct {
	Term string  `yaml:"term"`
	Note *string `yaml:"note"`
}

func (sloka Sloka) ActiveFootnotes() []Footnote {
	active := make([]Footnote, 0, len(sloka.Footnotes))
	for _, footnote := range sloka.Footnotes {
		if footnote.Note != nil && *footnote.Note != "" {
			active = append(active, footnote)
		}
	}
	return active
}

type Endnote struct {
	VerseLabel string
	Term       string
	Note       string
}

// ActiveEndnotes flattens the chapter's sloka footnotes into endnotes in
// reading order. The first note-bearing verse of each song carries a
// "Song v.N" label; later verses of the same song are labeled "v.N".
func (chapter Chapter) ActiveEndnotes() []Endnote {
	var endnotes []Endnote
	currentSong := ""
	for _, verse := range chapter.Content {
		for _, sloka := range verse.Slokas {
			active := sloka.ActiveFootnotes()
			if len(active) == 0 {
				continue
			}
			label := "v." + string(sloka.Number)
			if verse.Title != "" && verse.Title != currentSong {
				label = verse.Title + " " + label
				currentSong = verse.Title
			}
			for _, fn := range active {
				endnotes = append(endnotes, Endnote{
					VerseLabel: label,
					Term:       fn.Term,
					Note:       *fn.Note,
				})
			}
		}
	}
	return endnotes
}

// EndnoteChapter groups endnotes by chapter for book-level appendix.
type EndnoteChapter struct {
	Title string
	Songs []EndnoteSong
}

// EndnoteSong groups endnotes by song within a chapter.
type EndnoteSong struct {
	Title  string
	Verses []EndnoteVerse
}

// EndnoteVerse groups endnotes by verse within a song.
type EndnoteVerse struct {
	Number   string
	Endnotes []BookEndnote
}

// BookEndnote is a single endnote entry for the book appendix.
type BookEndnote struct {
	Term string
	Note string
}

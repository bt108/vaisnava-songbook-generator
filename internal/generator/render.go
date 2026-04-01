package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// glossGap returns the inter-gloss space to emit after entries[i].
// No space when the gloss ends with "-" (a compound join) or when the
// entry is the last of its line (a break follows, or the pada ends).
func glossGap(entries []Entry, i int) string {
	if strings.HasSuffix(entries[i].Gloss, "-") {
		return ""
	}
	if i+1 >= len(entries) || entries[i+1].Break {
		return ""
	}
	return `\,`
}

// BookData wraps Book with computed book-level endnote data for the template.
type BookData struct {
	Book
	EndnoteChapters []EndnoteChapter
}

func Generate(dataDir, templateDir, outDir string, stdout io.Writer) error {
	book, err := loadYAML[Book](filepath.Join(dataDir, "book.yaml"))
	if err != nil {
		return err
	}

	templates, err := parseTemplates(templateDir)
	if err != nil {
		return err
	}

	chaptersDir := filepath.Join(outDir, "chapters")
	if err := os.MkdirAll(chaptersDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// Collect book-level endnotes as chapters are loaded, grouped by song and verse.
	var endnoteChapters []EndnoteChapter

	for _, chapterRef := range book.Chapters {
		if !chapterRef.IsEnabled() {
			continue
		}
		if chapterRef.ID == "" {
			return fmt.Errorf("chapter entry missing id")
		}

		chapter, err := loadChapter(dataDir, chapterRef.ID, book.Fragments)
		if err != nil {
			return err
		}
		if err := renderFile(templates, "chapter.tex.tmpl", filepath.Join(chaptersDir, chapterRef.ID+".tex"), chapter); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "  %s\n", filepath.ToSlash(filepath.Join(outDir, "chapters", chapterRef.ID+".tex")))

		// Collect endnotes from this chapter, grouped by song then verse.
		var songs []EndnoteSong
		currentSongTitle := ""
		var currentVerses []EndnoteVerse
		var currentVerse *EndnoteVerse

		flushVerse := func() {
			if currentVerse != nil && len(currentVerse.Endnotes) > 0 {
				currentVerses = append(currentVerses, *currentVerse)
			}
			currentVerse = nil
		}

		flushSong := func() {
			flushVerse()
			if len(currentVerses) > 0 {
				songs = append(songs, EndnoteSong{
					Title:  currentSongTitle,
					Verses: currentVerses,
				})
			}
			currentVerses = nil
			currentSongTitle = ""
		}

		for _, verse := range chapter.Content {
			for _, sloka := range verse.Slokas {
				for _, fn := range sloka.Footnotes {
					if fn.Note == nil || *fn.Note == "" {
						continue
					}

					// New song?
					if verse.Title != currentSongTitle {
						flushSong()
						currentSongTitle = verse.Title
					}

					// New verse number?
					numStr := string(sloka.Number)
					if currentVerse == nil || currentVerse.Number != numStr {
						flushVerse()
						currentVerse = &EndnoteVerse{Number: numStr}
					}

					currentVerse.Endnotes = append(currentVerse.Endnotes, BookEndnote{
						Term: fn.Term,
						Note: *fn.Note,
					})
				}
			}
		}
		flushSong()

		if len(songs) > 0 {
			endnoteChapters = append(endnoteChapters, EndnoteChapter{
				Title: chapter.Title,
				Songs: songs,
			})
		}
	}

	bookData := BookData{Book: book, EndnoteChapters: endnoteChapters}
	if err := renderFile(templates, "book.tex.tmpl", filepath.Join(outDir, "book.tex"), bookData); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "  %s\n", filepath.ToSlash(filepath.Join(outDir, "book.tex")))
	return nil
}

func parseTemplates(templateDir string) (*template.Template, error) {
	functions := template.FuncMap{
		"trim":     strings.TrimSpace,
		"glossGap": glossGap,
	}

	templates, err := template.New("root").Funcs(functions).ParseGlob(filepath.Join(templateDir, "*.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	return templates, nil
}

func renderFile(templates *template.Template, templateName, path string, data any) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer file.Close()
	if err := templates.ExecuteTemplate(file, templateName, data); err != nil {
		return fmt.Errorf("render %s: %w", templateName, err)
	}
	return nil
}

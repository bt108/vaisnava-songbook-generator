package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	projectRoot := filepath.Join("..", "..")
	outDir := t.TempDir()
	var output bytes.Buffer

	if err := Generate(filepath.Join(projectRoot, "data"), filepath.Join(projectRoot, "templates"), outDir, &output); err != nil {
		t.Fatal(err)
	}

	mangalarati, err := os.ReadFile(filepath.Join(outDir, "chapters", "mangala-arati.tex"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(mangalarati)
	for _, expected := range []string{
		`\chapter{Maṅgala-ārati}`,
		`\begin{song}{Śrī Śrī Gurv-aṣṭaka}`,
		`\songheader{Śrī Śrī Gurv-aṣṭaka}{Viśvanātha Cakravartī Ṭhākura}`,
		`\begin{sloka}{}`,
		`\padaEntry{vande}{I praise}`,
	} {
		if !strings.Contains(content, expected) {
			t.Errorf("mangalarati.tex does not contain %q", expected)
		}
	}
	if strings.Contains(content, `\begin{sloka}{Śrī Śrī Gurv-aṣṭaka}`) {
		t.Error("song title must not render inside the sloka header")
	}
	if strings.Contains(content, `\blankfootnote{saṅghān}`) {
		t.Errorf("mangalarati.tex rendered a null footnote")
	}

	// chapter files must not contain Notes sections (endnotes moved to book appendix)
	if strings.Contains(content, `\section*{Notes}`) {
		t.Error("mangala-arati.tex must not contain a Notes section (endnotes moved to book appendix)")
	}
	if strings.Contains(content, `\footnoteEntry`) {
		t.Error("mangala-arati.tex must not contain footnote entries (moved to book appendix)")
	}

	// the song header above the first sloka shows the title; the running
	// header title on that page is suppressed by the empty \sectionmark{}
	// the first sloka emits (sty/gloss.tex), so the template no longer
	// applies the songfirst page style
	if strings.Contains(content, `\thispagestyle{songfirst}`) {
		t.Error("template must not emit \\thispagestyle{songfirst}: songfirst is handled by the empty running-head mark of the first sloka's \\sectionmark")
	}
	if !strings.Contains(content, `\songheader{Śrī Śrī Gurv-aṣṭaka}{Viśvanātha Cakravartī Ṭhākura}`) {
		t.Error("song header should render the title and author above the first sloka")
	}
	if !strings.Contains(content, `\begin{sloka}{}`+"\n"+`\songheader{Śrī Śrī Gurv-aṣṭaka}{Viśvanātha Cakravartī Ṭhākura}`) {
		t.Error("song header must render inside the first sloka's unbreakable box")
	}

	// pranama-prabhupada slokas 1 and 2 flow together: no \clearpage between them
	invocation, err := os.ReadFile(filepath.Join(outDir, "chapters", "invocation.tex"))
	if err != nil {
		t.Fatal(err)
	}
	invocationContent := string(invocation)
	if !regexp.MustCompile(`\\end\{sloka\}\s*\\begin\{sloka\}\{\}`).MatchString(invocationContent) {
		t.Error("pranama-prabhupada slokas 1 and 2 should share a page (no \\clearpage between)")
	}

	// default slokas flow together: no \clearpage directly after \end{sloka}
	if regexp.MustCompile(`\\end\{sloka\}\s*\\clearpage`).MatchString(content) {
		t.Error("default slokas must not clear the page after themselves")
	}

	book, err := os.ReadFile(filepath.Join(outDir, "book.tex"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(book), `\include{chapters/nrsimha}`) {
		t.Error("book.tex does not include enabled chapters")
	}

	// book appendix must contain endnotes grouped by chapter
	bookContent := string(book)
	notesIdx := strings.Index(bookContent, `\chapter*{End Notes}`)
	if notesIdx < 0 {
		t.Fatal("book.tex has no End Notes appendix chapter")
	}
	if !strings.Contains(bookContent, `\appendix`) {
		t.Error("book.tex must emit \\appendix before the End Notes chapter")
	}
	if !strings.Contains(bookContent, `\addcontentsline{toc}{chapter}{End Notes}`) {
		t.Error("book.tex must add End Notes to the table of contents")
	}
	if idx := strings.Index(bookContent, `\textsc{`); idx < 0 || idx < notesIdx {
		t.Errorf("book.tex must contain \\textsc endnote terms after the Notes chapter (first at %d, chapter at %d)", idx, notesIdx)
	}
	if strings.Contains(bookContent, `\textsuperscript{`) {
		t.Error("book.tex must not contain \\textsuperscript (numbered endnotes removed)")
	}
	if idx := strings.Index(bookContent, `\subsection*{`); idx < 0 || idx < notesIdx {
		t.Errorf("book.tex must contain \\subsection* song groupings after the Notes chapter (first at %d, chapter at %d)", idx, notesIdx)
	}

	nrsimha, err := os.ReadFile(filepath.Join(outDir, "chapters", "nrsimha.tex"))
	if err != nil {
		t.Fatal(err)
	}
	nrsimhaContent := string(nrsimha)
	if strings.Contains(nrsimhaContent, `\section*{Notes}`) {
		t.Error("nrsimha.tex has a Notes section but the chapter has no footnotes")
	}
	if strings.Contains(nrsimhaContent, `\footnoteEntry`) {
		t.Error("nrsimha.tex renders a footnote entry but the chapter has no footnotes")
	}

	// a sloka with no translation renders a centered no-translation rule
	// with no extra vertical space (nscenter has zero padding); gaura-arati's
	// two refrains are the only translation-less slokas
	gaura, err := os.ReadFile(filepath.Join(outDir, "chapters", "gaura-arati.tex"))
	if err != nil {
		t.Fatal(err)
	}
	gauraContent := string(gaura)
	if got := strings.Count(gauraContent, `\end{nscenter}`); got != 2 {
		t.Errorf("gaura-arati.tex has %d centered no-translation rules, want 2 (the two refrains)", got)
	}
	if !strings.Contains(gauraContent, "\\begin{nscenter}\n\\notransrule\n\\end{nscenter}") {
		t.Error("gaura-arati.tex refrains should render a centered \\notransrule in nscenter")
	}

	// prose chapters render their body text between the chapter heading and any songs
	guide, err := os.ReadFile(filepath.Join(outDir, "chapters", "how-to-use-this-book.tex"))
	if err != nil {
		t.Fatal(err)
	}
	guideContent := string(guide)
	for _, expected := range []string{
		`\chapter{How to Use This Book}`,
		`\section*{How the book is arranged}`,
		`\section*{Transliteration}`,
		`\section*{Pronunciation}`,
	} {
		if !strings.Contains(guideContent, expected) {
			t.Errorf("how-to-use-this-book.tex does not contain %q", expected)
		}
	}
	if strings.Contains(guideContent, `\begin{song}`) {
		t.Error("a prose-only chapter must not render any songs")
	}
}

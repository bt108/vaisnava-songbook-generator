# AGENTS.md

A Vaishnava prayer-songbook PDF generator. Go binary reads YAML → Go templates → LaTeX → LuaLaTeX → PDF.

## Prerequisites

- Go toolchain
- LuaLaTeX (`lualatex`) + `Charis SIL` font
- `Noto Serif Bengali` font — renders the Bengali letters (ড়, ড) in the how-to chapter's transliteration section
- `yamlfmt` — YAML formatter (run after editing `.yaml` files)

## Commands

```sh
make                # generate LaTeX from YAML, then lualatex (twice)
make generate       # YAML → LaTeX only (go run ./cmd/songbook)
make out/book.pdf   # LaTeX → PDF only (generate runs first as a make dependency)
make booklet         # impose out/book.pdf into an A5 saddle-stitched booklet (A4 landscape sheets)
make clean          # rm -rf out/
```

Tests: `go test ./internal/generator/ -v`

## Project layout

```
cmd/songbook/main.go   # CLI (flags: -data, -templates, -out)
internal/generator/    # types.go (data model), load.go (YAML loading), render.go (text/template rendering)
data/book.yaml         # title, intro, chapters list, book-level fragments
data/chapters/*.yaml   # per-chapter content
data/verse/*.yaml      # reusable verse definitions
templates/*.tex.tmpl   # 3 Go templates (book, chapter, sloka)
sty/pre.tex            # LaTeX preamble (packages, layout, chapterIntro, \chapterquote)
sty/gloss.tex          # verse rendering macros (\pada, \padaEntry, \ftrans, etc.)
sty/glossy.sty         # interlinear gloss 2-row table macro
out/                   # generated LaTeX + PDF (gitignored)
post-processing/       # booklet.tex for the saddle-stitched A5 booklet (A4 landscape)
```

## Content model (types in `internal/generator/types.go`)

- **Book** → title, intro, chapters list (with `enabled: true/false`), fragments
- **Chapter** → title, optional `intro` string (plain text or LaTeX), content list
- **Verse** → slokas (with padas of entries); optional inline fragments
- **Entry** → word + gloss + optional fragment/break

## Key conventions

- **Everything YAML-driven.** Edit `data/`, not Go code, to change content.
- **YAML style:**
  - **Text fields** (`intro`, `translation`, footnote `note`) → block scalar `|` (multi-line) or `|-` (single-line). Never single-quoted `'...'` strings — they mangle `''` into `'`.
  - **Entry pairs** (`entries:`, `footnotes:`, `fragments:`) → one-line flow mappings: `- {word: x, gloss: y}`, `- {break: true}`, `- {fragment: name}`, `- {term: x, note: y}`.
  - Everything else: `field:"value"` format, not `field: { key: val }` nesting.
- **yamlfmt** after every YAML edit — run `yamlfmt .` to normalize all.
- **IAST diacritics** for Sanskrit transliteration.
- **Chapters** enabled/disabled via `enabled: true/false` in `book.yaml`. Missing from the list = disabled.
- **File naming:**
  - Chapter ids and verse filenames are kebab-case + IAST, and match each other (e.g. chapter id `gaura-arati` ↔ verse `gaura-arati`). Chapter ids appear in `book.yaml`, `data/chapters/<id>.yaml`, and generated `out/chapters/<id>.tex`.
  - Verse families are prefixed by topic: `pranama-*`, `tulasi-*`, `mahamantra-*`. Use full names, never cryptic abbreviations (e.g. `pranama-prabhupada`, not `pranama-sp`).
- **`include:` references resolve against `data/verse/` only.**
- **Fragments** define reusable entry sequences. Prefer verse-level fragments over book-level.
- **Pada line breaks** via `break: true` entry in a pada's entries list.
- **Sloka `number`** is normally a plain integer and is suppressed in the output; non-numeric values (`First Refrain`) render as that sloka's header label.
- **Sloka `clearPage: true`** inserts a page break after that sloka (distinct from the pada-level `break: true`).
- **Verse `language:`** is `bn` (Bengali) or `sk` (Sanskrit). `bn` renders a gray superscript B on the song title (`\bengalimark` in `sty/gloss.tex`); set it on Bengali songs.
- **Sloka `footnotes:`** render as chapter-level endnotes in the `Notes` section at the end of the chapter (via `\footnoteEntry` in `sty/gloss.tex`), not under the verse. Follow the rules in `docs/footnotes-style.md` before writing them.
- **`\chapterquote{text}{source}`** — optional LaTeX command for styled quotes inside chapter intros. Italic + indented text, source right-aligned. Works only inside `chapterIntro` environment (which is used when a chapter has an `intro:` field). The source attribution order is: author (if any), work, reference, place (if any), date (if any) — e.g. `(Śrīla Prabhupāda, Śrīmad-Bhāgavatam 2.9.1, Tokyo, April 1972)`, `(Śrī Caitanya-caritāmṛta, Ādi-līlā 1.20)`, `(Letter to Trai, December 27, 1972, Bombay)`. Omit elements that don't apply (e.g. a work-only quote has no author).
- **Importing new songs** from `source-books/` (canonical source `ss.md`, txt, or the EPUB `Temple_Songs_*.html` files) → use the `/import-song` skill (`.opencode/skills/import-song/`), which covers the extract → YAML → gloss pipeline. Never improvise a verse file from scratch.
- **Fixing overfull hboxes** ("Overfull \hbox" in `out/book.log`, Sanskrit/gloss lines wider than the 340.2pt text block, or deciding whether to shorten a gloss vs. add a pada break) → use the `/resolve-overflow` skill (`.opencode/skills/resolve-overflow/`), which has the width model, the Sanskrit-bound/gloss-marginal/free-win classification, and the measurement procedure.

## Build notes

- `make out/book.pdf` sets `TEXINPUTS=.:../sty/` for `\input{pre}` and `\input{gloss}` to resolve.
- `make booklet` compiles post-processing/booklet.tex from post-processing/ (so
  `../out/book.pdf` resolves) with output to out/booklet.pdf. Pages are never
  scaled; the ~1 mm A4 leftover goes in the fold.
- LuaLaTeX runs twice (TOC, cross-references).
- `out/` is gitignored.

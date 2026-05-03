---
name: resolve-overflow
description: Use when a Sanskrit pada and/or its gloss overflows the line width in the songbook — "Overfull \hbox" warnings in out/book.log, pada lines sticking out past the margin, or deciding whether a gloss needs shortening / a pada needs a break. Classifies each overflowing pada as Sanskrit-bound, gloss-marginal, or a free-win break, then resolves via gloss shortening or pada breaks. Triggered by "overfull", "overflow", "too wide", "sticks out", "hbox", "gloss too long", "pada too long", "fix the overfull boxes".
---

# Resolve Horizontal Overflows (Sanskrit / gloss)

When a pada's words + glosses are wider than the text block, the centered
`\gloss` tabular line overflows the margin and LuaLaTeX reports
`Overfull \hbox (N.NNpt too wide)`. The fix is either to **shorten glosses**
(where possible) or to **split the pada with a break** (where the Sanskrit
alone is too wide). This skill teaches the exact width model, how to classify
each overflowing pada, and the fix hierarchy.

## The width model (know this cold)

- **Text width = 340.2pt** (text block x ∈ [56.7, 396.9], from `sty/pre.tex`
  geometry).
- **Each pada = ONE unbreakable centered line.** A pada is a `\gloss` 2-row
  tabular (word row `\Large\bfseries`, gloss row
  `\color{glossgray}\itshape\normalsize` — see `sty/glossy.sty` +
  `sty/gloss.tex`) inside `\begin{pada}` … `\end{pada}` in
  `templates/sloka.tex.tmpl`.
- **Column width = max(word, gloss).** Each entry is a `c`-column tabular, so
  the entry's width is the wider of its Sanskrit word and its gloss.
- **Inter-entry gap = `\,` kern at normalsize = 1.83pt** (THIN space). Emitted
  by `glossGap` in `internal/generator/render.go` — *no* gap after a gloss
  ending in `-` (compound join) or before a `break`/last entry. Compound words
  (trailing `-` in the word field) are pre-merged at load time in
  `internal/generator/load.go`, so they rarely appear as separate entries, but
  `glossGap` remains as a safety net for any unmerged compounds.
- **Pada width = Σ max(word_w, gloss_w) + 1.83 × (n_entries − 1).**
  Compare against **340.2pt**. Excess over 340.2 = the reported overfull size
  (the real build matches this model to ~0.04pt).

## Diagnose

1. Build and collect overfulls:
   ```sh
   make out/book.pdf
   grep -n 'Overfull \\hbox' out/book.log
   ```
2. **Ignore tiny overfulls**: < 1pt is rounding noise (e.g. 0.14pt from page
   output, 0.5pt). Only > 1pt needs fixing.
3. Isolate the chapter's section in the log: find the line
   `(./chapters/<id>.tex` and scan forward to the next chapter open.
4. Map each significant overfull to its pada: the log gives "in paragraph at
   lines L--M" — those are line numbers in the generated
   `out/chapters/<id>.tex`. Each pada renders as `\begin{pada}` … so walk the
   sloka blocks to find which sloka/pada overflows, then fix that pada's YAML
   in `data/verse/<id>.yaml`.

## Classify the overflowing pada (3 kinds)

Compute two numbers per pada (fragments expanded inline — see
`internal/generator/load.go` `resolveFragments`):

- **Sanskrit-only width**: Σ word widths + 1.83 × gaps (ignore glosses).
- **Full width**: Σ max(word, gloss) + 1.83 × gaps.

| Kind | Condition | Fix |
|---|---|---|
| **Sanskrit-bound** | Sanskrit-only > 340.2pt | Break is **mandatory** — no gloss shortening can ever fix it. Keep/add the break. (In Hari-nāma-ratnam, 35 of 44 overflowing padas were this kind.) |
| **Gloss-marginal** | Sanskrit-only fits, full > 340.2pt | **Rescuable by shortening glosses** — the glosses are what push it over. |
| **Free win** | Has a break but full unbroken width ≤ 340.2pt | **Remove the break** — it never needed one. No gloss edit needed. |

## Fix hierarchy (cheapest first)

1. **Remove free-win breaks.** Delete the `- {break: true}` entry from that
   pada's `entries:` in the YAML. Pure win, no content change.
2. **Shorten glosses** (gloss-marginal padas). A gloss replacement only helps
   if the **new gloss is narrower than the old gloss AND the word** (column =
   max). Shortening a gloss that is already narrower than its word does
   nothing — check `word_w` vs `gloss_w` first. Save = old_gloss −
   max(word, new_gloss).
3. **Keep/add breaks** for Sanskrit-bound padas: insert
   `- {break: true}` at a natural midpoint of the pada's `entries:`
   (prefer breaking before the refrain/fragment so it sits on its own line).
   Each break splits the pada into two rendered lines, both must fit.

**Global safety rule**: shortening a gloss can only narrow columns, never
widen — so a shortening is safe everywhere the word appears. Apply it
**consistently per word** (same Sanskrit word = same gloss across the file,
per gloss-spec rule 3). Use `replaceAll` edits or a word→gloss fix map so
you don't leave the same word glossed two ways.

**Break removal caution**: after removing a break, that pada renders on ONE
line again — re-verify it fits (step: full unbroken width ≤ 340.2) or the
overfull just moves.

## Gloss shortening candidates (measured)

Glosses below were measured at gloss style (`\itshape\normalsize`, Charis SIL
Compact) — see `references/measure.md` for the reproducible procedure.

| Old gloss | Width | Candidates (width) |
|---|---|---|
| indeed | 29.6 | **just** 16.0, truly 21.4, verily 25.1, surely 26.9 |
| lifespan | 34.9 | **life** 14.0 |
| remember | 44.7 | **recall** 25.0 |
| therefore | 40.4 | **hence** 25.7 |
| with heart | 46.2 | **heartfelt** 37.9, heart 23.5 |

Sanity anchor: word `hi` = 13.8pt, `harer` = 36.9, `nāmaiva` = 58.7,
`kevalam` = 57.8 — the refrain (`harer nāmaiva kevalam` = "the name of Hari
alone") is 36.9 + 58.7 + 57.8 + 2×1.83 = **157.1pt**, so a refrain-ended pada
has ~183pt of budget for the rest.

## Verify

1. `yamlfmt <file>` then `go test ./internal/generator/`.
2. `make out/book.pdf`; confirm **0 overfulls > 1pt in the chapter's section**
   of `out/book.log`.
3. Check the sloka line split didn't drift: count `\begin{pada}` blocks per
   sloka in `out/chapters/<id>.tex` (2-line slokas = 2 padas, etc.).
4. `make booklet` (A5 saddle-stitched booklet compiles cleanly).
5. Commit the YAML change with a message listing what was shortened/removed
   and the new 2/3/4-line sloka split.

## Worked example (Hari-nāma-ratnam)

44 overflowing padas → 35 Sanskrit-bound (breaks kept), 9 gloss-marginal
(rescued via `indeed→just`, `remember→recall`, `therefore→hence`,
`lifespan→life`, `with heart→heartfelt`), 1 free win (v5p2 break removed).
Result: 45 → 35 breaks; sloka split 62/31/7 → 71/23/6 (two/three/four-line);
chapter has 0 overfulls > 1pt. Committed as
`feat(content): shorten glosses to rescue 10 pada breaks in Hari-nāma-ratnam`.

## Related

- `import-song` skill + its `references/gloss-spec.md` — the width discipline
  that prevents overflows at glossing time (gloss ≈ word width, never > ~1.5×
  the Sanskrit letter count).
- `docs/footnotes-style.md` — footnote conventions (not related to width).

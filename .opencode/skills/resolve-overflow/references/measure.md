# Measuring word/gloss widths (reproducible)

The width model (`SKILL.md`) needs per-entry widths for **Sanskrit words** and
**glosses** as rendered by the real build. Never eyeball them — measure them.
The procedure below reproduces the exact values used to classify and fix
Hari-nāma-ratnam (e.g. `hi` 13.8pt, `harer` 36.9pt, `nāmaiva` 58.7pt,
`kevalam` 57.8pt; glosses `indeed` 29.6pt, `just` 16.0pt, `remember` 44.7pt,
`recall` 25.0pt).

## Why two styles

`sty/gloss.tex` sets the two tabular rows:

```tex
\gsSetFirstStyle{\Large\bfseries}                        % Sanskrit word row
\gsSetSecondStyle{\color{glossgray}\itshape\normalsize}  % gloss row
```

So words must be measured at `\Large\bfseries` and glosses at
`\itshape\normalsize` — measuring both at the same size gives wrong widths.

## Procedure

1. **Build a measure document** in a scratch dir (not in the repo), one line
   per candidate, in the correct style:

   ```tex
   \documentclass[11pt]{book}
   \input{pre}
   \pagestyle{empty}
   \begin{document}
   \noindent{\Large\bfseries sarveṣām}\par        % words
   \noindent{\Large\bfseries hi}\par
   \noindent{\normalsize\itshape indeed}\par      % glosses
   \noindent{\normalsize\itshape just}\par
   \end{document}
   ```

   `\input{pre}` pulls in the real preamble (fonts, sizes) — mandatory for
   accurate widths. Compile with the songbook `sty/` on the search path:

   ```sh
   TEXINPUTS=".:/home/t/Downloads/songbook/sty/:" lualatex -interaction=nonstopmode cand.tex
   ```

2. **Extract widths** from the PDF bounding boxes:

   ```sh
   pdftotext -bbox cand.pdf - | grep '<word' | python3 -c "
   import sys, re
   for line in sys.stdin:
       m = re.match(r'<word xMin=\"([\d.]+)\" yMin=\"([\d.]+)\" xMax=\"([\d.]+)\" yMax=\"([\d.]+)\">(.*)</word>', line.strip())
       if m:
           w = float(m.group(3)) - float(m.group(1))
           print(m.group(5), round(w, 1))
   "
   ```

   Width = xMax − xMin. Values match the real build's tabular column widths to
   ~0.04pt — good enough to predict whether a pada fits.

   **Known pitfall — prefer `\settowidth` over `pdftotext -bbox`**: the bbox
   method is unreliable for multi-word glosses (pdftotext splits them into
   separate `<word>` tokens, breaking per-string lookup) and for words with
   apostrophes (`Ligatures=TeX` turns `'` into `’`, so `madhurebhyo'pi` never
   matches). The robust alternative measures each string directly in LaTeX:

   ```tex
   \newlength\mw
   \settowidth{\mw}{\Large\bfseries madhurebhyo'pi}\typeout{W|madhurebhyo'pi|\the\mw}
   \settowidth{\mw}{\itshape\normalsize own-respective-duty}\typeout{G|own-respective-duty|\the\mw}
   ```

   then parse the `W|`/`G|` lines from the log. Two more gotchas:
   - **Same string, two styles**: a string that appears as *both* a word and a
     gloss (e.g. `karma` = 43.7pt as word, 29.6pt as gloss) must be stored
     under separate word/gloss namespaces — a single string-keyed dict silently
     overwrites one with the other (~14pt error).
   - **`\,` only in context**: when measuring a whole pada, wrap it in
     `\begin{glossy}` (that's where `\gloss` is defined) — `\gloss` is
     undefined outside the environment.

## Using the widths

For a pada (fragments expanded inline), column per entry = **max(word, gloss)**:

```
pada_width = Σ max(word_w, gloss_w) + 1.83 × (n_entries − 1)
```

- Text width = **340.2pt**. `pada_width − 340.2` = the overfull the real build
  reports.
- The inter-entry gap is `\,` (THIN) = **1.83pt at normalsize** — emitted
  by `glossGap` in `internal/generator/render.go`. Skipped after a
  trailing-hyphen gloss (compound join) and before a break / at the pada end.
  Compound words are pre-merged at load time in `internal/generator/load.go`,
  so unmerged compounds are rare, but `glossGap` remains as a safety net.

### Classification decision

- **Sanskrit-only** (glosses set to zero): `Σ word_w + gaps > 340.2`
  → **Sanskrit-bound**: a break is mandatory, glosses can't help.
- **Full** (`Σ max(word, gloss) + gaps > 340.2`), Sanskrit-only ≤ 340.2
  → **gloss-marginal**: rescuable by shortening the glosses that exceed their
  words.
- A pada that **has a break but full width ≤ 340.2** → **free win**: remove
  the break, no other change.

### How much does a gloss swap save?

`save = old_gloss_w − max(word_w, new_gloss_w)` (column takes the max; a new
gloss narrower than the word saves nothing beyond making the column equal the
word). Sum the saves over the padas entries; if `Σ saves ≥ excess`, the pada
fits unbroken.

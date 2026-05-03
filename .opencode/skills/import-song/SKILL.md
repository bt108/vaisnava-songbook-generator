---
name: import-song
description: Use when importing a Vaishnava song/verse into the songbook from an external source file (markdown, txt, or html/epub — primarily the Sacred Song Symphony source-books/ss.md), building the data/verse/*.yaml datafile — extracting the song and its author, taking translations verbatim, splitting Sanskrit into padas/entries, and filling glosses. Triggered by requests like "import song X from file Y".
---

# Import Song

Import a song from an external source into the songbook's YAML data model. The
source is almost always the **Sacred Song Symphony** (`source-books/ss.md`,
referred to below as "ss.md") — otherwise a **plain-text file** (e.g.
`misc/song.txt`) or an **HTML file** (e.g. the EPUB `Temple_Songs_*.html` files
in `misc/.../`). Each source has its own layout; this skill teaches the one
pipeline that handles all of them.

## The pipeline

1. **Locate** the song in the source file (by `Song Name:` header in txt,
   `<div class="head1">` heading in HTML, or the index/page markers in ss.md).
2. **Extract** the raw song block: metadata, numbered Sanskrit stanzas, and the
   translation section.
3. **Build the YAML skeleton** (`data/verse/<id>.yaml`): title, **author**,
   slokas, padas, entries with Sanskrit words split per the reference —
   glosses left **empty**.
4. **Split Sanskrit** into entries (compounds at `-`, padas at line breaks;
   from ss.md use the **standalone** version for pada breaks, word split
   points, and spelling — never the word-for-word rendition, which is used
   only for glosses and translations).
5. **Fill glosses** (Layer 1, per entry) following
   `references/gloss-spec.md`.
6. **Verify** against the source (pada count, split points, translations).
7. **`yamlfmt`** the file.

## When to Use

- "import song X from file Y"
- Adding a new verse file to `data/verse/` from ss.md, a txt, or an html
  source
- Rebuilding/fixing an existing verse file against its source

Do NOT use for: editing chapter structure, intro text, LaTeX templates, or
`sty/` files.

## Mandatory style rules (apply to every file you touch)

- `entries:`, `fragments:`, `footnotes:` are **one-line flow mappings**:
  `- {word: x, gloss: y}`, `- {break: true}`, `- {fragment: name}`,
  `- {term: x, note: y}`.
- All text fields (`translation`, `note`) are **block scalars** (`|` or `|-`),
  never single-quoted strings.
- Run `yamlfmt .` (or `yamlfmt <file>`) after editing any YAML.
- Verse-level `translation` is **verbatim from the source** except for the
  mechanical EPUB fixes in Step 5 item 1 (digit substitution, period/comma
  mixups, missing capitalization) — never rewritten beyond those.
- Never touch: filenames, `fragment:` references, `footnotes` content,
  `bar` values.

---

## Step 1 — Locate the song

### TXT source (e.g. `source-books/song.txt`)

Songs are delimited by a header line:

```
Song Name: Ar Keno Maya Jale
Author: Bhaktivinoda Thakura
Book Name: Gitavali (Section: Sreyo-Nirnaya Song 2)
(1)
ār keno māyā-jāle poḍitecho, jīva-mīn
nāhi jāno baddha ho’ye ro’be tumi ciro-din
...
TRANSLATION
1) O fish like soul, ...
```

- Search for `^Song Name: <name>` (case-insensitive).
- The song block runs from that header until the next `Song Name:` header.
- Header fields: `Song Name:`, `Official Name:`, `Author:`, `Book Name:`.
  `Official Name:` may be the canonical Sanskrit title (prefer it for the YAML
  `title:` when present and clearly canonical; otherwise use `Song Name:`).
- Stanzas are numbered `(1)`, `(2)`, ... in order.
- The `TRANSLATION` line begins the translation section; translations are
  numbered `1)`, `2)`, ... matching stanza numbers.
- **Ignore** page-number noise lines like `(19)` and `ISKCON desire tree` —
  they interrupt the block between stanzas. A stanza number is a number that
  actually starts a Sanskrit stanza (followed by verse text); page numbers are
  isolated `(NNN)` lines that reset or jump.

### HTML/EPUB source (e.g. `source-books/A. C. Bhaktivedanta Swami Prabhupada - Temple Song Book-Bhaktivedanta Book Trust (1990)/Temple_Songs_NNNN.html`)

- All sources live under `source-books/` — the EPUB extraction is the
  directory `A. C. Bhaktivedanta Swami Prabhupada - Temple Song Book-Bhaktivedanta Book Trust (1990)/` containing the `Temple_Songs_NNNN.html` files (plus `content.opf`, `toc.ncx`, etc.).

- Each song usually lives in one `Temple_Songs_NNNN.html` file.
- Search the directory for the song name in `<div class="head1">` (the song
  title) — or in the `<h3>` block / TOC links.
- Structure of a song:

```html
<h3 id="calibre_toc_7" ...><div class="calibre9">Śrī Śrī Gurv-aṣṭaka </div></h3>
<div class="head1">Śrī Śrī Gurv-aṣṭaka</div>
<div class="head3">Eight Prayers to the Guru<br/> by Śrīla Viśvanātha Cakravartī Ṭhākura</div>
<div class="purport">...intro/overview...</div>

<div class="verse">(1)<br/> saṁsāra-dāvānala-līḍha-loka-<br/>
trāṇāya kāruṇya-ghanāghanatvam<br/>
prāptasya kalyāṇa-guṇārṇavasya<br/>
vande guroḥ śrī-caraṇāravindam </div>

<div class="purport">The spiritual master is receiving...</div>
```

- `<div class="head1">` = song title (this is the YAML `title:`).
- Each `<div class="verse">` = one sloka. The leading `(N)` (or a label like
  `(First Refrain)`) is the `number:`. Each `<br/>`-separated segment is one
  **pāda (visual line)**.
- The `<div class="purport">` following a verse (or group of verses) is the
  verse-level translation. Take it **verbatim** (modulo the mechanical EPUB
  fixes in Step 5 item 1 — digit substitution, period/comma mixups, missing
  capitalization).
- Optional first `<div class="purport">` before any verse = song intro; not a
  verse translation.

Extraction helper (Python) — from the EPUB directory:

```python
import re, html

with open(f'{epub_dir}/Temple_Songs_NNNN.html') as f:
    content = f.read()

verses = re.findall(r'<div class="verse">(.*?)</div>', content, re.DOTALL)
for v in verses:
    v = re.sub(r'<br\s*/?>', '', v, re.IGNORECASE)
    v = html.unescape(v)
    v = re.sub(r'<[^>]+>', '', v)
    v = re.sub(r'\s+', ' ', v).strip()
    v = re.sub(r'^\(\d+\)\s*', '', v)
    v = re.sub(r'(\S-)\s+(\S)', r'\1\2', v)
    print(v)
```

Caveat: the `.*?` regex is non-greedy and assumes no nested tags inside the
`verse`/`purport` divs — true for the Temple Song Book EPUB, but if a div
contains nested markup, strip tags first. Also note the helper flattens
`<br/>` away: recover pada (visual line) boundaries from the raw HTML's
`<br/>` segments, not from this flattened output.

### Sacred Song Symphony source (`source-books/ss.md`) — the primary source

This is the **preferred source** for new imports. It is a markdown export of
the book *Sacred Song Symphony* (grep for the song name — titles appear in the
index and as section headings). Its special structure:

- **Every song appears twice.** Find the index/TOC entry first: it lists the
  song with a **page number** and an **xref** (e.g. `469↔454`). The song is
  then rendered twice:
  1. **Standalone version** (song alone, early pages) — the visual layout of
     the song as sung: the source of truth for **pada breaks**, **word
     split points**, and **spelling**.
  2. **Word-for-word version** (glossed, later pages) — used only for
     **glosses** and the verse translations. Its splits and spellings must
     NOT be copied into the YAML when they differ from the standalone.
- **Author** comes from the standalone section heading. The pattern is:
  `[Page N]` line, then the **title** line, then the **author** line (e.g.
  `Dāmodarāṣṭaka` / `Vyāsadeva (Satyavrata Muni)`, `Gurvaṣṭaka` /
  `Viśvanātha Cakravartī Ṭhākura`). Take the author line verbatim as the
  YAML `author:` value. If the title line is directly followed by the verse
  text (no author line — common for compiled sections like the Maṅgalācaraṇa
  prayers, Tulasī prayers, and the pañca-tattva mantra), use `unknown`.
  Also use `unknown` when the song has no matching ss.md entry at all (e.g.
  `maha-prasade-govinde.yaml` is not in ss.md) or the author would have to
  be guessed from tradition rather than read from the source.
- Navigate with the page markers: pages appear as `[Page 123]` lines. The
  two renditions of a song are ~340 pages apart (the xref gives both).
- **Splits and spelling come from the standalone version — never the
  word-for-word rendition.** If the standalone breaks a line the word-for-word
  version keeps together, follow the standalone (see
  `sad-gosvami-astakam.yaml`: its padas match the standalone layout, e.g.
  breaking after `parau`). The word-for-word version may rejoin or re-break
  differently and may spell words differently — always prefer the standalone's
  splits, spellings, and sandhi forms, and **adopt the standalone's spelling in
  the data file** even when the existing data file differs (e.g. change
  data-file `hoy` to ss.md's `haya`, `jāi` to `yāi`). The word-for-word
  version is used only for glosses and translations.
- **ss.md is not an absolute source of truth.** It is a digitized book and can
  itself contain typos, OCR artifacts, and inconsistent spellings. If you
  suspect a typo in ss.md (e.g. a word that contradicts how it appears
  elsewhere in ss.md, in other editions, or on the internet — like the
  anusvāra on `dalāyatākṣam` that appears only in the word-for-word rendition
  of govindam while every online source writes plain `m`), try to figure out
  the truth: cross-check the internet (Vedabase, Gitabase, other songbooks)
  and **report your finding to the user**, applying the verified correction
  only after discussing it. Never silently copy a suspected typo.
- **Reproduce wraps exactly, don't "fix" them.** If the standalone version
  wraps a word mid-word *without* a hyphen (e.g. `bhajanā` / `nandena` from
  `bhajanānandena`, `gaṇā` / `kīrṇe` from `gaṇākīrṇe`), split the entry at
  that point and gloss the halves with their literal senses (`bhajanā` =
  worship, `nandena` = joy). Keep the split verbatim — do not rejoin the
  word. List such splits in your report so the user can see them.

---

## Step 2 — Build the YAML skeleton

File: `data/verse/<kebab-case-id>.yaml`. ID conventions:

- Kebab-case, IAST where applicable (e.g. `gaura-arati`, `pranama-prabhupada`,
  `maha-prasade-govinde`).
- Prefix families by topic: `pranama-*`, `tulasi-*`, `mahamantra-*`.
- Match the chapter id if this song is a chapter (e.g. chapter `gaura-arati` ↔
  verse `gaura-arati`).

Skeleton (glosses empty; verse translation verbatim):

```yaml
title: <from source>
author: <from source, or unknown>
slokas:
  - number: 1
    padas:
      - entries:
          - {word: ...-, gloss: ""}
          - {word: ..., gloss: ""}
    translation: |-          # verbatim from source
      ...
    footnotes: []
```

Notes:

- `title:` comes verbatim from the source heading — except **txt sources,
  which often write titles in ASCII** (no diacritics, e.g. `Ar Keno Maya
  Jale`). Convert to proper IAST for the YAML title (`Ār Kena Māyā Jāle`);
  HTML/EPUB headings are already IAST.
- `author:` comes from the standalone section heading (ss.md) — the line
  right after the title — or `unknown` when the source gives no author (see
  Step 1, ss.md section). Never invent an author from tradition.
- `number:` may be `1`, `2`, ... or a label like `First Refrain` (string).
- `clearPage: true` only when the source/design requires a page break (see
  existing files; default is absent/false).
- If the verse has **no translation** in the source (e.g. a refrain), omit the
  `translation:` field entirely — the generator renders `———` instead.
- `footnotes:` starts as `[]`. Follow `docs/footnotes-style.md` only when
  adding real etymological notes (not part of a bare import).

---

## Step 3 — Split Sanskrit into entries

From the reference text (the EPUB verse or txt stanza), one `word:` per token:

- `-` in the reference = split point. Each segment becomes its own YAML
  `word:` entry. Non-final parts keep a trailing `-`.
  ```python
  expected = []
  for chunk in ref_text.split():
      parts = chunk.split('-')
      for i, part in enumerate(parts):
          expected.append(part + '-' if i < len(parts) - 1 else part)
  ```
- Match **sandhi forms exactly** as written in the source (`nṛsiṁho` not
  `nṛsiṁhaḥ`, `bahir` not `bahiḥ`). Do not "fix" them to dictionary forms.
- **Apparent source typos** (EPUB quirks like `punyaṁ` for `puṇyaṁ`,
  `idarṅ` for `idam`): if you are **highly confident** of the correct IAST
  form, fix it. If you are **not** highly confident, keep the source form
  and **ask the user at the end** of the import (list it in your report
  with your best guess) for the correct form.
- Include trailing punctuation as part of the word entry (`loy,` not `loy`).
  Quote when needed: `{word: 'loy,', gloss: ''}`. YAML flow mappings with `,`
  or `:` inside values **must** be quoted.
- Compound split points (the `-` boundaries) come from the reference, never
  invented. If the source shows `saṁsāra-dāvānala-līḍha-loka-` on one visual
  line, keep those splits.

### Pada (line) grouping

Each visual line of the reference is a pāda. Group words into YAML `padas` so:

- **Visual lines = number of `padas` + number of `{break: true}` entries.**
- Each pāda's words go in a single pada's `entries`. If one YAML pada spans
  two visual pādas (e.g. a long line that must stay together), insert
  `{break: true}` between the pādas' words.
- **From ss.md: use the *standalone* version for visual-line boundaries,
  split points, and spelling.** The word-for-word version may group lines
  differently and spell words differently; the standalone layout is what the
  reader sees and what the PDF should match. Never split or spell as in the
  word-for-word rendition.
- **Never split a compound across a pada/`break: true` boundary** unless the
  reference itself shows the trailing `-` at end of line (e.g. `loka-`).
  Exception: ss.md wraps a word mid-word *without* a hyphen (`bhajanā` /
  `nandena`) — split there too, verbatim (see Step 1, ss.md section).
- Consecutive short pādas may be merged into one pada as long as no compound
  is split.

### Recurring phrases → `fragments:`

If the same multi-word phrase repeats across slokas (e.g. `vande guroḥ
śrī-caraṇāravindam` at the end of every verse), hoist it to a top-level
`fragments:` block, gloss it **once**, and reference it in padas via
`- {fragment: <name>}` (no gloss field). See `gurvastakam.yaml` for the model.
Only do this when the phrase genuinely recurs; otherwise keep entries inline.
If the recurring phrase appears with an apparent source typo in *one* verse
(e.g. a refrain spelled `jagannāthah` instead of `jagannāthaḥ`): if you are
highly confident it's a typo, use the fragment everywhere with the corrected
form. If not confident, keep the fragment for the well-formed occurrences,
inline that one occurrence as ordinary entries, and ask the user at the end.

**Fragments may contain `{break: true}` entries** — the generator splices
fragment entries inline verbatim (breaks included; see
`resolveFragments` in `internal/generator/generator.go`), so a fragment can
span two visual lines. If the standalone layout splits a recurring phrase
across a line break (e.g. the Ṣaḍ-gosvāmy-aṣṭaka refrain `vande
rūpa-sanātanau raghu-yugau` / `śrī-jīva-gopālakau`), put the `{break: true}`
*inside* the fragment and reference it once — do **not** create two
fragments (`vanderupa1`/`vanderupa2` was a workaround and is unnecessary).
See `sad-gosvami-astakam.yaml` for the model.

---

## Step 4 — Fill glosses (Layer 1)

Follow `references/gloss-spec.md` exactly. This is the heart of the import.

- Only `gloss:` (per entry) is filled. Nothing else changes.
- **Width discipline (mandatory):** each gloss must be about the same width as
  the Sanskrit word (counting letters **and spaces**), capped at ~1.5× the
  letter count — never cross it by much. Keep glosses short; drop filler
  ("the", "to", "in", "of") where the sense survives. Full rule in
  `references/gloss-spec.md`.
- **`nam-` root is always `I bow`:** `namo`, `namas`, `namaha`, `namaste`,
  `namāmi` are glossed `I bow` in every song, now and future.
- **Pronouns are flat:** `te` → `you`, `me` → `me` — never `to you`, `unto me`.
- Verse-level `translation:` is verbatim from the source — never rewritten.
  The only allowed changes are the mechanical EPUB fixes in Step 5 item 1
  (digit substitution, period/comma mixups, missing capitalization).
- Fragments are already glossed (from step 3) — never re-glossed at use site.
- Check the file's own `footnotes:` for term meanings and stay consistent.

---

## Step 5 — Verify against the source

Checklist (this is the verify-fix pass):

1. **Verse translations** match the source `<div class="purport">` (or txt
   `N)` translation) — ignoring curly vs straight quote differences and
   double spaces (EPUB may show ` ārati ` with extra spaces; YAML uses
   single). Mechanical EPUB artifacts are **fixed automatically**, not
   copied: digit-substitution (`0 Lord` → `O Lord`), period/comma mixups
   (`wife. as` → `wife, as`), and missing capitalization after periods.
   Verify these fixes were applied and nothing else changed.
2. **Sanskrit words** match the reference exactly — sandhi forms, split
   points, trailing punctuation — except for confident typo fixes per
   Step 3.
3. **Pada count**: visual lines = padas + breaks. No compound split across a
   pada boundary.
4. **Title** matches the source heading (`<div class="head1">` or txt header).
5. **Author** matches the ss.md standalone heading's author line, or is
   `unknown` when no author line exists (compiled sections) or the song is
   absent from ss.md.
6. Glosses consistent across the file and with existing files for the same
   stems (`vande` → `I praise`, `guroḥ` → `of guru` everywhere).
7. Style: one-line flow mappings, block scalars, no single-quoted text fields.
8. Run `yamlfmt <file>` and confirm clean.
9. **Case-table examples** — the ``Basic Sanskrit and Bengali grammar'' section of
   `data/chapters/how-to-use-this-book.yaml` keeps one real example per Sanskrit
   case. If any case currently has no example and this song contains a clean
   example of it (a word whose case is unambiguous, ideally from the song's own
   glosses), add it to the table's ``In this book'' / ``Gloss'' columns and
   mention the addition in your report. Do not replace an existing example with
   one from the new song.

### Asking the user

If any apparent source typo was kept because you were not highly confident of
the correct form, the import is **not final**: list each kept typo (source
form + your best guess + location) in your report and ask the user for the
correct form before the file is considered done. Same for any Sanskrit word
whose form is genuinely uncertain — don't guess silently.

---

## Common mistakes

| Mistake | Fix |
|---|---|
| Over-splitting compounds beyond the reference's `-` boundaries | Split only where the reference hyphenates (ss.md also splits mid-word *without* a hyphen — reproduce those verbatim) |
| Dictionary sandhi forms instead of the pāda-pāṭha in the reference | Copy the reference form exactly |
| Missing trailing punctuation in word entries | `loy,` not `loy`; quote if needed |
| Inventing an author from tradition instead of the source | Take the author line verbatim from the ss.md standalone heading; use `unknown` when there is none |
| Glossing a `fragment:` reference | `{fragment: name}` has no gloss field |
| Creating two fragments when a refrain spans two visual lines | Put `{break: true}` inside one fragment and reference it once |
| Rewriting the verse-level translation | Verbatim from the source — only the mechanical EPUB fixes in Step 5 item 1 are allowed |
| Keeping obvious EPUB typos (`0 Lord`, `punyaṁ`) | Fix automatically when highly confident; ask the user when not |
| Splitting a compound across a pāda boundary | Keep compounds on one line unless source shows trailing `-` |
| Single-quoted YAML strings for text fields | Use block scalars (`\|`, `\|-`) |
| Confusing EPUB double spacing with a real mismatch | Collapse to single spaces; not a mismatch |
| Editing curly to straight quotes | Don't — quote style is irrelevant |
| Using word-for-word split points, spelling, or pada breaks from ss.md | Use the *standalone* version for visual lines, split points, and spelling; the word-for-word rendition is used only for glosses/translations |
| Overlong gloss that overflows the Sanskrit word's width | Keep it ≈ the word's width (cap ~1.5× letters, counting spaces); drop filler (`the`, `to`, `in`, `of`) where the sense survives |
| Glossing the `nam-` root as "obeisances"/"homage"/"to You" | Always `I bow` (`namo`, `namas`, `namaha`, `namaste`, `namāmi`), in every song |
| Glossing pronouns with prepositions (`to you`, `unto me`, `of you`) | Flat pronoun only: `te` → `you`, `me` → `me`, `tava` → `your` |

## References

- `references/gloss-spec.md` — Layer 1 (gloss) rules, with the consistency
  and tradition-verification guidance.

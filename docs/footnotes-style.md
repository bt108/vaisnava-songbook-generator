# Footnote Style Guide

How to write the etymological footnotes on Sanskrit verse terms. Reference
model: `data/verse/gurvastakam.yaml`.

## Structure

A sloka may carry a `footnotes:` list under its own keys (sibling of `padas:`
and `translation:`):

```yaml
- number: 5
  padas:
    - entries:
        - {word: apāra-, gloss: unlimited-}
  translation: |-
    ...
  footnotes:
    - {term: apāra, note: "fr. \\sa{a} (without) + \\sa{pāra} (far shore); a sea whose far shore is never reached — boundless."}
    - {term: mādhurya, note: "fr. \\sa{madhu} (honey); sweetness arrived at through honey — 'honeyed.'"}
```

- `term:` is the headword, in IAST, as it appears in the verse text
  (compound strips and case endings dropped: `sākṣād-` → `sākṣāt`,
  `śāstrair` → `śāstrair`).
- `note:` is the body. Double-quoted YAML string. The note body renders in the
  book-level **End Notes** appendix at the end of the book, grouped by chapter →
  song → verse. Notes are never shown under the verse itself.

## Ordering

List footnotes **in the order the terms appear in the verse text** (pada 1
first, then pada 2). The block renders top-to-bottom in list order.

## When to footnote

Only add a footnote when it earns its space:

- **Etymology that adds depth** — the root is non-obvious or the semantic
  chain is illuminating (e.g. `prasāda` "that which sits in front → clearness →
  free gift → remnants").
- **Lesser-known origin of a common word** (e.g. `guṇa` "thread, strand,
  bow-string → quality").
- **No duplicates** across the whole book. A term is footnoted once. If you
  want to preserve the term in the list without rendering it, use
  `{term: X, note: null}` (null notes are skipped by the generator).
- Skip when the word is self-evident or the etymology adds nothing.

## Note conventions

- Sanskrit in the note body is marked `\sa{...}`. In the YAML double-quoted
  string the backslash is doubled: `\\sa{karuṇā}`.
- Root derivation opens with `fr. \sa{...}` and **the gloss always goes in
  parens right after its `\sa{...}`**, both for single roots and for each part
  of a `+` compound:
  - `"fr. \\sa{śās} (to teach, rule, correct); śāstra, 'that which instructs and rectifies.'"`
  - `"fr. \\sa{a} (without) + \\sa{pāra} (far shore); ..."`
- "Literally" is the shorthand **`lit.`**: `"lit. 'thread, strand, bow-string' → 'quality'; ..."`.
- **`cf.` only ever points to another footnote in the same chapter's endnotes**
  (e.g. the `bhakta` note's "cf. FN. to `\sa{bhāj}`"). Notes are now collected
  globally but remain organized by chapter in the End Notes appendix — a
  cross-chapter `cf.` can't be easily followed by the reader, so don't use it;
  say what you mean inline instead.
- The semantic chain uses `→` (e.g. `delight → love`).
- Compound words are broken into their parts with the sense of each, plus the
  combined sense when it is not obvious:
  - `"fr. \\sa{ni} (in, within) + \\sa{kuñja} (bower, thicket); a secluded grove."`
- When a word has multiple roots or senses, pick the sense active in the verse.
- **Quoting** uses LaTeX conventions — a straight `'` renders as a closing
  quote. Use backtick + apostrophe for single quotes and double-backtick +
  double-apostrophe for double quotes:
  - `` `quality', `honeyed.' `` → single quotes
  - `` ``that which one walks with.'' `` → double quotes
  - Possessive apostrophes (`soul's`) stay a plain `'`.
- Glosses in the pada entries stay **literal** (singular form if that is the
  literal value), e.g. `sandhyaṁ` glosses as `junction`, not `junctions` —
  the plural comes from the compound `tri-sandhyaṁ`.

## Good vs. bad

Good:
```yaml
- {term: guṇa, note: "lit. 'thread, strand, bow-string' → 'quality'; a quality as a strand woven into a being."}
```

Bad:
```yaml
- {term: guṇa, note: "quality, attribute"}                 # no depth
- {term: guṇa, note: "cf. Sāṅkhya"}                        # not a book footnote — reader can't follow
- {term: rati, note: "love"}                               # self-evident, no root
```

## Adding a footnote

To add a footnote to a verse, append a `footnotes:` list to the sloka in its
YAML file (under `data/verse/`). Each entry is a flow mapping with `term` and
`note`:

```yaml
footnotes:
  - {term: bhakti, note: "fr. \\sa{bhaj} (to serve, worship); loving devotion."}
```

The file can live in any chapter — footnotes from all chapters are collected into
the End Notes appendix at the end of the book.

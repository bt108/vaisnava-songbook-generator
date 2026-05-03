# Interlinear Gloss Spec

Applies to songbook datafiles: `title:`, optional top-level `fragments:`
(recurring phrases, glossed once), and `slokas:` → `padas:` → `entries:`.

```yaml
title: "..."
fragments:                          # top-level, defines recurring phrases ONCE
  vandeguro:
    - {word: vande, gloss: I praise}        # already filled — read-only
    - {word: guroḥ, gloss: of guru}
    - {word: śrī-caraṇāravindam, gloss: the radiant-feet-lotus}
slokas:
  - number: 1
    padas:
      - entries:
          - {word: saṁsāra-, gloss: ""}
          - {word: dāvānala-, gloss: ""}
          - {word: trāṇāya, gloss: ""}
          - {fragment: vandeguro}   # no gloss field — recurring fixed phrase
    translation: |-                 # <- guru's full verse translation, verbatim, DO NOT TOUCH
      ...
    footnotes:
      - {term: ..., note: ...}      # read-only reference — use for meaning, don't edit
```

**Style rules (mandatory):** entries, fragments, and footnotes are **one-line
flow mappings** (`{word: x, gloss: y}`, `{break: true}`, `{fragment: name}`,
`{term: x, note: y}`). All text fields (`translation`, `note`) are **block
scalars** (`|` or `|-`), never single-quoted strings. Run `yamlfmt` on any
file you edit.

Only one field is ever edited: **`gloss`** (per entry, Layer 1). Nothing else
in the file is touched.

---

## Layer 1 — `gloss` (per entry)

**Job:** the literal sense of this one entry, nothing more.

**Width discipline (mandatory).** The gloss is typeset in a smaller font than
the Sanskrit word, so it must be kept short. Aim for the gloss to be about the
**same width as the Sanskrit word** — counting letters **and spaces** — and
never cross roughly **1.5× the number of Sanskrit letters**. Really try not to
cross it at all, at least not by much. When a short gloss and a long one both
work, take the short one: prefer single words, drop filler words ("the", "to",
"in", "of") wherever the sense survives without them, and let the verse
translation carry the nuance. Rule of thumb: a 5-letter word caps around 7-8
gloss characters, a 10-letter word around 15. The `nam-` root is the model:
`namo` / `namas` / `namaha` / `namaste` / `namāmi` → **`I bow`** — short and
uniform, in every song, now and future.

1. **One gloss per entry — but check whether the compound is already split at
   the data level or not.**
   - **Split compounds** (common case): the data already breaks the compound
     into separate consecutive entries, each with its own trailing hyphen in
     `word` (e.g. `saṁsāra-`, `dāvānala-`, `kāruṇya-`). Gloss each entry with
     its own single plain sense — do **not** add a hyphen inside the gloss
     text; the structural split already does that job.
   - **Unsplit compounds** (happens when sandhi fuses the members into one
     written word and the data was not broken up further — e.g.
     `prahlādāhlāda-` = *Prahlāda* + *āhlāda*; or `ghanāghanatvam` = *ghana* +
     *aghana* + *tvam*, "simultaneously solid and non-solid," per its
     footnote): the gloss **must** carry all the components itself, hyphenated
     within the single gloss string (e.g. `Prahlāda-joy`,
     `solid-nonsolid-ness`). This is the one place a hyphen belongs inside the
     gloss text.
   - Rule of thumb: if the compound's members already appear as separate
     `entries`, don't hyphenate the gloss. If a single `word` entry is itself
     doing the work of multiple compound members (recognizable by sandhi
     fusing two recognizable stems into one written word with no entry break),
     hyphenate inside that entry's gloss.
   - **Trailing-hyphen echo:** if the Sanskrit `word` itself ends in a hyphen
     (i.e. it continues into the next entry), the `gloss` also ends in a
     hyphen — e.g. `{word: "vakṣaḥ-", gloss: "chest-"}`,
     `{word: "prahlādāhlāda-", gloss: "Prahlāda-joy-"}`. This applies
     regardless of whether the gloss is a single plain sense (split-compound
     case) or itself internally hyphenated (unsplit-compound case) — the
     trailing hyphen is a separate, mechanical echo of the word's own trailing
     hyphen, not related to the internal-hyphenation rule above. A `word` with
     no trailing hyphen gets a `gloss` with no trailing hyphen.
2. **Give the core sense, not a grammatical sentence.** Skip filler case-words
   ("of," "to," "in") unless genuinely needed to not mislead. Nouns stay
   nouns, verb roots stay verb-ish — don't build full clauses here.
   - **Pronouns keep their own core word, not an inserted preposition.** A
     pronoun like `te` or `tava` already *is* the case relationship
     (dative/genitive), not a case-marker glued onto a noun — so gloss it with
     its natural single-word English form (`te` → `you`, `tava` → `your`,
     `me` → `me`), not a prepositional phrase (`to you`, `of you`, `unto me`).
     Apply this **always**, even where the translation reads "of yours" or
     "to you" — the gloss is just the pronoun, which also keeps it short
     enough for the width rule. What gets dropped per this rule is a
     *preposition layered on top* of a noun/pronoun's core sense, not the
     pronoun's own basic meaning.
   - **The `nam-` root is always `I bow`.** `namo`, `namas`, `namaha`,
     `namaste`, `namāmi` — every form of the root *nam* ("to bow, to offer
     obeisances") — is glossed **`I bow`** every time it appears, in every
     song, now and future. It is short enough to fit the width rule and stays
     uniform across the whole songbook. (A fused compound like `namāmīśvaraṁ`
     = *namāmi* + *īśvaram* glosses `I bow to the Lord`.)
   - **Bengali locative `-e`** (many songs are Bengali, e.g. Bhaktivinoda
     Ṭhākura): gloss the word with its locative sense `in/at the X` —
     `jāle` → `in the net`, `mane` → `in the mind`. Like the pronoun rule,
     this is the word's own inflection, not an added case-word.
   - **Preserve the Sanskrit's own grammatical number**, even where natural
      English would pluralize (e.g. `nakham`, grammatically singular, glosses
      as `nail` even though hands have several nails).
3. **Consistency is mandatory.** Every recurring word/root (`vande`, `guroḥ`,
   `śrī-`, `caraṇāravindam`, common epithets) gets the **same English gloss
   every time it appears**, across the whole file and across files. Before
   glossing a word, check whether it (or its stem) already appears elsewhere
   in the file and match that choice.
4. **Use the footnotes as the authority when present.** If a `footnotes` entry
   already defines a term's nuance (e.g. `kāruṇya` = "deep compassion, a step
   deeper than karuṇā"), the gloss should reflect that same sense in short
   form — don't contradict the footnote.
5. **`fragments` block (top-level) is read-only.** It defines each recurring
   phrase once, already glossed — never edit it.
6. **Pada-level `fragment: name` references are never glossed themselves**
   (they have no `gloss` field) — leave as-is.
7. **If a fragment's phrase is ever also spelled out inline** as ordinary
   entries elsewhere (rather than via `fragment: name`), the individual
   glosses given there must add up to the same meaning as the fragment block's
   combined gloss — check the `fragments` block first and stay consistent with
   it, even if the local split of words differs (e.g. the fragment glosses
   `śrī-caraṇāravindam` as one unit "the radiant-feet-lotus", but a verse
   might spell it as two separate entries `śrī-` and `caraṇāravindam` — gloss
   those two so together they still convey "radiant" + "feet-lotus", not a
   different phrasing).
8. **Proper nouns / untranslated technical terms** (deity names, `bhakti`,
   etc.): leave in transliteration as the gloss (i.e. gloss = the term itself
   or its standard epithet-translation, per existing footnote if one exists),
   not a paraphrase.
   - **Name vs. descriptive compound.** Some words are etymologically
     transparent compounds (`narasiṁha`/`nṛsiṁha` = man+lion) but function as
     the deity's actual invoked name throughout the tradition (found in
     name-lists, used as the direct term of address) — these stay
     transliterated per this rule even though you could translate the parts
     (e.g. `Nṛsiṁha`, `Keśava`, `Hari`, `Kṛṣṇa`). Other words are structurally
     identical compounds of ordinary translatable roots but are being used
     *descriptively* in this specific line to unpack a quality, form, or
     feature rather than to name/address someone (e.g. `narahari-` inside
     `dhṛta-narahari-rūpa`, "bearing the man-lion form" — describing what the
     form consists of; or `jagadīśa` = "world-lord" as a fresh epithet of
     praise rather than the primary name being invoked) — these get translated
     and hyphenated per rule 1 instead. Rule of thumb: is this word
     *naming/addressing* the entity, or *characterizing* it? Naming →
     transliterate; characterizing → translate.
   - **Multiple stems for the same deity are not normalized to one spelling.**
     If a deity is referred to by two or more different (but synonymous)
     Sanskrit stems across the file or across the songbook (e.g. `narasiṁha`
     and `nṛsiṁha`, both "man-lion," referring to the same avatar), gloss each
     stem with its own matching transliteration — don't force them to match
     each other. Consistency (rule 3) applies *within* each stem family (every
     `nṛsiṁha`-form gets `Nṛsiṁha`, every `narasiṁha`-form gets `Narasiṁha`),
     not across genuinely different stems.
9. **No punctuation, no capitalization** except for proper nouns. (Standard
   English orthography for the pronoun "I" isn't an exception worth avoiding —
   it stays capitalized as always.)

---

## Verifying against tradition

Many songs in this book are famous, widely-published verses (classical
stotras, well-known praṇāma mantras) with established word-for-word (padārtha)
breakdowns already circulating in print and online. Before glossing, especially
for compounds whose split isn't obvious from the words alone, search for and
cross-check these traditional sources. Use them as **verification of meaning**,
not as a template to copy wholesale — the traditional sources often keep the
case-marking prepositions ("unto You," "of Hiraṇyakaśipu") that this spec's
Layer 1 rule 2 deliberately strips out, and may combine multiple compound
members into one translated phrase where this schema's `entries` split them
apart. Confirm the underlying sense against tradition, then apply this spec's
own minimalism and structure rules on top of it.

---

## What is never touched

- `title`, `number`, `word`, `fragment`
- top-level `fragments` block (already glossed — use its values, don't edit
  them)
- verse-level `translation` (guru's — verbatim from source; only the
  mechanical EPUB fixes from SKILL.md Step 5 item 1 — digit substitution,
  period/comma mixups, missing capitalization — may be applied)
- `footnotes` (read-only reference — use, don't edit)
- Overall YAML structure, key order, formatting

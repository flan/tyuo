# tyuo

tyuo is a Markov-chain-based chatterbot, inspired by
[Jason Hutchens' MegaHAL](https://en.wikipedia.org/wiki/MegaHAL).

More specifically, tyuo is a rework of a hack against the original C-based MegaHAL,
called yuo, based on a limited understanding of how it worked and very rudimentary
knowledge of the language.

This implementation is almost entirely clean-room, save for the "surprise" calculation,
which was preserved for use either as a tie-breaker against tyo's own scoring model or
for more-faithful emulation of its inspiration.

Notably, this design diverges by concurrently making use of a mix of n-grams, from di-
through quintgram, trading time-to-novel production and search-space complexity against
relative quality of output.

It also deviates in that it was written entirely with a chatroom use-case in mind, rather
than evaluation by a critical judge.


## usage

This system can run as a daemon or in a terminal or however you want.
Run with `-help` to see what options exist. The most significant one is `-data-dir`.

Within the specified directory, a few files need to exist:

```
$DATA_DIR
├── contexts
│   └── <context-id>.json
└── languages
    ├── <language>.banned
    └── <language>.boring
```

Samples for all of these can be found in the `data/` directory within this project.

`boring` words are those that serve a purely functional role in framing an inquiry.
When someone wants to ask tyuo a question, picking these as keywords will almost certainly
lead to a response they won't want, so they're discounted as choices when parsing the input.

`banned` substrings are parts of words that should never be uttered by the chatbot.
This list protects both against tyuo learning phrases containing the words (avoiding undesired
n-gram paths and dictionary entries) and producing phrases containing the words if the list is
modified after learning has occurred. There is also a per-context mechanism to extend bans on
a local scale.

The design supports multiple languages, but at present, only English has been implemented.
Take a look at `tyuo/logic/language/english.go` if you want to get started on implementing another;
you have a great deal of freedom in deciding how and what your language needs to handle.
English, for example, does a fair bit of massaging to correct for unambiguous errors, but all
of its extensions are local to its parsing and formatting functions.

`<context-id>.json` is given the ID for a context as its name and contains details about how that
context should operate, all the knobs and dials that can be tweaked to affect how it learns and speaks.

```javascript
{
    /* the language to use in this context */
    "Language": "english",
    
    "Ngrams": {
        /* digrams are the simplest and fastest transition model; using them will
         * produce behaviour that is often novel, sometimes insightful,
         * frequently deranged, particularly as learning progresses
         * 
         * it's pretty random and will only resemble native speech by chance
         */
        "Digrams": false,
        /* trigrams are a fairly middle-ground option, producing relevant
         * observations with some regularity and having about as much
         * sentence-structure correctness as a machine-translation between
         * languages with no common ancestry
         */
        "Trigrams": true,
        /*
         * quadgrams are a reasonably stable choice for production of "how do you
         * do, fellow humans" responses, being well-formed for the most part, but
         * closely reflecting observed input: a lot of data will need to be
         * learned before novel structures will be produced with any regularity
         * and search-spaces will sometimes be exhausted while finding viable
         * paths
         */
        "Quadgrams": true,
        /* quintgrams (and anything above them) will rarely deviate from mimicing
         * what was learned; occasional novel productions are possible, but it
         * will not be uncommon to see near-verbatim recreations of input data for
         * the context's entire lifespan
         */
        "Quintgrams": false
    },
    
    "Learning": {
        /* how long, in tokens, input needs to be before learning will occur;
         * it is automatically fed to any enabled n-gram structures that
         * can accomodate the given length
         */
        "MinTokenCount": 6,
        /* the number of runes allowed within any single token,
         * used to prevent over-hyphenated compounds that will only
         * ever be seen a handful of times from cluttering the database
         * 
         * for English, 13 should capture most words; 15 is probably the
         * highest you'd ever want, and 10 is still pretty reasonable
         */
        "MaxTokenLength": 13,

        /* how long to hold on to n-gram structures, in seconds
         * 
         * this allows unused paths, like those created by typos or
         * short-lived memes, to expire
         */
        "MaxAge": 31536000,

        /* the number of dictionary occurrences or transitions at which
         * to trigger rescale logic, which eliminates obsolete entries and
         * keeps the numbers in check
         */
        "RescaleThreshold": 1000,
        /* the divisor for rescaling; this affects how frequently it happens
         * and how long rare entries hang around
         */
        "RescaleDecimator": 3
    },
    "Production": {
        /* the maximum number of searches to conduct simultaneously,
         * used to limit resource usage when doing long recursive explorations
         */
        "MaxParallelOperations": 8,
        
        /* the number of keytokens or terminals to choose before starting a search
         *
         * this is the most influential factor in determining how many productions
         * are generated and should be turned up if tyuo is failing to produce output
         * or down if it's consuming too many resources and overproducing options
         */
        "TokensInitial": 2,
        /* how many paths to explore from the initial token, in both directions
         * 
         * this is the second-most inluential factor in quantity of options produced
         * and should be increased if there isn't enough originality in tyuo's
         * utterances or decreased if resources are constrained
         */
        "SearchBranchesInitial": 4,
        /* how many paths to explore from bounaries, in both directions
         * 
         * this regulates the case where tyuo will attempt to build a production
         * that starts or ends with the initial token, controlling how many paths
         * will be explored from that origin; like the others, more means more
         * diversity, but also more resource-usage
         */
        "SearchBranchesFromBoundaryInitial": 2,
        /* this controls the breadth component of each recursive search
         * 
         * tyuo attempts to complete at least one viable chain from each of its
         * origin points and will only explore as many child-paths, from each node,
         * as is required to satisfy that requirement; increasing this number improves
         * the likelihood of success at low cost, but it will increase the time required
         * to complete a search that can't possibly succeed
         */
        "SearchBranchesChildren": 2,

        /* the minimum number of tokens that need to be present in a production for it
         * to be a viable candidate as an utterance
         */
        "MinLength": 5,
        /* the upper limit on how long a production can be
         * 
         * this prevents searches from being infinitely deep and constrains productions
         * so they aren't too long to be amusing quips.
         */
        "MaxLength": 30,
        /* the likelihood of stopping production, upon finding a terminal transition,
         * before reaching the target range
         */
        "StopProbability": 0.25,

        /* the minimum desired length of a production */
        "TargetMinLength": 8,
        /* the maximum desired length of a production */
        "TargetMaxLength": 16,
        /* the likelihood of stopping production, upon finding a terminal transition,
         * after reaching the target range
         */
        "TargetStopProbability": 0.375,

        /* when tyuo learns a token, it keeps track of how often it appears in what
         * it considers "base" form, stripped of diacritics and rendered in lower-case,
         * as much as possible based on Unicode decomposition rules
         * 
         * if a token is represented in its base form at least this often, as a
         * percentage where 1.0 == 100%, then its base form will be considered the
         * generic representation and the language may capitalise it or perform other
         * transformations as appropriate
         * 
         * if that isn't the case, then the most commonly occurring variant is chosen,
         * allowing tyuo to mimic understanding of proper nouns and other orthographic
         * patterns
         */
        "BaseRepresentationThreshold": 0.9,
        
        /* tyuo incorporates MegaHAL's "surprise" scoring mechanism, in addition to its own
         * 
         * this algorithm is intended to encourage selection of more-novel productions when
         * there are multiple choices available
         * 
         * turning on either direction will cause surprise to be calculated using the
         * highest-order-enabled n-gram model, at the cost of doing a lookup for all
         * applicable n-grams; this isn't usually too bad, due to cache locality, but
         * it may be an expense you won't want to incur, especially if you don't select
         * based on surprise in your application
         */
        "CalculateSurpriseForward": true,
        "CalculateSurpriseReverse": true
    }
}
```

With a config file in place, when tyuo is asked to interact with the corresponding context, it will load the config
and, if necessary, create a new `<context-id>.sqlite3` file in the same directory, which hosts its dictionaries and
n-grams.


### interaction
See `scripts/` for a few toy Python scripts that demonstrate how to interact with this system in a debug capacity.


## dependencies

You may need to grab these with `go get` to build this project. There shouldn't be any special versioning requirements.

First-party dependencies; you shouldn't have any concerns about introducing these:

- `golang.org/x/text/cases`
- `golang.org/x/text/language`
- `golang.org/x/text/runes`
- `golang.org/x/text/transform`
- `golang.org/x/text/unicode/norm`

Third-party dependencies; these are pretty mature and well-vetted:

- `github.com/juju/loggo`
  - logging
- `github.com/mattn/go-sqlite3`
  - SQLite3 driver


## theory of operation

At its core, tyuo is an implementation of an [n-gram](https://en.wikipedia.org/wiki/N-gram)-based
[Markov chain](https://en.wikipedia.org/wiki/Markov_chain) model.
It supports digrams, trigrams, quadgrams, and quintgrams.

When generating productions, an arbitrary keytoken is chosen from the input-set and then both forward and backward
walks occur to reach a terminal state, the end of some other previously observed input. After that, the production
fragment is fed into the same process in the opposite direction, this time beginning with its tail, hopefully ensuring
that the full path from front to back is of consistent quality. (This is in contract to MegaHAL, which does two
independent walks from the keytoken, then joins them together)

If no productions can be assembled in this manner, then forward and backward walks are conducted from arbitrary
terminal-initial states, attempting to choose granches that include a keytoken along the way. This is so tyuo almost
always has something to say, even if it ends up being completely irrelevant.

To support this flow, the highest-order n-gram model is tried first in each step, successively falling back to the lowest
if no n-grams match the head/tail in its current state; if no matching n-grams can be found, then that production fails
and it's up to the other searches to yeild usable results.

After all productions have been generated, they are scored and everything that meets minimum viable criteria is presented
to the caller, with scoring results, so it can pick the most suitable option for its needs.


### implementation details

tyuo uses an isolated SQLite3 database to serve each context. There is no common memory or overlap between them,
though it does share language-level banned and boring lists.

As much as reasonably possible, tyuo will not hold any information in the database in memory, to reduce its process
footprint when idle, which is likely to be close to 100% of the time, given that its operations tend to be on the order
of 30ms on a Ryzen 3700X.

One exception to this is banned tokens. It's just much more efficient to hold that relatively small set of strings and IDs in
memory to quickly filter out unwanted input without unnecessarily hitting the database.

API accesses are internally subject to a read/write-lock, allowing for threadsafe use by any number of callers.

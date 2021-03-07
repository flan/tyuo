//Ban-checking should probably do a global lcase-check against all banned tokens,
//rather than tokenising first, since 100 strstrs is probably cheaper than
//700 + related copying, even if it's over a longer input

//Also, before returning the list of candidates to the caller, run them through
//another ban-pass and eliminate anything that fails the check
//This is an unfortunate necessity of having a language-level ban-list:
//something offensive could have made its way into the dictionary before that
//applied, then not been explicitly banned  within the context
//...or maybe, on-load, do a substring match of all language-level banned
//tokens and ban any resulting IDs in the runtime memory
//That would probably be much more performant.


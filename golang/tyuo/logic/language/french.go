package language
//When parsing and formatting French, things like punctuation have different
//spacing rules from English

//It probably makes sense to just eat the byte used by a space and store the token
//in printable form, rather than reprocess everything at presentation time

//French has similar roles to English for apostrophes, but it also has mappings
//for accents; if an accent is found, treat that as a variant form, too, but consider
//it to be the same as a capital


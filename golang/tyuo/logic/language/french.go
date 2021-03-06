//When parsing and formatting French, things like punctuation have different
//spacing rules from English

//It probably makes sense to just eat the byte used by a space and store the token
//in printable form, rather than reprocess everything at presentation time

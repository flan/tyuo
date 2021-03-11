package language
import (
    "regexp"
    
    "github.com/flan/tyuo/context"
)

var nullToken = context.ParsedToken{
    Base: "",
    Variant: "",
}

const hyphen = '-'
const apostrophe = '\''

type punctuationCompressor struct {
    punc string
    re *regexp.Regexp
}
var punctuationCompression = []punctuationCompressor{
    punctuationCompressor{
        punc: "…",
        re: regexp.MustCompile("^\\.\\.+$"),
    },
    punctuationCompressor{
        punc: "‼",
        re: regexp.MustCompile("^!!+$"),
    },
    punctuationCompressor{
        punc: "⁇",
        re: regexp.MustCompile("^\\?\\?+$"),
    },
    punctuationCompressor{
        punc: "—",
        re: regexp.MustCompile("^--$"),
    },
    punctuationCompressor{
        punc: "⁈",
        re: regexp.MustCompile("^[\\?!]{2,}$"),
    },
}

//when digesting input, convert these into common tokens for interoperability
//when formatting output, the language's rules can map them back to what they should be
var punctuationEquivalents = map[rune]rune{
    '。': '.',
    '、': ',',
    '？': '?',
    '！': '!',
}

var punctuation runeset
func init() {
    punctuation = make(runeset, len(context.PunctuationIdsByToken) + len(punctuationEquivalents) + 1)
    for _, puncSpec := range context.Punctuation {
        punctuation[puncSpec.GetRepr()] = voidInstance
    }
    for punc, _ := range punctuationEquivalents {
        punctuation[punc] = voidInstance
    }
    punctuation['-'] = voidInstance //special case for '--'
}

func punctuationDigest(punc string) (context.ParsedToken) {
    if len(punc) == 1 {
        if punc[0] != '-' { //only valid when doubled
            return context.ParsedToken{
                Base: punc,
                Variant: punc,
            }
        }
    } else {
        for _, pc := range punctuationCompression {
            if pc.re.MatchString(punc) {
                return context.ParsedToken{
                    Base: pc.punc,
                    Variant: pc.punc,
                }
            }
        }
    }
    return nullToken
}

//processes punctuation at the start and end of a token, returning whatever remains
func punctuationDissect(token []rune) (context.ParsedToken, context.ParsedToken, []rune, bool) {
    var punctuationBefore = nullToken
    var punctuationAfter = nullToken
    
    currentToken := make([]rune, 0, 3)
    //forward search
    for _, r := range token {
        if _, isPunctuation := punctuation[r]; isPunctuation {
            if equiv, equivDefined := punctuationEquivalents[r]; !equivDefined {
                currentToken = append(currentToken, r)
            } else {
                currentToken = append(currentToken, equiv)
            }
        } else {
            break
        }
    }
    if len(currentToken) > 0 {
        token = token[len(currentToken):]
        parsedToken := punctuationDigest(string(currentToken))
        if parsedToken != nullToken {
            punctuationBefore = parsedToken
        } else {
            return punctuationBefore, punctuationAfter, token, false
        }
        currentToken = make([]rune, 0, 3)
    }
    
    //reverse search
    for i := len(token) - 1; i >= 0; i-- {
        r := token[i]
        if _, isPunctuation := punctuation[r]; isPunctuation {
            if equiv, equivDefined := punctuationEquivalents[r]; !equivDefined {
                currentToken = append(currentToken, r)
            } else {
                currentToken = append(currentToken, equiv)
            }
        } else {
            break
        }
    }
    if len(currentToken) > 0 {
        token = token[:len(token) - len(currentToken)]
        
        //NOTE: none of the supported punctuation cases about ordering, so
        //this should be unnecessary, but reverse the slice for consistency
        for i, j := 0, len(currentToken) - 1; i < j; i, j = i + 1, j - 1 {
            currentToken[i], currentToken[j] = currentToken[j], currentToken[i]
        }
        
        parsedToken := punctuationDigest(string(currentToken))
        if parsedToken != nullToken {
            punctuationAfter = parsedToken
        } else {
            return punctuationBefore, punctuationAfter, token, false
        }
    }
    
    return punctuationBefore, punctuationAfter, token, true
}

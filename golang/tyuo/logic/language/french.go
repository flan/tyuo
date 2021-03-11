package language
import (
    "github.com/flan/tyuo/context"
    
    "golang.org/x/text/transform"
)

var frenchCharacters = runeset{
    'a': voidInstance,
    'b': voidInstance,
    'c': voidInstance,
    'd': voidInstance,
    'e': voidInstance,
    'f': voidInstance,
    'g': voidInstance,
    'h': voidInstance,
    'i': voidInstance,
    'j': voidInstance,
    'k': voidInstance,
    'l': voidInstance,
    'm': voidInstance,
    'n': voidInstance,
    'o': voidInstance,
    'p': voidInstance,
    'q': voidInstance,
    'r': voidInstance,
    's': voidInstance,
    't': voidInstance,
    'u': voidInstance,
    'v': voidInstance,
    'w': voidInstance,
    'x': voidInstance,
    'y': voidInstance,
    'z': voidInstance,
    'A': voidInstance,
    'B': voidInstance,
    'C': voidInstance,
    'D': voidInstance,
    'E': voidInstance,
    'F': voidInstance,
    'G': voidInstance,
    'H': voidInstance,
    'I': voidInstance,
    'J': voidInstance,
    'K': voidInstance,
    'L': voidInstance,
    'M': voidInstance,
    'N': voidInstance,
    'O': voidInstance,
    'P': voidInstance,
    'Q': voidInstance,
    'R': voidInstance,
    'S': voidInstance,
    'T': voidInstance,
    'U': voidInstance,
    'V': voidInstance,
    'W': voidInstance,
    'X': voidInstance,
    'Y': voidInstance,
    'Z': voidInstance,
    'ç': voidInstance,
    'Ç': voidInstance,
    'é': voidInstance,
    'É': voidInstance,
    'â': voidInstance,
    'Â': voidInstance,
    'ê': voidInstance,
    'Ê': voidInstance,
    'î': voidInstance,
    'Î': voidInstance,
    'ô': voidInstance,
    'Ô': voidInstance,
    'û': voidInstance,
    'Û': voidInstance,
    'à': voidInstance,
    'À': voidInstance,
    'è': voidInstance,
    'È': voidInstance,
    'ì': voidInstance,
    'Ì': voidInstance,
    'ò': voidInstance,
    'Ò': voidInstance,
    'ù': voidInstance,
    'Ù': voidInstance,
    'ë': voidInstance,
    'Ë': voidInstance,
    'ï': voidInstance,
    'Ö': voidInstance,
    'ü': voidInstance,
    'Ü': voidInstance,
    apostrophe: voidInstance,
    hyphen: voidInstance,
}

const frenchConsecutiveVowelLimit = 3
const frenchConsecutiveConsonantLimit = 4
var frenchVowelsNormalised = runeset{
    'a': voidInstance,
    'e': voidInstance,
    'i': voidInstance,
    'o': voidInstance,
    'u': voidInstance,
}

var frenchLanguageDefinition = languageDefinition{
    delimiter: ' ',
    characters: frenchCharacters,
    
    digestToken: func (token []rune, normaliser *transform.Transformer) ([]context.ParsedToken, bool) {
        tokens := make([]context.ParsedToken, 0, 2)
        punctuationBefore, punctuationAfter, token, learnable := punctuationDissect(token)
        if !learnable {
            return tokens, false
        }
        
        if punctuationBefore != nullToken {
            tokens = append(tokens, punctuationBefore)
        }
        
        if len(token) > 0 { //not fully digested
            //punctuationDissect will deal with leading/trailing hyphens, so just check for apostrophes
            if token[0] == apostrophe || token[len(token) - 1] == apostrophe {
                //a learnable token can't be bounded by an apostrophe, since it might be a quotation mark
                return tokens, false
            }
            
            containsPunctuation := false
            for _, r := range token {
                if _, isCharacter := frenchCharacters[r]; !isCharacter {
                    return tokens, false
                }
                if r == apostrophe || r == hyphen {
                    if containsPunctuation { //allow at most one punctuation-mark, to limit abuse
                        return tokens, false
                    }
                    containsPunctuation = true
                }
            }
            
            //get the normalised base-form
            variant := string(token)
            base, _, err := transform.String(*normaliser, variant)
            if err != nil {
                logger.Warningf("unable to normalise token %s: %s", variant, err)
                return tokens, false
            }
            
            //check to make sure there aren't too many vowels or consonants clumped together; this likely indicates gibberish
            vowelCount := 0
            consonantCount := 0
            for _, r := range base {
                if _, isVowel := frenchVowelsNormalised[r]; isVowel {
                    vowelCount++
                    if vowelCount > frenchConsecutiveVowelLimit {
                        return tokens, false
                    }
                    consonantCount = 0
                } else {
                    consonantCount++
                    if consonantCount > frenchConsecutiveConsonantLimit {
                        return tokens, false
                    }
                    vowelCount = 0
                }
            }
            
            tokens = append(tokens, context.ParsedToken{
                Base: base,
                Variant: variant,
            })
        }
        
        if punctuationAfter != nullToken {
            tokens = append(tokens, punctuationAfter)
        }
        
        return tokens, true
    },
}


//When parsing and formatting French, things like punctuation have different
//spacing rules from English

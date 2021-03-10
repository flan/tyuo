package language
import (
    "github.com/flan/tyuo/context"
)


//WARNING Disregard all other notes on punctuation and the encoding thereof.
//punctuation will be encoded as standalone tokens, not part of the word to
//which it's attached.
//additionally, this means they won't be chosen as keywords, which sidesteps the
//problem of something like "tyuo, desu?" not matching "desu" in other places.
//Compound punctuation won't be learnable if it's not in the recognised set.

func digestToken(
    token []rune,
    language *language,
) ([]context.ParsedToken, bool) {
    parsedTokens := parseSymbol(token)
    if len(parsedTokens) > 0 {
        return parsedTokens, true
    }
    
    //asks the language if it's a valid word
    //response is zero or more parsed tokens; two typically means it was a word with
    //some punctuation at the end
    
    //the language's processor is responsible for resolving whether the use of punctuation was appropriate
}

func lex(
    input string,
    maxTokenLength int,
    language *language,
) ([]context.ParsedToken, bool) {
    delimiter := language.delimiter
    characters := language.characters
    
    tokens := make([]context.ParsedToken, 0, 16)
    
    var inputLearnable bool = true
    var currentToken []rune = make([]rune, 0, maxTokenLength)
    var currentTokenValid bool = true
    for _, r := range input {
        if r == delimiter {
            if len(currentToken) > 0 {
                if currentTokenValid {
                    digestedTokens, learnable := digestToken(currentToken, language)
                    if len(digestedTokens) > 0 {
                        tokens = append(tokens, digestedTokens...)
                    }
                    if !learnable {
                        inputLearnable = false
                    }
                }
                currentToken = make([]rune, 0, maxTokenLength)
            }
            currentTokenValid = true
        }
        
        if _, isCharacter := characters[r]; !isCharacter {
            if _, isPunctuation := punctuation[r]; !isPunctuation {
                if _, isSymbolRune := symbolRunes[r]; !isSymbolRune {
                    currentTokenValid = false
                    inputLearnable = false
                    continue
                }
            }
        }
        
        if currentTokenValid {
            if len(currentToken) < maxTokenLength {
                currentToken = append(currentToken, r)
            } else {
                currentTokenValid = false
                inputLearnable = false
            }
        }
    }
    if currentTokenValid && len(currentToken) > 0 {
        digestedTokens, learnable := digestToken(currentToken, language)
        if len(digestedTokens) > 0 {
            tokens = append(tokens, digestedTokens...)
        }
        if !learnable {
            inputLearnable = false
        }
    }
    
    return tokens, inputLearnable
}

func Parse(input string, ctx *context.Context) ([]context.ParsedToken, bool) {
    
    //returns the usable tokens read from input and a boolean value indicating
    //whether the input was learnable or not
}

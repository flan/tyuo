package language
import (
    "github.com/flan/tyuo/context"
    
    "golang.org/x/text/transform"
)


func digestToken(
    token []rune,
    digester func([]rune, *transform.Transformer)([]context.ParsedToken, bool),
    normaliser *transform.Transformer,
) ([]context.ParsedToken, bool) {
    parsedTokens := parseSymbol(token)
    if len(parsedTokens) > 0 {
        return parsedTokens, true
    }
    
    return digester(token, normaliser)
}

func lex(
    input string,
    learn bool,
    maxTokenLength int,
    language *languageDefinition,
) ([]context.ParsedToken, bool) {
    delimiter := language.delimiter
    characters := language.characters
    
    normaliser := context.MakeStringNormaliser()
    digester := language.digestToken
    
    tokens := make([]context.ParsedToken, 0, 16)
    
    var inputLearnable bool = true
    var currentToken []rune = make([]rune, 0, maxTokenLength)
    var currentTokenValid bool = true
    for _, r := range input {
        if r == delimiter {
            if len(currentToken) > 0 {
                if currentTokenValid {
                    digestedTokens, learnable := digestToken(currentToken, digester, normaliser)
                    if len(digestedTokens) > 0 {
                        tokens = append(tokens, digestedTokens...)
                    }
                    if !learnable {
                        if learn {
                            return nil, false
                        }
                        inputLearnable = false
                    }
                }
                currentToken = make([]rune, 0, maxTokenLength)
            }
            currentTokenValid = true
            continue
        }
        
        if _, isCharacter := characters[r]; !isCharacter {
            if _, isPunctuation := punctuation[r]; !isPunctuation {
                if _, isSymbolRune := symbolRunes[r]; !isSymbolRune {
                    if learn {
                        return nil, false
                    }
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
                if learn {
                    return nil, false
                }
                currentTokenValid = false
                inputLearnable = false
            }
        }
    }
    if currentTokenValid && len(currentToken) > 0 {
        digestedTokens, learnable := digestToken(currentToken, digester, normaliser)
        if len(digestedTokens) > 0 {
            tokens = append(tokens, digestedTokens...)
        }
        if !learnable {
            if learn {
                return nil, false
            }
            inputLearnable = false
        }
    }
    
    return tokens, inputLearnable
}

func Parse(input string, learn bool, ctx *context.Context) ([]context.ParsedToken, bool) {
    var lang *languageDefinition
    switch ctx.GetLanguage() {
        case context.LanguageEnglish:
            lang = &englishLanguageDefinition
        case context.LanguageFrench:
            lang = &frenchLanguageDefinition
        default:
            logger.Errorf("unrecognised language: %s", ctx.GetLanguage())
            return make([]context.ParsedToken, 0), false
    }
    parsedTokens, learnable := lex(input, learn, ctx.GetMaxTokenLength(), lang)
    
    if learnable {
        //one final pass over tokens to make sure there's no consecutive punctuation (which is all single-token strings now)
        previousTokenIsPunctuation := false
        for _, token := range parsedTokens {
            _, isPunctuation := context.PunctuationIdsByToken[token.Base]
            if isPunctuation && previousTokenIsPunctuation {
                return parsedTokens, false
            }
            previousTokenIsPunctuation = isPunctuation
        }
    }
    
    return parsedTokens, learnable
}

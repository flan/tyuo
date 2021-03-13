package language
import (
    "github.com/flan/tyuo/context"
)

func Format(production []int, dictionaryTokens map[int]context.DictionaryToken, ctx *context.Context) (string) {
    lang := getLanguageDefinition(ctx.GetLanguage())
    if lang == nil {
        return ""
    }
    
    return lang.formatUtterance(production, dictionaryTokens, ctx.GetProductionBaseRepresentationThreshold())
}

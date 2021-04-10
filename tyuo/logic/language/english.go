package language
import (
    "github.com/flan/tyuo/context"
    
    "golang.org/x/text/transform"
    
    "strings"
    "unicode"
)

var englishCharacters = runeset{
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
    'é': voidInstance,
    'É': voidInstance,
    'ï': voidInstance,
    apostrophe: voidInstance,
    hyphen: voidInstance,
}

const englishConsecutiveVowelLimit = 3
const englishConsecutiveConsonantLimit = 6
var englishVowelsNormalised = runeset{
    'a': voidInstance,
    'e': voidInstance,
    'i': voidInstance,
    'o': voidInstance,
    'u': voidInstance,
}


//people are very lazy, but with contractions, this can lead to divergent paths for "didn't" and "didnt",
//which is very undesireable; leave the common usage intact when displaying, but map some words to their
//correct form for identification consistency
var englishCorrections = map[string]string{
    "aint": "ain't",
    "arent": "aren't",
    "cant": "can't", //"cant" is a very obscure word
    "couldnt": "couldn't",
    "couldve": "could've",
    "didnt": "didn't",
    "doesnt": "doesn't",
    "dont": "don't",
    "hadnt": "hadn't",
    "havent": "haven't",
    "hed": "he'd",
    "hes": "he's",
    //"id": "i'd", //"id" comes up sometimes
    "im": "i'm",
    "ima": "i'mma",
    "i'ma": "i'mma",
    "imma": "i'mma",
    "isnt": "isn't",
    "ive": "i've",
    "mightve": "might've",
    "mustnt": "mustn't",
    "mustve": "must've",
    //"shed": "she'd", //the word "shed" isn't exactly obscure
    "shes": "she's",
    "souldnt": "shouldn't",
    "souldve": "should've",
    "hte": "the",
    "teh": "the",
    "their's": "theirs",
    "theres": "there's",
    "theyre": "they're",
    "theyd": "they'd",
    "theyll": "they'll",
    "theyve": "they've",
    "wasnt": "wasn't",
    "wernt": "weren't",
    "werent": "weren't",
    "weve": "we've",
    "wheres": "where's",
    "whos": "who's",
    "wont": "won't", //"wont" is a very obscure word
    "wouldnt": "wouldn't",
    "yall": "y'all",
    "youll": "you'll",
}

type englishErrorWordFragment struct {
    correct string
    incorrect []string
}
//and handle some commonly misspelled word-fragments
var englishErrorWordFragments = []englishErrorWordFragment{
    englishErrorWordFragment{
        correct: "acceptabl",
        incorrect: []string{"acceptibl"},
    },
    englishErrorWordFragment{
        correct: "accidental",
        incorrect: []string{"accidentall", "accidentl"},
    },
    englishErrorWordFragment{
        correct: "accommodat",
        incorrect: []string{"accomodat", "acommodat"},
    },
    englishErrorWordFragment{
        correct: "achiev",
        incorrect: []string{"acheiv"},
    },
    englishErrorWordFragment{
        correct: "acknowledg",
        incorrect: []string{"acknowleg", "aknowledg"},
    },
    englishErrorWordFragment{
        correct: "aggress",
        incorrect: []string{"agress"},
    },
    englishErrorWordFragment{
        correct: "almost",
        incorrect: []string{"allmost"},
    },
    englishErrorWordFragment{
        correct: "annual",
        incorrect: []string{"anual"},
    },
    englishErrorWordFragment{
        correct: "apparent",
        incorrect: []string{"apparant", "aparent", "apparrent", "aparrent"},
    },
    englishErrorWordFragment{
        correct: "arctic",
        incorrect: []string{"artic"},
    },
    englishErrorWordFragment{
        correct: "argument",
        incorrect: []string{"arguement"},
    },
    englishErrorWordFragment{
        correct: "atheist",
        incorrect: []string{"athiest", "athist"},
    },
    
    
    englishErrorWordFragment{
        correct: "barbecue",
        incorrect: []string{"bbq", "barbeque", "barbequeue"},
    },
    englishErrorWordFragment{
        correct: "because",
        incorrect: []string{"beatiful"},
    },
    englishErrorWordFragment{
        correct: "beginning",
        incorrect: []string{"begining"},
    },
    englishErrorWordFragment{
        correct: "belie",
        incorrect: []string{"belei"},
    },
    englishErrorWordFragment{
        correct: "business",
        incorrect: []string{"buisness"},
    },
    
    
    englishErrorWordFragment{
        correct: "calendar",
        incorrect: []string{"calender"},
    },
    englishErrorWordFragment{
        correct: "categor",
        incorrect: []string{"catagor"},
    },
    englishErrorWordFragment{
        correct: "cemetery",
        incorrect: []string{"cemetary", "cematery"},
    },
    englishErrorWordFragment{
        correct: "congratulat",
        incorrect: []string{"congradulat"},
    },
    englishErrorWordFragment{
        correct: "conscious",
        incorrect: []string{"concious", "consious"},
    },
    englishErrorWordFragment{
        correct: "controvers",
        incorrect: []string{"contravers"},
    },
    englishErrorWordFragment{
        correct: "decei",
        incorrect: []string{"decie"},
    },
    englishErrorWordFragment{
        correct: "definit",
        incorrect: []string{"definat"},
    },
    englishErrorWordFragment{
        correct: "desper",
        incorrect: []string{"despar"},
    },
    englishErrorWordFragment{
        correct: "differ",
        incorrect: []string{"diffr"},
    },
    
    
    englishErrorWordFragment{
        correct: "embarrass",
        incorrect: []string{"embarass"},
    },
    englishErrorWordFragment{
        correct: "existen",
        incorrect: []string{"existan"},
    },
    englishErrorWordFragment{
        correct: "experien",
        incorrect: []string{"experian"},
    },
    
    
    englishErrorWordFragment{
        correct: "foreign",
        incorrect: []string{"foriegn"},
    },
    englishErrorWordFragment{
        correct: "friend",
        incorrect: []string{"freind"},
    },
    englishErrorWordFragment{
        correct: "fulfil",
        incorrect: []string{"fullfil", "fulfill"},
    },
    
    
    englishErrorWordFragment{
        correct: "gauge",
        incorrect: []string{"guage"},
    },
    englishErrorWordFragment{
        correct: "guida",
        incorrect: []string{"guide"},
    },
    
    
    englishErrorWordFragment{
        correct: "harass",
        incorrect: []string{"harrass"},
    },
    englishErrorWordFragment{
        correct: "heroes",
        incorrect: []string{"heros"},
    },
    englishErrorWordFragment{
        correct: "hygien",
        incorrect: []string{"hygen", "hygein"},
    },
    englishErrorWordFragment{
        correct: "hypocri",
        incorrect: []string{"hipocrit", "hippocrit"},
    },
    
    
    englishErrorWordFragment{
        correct: "ignoran",
        incorrect: []string{"ignoren"},
    },
    englishErrorWordFragment{
        correct: "independent",
        incorrect: []string{"independant"},
    },
    englishErrorWordFragment{
        correct: "indispensabl",
        incorrect: []string{"indispensibl"},
    },
    englishErrorWordFragment{
        correct: "inoculat",
        incorrect: []string{"innoculat"},
    },
    
    
    englishErrorWordFragment{
        correct: "jewelry",
        incorrect: []string{"jewelery"},
    },
    englishErrorWordFragment{
        correct: "judgment",
        incorrect: []string{"judgement"},
    },
    
    
    englishErrorWordFragment{
        correct: "kernel",
        incorrect: []string{"kernal"},
    },
    
    
    englishErrorWordFragment{
        correct: "necessar",
        incorrect: []string{"neccessar"},
    },
    englishErrorWordFragment{
        correct: "niece",
        incorrect: []string{"neice"},
    },
    englishErrorWordFragment{
        correct: "notice",
        incorrect: []string{"notica"},
    },
    
    
    englishErrorWordFragment{
        correct: "occasion",
        incorrect: []string{"occassion"},
    },
    englishErrorWordFragment{
        correct: "occurre",
        incorrect: []string{"occurra", "occure"},
    },
    englishErrorWordFragment{
        correct: "omission",
        incorrect: []string{"ommision", "omision"},
    },
    
    
    englishErrorWordFragment{
        correct: "pastime",
        incorrect: []string{"passtime", "pasttime"},
    },
    englishErrorWordFragment{
        correct: "personnel",
        incorrect: []string{"personell", "personel"},
    },
    englishErrorWordFragment{
        correct: "possess",
        incorrect: []string{"posess", "posses"},
    },
    englishErrorWordFragment{
        correct: "potatoes",
        incorrect: []string{"potatos"},
    },
    englishErrorWordFragment{
        correct: "privilege",
        incorrect: []string{"privelege", "priviledge"},
    },
    
    
    englishErrorWordFragment{
        correct: "publicly",
        incorrect: []string{"publically"},
    },
    
    
    englishErrorWordFragment{
        correct: "quarantine",
        incorrect: []string{"quarentine"},
    },
    englishErrorWordFragment{
        correct: "queue",
        incorrect: []string{"que"},
    },
    
    
    englishErrorWordFragment{
        correct: "receive",
        incorrect: []string{"recieve"},
    },
    englishErrorWordFragment{
        correct: "receipt",
        incorrect: []string{"reciept"},
    },
    englishErrorWordFragment{
        correct: "recommend",
        incorrect: []string{"recomend", "reccommend"},
    },
    englishErrorWordFragment{
        correct: "relevan",
        incorrect: []string{"releven"},
    },
    
    
    englishErrorWordFragment{
        correct: "restaurant",
        incorrect: []string{"restarant", "restaraunt"},
    },
    englishErrorWordFragment{
        correct: "rhythm",
        incorrect: []string{"rythm", "rythem"},
    },
    
    
    englishErrorWordFragment{
        correct: "separate",
        incorrect: []string{"seperate"},
    },
    englishErrorWordFragment{
        correct: "speech",
        incorrect: []string{"speach"},
    },
    
    
    englishErrorWordFragment{
        correct: "surpris",
        incorrect: []string{"supris"},
    },
    
    
    englishErrorWordFragment{
        correct: "tomatoes",
        incorrect: []string{"tomatos"},
    },
    
    
    englishErrorWordFragment{
        correct: "tomorrow",
        incorrect: []string{"tommorow", "tommorrow"},
    },
    
    
    englishErrorWordFragment{
        correct: "vacuum",
        incorrect: []string{"vaccuum", "vaccum", "vacume"},
    },
    
    
    englishErrorWordFragment{
        correct: "weird",
        incorrect: []string{"wierd"},
    },
    
    
    englishErrorWordFragment{
        correct: "zeroes",
        incorrect: []string{"zeros"},
    },
}


var englishLanguageDefinition = languageDefinition{
    delimiter: ' ',
    characters: englishCharacters,
    
    digestToken: func(token []rune, normaliser *transform.Transformer) ([]context.ParsedToken, bool) {
        tokens := make([]context.ParsedToken, 0, 2)
        punctuationBefore, punctuationAfter, token, learnable := punctuationDissect(token)
        if !learnable {
            return nil, false
        }
        
        if punctuationBefore != nullToken {
            tokens = append(tokens, punctuationBefore)
        }
        
        if len(token) > 0 { //not fully digested
            //punctuationDissect will deal with leading/trailing hyphens, so just check for apostrophes
            if token[0] == apostrophe || token[len(token) - 1] == apostrophe {
                //a learnable token can't be bounded by an apostrophe, since it might be a quotation mark
                return nil, false
            }
            
            containsPunctuation := false
            for _, r := range token {
                if _, isCharacter := englishCharacters[r]; !isCharacter {
                    return nil, false
                }
                if r == apostrophe || r == hyphen {
                    if containsPunctuation { //allow at most one punctuation-mark, to limit abuse
                        return nil, false
                    }
                    containsPunctuation = true
                }
            }
            
            //get the normalised base-form
            variant := string(token)
            base, _, err := transform.String(*normaliser, variant)
            if err != nil {
                logger.Warningf("unable to normalise token %s: %s", variant, err)
                return nil, false
            }
            
            //check to make sure there aren't too many vowels or consonants clumped together; this likely indicates gibberish
            vowelCount := 0
            consonantCount := 0
            for _, r := range base {
                if _, isVowel := englishVowelsNormalised[r]; isVowel {
                    vowelCount++
                    if vowelCount > englishConsecutiveVowelLimit {
                        return nil, false
                    }
                    consonantCount = 0
                } else {
                    consonantCount++
                    if consonantCount > englishConsecutiveConsonantLimit {
                        return nil, false
                    }
                    vowelCount = 0
                }
            }
            
            //see if it's a word with a direct correction
            if correctedForm, defined := englishCorrections[base]; defined {
                base = correctedForm
            } else { //account for common spelling errors to reduce n-gram spread
                for _, eewf := range englishErrorWordFragments {
                    replacementMade := false
                    for _, incorrect := range eewf.incorrect {
                        if strings.Contains(base, incorrect) {
                            base = strings.Replace(base, incorrect, eewf.correct, 1)
                            replacementMade = true
                            break
                        }
                    }
                    if replacementMade {
                        break
                    }
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
    
    formatUtterance: func(production []int, dictionaryTokens map[int]context.DictionaryToken, baseRepresentationThreshold float32) (string) {
        var output strings.Builder
        
        startOfSentence := true
        spaceRequired := false
        for i, id := range production {
            spaceRequired = true
            
            //see if it's punctuation
            if punctuation, defined := context.PunctuationTokensById[id]; defined {
                switch punctuation {
                    case "…":
                        output.WriteString(punctuation)
                        if i != 0 {
                            spaceRequired = false
                        }
                    case ".", "?", "!", "⁈", "‼", "⁇":
                        output.WriteString(punctuation)
                        startOfSentence = true
                    case "—", "&":
                        output.WriteByte(' ')
                        output.WriteString(punctuation)
                    default:
                        output.WriteString(punctuation)
                }
                continue
            }
            
            //since whatever's left is either a word or a symbol, add a space if it makes sense to do so
            if spaceRequired && i > 0 {
                output.WriteByte(' ')
            }
            
            //see if it's a symbol
            if symbol, defined := context.SymbolsTokensById[id]; defined {
                output.WriteString(symbol)
                continue
            }
            
            //it must be a word
            if dictionaryToken, defined := dictionaryTokens[id]; defined {
                representation, isBase := dictionaryToken.Represent(baseRepresentationThreshold)
                if isBase && startOfSentence {
                    var head rune
                    var tail []rune
                    for j, r := range representation {
                        if j == 0 {
                            head = r
                        } else {
                            tail = append(tail, r)
                        }
                    }
                    output.WriteRune(unicode.ToUpper(head))
                    output.WriteString(string(tail))
                } else {
                    output.WriteString(representation)
                }
                startOfSentence = false
            } else {
                logger.Errorf("unable to resolve dictionary token for %d; the database is inconsistent", id)
                output.Reset()
                break
            }
        }
        
        return output.String()
    },
}

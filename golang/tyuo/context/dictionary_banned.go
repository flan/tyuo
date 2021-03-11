package context
import (
    "bufio"
    "os"
    "strings"
    
    "golang.org/x/text/transform"
)

type bannedToken struct {
    baseRepresentation string
    dictionaryId int
}


func processBannedSubstrings(listPath string) ([]string, error) {
    file, err := os.Open(listPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    normaliser := MakeStringNormaliser()
    
    output := make([]string, 0)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        substring, _, err := transform.String(*normaliser, strings.TrimSpace(scanner.Text()))
        if err != nil {
            return nil, err
        }
        if len(substring) > 0 {
            output = append(output, substring)
        }
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    logger.Debugf("loaded %d language-level banned substrings", len(output))
    return output, nil
}


type bannedDictionary struct {
    database *database

    //words from database
    bannedTokens []bannedToken
    bannedIds map[int]void

    //tokens from the list
    bannedSubstringsGeneric []string
    bannedIdsGeneric map[int]void
}
func prepareBannedDictionary(
    database *database,
    bannedSubstringsGeneric []string,
) (*bannedDictionary, error) {
    bannedTokens := make([]bannedToken, 0)
    bannedIds := make(map[int]void)
    if bts, err := database.bannedLoadBannedTokens(nil); err == nil {
        for _, bt := range bts {
            bannedTokens = append(bannedTokens, bt)
            if bt.dictionaryId != undefinedDictionaryId {
                bannedIds[bt.dictionaryId] = voidInstance
            }
        }
    } else {
        return nil, err
    }
    logger.Debugf("loaded %d banned tokens", len(bannedTokens))
    logger.Debugf("loaded %d banned IDs", len(bannedIds))

    //enumerate the IDs of anything in the dictionary that predated additions to the language-level ban-list
    bannedIdsGeneric := make(map[int]void)
    if bvs, err := database.dictionaryEnumerateTokensBySubstring(bannedSubstringsGeneric); err == nil {
        for _, bannedId := range bvs {
            bannedIdsGeneric[bannedId] = voidInstance
        }
    } else {
        return nil, err
    }
    logger.Debugf("identified %d IDs mapped to banned language-level tokens", len(bannedIdsGeneric))

    return &bannedDictionary{
        database: database,

        bannedTokens: bannedTokens,
        bannedIds: bannedIds,

        bannedSubstringsGeneric: bannedSubstringsGeneric,
        bannedIdsGeneric: bannedIdsGeneric,
    }, nil
}
func (bd *bannedDictionary) ban(substrings stringset) (error) {
    normaliser := MakeStringNormaliser()
    
    bannedSubstrings := make([]string, 0, len(substrings))
    for substring := range substrings {
        normalisedSubstring, _, err := transform.String(*normaliser, strings.TrimSpace(substring))
        if err != nil {
            return err
        }
        if len(normalisedSubstring) == 0 {
            continue
        }

        alreadyBanned := false
        for _, bt := range bd.bannedTokens {
            if bt.baseRepresentation == normalisedSubstring {
                alreadyBanned = true
                break
            }
        }
        if !alreadyBanned {
            bannedSubstrings = append(bannedSubstrings, normalisedSubstring)
        }
    }
    if len(bannedSubstrings) == 0 {
        return nil
    }
    logger.Infof("banning %d substrings: %v...", len(bannedSubstrings), bannedSubstrings)
    
    if newlyBannedTokens, err := bd.database.bannedBanSubstrings(bannedSubstrings); err == nil {
        for _, bt := range newlyBannedTokens {
            bd.bannedTokens = append(bd.bannedTokens, bt)
            if bt.dictionaryId != undefinedDictionaryId {
                if bt.dictionaryId != undefinedDictionaryId {
                    bd.bannedIds[bt.dictionaryId] = voidInstance
                }
            }
        }
    } else {
        return err
    }
logger.Criticalf("%v", bd.bannedTokens)
    return nil
}
func (bd *bannedDictionary) unban(substrings stringset) (error) {
    normaliser := MakeStringNormaliser()
    
    bannedTokenIndexes := make([]int, 0, len(substrings))
    bannedSubstrings := make([]string, 0, len(substrings))
    for substring := range substrings {
        normalisedSubstring, _, err := transform.String(*normaliser, strings.TrimSpace(substring))
        if err != nil {
            return err
        }
        if len(normalisedSubstring) == 0 {
            continue
        }

        for idx, bt := range bd.bannedTokens {
            if bt.baseRepresentation == normalisedSubstring {
                bannedTokenIndexes = append(bannedTokenIndexes, idx)
                bannedSubstrings = append(bannedSubstrings, normalisedSubstring)
                break
            }
        }
    }
    if len(bannedSubstrings) == 0 {
        return nil
    }
    logger.Infof("unbanning %d substrings: %v...", len(bannedSubstrings), bannedSubstrings)

    if err := bd.database.bannedUnbanTokens(bannedSubstrings); err != nil {
        return err
    }

    for i, idx := range bannedTokenIndexes {
        delete(bd.bannedIds, bd.bannedTokens[idx].dictionaryId)

        //move an element from the tail over the one to be removed
        bd.bannedTokens[idx] = bd.bannedTokens[len(bd.bannedTokens) - (1 + i)]
    }
    //cut the tail
    bd.bannedTokens = bd.bannedTokens[:len(bd.bannedTokens) - len(bannedTokenIndexes)]
    
logger.Criticalf("%v", bd.bannedTokens)
    return nil
}
func (bd *bannedDictionary) containsBannedToken(s string) (bool) {
    for _, bs := range bd.bannedSubstringsGeneric {
        if strings.Contains(s, bs) {
            return true
        }
    }
    for _, bt := range bd.bannedTokens {
        if strings.Contains(s, bt.baseRepresentation) {
            return true
        }
    }
    return false
}
func (bd *bannedDictionary) getIdBannedStatus(id int) (bool) {
    if _, defined := bd.bannedIds[id]; defined {
        return true
    }
    if _, defined := bd.bannedIdsGeneric[id]; defined {
        return true
    }
    return false
}
func (bd *bannedDictionary) getIdsBannedStatus(ids intset) (map[int]bool) {
    results := make(map[int]bool, len(ids))
    for id := range ids {
        results[id] = bd.getIdBannedStatus(id)
    }
    return results
}

package context
import (
    "bufio"
    "os"
    "strings"
)

type bannedToken struct {
    baseRepresentation string
    dictionaryId int
}


func processBannedTokens(listPath string) ([]string, error) {
    file, err := os.Open(listPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    output := make([]string, 0)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        token := strings.ToLower(strings.TrimSpace(scanner.Text()))
        if len(token) > 0 {
            output = append(output, token)
        }
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    logger.Debugf("loaded {} language-level banned tokens", len(output))
    return output, nil
}


type bannedDictionary struct {
    database *database

    //words from database
    bannedTokens []bannedToken
    bannedIds map[int]void

    //tokens from the list
    bannedTokensGeneric []string
    bannedIdsGeneric map[int]void
}
func prepareBannedDictionary(
    database *database,
    bannedTokensGeneric []string,
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
    logger.Debugf("loaded {} banned tokens", len(bannedTokens))
    logger.Debugf("loaded {} banned IDs", len(bannedIds))

    //enumerate the IDs of anything in the dictionary that predated additions to the language-level ban-list
    bannedIdsGeneric := make(map[int]void)
    if bvs, err := database.dictionaryEnumerateTokensBySubstring(bannedTokensGeneric); err == nil {
        for _, bannedId := range bvs {
            bannedIdsGeneric[bannedId] = voidInstance
        }
    } else {
        return nil, err
    }
    logger.Debugf("identified {} IDs mapped to banned language-level tokens", len(bannedIdsGeneric))

    return &bannedDictionary{
        database: database,

        bannedTokens: bannedTokens,
        bannedIds: bannedIds,

        bannedTokensGeneric: bannedTokensGeneric,
        bannedIdsGeneric: bannedIdsGeneric,
    }, nil
}
func (bd *bannedDictionary) ban(tokens stringset) (error) {
    bannedTokens := make([]string, 0, len(tokens))
    for token := range tokens {
        lcaseToken := strings.ToLower(token)

        alreadyBanned := false
        for _, bt := range bd.bannedTokens {
            if bt.baseRepresentation == lcaseToken {
                alreadyBanned = true
                break
            }
        }
        if !alreadyBanned {
            bannedTokens = append(bannedTokens, lcaseToken)
        }
    }
    if len(bannedTokens) == 0 {
        return nil
    }

    if newlyBannedTokens, err := bd.database.bannedBanTokens(bannedTokens); err == nil {
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
    return nil
}
func (bd *bannedDictionary) unban(tokens stringset) (error) {
    bannedTokenIndexes := make([]int, 0, len(tokens))
    bannedTokens := make([]string, 0, len(tokens))
    for token := range tokens {
        lcaseToken := strings.ToLower(token)

        for idx, bt := range bd.bannedTokens {
            if bt.baseRepresentation == lcaseToken {
                bannedTokenIndexes = append(bannedTokenIndexes, idx)
                bannedTokens = append(bannedTokens, lcaseToken)
                break
            }
        }
    }
    if len(bannedTokens) == 0 {
        return nil
    }

    if err := bd.database.bannedUnbanTokens(bannedTokens); err != nil {
        return err
    }

    for i, idx := range bannedTokenIndexes {
        delete(bd.bannedIds, bd.bannedTokens[idx].dictionaryId)

        //move an element from the tail over the one to be removed
        bd.bannedTokens[idx] = bd.bannedTokens[len(bd.bannedTokens) - (1 + i)]
    }
    //cut the tail
    bd.bannedTokens = bd.bannedTokens[:len(bd.bannedTokens) - len(bannedTokenIndexes)]

    return nil
}
func (bd *bannedDictionary) containsBannedToken(s string) (bool) {
    for _, bt := range bd.bannedTokensGeneric {
        if strings.Contains(s, bt) {
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

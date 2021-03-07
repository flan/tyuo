package context
import (
    "database/sql"
    "encoding/json"
    "flag"
    "fmt"
    "github.com/4kills/go-zlib"
    "os"
    "path/filepath"
    "strings"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

var dbDebug = flag.Bool("db-debug", false, "whether to use database debugging features; should usually be false")

func prepareSqliteArrayParams(start int, count int) (string) {
    arrayParams := make([]string, count)
    for i := 0; i < count; i++ {
        arrayParams[i] = fmt.Sprintf("?%d", start + i)
    }
    return strings.Join(arrayParams, ",")
}

func stringSliceToInterfaceSlice(s []string) ([]interface{}) {
    output := make([]interface{}, len(s))
    for i, v := range(s) {
        output[i] = v
    }
    return output
}
func intSliceToInterfaceSlice(s []int) ([]interface{}) {
    output := make([]interface{}, len(s))
    for i, v := range(s) {
        output[i] = v
    }
    return output
}
func stringSetToInterfaceSlice(s map[string]bool) ([]interface{}) {
    output := make([]interface{}, 0, len(s))
    for k, _ := range(s) {
        output = append(output, k)
    }
    return output
}
func intSetToInterfaceSlice(s map[int]bool) ([]interface{}) {
    output := make([]interface{}, 0, len(s))
    for k, _ := range(s) {
        output = append(output, k)
    }
    return output
}




type Database struct {
    connection *sql.DB
}
func prepareDatabase(
    dbPath string,
) (*Database, error) {
    connection, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    //SQLite databases should only be opened once per process, so disable Go's pooling
    connection.SetMaxOpenConns(1)
    
    logger.Debugf("preparing database structures...");
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS dictionary (
        caseInsensitiveRepresentation TEXT NOT NULL UNIQUE,
        id INTEGER NOT NULL PRIMARY KEY,
        caseInsensitiveOccurrences INTEGER NOT NULL,
        capitalisedFormsJSON TEXT
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS dictionary_banned (
        caseInsensitiveRepresentation TEXT NOT NULL PRIMARY KEY
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS terminals (
        dictionaryId INTEGER NOT NULL PRIMARY KEY,
        lastObservedForward INTEGER, --UNIX timestamp
        lastObservedReverse INTEGER, --UNIX timestamp

        FOREIGN KEY(dictionaryId)
        REFERENCES dictionary(id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    
    //for n-grams, the JSON structure will never be empty, since there
    //has to be at least one transition for a write to occur
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS digrams_forward (
        dictionaryIdFirst INTEGER NOT NULL,
        transitionsJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst),
        FOREIGN KEY(dictionaryIdFirst)
        REFERENCES dictionary(id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS digrams_reverse (
        dictionaryIdFirst INTEGER NOT NULL,
        transitionsJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst),
        FOREIGN KEY(dictionaryIdFirst)
        REFERENCES dictionary(id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS trigrams_forward (
        dictionaryIdFirst INTEGER NOT NULL,
        dictionaryIdSecond INTEGER NOT NULL,
        transitionsJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst, dictionaryIdSecond),
        FOREIGN KEY(dictionaryIdFirst, dictionaryIdSecond)
        REFERENCES dictionary(id, id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS trigrams_reverse (
        dictionaryIdFirst INTEGER NOT NULL,
        dictionaryIdSecond INTEGER NOT NULL,
        transitionsJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst, dictionaryIdSecond),
        FOREIGN KEY(dictionaryIdFirst, dictionaryIdSecond)
        REFERENCES dictionary(id, id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS quadgrams_forward (
        dictionaryIdFirst INTEGER NOT NULL,
        dictionaryIdSecond INTEGER NOT NULL,
        dictionaryIdThird INTEGER NOT NULL,
        transitionsJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird),
        FOREIGN KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird)
        REFERENCES dictionary(id, id, id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS quadgrams_reverse (
        dictionaryIdFirst INTEGER NOT NULL,
        dictionaryIdSecond INTEGER NOT NULL,
        dictionaryIdThird INTEGER NOT NULL,
        transitionsJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird),
        FOREIGN KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird)
        REFERENCES dictionary(id, id, id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS quintgrams_forward (
        dictionaryIdFirst INTEGER NOT NULL,
        dictionaryIdSecond INTEGER NOT NULL,
        dictionaryIdThird INTEGER NOT NULL,
        dictionaryIdFourth INTEGER NOT NULL,
        transitionsJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird, dictionaryIdFourth),
        FOREIGN KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird, dictionaryIdFourth)
        REFERENCES dictionary(id, id, id, id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS quintgrams_reverse (
        dictionaryIdFirst INTEGER NOT NULL,
        dictionaryIdSecond INTEGER NOT NULL,
        dictionaryIdThird INTEGER NOT NULL,
        dictionaryIdFourth INTEGER NOT NULL,
        transitionsJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird, dictionaryIdFourth),
        FOREIGN KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird, dictionaryIdFourth)
        REFERENCES dictionary(id, id, id, id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    
    logger.Debugf("preparing database pragma...");
    //while foreign keys are declared in the structure, because tokens are never
    //removed from the database, their enforcement is unnecessary outside of debugging
    foreignKeys := "OFF"
    if *dbDebug {
        foreignKeys = "ON"
    }
    if _, err = connection.Exec(fmt.Sprintf("PRAGMA foreign_keys = %s", foreignKeys)); err != nil {
        connection.Close()
        return nil, err
    }
    //since only lower-case matches occur, make comparisons more efficient
    if _, err = connection.Exec("PRAGMA case_sensitive_like = TRUE"); err != nil {
        connection.Close()
        return nil, err
    }

    return &Database{
        connection: connection,
    }, nil
}
func (db *Database) Close() (error) {
    return db.connection.Close()
}




func deserialiseCapitalisedFormsJSON(data *sql.NullString) (map[string]int) {
    if data.Valid {
        var buffer map[string]int = nil
        if err := json.Unmarshal([]byte(data.String), &buffer); err == nil {
            return buffer
        } else { //some sort of database corruption, almost certainly due to misuse
            logger.Warningf("unable to deserialise capitalisedFormsJSON; reinitialising state: %s", err)
        }
    }
    return make(map[string]int)
}
func serialiseCapitalisedFormsJSON(data map[string]int) (interface{}) {
    if len(data) == 0 {
        return nil
    }
    
    if buffer, err := json.Marshal(data); err == nil {
        return buffer
    } else {
        logger.Warningf("unable to serialise capitalisedFormsJSON; reinitialising state: %s", err)
        return nil
    }
}
func (db *Database) dictionaryEnumerateTokensBySubstring(tokens []string) (map[string]int, error) {
    if stmt, err := db.connection.Prepare(`
    SELECT
        caseInsensitiveRepresentation,
        id
    FROM
        dictionary
    WHERE
        caseInsensitiveRepresentation LIKE ?1
    `); err == nil {
        defer stmt.Close()
        
        output := make(map[string]int)
        for _, token := range tokens {
            if rows, err := stmt.Query(fmt.Sprintf("%%%s%%", token)); err == nil {
                defer rows.Close()
                for rows.Next() {
                    var cir string
                    var did int
                    if err:= rows.Scan(&cir, &did); err == nil {
                        output[cir] = did
                    } else {
                        return nil, err
                    }
                }
            } else {
                return nil, err
            }
        }
        return output, nil
    } else {
        return nil, err
    }
}
func processDictionaryRows(maxCount int, rows *sql.Rows) ([]dictionaryToken, error) {
    output := make([]dictionaryToken, 0, maxCount)
    for rows.Next() {
        var cir string
        var did int
        var cio int
        var cfj sql.NullString
        if err:= rows.Scan(&cir, &did, &cio, &cfj); err == nil {
            output = append(output, dictionaryToken{
                id: did,
                caseInsensitiveRepresentation: cir,
                caseInsensitiveOccurrences: cio,
                capitalisedForms: deserialiseCapitalisedFormsJSON(&cfj),
            })
        } else {
            return nil, err
        }
    }
    return output, nil
}
func (db *Database) dictionaryGetTokensByToken(tokens map[string]bool) ([]dictionaryToken, error) {
    if len(tokens) == 0 {
        return make([]dictionaryToken, 0), nil
    }
    
    query := fmt.Sprintf(`
    SELECT
        caseInsensitiveRepresentation,
        id,
        caseInsensitiveOccurrences,
        capitalisedFormsJSON
    FROM
        dictionary
    WHERE
        caseInsensitiveRepresentation IN (%s)
    LIMIT %d
    `, prepareSqliteArrayParams(1, len(tokens)), len(tokens))
    
    if rows, err := db.connection.Query(
        query,
        stringSetToInterfaceSlice(tokens)...,
    ); err == nil {
        defer rows.Close()
        return processDictionaryRows(len(tokens), rows)
    } else {
        return nil, err
    }
}
func (db *Database) dictionaryGetTokensById(ids map[int]bool) ([]dictionaryToken, error) {
    if len(ids) == 0 {
        return make([]dictionaryToken, 0), nil
    }
    
    query := fmt.Sprintf(`
    SELECT
        caseInsensitiveRepresentation,
        id,
        caseInsensitiveOccurrences,
        capitalisedFormsJSON
    FROM
        dictionary
    WHERE
        id IN (%s)
    LIMIT %d
    `, prepareSqliteArrayParams(1, len(ids)), len(ids))
    
    if rows, err := db.connection.Query(
        query,
        intSetToInterfaceSlice(ids)...,
    ); err == nil {
        defer rows.Close()
        return processDictionaryRows(len(ids), rows)
    } else {
        return nil, err
    }
}
func (db *Database) dictionarySetTokens(tokens []*dictionaryToken) (error) {
    if len(tokens) == 0 {
        return nil
    }
    
    tx, err := db.connection.Begin()
    if err != nil {
        return err
    }
    
    if stmt, err := tx.Prepare(`
    INSERT INTO dictionary(
        caseInsensitiveRepresentation,
        id,
        caseInsensitiveOccurrences,
        capitalisedFormsJSON
    ) VALUES (?1, ?2, ?3, ?4)
    ON CONFLICT(id) DO UPDATE SET
        caseInsensitiveOccurrences = ?5,
        capitalisedFormsJSON = ?6
    `); err == nil {
        for _, token := range tokens {
            cfj := serialiseCapitalisedFormsJSON(token.capitalisedForms)
            
            if _, err = stmt.Exec(
                token.caseInsensitiveRepresentation,
                token.id,
                token.caseInsensitiveOccurrences,
                cfj,
                //update-case:
                token.caseInsensitiveOccurrences,
                cfj,
            ); err != nil {
                if e := stmt.Close(); e != nil {
                    logger.Warningf("unable to close statement: %s", e)
                }
                if e := tx.Rollback(); e != nil {
                    logger.Warningf("unable to roll-back transaction: %s", e)
                }
                return err
            }
        }
        stmt.Close()
    } else {
        if e := tx.Rollback(); e != nil {
            logger.Warningf("unable to roll-back transaction: %s", e)
        }
        return err
    }
    return tx.Commit()
}
func (db *Database) dictionaryGetNextIdentifier() (int, error) {
    var maxIdentifier = undefinedDictionaryId //lowest allowable identifier, used to initialise dictionaries
    const query = "SELECT MAX(id) FROM dictionary"
    row := db.connection.QueryRow(query)
    if err := row.Scan(&maxIdentifier); err != nil {
        if err != sql.ErrNoRows {
            return maxIdentifier, err
        }
    }
    return maxIdentifier + 1, nil
}




func (db *Database) bannedLoadBannedTokens(
    tokenSubset []string,
) ([]bannedToken, error) {
    query := `
    SELECT
        banned.caseInsensitiveRepresentation,
        dict.id
    FROM
        dictionary_banned AS banned
    LEFT JOIN dictionary AS dict ON
        banned.caseInsensitiveRepresentation = dict.caseInsensitiveRepresentation
    `
    
    if len(tokenSubset) > 0 {
        query += fmt.Sprintf(
            "WHERE banned.caseInsensitiveRepresentation IN (%s) LIMIT %d",
            prepareSqliteArrayParams(1, len(tokenSubset)),
            len(tokenSubset),
        )
    }
    if rows, err := db.connection.Query(
        query,
        stringSliceToInterfaceSlice(tokenSubset)...,
    ); err == nil {
        defer rows.Close()
        
        output := make([]bannedToken, 0)
        for rows.Next() {
            var cir string
            var did int
            if err:= rows.Scan(&cir, &did); err == nil {
                output = append(output, bannedToken{
                    caseInsensitiveRepresentation: cir,
                    dictionaryId: did,
                })
            } else {
                return nil, err
            }
        }
        return output, nil
    } else {
        return nil, err
    }
}
func (db *Database) bannedBanTokens(tokens []string) ([]bannedToken, error) {
    tx, err := db.connection.Begin()
    if err != nil {
        return nil, err
    }
    
    const query = `
    INSERT INTO
        dictionary_banned(caseInsensitiveRepresentation)
    VALUES (?1)
    ON CONFLICT DO NOTHING
    `
    if stmt, err := tx.Prepare(query); err == nil {
        for _, token := range tokens {
            if _, err = stmt.Exec(token); err != nil {
                break
            }
        }
        if e := stmt.Close(); e != nil {
            logger.Warningf("unable to close statement: %s", e)
        }
    }
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    if err = tx.Commit(); err != nil {
        return nil, err
    }
    return db.bannedLoadBannedTokens(tokens);
}
func (db *Database) bannedUnbanTokens(tokens []string) (error) {
    query := fmt.Sprintf(`
    DELETE FROM
        dictionary_banned
    WHERE caseInsensitiveRepresentation IN (%s)
    `, prepareSqliteArrayParams(1, len(tokens)))
    
    _, err := db.connection.Exec(query, stringSliceToInterfaceSlice(tokens)...)
    return err
}




func (db *Database) terminalsGetTerminals(ids map[int]bool, oldestAllowedTime int64) (map[int]Terminal, error) {
    if len(ids) == 0 {
        return make(map[int]Terminal, 0), nil
    }
    
    remainingIds := make(map[int]bool, len(ids))
    for id := range ids {
        remainingIds[id] = false
    }
    
    query := fmt.Sprintf(`
    SELECT
        dictionaryId,
        lastObservedForward,
        lastObservedReverse
    FROM
        terminals
    WHERE
        id IN (%s)
    LIMIT %d
    `, prepareSqliteArrayParams(1, len(ids)), len(ids))
    
    if rows, err := db.connection.Query(
        query,
        intSetToInterfaceSlice(ids)...,
    ); err == nil {
        defer rows.Close()
        results := make(map[int]Terminal, len(ids))
        for rows.Next() {
            var did int
            var lof, lor sql.NullInt64
            if err:= rows.Scan(&did, &lof, &lor); err == nil {
                results[did] = Terminal{
                    dictionaryId: did,
                    
                    Forward: lof.Valid && lof.Int64 > oldestAllowedTime,
                    Reverse: lor.Valid && lor.Int64 > oldestAllowedTime,
                }
                delete(remainingIds, did)
            } else {
                return nil, err
            }
        }
        for did := range remainingIds {
            results[did] = Terminal{
                dictionaryId: did,
                
                Forward: false,
                Reverse: false,
            }
        }
        return results, nil
    } else {
        return nil, err
    }
}
func (db *Database) terminalsSetStatus(terminals []*Terminal) (error) {
    if len(terminals) == 0 {
        return nil
    }
    
    currentTime := time.Now().Unix()
    
    tx, err := db.connection.Begin()
    if err != nil {
        return err
    }
    
    //it should be quite rare that the same symbol is both the forward- and
    //reverse-terminal while learning, so the logic was simplified a bit, opting
    //to double-execute if that case happens
    if stmtForward, err := tx.Prepare(`
    INSERT INTO terminals(
        dictionaryId,
        lastObservedForward,
        lastObservedReverse
    ) VALUES (?1, ?2, NULL)
    ON CONFLICT(dictionaryId) DO UPDATE SET
        lastObservedForward = ?3
    `); err == nil {
        if stmtReverse, err := tx.Prepare(`
        INSERT INTO terminals(
            dictionaryId,
            lastObservedForward,
            lastObservedReverse
        ) VALUES (?1, NULL, ?2)
        ON CONFLICT(dictionaryId) DO UPDATE SET
            lastObservedReverse = ?3
        `); err == nil {
            for _, terminal := range terminals {
                if terminal.Forward {
                    _, err = stmtForward.Exec(
                        terminal.dictionaryId,
                        currentTime,
                        currentTime,
                    )
                }
                if err == nil && terminal.Reverse {
                    _, err = stmtReverse.Exec(
                        terminal.dictionaryId,
                        currentTime,
                        currentTime,
                    )
                }
                if err != nil {
                    if e := stmtForward.Close(); e != nil {
                        logger.Warningf("unable to close forward statement: %s", e)
                    }
                    if e := stmtReverse.Close(); e != nil {
                        logger.Warningf("unable to close reverse statement: %s", e)
                    }
                    if e := tx.Rollback(); e != nil {
                        logger.Warningf("unable to roll-back transaction: %s", e)
                    }
                    return err
                }
            }
            stmtReverse.Close()
        }
        stmtForward.Close()
    } else {
        if e := tx.Rollback(); e != nil {
            logger.Warningf("unable to roll-back transaction: %s", e)
        }
        return err
    }
    return tx.Commit()
}
//this provides starting-point candidates for doing a forward- or reverse-
//random-walk, in the event that a keyword-oriented walk fails.
func (db *Database) terminalsGetStarters(count int, forward bool, oldestAllowedTime int64) ([]int, error) {
    if count <= 0 {
        return make([]int, 0), nil
    }
    
    direction := "Reverse" //a reverse-terminal is position 0 in production
    if !forward {
        direction = "Forward"
    }
    
    query := fmt.Sprintf(`
    SELECT
        dictionaryId
    FROM
        terminals
    WHERE
        lastObserved%s > ?1
    ORDER BY RANDOM()
    LIMIT %d
    `, direction, count)
    
    if rows, err := db.connection.Query(
        query,
        oldestAllowedTime,
    ); err == nil {
        defer rows.Close()
        
        results := make([]int, 0, count)
        for rows.Next() {
            var did int
            if err:= rows.Scan(&did); err == nil {
                results = append(results, did)
            } else {
                return nil, err
            }
        }
        return results, nil
    } else {
        return nil, err
    }
}




func deserialiseTransitionsJSONZLIB(data []byte, oldestAllowedTime int64) (map[int]transitionSpec) {
    reader, err := zlib.NewReader(nil)
    if err != nil {
        panic(err) //inconsistency in zlib library between reader and writer;
                   //make them do the same thing
                   //(if this fails, the environment is unusable)
    }
    defer reader.Close()
    _, decompressed, err := reader.ReadBuffer(data, nil)
    if err != nil {
        logger.Warningf("unable to decompress zlib data; reinitialising state: %s", err)
        return make(map[int]transitionSpec, 1)
    }
    
    var buffer [][3]int = nil
    if err := json.Unmarshal([]byte(decompressed), &buffer); err == nil {
        output := make(map[int]transitionSpec, len(buffer) + 1)
        for _, tspec := range buffer {
            lastObserved := int64(tspec[2])
            if lastObserved > oldestAllowedTime {
                output[tspec[0]] = transitionSpec{
                    occurrences: tspec[1],
                    lastObserved: lastObserved,
                }
            }
        }
        return output
    } else { //some sort of database corruption, almost certainly due to misuse
        logger.Warningf("unable to deserialise transitions; reinitialising state: %s", err)
    }
    return make(map[int]transitionSpec, 1)
}
func serialiseTransitionsJSONZLIB(specs map[int]transitionSpec) ([]byte) {
    destructuredData := make([][3]int, 0, len(specs))
    for did, tspec := range specs {
        destructuredData = append(destructuredData, [3]int{did, tspec.occurrences, int(tspec.lastObserved)})
    }
    
    if buffer, err := json.Marshal(destructuredData); err == nil {
        writer := zlib.NewWriter(nil)
        defer writer.Close()
        if compressed, err := writer.WriteBuffer(buffer, nil); err == nil {
            return compressed
        } else {
            logger.Warningf("unable to compress transitions; reinitialising state: %s", err)
        }
    } else {
        logger.Warningf("unable to serialise transitions; reinitialising state: %s", err)
    }
    return []byte{120, 156, 139, 142, 5, 0, 1, 21, 0, 185} //zlib-compressed empty JSON array
}

func ngramsGetDirectionString(forward bool) (string) {
    if forward {
        return "forward"
    }
    return "reverse"
}

func (db *Database) digramsGet(specs map[DigramSpec]bool, forward bool, oldestAllowedTime int64) (map[DigramSpec]Digram, error) {
    if len(specs) == 0 {
        return make(map[DigramSpec]Digram, 0), nil
    }
    
    if stmt, err := db.connection.Prepare(fmt.Sprintf(`
    SELECT
        transitionsJSONZLIB
    FROM
        digrams_{}
    WHERE
        dictionaryIdFirst = ?1
    LIMIT 1
    `, ngramsGetDirectionString(forward))); err == nil {
        defer stmt.Close()
        
        output := make(map[DigramSpec]Digram, len(specs))
        for spec := range specs {
            var transitionsJSONZLIB []byte
            row := stmt.QueryRow(spec.DictionaryIdFirst)
            if err := row.Scan(&transitionsJSONZLIB); err == nil {
                output[spec] = Digram{
                    Transitions: prepareTransitions(deserialiseTransitionsJSONZLIB(transitionsJSONZLIB, oldestAllowedTime)),
                    
                    dictionaryIdFirst: spec.DictionaryIdFirst,
                }
            } else if err != sql.ErrNoRows {
                return nil, err
            } else {
                output[spec] = Digram{
                    Transitions: prepareTransitionsEmpty(),
                    
                    dictionaryIdFirst: spec.DictionaryIdFirst,
                }
            }
        }
        return output, nil
    } else {
        return nil, err
    }
}
func (db *Database) digramsSet(digrams []Digram, forward bool) (error) {
    if len(digrams) == 0 {
        return nil
    }
    
    tx, err := db.connection.Begin()
    if err != nil {
        return err
    }
    
    if stmt, err := tx.Prepare(fmt.Sprintf(`
    INSERT INTO digrams_%s(
        dictionaryIdFirst,
        transitionsJSONZLIB
    ) VALUES (?1, ?2)
    ON CONFLICT(dictionaryIdFirst) DO UPDATE SET
        transitionsJSONZLIB = ?3
    `, ngramsGetDirectionString(forward))); err == nil {
        for _, digram := range digrams {
            transitionsJSONZLIB := serialiseTransitionsJSONZLIB(digram.Transitions.transitions)
            if _, err = stmt.Exec(
                digram.dictionaryIdFirst,
                transitionsJSONZLIB,
                transitionsJSONZLIB,
            ); err != nil {
                if e := stmt.Close(); e != nil {
                    logger.Warningf("unable to close statement: %s", e)
                }
                if e := tx.Rollback(); e != nil {
                    logger.Warningf("unable to roll-back transaction: %s", e)
                }
                return err
            }
        }
        stmt.Close()
    } else {
        if e := tx.Rollback(); e != nil {
            logger.Warningf("unable to roll-back transaction: %s", e)
        }
        return err
    }
    return tx.Commit()
}


func (db *Database) trigramsGetOnlyFirst(dictionaryIdFirst int, count int, forward bool, oldestAllowedTime int64) ([]Trigram, error) {
    if rows, err := db.connection.Query(fmt.Sprintf(`
    SELECT
        dictionaryIdSecond,
        transitionsJSONZLIB
    FROM
        trigrams_{}
    WHERE
        dictionaryIdFirst = ?1
    ORDER BY RANDOM()
    LIMIT %d
    `, ngramsGetDirectionString(forward), count)); err == nil {
        defer rows.Close()
        
        output := make([]Trigram, 0, count)
        for rows.Next() {
            var dictionaryIdSecond int
            var transitionsJSONZLIB []byte
            if err:= rows.Scan(&dictionaryIdSecond, &transitionsJSONZLIB); err == nil {
                output = append(output, Trigram{
                    Transitions: prepareTransitions(deserialiseTransitionsJSONZLIB(transitionsJSONZLIB, oldestAllowedTime)),
                    
                    dictionaryIdFirst: dictionaryIdFirst,
                    dictionaryIdSecond: dictionaryIdSecond,
                })
            } else {
                return nil, err
            }
        }
        return output, nil
    } else {
        return nil, err
    }
}

/*
type TrigramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
}
type Trigram struct {
    transitions
    
    DictionaryIdFirst int
    DictionaryIdSecond int
}

type QuadgramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
    DictionaryIdThird int
}
type Quadgram struct {
    transitions
    
    DictionaryIdFirst int
    DictionaryIdSecond int
    DictionaryIdThird int
}

type QuintgramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
    DictionaryIdThird int
    DictionaryIdFourth int
}
type Quintgram struct {
    transitions
    
    DictionaryIdFirst int
    DictionaryIdSecond int
    DictionaryIdThird int
    DictionaryIdFourth int
}
*/













type DatabaseManager struct {
    dbDir string
    
    databases map[string]*Database
}
func prepareDatabaseManager(dbDir string) (*DatabaseManager) {
    return &DatabaseManager{
        dbDir: dbDir,
        
        databases: make(map[string]*Database),
    }
}
func (dbm *DatabaseManager) Close() {
    logger.Debugf("closing databases...")
    for _, database := range dbm.databases {
        database.Close();
    }
    dbm.databases = make(map[string]*Database)
}
func (dbm *DatabaseManager) idToPath(contextId string) (string) {
    return filepath.Join(dbm.dbDir, contextId + ".sqlite3")
}
func (dbm *DatabaseManager) Create(contextId string) (error) {
    logger.Infof("creating database {}...", contextId)
    if db, err := prepareDatabase(dbm.idToPath(contextId)); err == nil {
        return db.Close()
    } else {
        return err
    }
}
func (dbm *DatabaseManager) Drop(contextId string) (error) {
    logger.Infof("dropping database {}...", contextId)
    if database, defined := dbm.databases[contextId]; defined {
        if err := database.Close(); err != nil {
            //it'll be referenced anyway, so this isn't critical
            logger.Warningf("unable to close database %s: %s", contextId, err)
        }
        delete(dbm.databases, contextId)
    }
    return os.Remove(dbm.idToPath(contextId))
}
func (dbm *DatabaseManager) Load(contextId string) (*Database, error) {
    logger.Infof("loading database {}...", contextId)
    dbPath := dbm.idToPath(contextId)
    if _, err := os.Stat(dbPath); err != nil {
        return nil, err
    }
    return prepareDatabase(dbPath)
}

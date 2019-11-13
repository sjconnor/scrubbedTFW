package spell

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/svipul/tinfoil-wizard/errorutil"
	"github.com/svipul/tinfoil-wizard/httputil"
)

// Spell represents a single spell entity's complete properties
// corresponds to the Spells table properties
type Spell struct {
	Key      int64  `json:"SpellKey"`
	Name     string `json:"SpellName"`
	MSchool  string `json:"MagicSchool"`
	Level    int64  `json:"Level"`
	Range    string `json:"SpellRange"`
	Duration string `json:"Duration"`
	Concen   int64  `json:"RequiresConcentration"`
	CTime    string `json:"CastingTime"`
	VComp    int64  `json:"VerbalComponent"`
	SComp    int64  `json:"SomaticComponent"`
	MComp    string `json:"MaterialComponent"`
	Desc     string `json:"Description"`
}

// CharacterSpell represents a specific character's known spells
// corresponds to the CharacterKnownSpells table properties
type CharacterSpell struct {
	Key      int64  `json:"SpellKey"`
	Name     string `json:"SpellName"`
	Level    int64  `json:"Level"`
	Prep     int64  `json:"Prepared"`
	Conc     int64  `json:"HoldingConcentration"`
	Range    string `json:"SpellRange"`
	CTime    string `json:"CastingTime"`
	Concen   int64  `json:"RequiresConcentration"`
	Desc     string `json:"Description"`
	MSchool  string `json:"MagicSchool"`
	Duration string `json:"Duration"`
	MComp    string `json:"MaterialComponent"`
}

type CharacterSpellPatch struct {
	Prep int64 `json:"Prepared"`
}

// TODO description
type SpellClient interface {
	AllSpells(w http.ResponseWriter, r *http.Request)
	HandleSpell(w http.ResponseWriter, r *http.Request)
	CharacterSpells(w http.ResponseWriter, r *http.Request)
}

type client struct {
	db *sql.DB
}

// NewSpellClient TODO description
func NewSpellClient(db *sql.DB) SpellClient {
	return &client{db: db}
}

/******************************************************************************
** HANDLERS:
	** AllSpells: ENDPOINT "/api/spells"
	** HandleSpell: ENDPOINT "/api/spells/{sID}"
	** CharacterSpells: ENDPOINT "/api/characters/{cID}/spells"
	** HandleCharacterSpell: ENDPOINT: "/api/characters/{cID}/spells/{sID}"
******************************************************************************/

// AllSpells: GET spells from Spells table
// ENDPOINT "/api/spells"
//
func (c *client) AllSpells(w http.ResponseWriter, r *http.Request) {

	log.Print("AllSpells endpoint")

	switch r.Method {

	case http.MethodGet: // get all spells from Spells table

		// If the Character ID is included as a parameter, only get the list of
		// all spells that the character can learn at the current level that
		// they don't already know
		if cID, ok := r.URL.Query()["characterid"]; ok {
			log.Print("GET on AllSpells for specific character: " + cID[0])

			// The query to return all spells that a chacaracter can learn for their class at
			// their specific level that they don't already know
			var query = `SELECT Spells.SpellKey, Spells.SpellName, Spells.Description, Spells.Duration, Spells.MagicSchool, Spells.Level, Spells.CastingTime, Spells.SpellRange FROM Spells JOIN ClassSpells ON Spells.SpellKey = ClassSpells.SpellKey
			JOIN Characters ON ClassSpells.ClassKey = Characters.ClassKey WHERE Characters.CharacterKey = ` + cID[0] + ` AND
			Spells.Level <= (SELECT cs.MaxSpellLevel FROM ClassSlotsPerLevel cs 
				JOIN Characters c ON cs.ClassKey = c.ClassKey AND cs.ClassLevel = c.ClassLevel WHERE c.CharacterKey = ` + cID[0] + `) AND
			Spells.SpellKey NOT IN (SELECT SpellKey FROM CharacterKnownSpells WHERE UserCharacterKey = ` + cID[0] + `) ORDER BY Spells.SpellName`

			//query for all spells available for a given spell's spell class
			ret, err := c.db.Query(query)
			defer ret.Close()
			if err != nil { // new error logging method suggested by Shannon
				log.Printf("mysql: could not access Spells table: %v", err)
				httputil.ErrorResponse(w, errorutil.New(http.StatusInternalServerError, "internal error"))
				return
			}

			allSpells := []*Spell{} // fill this with data returned by query

			// scan each query row and append to allSpells
			for ret.Next() {

				var spell Spell

				err = ret.Scan(&spell.Key, &spell.Name, &spell.Desc, &spell.Duration, &spell.MSchool, &spell.Level, &spell.CTime, &spell.Range)
				if err != nil {
					log.Printf("\nCould not parse spell: %v", err)
					httputil.ErrorResponse(w, errorutil.New(http.StatusInternalServerError, "internal error"))
					return
				}

				allSpells = append(allSpells, &spell)

			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(allSpells)
		} else { // Just return a list of all spells in the database
			log.Print("GET on AllSpells")

			// query for all entities in Spells joined ClassSpells
			ret, err := c.db.Query(`SELECT Spells.SpellKey, Spells.SpellName, Spells.Description, Spells.Level
			FROM Spells JOIN ClassSpells
			ON Spells.SpellKey = ClassSpells.SpellKey ORDER BY SpellName;`)
			defer ret.Close() // close the connection when done!
			if err != nil {   // new error logging method suggested by Shannon
				log.Printf("mysql: could not access Spells table: %v", err)
				httputil.ErrorResponse(w, errorutil.New(http.StatusInternalServerError, "internal error"))
				return
			}

			allSpells := []*Spell{} // fill this with data returned by query

			// scan each query row and append to allSpells
			for ret.Next() {

				var spell Spell

				err = ret.Scan(&spell.Key, &spell.Name, &spell.Desc, &spell.Level)
				if err != nil {
					log.Printf("\nCould not parse spell: %v", err)
					httputil.ErrorResponse(w, errorutil.New(http.StatusInternalServerError, "internal error"))
					return
				}

				allSpells = append(allSpells, &spell)

			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(allSpells)
		}

	default: // TODO handle this thoroughly

		w.WriteHeader(http.StatusForbidden) // return 403

	}

}

// HandleSpell: GET, PATCH, DELETE on single spell in Spells table
// ENDPOINT "/api/spells/{sID}"
func (c *client) HandleSpell(w http.ResponseWriter, r *http.Request) {

	log.Printf("Spell endpoint")

	thisSpell := new(Spell)

	// get the ID out of the request variables
	vars := mux.Vars(r)
	sID := vars["sID"]
	var err error
	if thisSpell.Key, err = strconv.ParseInt(sID, 10, 64); err != nil {
		log.Printf("\nSpell id (handlespell) ParseInt fail %v", err)
		httputil.ErrorResponse(w, errorutil.Newf(http.StatusBadRequest, "%q is not a valid spell id", sID))
		return
	}

	switch r.Method {

	case http.MethodGet: // get spell by ID

		if err := c.getSpell(thisSpell); err != nil {
			httputil.ErrorResponse(w, err)
			return
		}

	case http.MethodPatch: // edit spell by ID

		// fill struct with properties from the request body
		c.patchSpell(thisSpell)           // patch the spell in the db
		w.WriteHeader(http.StatusCreated) // return 201 Created?

		return

	case http.MethodDelete: // delete spell by ID

		c.deleteSpell(thisSpell.Key)
		w.WriteHeader(http.StatusNoContent) // return 204 No Content?

		return // TODO handle end of delete differently?

	default: // TODO handle this thoroughly

		w.WriteHeader(http.StatusForbidden) // return 403
		return
	}

	// make it json before writing back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(thisSpell)

}

// CharacterSpells: GET, POST spells from/to CharacterKnownSpells table
// ENDPOINT "/api/characters/{cID}/spells"
// ENDPOINT "/api/characters/{cID}/spells/{sID}"
func (c *client) CharacterSpells(w http.ResponseWriter, r *http.Request) {
	log.Printf("Character Spells Endpoint")

	vars := mux.Vars(r)
	var cID = int64(-1)
	var sID = int64(-1)
	var err error
	varsCID := vars["cID"]
	// parse the cID out of url params
	if cID, err = strconv.ParseInt(varsCID, 10, 64); err != nil {
		log.Printf("Character id ParseInt fail %v", err)
		httputil.ErrorResponse(w, errorutil.Newf(http.StatusBadRequest, "%q is not a valid character id", varsCID))
		return
	}

	switch r.Method {
	case http.MethodGet: // get all spells for the character
		// Added order statement to get spells in order by name
		var query = "SELECT s.SpellKey, s.SpellName, s.RequiresConcentration, s.Description, s.Duration, s.MagicSchool, s.Level, s.CastingTime, s.SpellRange, ks.Prepared, ks.HoldingConcentration" +
			" FROM CharacterKnownSpells ks JOIN Spells s ON ks.SpellKey = s.SpellKey" +
			" WHERE ks.UserCharacterKey = " + strconv.FormatInt(cID, 10) + " ORDER BY s.SpellName"

		ret, err := c.db.Query(query)
		defer ret.Close() // close the connection when done!
		if err != nil {   // new error logging method suggested by Shannon
			log.Printf("mysql: could not access Spells table: %v", err)
			httputil.ErrorResponse(w, errorutil.New(http.StatusInternalServerError, "internal error"))
			return
		}

		// Object for storing scanned spells
		characterSpells := []*CharacterSpell{}

		// Loop through the return SQL rows
		for ret.Next() {
			var spell CharacterSpell
			err = ret.Scan(&spell.Key, &spell.Name, &spell.Concen, &spell.Desc, &spell.Duration, &spell.MSchool, &spell.Level, &spell.CTime, &spell.Range, &spell.Prep, &spell.Conc)
			if err != nil { // new error logging method suggested by Shannon
				log.Printf("\nmysql: could not access Spells table: %v", err)
				httputil.ErrorResponse(w, errorutil.New(http.StatusInternalServerError, "internal error"))
				return
			}

			characterSpells = append(characterSpells, &spell)

		}

		//Return the Spells as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(characterSpells)

	case http.MethodPut:
		log.Printf("Adding spell")

		// We only want to parse this here because the spell ID can be null
		// when grabbing all character spells
		varsSID := vars["sID"]
		if sID, err = strconv.ParseInt(varsSID, 10, 64); err != nil {
			log.Printf("\nSpell id characterspells(get) ParseInt fail %v", err)
			httputil.ErrorResponse(w, errorutil.Newf(http.StatusBadRequest, "%q is not a valid spell id", varsSID))
			return
		}

		var spell Spell
		spell.Key = sID
		c.getSpell(&spell)

		prepped := int64(0)
		if spell.Level == 0 {
			prepped = int64(1)
		}

		// Query for inserting the spell
		var query = "INSERT INTO CharacterKnownSpells (UserCharacterKey, SpellKey, Prepared) VALUES (" +
			strconv.FormatInt(cID, 10) + ", " + strconv.FormatInt(sID, 10) + ", " + strconv.FormatInt(prepped, 10) + ");"

		// Run the query against the DB
		if _, err := c.db.Exec(query); err != nil { // new error logging method suggested by Shannon
			log.Printf("mysql: could not access Spells table: %v", err)
			httputil.ErrorResponse(w, errorutil.New(http.StatusInternalServerError, "internal error"))
			return
		}

	case http.MethodPatch:
		// We only want to parse this here because the spell ID can be null
		// when grabbing all character spells
		varsSID := vars["sID"]
		if sID, err = strconv.ParseInt(varsSID, 10, 64); err != nil {
			log.Printf("Spell id characterspells(patch) ParseInt fail %v", err)
			httputil.ErrorResponse(w, errorutil.Newf(http.StatusBadRequest, "%q is not a valid spell id", varsSID))
			return
		}
		csp := &CharacterSpellPatch{}

		b, err := ioutil.ReadAll(r.Body) // byte-ify request body
		if err != nil {
			log.Printf("couldn't read patchSpell request body: %v", err)
			httputil.ErrorResponse(w, errorutil.New(http.StatusBadRequest, "invalid patch body"))
			return
		}

		if err := json.Unmarshal(b, csp); err != nil {
			log.Printf("couldn't unmarshall patch spell body: %v", err)
			httputil.ErrorResponse(w, errorutil.New(http.StatusBadRequest, "unparsable patch body"))
			return
		}

		var query = `UPDATE CharacterKnownSpells SET Prepared = ? WHERE UserCharacterKey = ? AND SpellKey = ?`
		if _, err := c.db.Exec(query, &csp.Prep, cID, sID); err != nil {
			log.Printf("couldn't update CharacterKnownSpells: %v", err)
			httputil.ErrorResponse(w, errorutil.New(http.StatusInternalServerError, "internal error"))
			return
		}
	}

}

/******************************************************************************
** HELPERS:
	** getSpell		WORKS
	** postSpell	STUB ONLY
	** patchSpell	STUB ONLY
	** deleteSpell	STUB ONLY
******************************************************************************/

// getSpell queries the database by spell ID and fills the struct
func (c *client) getSpell(thisSpell *Spell) error {

	// query for spell by ID in thisSpell.Key
	query := "SELECT * FROM Spells WHERE SpellKey = " + strconv.FormatInt(thisSpell.Key, 10) + ";"
	ret, err := c.db.Query(query)
	defer ret.Close() // close the connection when done!
	if err != nil {
		log.Printf("\nmysql: could not access Spells table: %v", err)
		return errorutil.New(http.StatusInternalServerError, "internal error")
	}

	// scan query results into struct
	for ret.Next() {
		err = ret.Scan(&thisSpell.Key, &thisSpell.Name,
			&thisSpell.MSchool, &thisSpell.Level,
			&thisSpell.Range, &thisSpell.Duration,
			&thisSpell.Concen, &thisSpell.CTime,
			&thisSpell.VComp, &thisSpell.SComp,
			&thisSpell.MComp, &thisSpell.Desc)
		if err != nil {
			log.Printf("\nCould not parse spell: %v", err)
			return errorutil.New(http.StatusInternalServerError, "internal error")
		}
	}
	return nil
}

// postSpell inserts the spell into the database
func (*client) postSpell(thisSpell *Spell) {

	// TODO

	return
}

// patchSpell updates the spell with this ID
func (*client) patchSpell(thisSpell *Spell) {

	// TODO

	return
}

// deleteSpell deletes the spell in the DB with the received ID
func (*client) deleteSpell(ID int64) {

	// TODO

	return

}

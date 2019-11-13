package character

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"firebase.google.com/go/auth"

	firebase "firebase.google.com/go"

	"github.com/gorilla/mux"
	"github.com/svipul/tinfoil-wizard/class"
	"github.com/svipul/tinfoil-wizard/errorutil"
	"github.com/svipul/tinfoil-wizard/httputil"
)

// Character represents a single Character entity's properties
type Character struct {
	ID                   int64           `json:"id"`
	Name                 string          `json:"name"`
	Class                *class.Class    `json:"class,omitempty"`
	Subclass             *class.Subclass `json:"subclass,omitempty"`
	Level                int             `json:"level,omitempty"`
	AbilityScore         int             `json:"abilityScore,omitempty"`
	PortraitPath         string          `json:"portraitPath,omitempty"`
	Concentrating        int64           `json:"concentrating"`
	Level1SlotsRemaining int64           `json:"Level1SlotsRemaining"`
	Level2SlotsRemaining int64           `json:"Level2SlotsRemaining"`
	Level3SlotsRemaining int64           `json:"Level3SlotsRemaining"`
	Level4SlotsRemaining int64           `json:"Level4SlotsRemaining"`
	Level5SlotsRemaining int64           `json:"Level5SlotsRemaining"`
	Level6SlotsRemaining int64           `json:"Level6SlotsRemaining"`
	Level7SlotsRemaining int64           `json:"Level7SlotsRemaining"`
	Level8SlotsRemaining int64           `json:"Level8SlotsRemaining"`
	Level9SlotsRemaining int64           `json:"Level9SlotsRemaining"`
	CastSuccess          bool            `json:"CastSuccess"`
}

// SetLevel updates the total remaining spell slots by level
func (c *Character) SetLevel(level int) {
	if level <= 0 || level > 20 {
		return
	}
	delta := []int64{}
	for i := 0; i < 10; i++ {
		delta = append(delta, int64(class.GetSpellSlots(level, i)-class.GetSpellSlots(c.Level, i)))
	}

	c.Level1SlotsRemaining += delta[1]
	if c.Level1SlotsRemaining < 0 {
		c.Level1SlotsRemaining = 0
	}
	c.Level2SlotsRemaining += delta[2]
	if c.Level2SlotsRemaining < 0 {
		c.Level2SlotsRemaining = 0
	}
	c.Level3SlotsRemaining += delta[3]
	if c.Level3SlotsRemaining < 0 {
		c.Level3SlotsRemaining = 0
	}
	c.Level4SlotsRemaining += delta[4]
	if c.Level4SlotsRemaining < 0 {
		c.Level4SlotsRemaining = 0
	}
	c.Level5SlotsRemaining += delta[5]
	if c.Level5SlotsRemaining < 0 {
		c.Level5SlotsRemaining = 0
	}
	c.Level6SlotsRemaining += delta[6]
	if c.Level6SlotsRemaining < 0 {
		c.Level6SlotsRemaining = 0
	}
	c.Level7SlotsRemaining += delta[7]
	if c.Level7SlotsRemaining < 0 {
		c.Level7SlotsRemaining = 0
	}
	c.Level8SlotsRemaining += delta[8]
	if c.Level8SlotsRemaining < 0 {
		c.Level8SlotsRemaining = 0
	}
	c.Level9SlotsRemaining += delta[9]
	if c.Level9SlotsRemaining < 0 {
		c.Level9SlotsRemaining = 0
	}

	c.Level = level
}

// PostCharBody TODO description
type PostCharBody struct {
	Char     *Character `json:"character"`
	Portrait string     `json:"portrait"`
}

// CharacterClient TODO description
type CharacterClient interface {
	AllCharacters(w http.ResponseWriter, r *http.Request)
	HandleCharacter(w http.ResponseWriter, r *http.Request)
	//HandleCharacterSlots(w https.ResponseWriter, r *http.Request)
}

// TODO description
type client struct {
	db  *sql.DB
	ac  *auth.Client
	clc class.Client
}

// TODO description
var (
	imgExt = map[string]string{
		"jpeg": "jpg",
	}
)

// TODO description
func NewClient(ctx context.Context, db *sql.DB, clc class.Client) (CharacterClient, error) {
	c := &client{db: db, clc: clc}
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

// TODO describe/comment
func (c *client) init(ctx context.Context) error {
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: "tinfoil-wizard"})
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	ac, err := app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error initializing auth: %v", err)
	}
	c.ac = ac
	return nil
}

/******************************************************************************
** HANDLERS:
** AllCharacters: ENDPOINT: "/api/characters"
** HandleCharacter: ENDPOINT: "/api/characters/{cID}"
******************************************************************************/

// AllCharacters: GET, POST characters by user (see uID functions)
// ENDPOINT: "/api/characters"
func (c *client) AllCharacters(w http.ResponseWriter, r *http.Request) {

	log.Print("AllCharacters endpoint")

	// get user's UID
	uID, err := c.requestUID(r)
	if err != nil {
		httputil.ErrorResponse(w, err)
		return
	}

	switch r.Method {

	case http.MethodGet: // get the characters for this user

		chars, err := c.getCharacters(uID)

		if err != nil {
			httputil.ErrorResponse(w, err)
			return
		}

		httputil.JSONResponse(w, http.StatusOK, chars)

	case http.MethodPost: // add a character for this user

		char, err := c.postCharacter(uID, r)

		if err != nil {
			httputil.ErrorResponse(w, err)
			return
		}

		httputil.JSONResponse(w, http.StatusCreated, char)

	default: // TODO handle this better?

		httputil.ErrorResponse(w, errorutil.New(http.StatusForbidden, "invalid request"))

	}
}

// HandleCharacter: GET, PATCH, DELETE character from Characters table
// ENDPOINT: "/api/characters/{cID}"
func (c *client) HandleCharacter(w http.ResponseWriter, r *http.Request) {

	log.Print("Characters endpoint")

	// NOTE: must set the content type BEFORE setting return code
	w.Header().Set("Content-Type", "application/json")

	// struct to fill - used by all methods
	thisCharacter := new(Character)

	// get the ID out of the request variables
	vars := mux.Vars(r)

	// set the character ID in the struct so it is easily accessible
	varsCID := vars["cID"]
	cID, err := strconv.ParseInt(varsCID, 10, 64)
	if err != nil {
		log.Printf("\ncharacter id ParseInt fail %v", err)
		httputil.ErrorResponse(w, errorutil.Newf(http.StatusBadRequest, "%s isn't a valid character id", varsCID))
		return
	}
	thisCharacter.ID = cID

	// get user's UID
	uID, err := c.requestUID(r)
	if err != nil {
		httputil.ErrorResponse(w, err)
		return
	}

	log.Printf("\nServing request: method %s, path %s", r.Method, r.URL.Path)
	switch r.Method {

	case http.MethodGet:
		// fill struct with character info from database
		if err := c.getCharacter(thisCharacter, uID); err != nil {
			httputil.ErrorResponse(w, err)
			return
		}
		// set status code and write result
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(thisCharacter)

	case http.MethodPatch:
		if err := c.patchCharacter(cID, uID, r); err != nil {
			httputil.ErrorResponse(w, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)

	case http.MethodDelete:
		if err := c.deleteCharacter(cID, uID); err != nil {
			log.Printf("\ndeleteCharacter failed: %v", err)
			httputil.ErrorResponse(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default: // TODO handle this thoroughly

		w.WriteHeader(http.StatusForbidden) // return 403

	}
}

/******************************************************************************
** HELPERS:
** requestUID		WORKS
** getCharacters	WORKS
** postCharacter	WORKS
** saveCharacter	WORKS
** savePortrait		WORKS
** patchCharacter	STUB ONLY
** deleteCharacter	WORKS
******************************************************************************/

// requestUID gets the user's Firebase UID from the supplied authentication JWT.
func (c *client) requestUID(r *http.Request) (string, error) {
	jwt := r.Header.Get("Authorization")
	tok, err := c.ac.VerifyIDToken(r.Context(), jwt)
	if err != nil {
		fmt.Printf("couldn't verify token: %v", err)
		return "", errorutil.New(http.StatusUnauthorized, "request must include valid JWT")
	}
	return tok.UID, nil
}

// getCharacters retrieves all of a user's characters from the DB, based on their UID.
func (c *client) getCharacters(uid string) ([]*Character, error) {
	rows, err := c.db.Query(
		`SELECT ch.CharacterKey, ch.Name, ch.ClassLevel,  ch.PortraitPath, ch.ClassKey, IFNULL(ch.SubclassKey, 0)
		FROM Characters ch
		WHERE ch.OwnerID = ?`,
		uid)
	defer rows.Close() // close the connection when done!
	if err != nil {
		log.Printf("couldn't get characters for uid %s: %v", uid, err)
		return nil, errorutil.New(http.StatusInternalServerError, "internal error")
	}

	var chars []*Character
	for rows.Next() {
		ch := &Character{}
		var clID, scID int
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.Level, &ch.PortraitPath, &clID, &scID); err != nil {
			log.Printf("couldn't read row: %v", err)
			continue
		}
		cl, err := c.clc.ClassByID(clID, false)
		if err != nil {
			log.Printf("couldn't get class for character (ID: %d): %v", ch.ID, err)
			return nil, err
		}
		ch.Class = cl
		if scID != 0 {
			sc, err := c.clc.SubclassByID(scID)
			if err != nil {
				log.Printf("couldn't get subclass for character (ID: %d): %v", ch.ID, err)
				return nil, err
			}
			ch.Subclass = sc
		}
		chars = append(chars, ch)
	}
	return chars, nil
}

// postCharacter creates a new character in the app, including a record of the character in
// the app's database and a portrait file to be displayed in the app's UI.
func (c *client) postCharacter(uid string, r *http.Request) (*Character, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("couldn't read postCharacter request body: %v", err)
		return nil, errorutil.New(http.StatusBadRequest, "invalid post body")
	}
	pcb := &PostCharBody{}
	if err := json.Unmarshal(b, pcb); err != nil {
		log.Printf("couldn't unmarshall postCharacter body: %v", err)
		return nil, errorutil.New(http.StatusBadRequest, "unparsable request body")
	}
	tx, err := c.db.Begin()
	if err != nil {
		log.Printf("couldn't create transaction %v", err)
		return nil, errorutil.New(http.StatusInternalServerError, "internal error")
	}
	if err := c.saveCharacter(tx, uid, pcb.Char); err != nil {
		return nil, err
	}
	if err := c.savePortrait(tx, pcb.Portrait, pcb.Char); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		log.Printf("error committing postCharacter transaction: %v", err)
		return nil, errorutil.New(http.StatusInternalServerError, "internal error")
	}
	return pcb.Char, nil
}

// saveCharacter saves the supplied character to the Characters table.
func (c *client) saveCharacter(tx *sql.Tx, ownerID string, ch *Character) error {

	var scID *int

	if ch.Subclass != nil {
		scID = &ch.Subclass.ID
	}
	rows, err := tx.Query(
		`SELECT Level1Slots
		 FROM ClassSlotsPerLevel
		 WHERE ClassLevel = 1 AND ClassKey = ?`,
		ch.Class.ID)
	if err != nil {
		log.Printf("\ncouldn't get ClassSlotsPerLevel data from DB: %v", err)
		return errorutil.New(http.StatusInternalServerError, "couldn't create character")
	}
	var startingSlots int
	if ok := rows.Next(); !ok {
		log.Printf("\nClassSlotsPerLevel returned empty row")
		return errorutil.New(http.StatusInternalServerError, "couldn't create character")
	}
	if err := rows.Scan(&startingSlots); err != nil {
		log.Printf("\nClassSlotsPerLevel rows scan error: %v", err)
		return errorutil.New(http.StatusInternalServerError, "couldn't create character")
	}
	rows.Close()

	r, err := tx.Exec(`INSERT INTO Characters (OwnerID, Name, ClassKey, ClassLevel, SCAbilityScore, SubclassKey, Level1SlotsRemaining)
		VALUES (?, ?, ?, ?, ?, ?, ?);`,
		ownerID, ch.Name, ch.Class.ID, ch.Level, ch.AbilityScore, scID, startingSlots)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("couldn't rollback saveCharacter changes. uh oh: %v", err)
		}
		log.Printf("couldn't save character to db: %v", err)
		return errorutil.New(http.StatusInternalServerError, "couldn't create character")
	}
	id, err := r.LastInsertId()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("couldn't rollback saveCharacter changes: %v", err)
		}
		log.Printf("couldn't retrieve character ID: %v", err)
		return errorutil.New(http.StatusInternalServerError, "couldn't create character")
	}
	ch.ID = id
	return nil
}

// savePortrait creates an image file for a character portrait, based on a supplied int64 encoded string,
// then updates the Character table with the portrait's path.
func (c *client) savePortrait(tx *sql.Tx, fullData string, ch *Character) error {
	if !strings.HasPrefix(fullData, "data:image/jpeg;base64,") {
		return errorutil.New(http.StatusBadRequest, "invalid portrait")
	}

	pData := fullData[len("data:image/jpeg;base64,"):]

	// Calculate a unique number based on the image data
	h := sha1.New()
	h.Write([]byte(pData))
	csum := h.Sum(nil)
	extension := "jpg"

	ch.PortraitPath = fmt.Sprintf("/portraits/%d.%x.%s", ch.ID, csum, extension)
	p, err := base64.StdEncoding.DecodeString(pData)
	if err != nil {
		log.Printf("couldn't decode portrait: %v", err)
		return errorutil.New(http.StatusBadRequest, "invalid portrait")
	}
	if err := ioutil.WriteFile(ch.PortraitPath[1:], p, os.ModePerm); err != nil {
		log.Printf("couldn't write portrait file: %v", err)
		return errorutil.New(http.StatusInternalServerError, "internal error")
	}
	if _, err := tx.Exec(
		`UPDATE Characters
		 SET PortraitPath = ?
		 WHERE CharacterKey = ?`,
		ch.PortraitPath, ch.ID); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("couldn't rollback character add: %v", err)
		}
		fmt.Printf("couldn't save character portrait path: %v", err)
		return errorutil.New(http.StatusInternalServerError, "internal error")
	}

	if matches, err := filepath.Glob(fmt.Sprintf("portraits/%d.*.*", ch.ID)); err == nil {
		for _, m := range matches {
			if m != ch.PortraitPath[1:] {
				os.Remove(m)
			}
		}
	}
	return nil
}

// getCharacter queries the database by character ID and fills the struct
func (c *client) getCharacter(thisCharacter *Character, uID string) error {
	authorized, err := c.isOwner(thisCharacter.ID, uID)
	if err != nil {
		return err
	}
	if !authorized {
		return errorutil.New(http.StatusForbidden, "can't delete someone else's character...")
	}

	// get the character by ID
	var clID, scID int
	// TODO if the default for Concentrated is 0, the api call to spells/0 will return Acid Splash (key 1) - why?
	if err := c.db.QueryRow(`SELECT ch.Name, ch.ClassLevel, ch.SCAbilityScore, ch.PortraitPath, ch.ClassKey, IFNULL(ch.SubclassKey, 0), IFNULL(ch.ConcentratedSpell, 0),
		ch.Level1SlotsRemaining, ch.Level2SlotsRemaining, ch.Level3SlotsRemaining, ch.Level4SlotsRemaining, ch.Level5SlotsRemaining, ch.Level6SlotsRemaining,
		ch.Level7SlotsRemaining, ch.Level8SlotsRemaining, ch.Level9SlotsRemaining, ch.CastSuccess
		FROM Characters ch
		WHERE ch.CharacterKey = ?`, thisCharacter.ID).Scan(&thisCharacter.Name, &thisCharacter.Level,
		&thisCharacter.AbilityScore, &thisCharacter.PortraitPath,
		&clID, &scID, &thisCharacter.Concentrating, &thisCharacter.Level1SlotsRemaining, &thisCharacter.Level2SlotsRemaining, &thisCharacter.Level3SlotsRemaining,
		&thisCharacter.Level4SlotsRemaining, &thisCharacter.Level5SlotsRemaining, &thisCharacter.Level6SlotsRemaining, &thisCharacter.Level7SlotsRemaining,
		&thisCharacter.Level8SlotsRemaining, &thisCharacter.Level9SlotsRemaining, &thisCharacter.CastSuccess); err != nil {
		log.Printf("mysql: could not access Character: %v", err)
		return errorutil.New(500, "internal error")
	}
	cl, err := c.clc.ClassByID(clID, false)
	if err != nil {
		log.Printf("couldn't get class for character (ID: %d): %v", thisCharacter.ID, err)
		return err
	}
	thisCharacter.Class = cl
	if scID != 0 {
		sc, err := c.clc.SubclassByID(scID)
		if err != nil {
			log.Printf("couldn't get subclass for character (ID: %d): %v", thisCharacter.ID, err)
			return err
		}
		thisCharacter.Subclass = sc
	}
	log.Print("\n\tfilled struct: ", thisCharacter)

	return nil
}

func (c *Character) updateCharacter(tx *sql.Tx) error {
	var nullable sql.NullInt64
	nullable.Valid = false // assume INVALID integer

	if c.Concentrating != 0 {
		nullable.Valid = true
		nullable.Int64 = c.Concentrating
	}
	query := `UPDATE Characters SET
		Name = ?,
		ClassLevel = ?,
		SCAbilityScore = ?,
		ConcentratedSpell = ?,
		Level1SlotsRemaining = ?,
		Level2SlotsRemaining = ?,
		Level3SlotsRemaining = ?,
		Level4SlotsRemaining = ?,
		Level5SlotsRemaining = ?,
		Level6SlotsRemaining = ?,
		Level7SlotsRemaining = ?,
		Level8SlotsRemaining = ?,
		Level9SlotsRemaining = ?,
        CastSuccess = ?
		WHERE CharacterKey = ?
		`
	if _, err := tx.Exec(query,
		c.Name,
		c.Level,
		c.AbilityScore,
		nullable,
		c.Level1SlotsRemaining,
		c.Level2SlotsRemaining,
		c.Level3SlotsRemaining,
		c.Level4SlotsRemaining,
		c.Level5SlotsRemaining,
		c.Level6SlotsRemaining,
		c.Level7SlotsRemaining,
		c.Level8SlotsRemaining,
		c.Level9SlotsRemaining,
		c.CastSuccess,
		c.ID,
	); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("couldn't rollback patchCharacter changes: %v", err)
		}
		log.Printf("couldn't modify character info: %v", err)
		return errorutil.New(http.StatusInternalServerError, "couldn't modify character")
	}
	return nil
}

// patch the supplied character to the Characters table.
func (c *client) patchCharacter(cID int64, uID string, r *http.Request) error {
	authorized, err := c.isOwner(cID, uID)
	if err != nil {
		return err
	}
	if !authorized {
		return errorutil.New(http.StatusForbidden, "can't delete someone else's character...")
	}

	pcb := &PostCharBody{}

	// unmarshal request body into a post char body type
	b, err := ioutil.ReadAll(r.Body) // byte-ify request body
	if err != nil {
		log.Printf("couldn't read patchCharacter request body: %v", err)
		return errorutil.New(http.StatusBadRequest, "invalid patch body")
	}

	if err := json.Unmarshal(b, pcb); err != nil {
		log.Printf("couldn't unmarshall patchCharacter body: %v", err)
		return errorutil.New(http.StatusBadRequest, "unparsable patch body")
	}

	tx, err := c.db.Begin()
	if err != nil {
		log.Printf("couldn't create transaction %v", err)
		return errorutil.New(http.StatusInternalServerError, "internal error")
	}

	// if character info included in patch, save it
	if pcb.Char != nil {
		log.Printf("patch includes character stats...")

		curChar := Character{ID: pcb.Char.ID}
		if err := c.getCharacter(&curChar, uID); err != nil {
			return err
		}

		// Secret handshake to do a long rest
		if pcb.Char.Level1SlotsRemaining == -1 {
			curChar.Level1SlotsRemaining = class.GetSpellSlots(curChar.Level, 1)
			curChar.Level2SlotsRemaining = class.GetSpellSlots(curChar.Level, 2)
			curChar.Level3SlotsRemaining = class.GetSpellSlots(curChar.Level, 3)
			curChar.Level4SlotsRemaining = class.GetSpellSlots(curChar.Level, 4)
			curChar.Level5SlotsRemaining = class.GetSpellSlots(curChar.Level, 5)
			curChar.Level6SlotsRemaining = class.GetSpellSlots(curChar.Level, 6)
			curChar.Level7SlotsRemaining = class.GetSpellSlots(curChar.Level, 7)
			curChar.Level8SlotsRemaining = class.GetSpellSlots(curChar.Level, 8)
			curChar.Level9SlotsRemaining = class.GetSpellSlots(curChar.Level, 9)
			curChar.CastSuccess = true
		}

		// Secret handshake to do an arcane recovery
		if pcb.Char.Level1SlotsRemaining == -2 {
			if class.GetSpellSlots(curChar.Level, 9) > 0 && curChar.Level9SlotsRemaining != class.GetSpellSlots(curChar.Level, 9) {
				curChar.Level9SlotsRemaining++
			} else if class.GetSpellSlots(curChar.Level, 8) > 0 && curChar.Level8SlotsRemaining != class.GetSpellSlots(curChar.Level, 8) {
				curChar.Level8SlotsRemaining++
			} else if class.GetSpellSlots(curChar.Level, 7) > 0 && curChar.Level7SlotsRemaining != class.GetSpellSlots(curChar.Level, 7) {
				curChar.Level7SlotsRemaining++
			} else if class.GetSpellSlots(curChar.Level, 6) > 0 && curChar.Level6SlotsRemaining != class.GetSpellSlots(curChar.Level, 6) {
				curChar.Level6SlotsRemaining++
			} else if class.GetSpellSlots(curChar.Level, 5) > 0 && curChar.Level5SlotsRemaining != class.GetSpellSlots(curChar.Level, 5) {
				curChar.Level5SlotsRemaining++
			} else if class.GetSpellSlots(curChar.Level, 4) > 0 && curChar.Level4SlotsRemaining != class.GetSpellSlots(curChar.Level, 4) {
				curChar.Level4SlotsRemaining++
			} else if class.GetSpellSlots(curChar.Level, 3) > 0 && curChar.Level3SlotsRemaining != class.GetSpellSlots(curChar.Level, 3) {
				curChar.Level3SlotsRemaining++
			} else if class.GetSpellSlots(curChar.Level, 2) > 0 && curChar.Level2SlotsRemaining != class.GetSpellSlots(curChar.Level, 2) {
				curChar.Level2SlotsRemaining++
			} else if class.GetSpellSlots(curChar.Level, 1) > 0 && curChar.Level1SlotsRemaining != class.GetSpellSlots(curChar.Level, 1) {
				curChar.Level1SlotsRemaining++
			}
			pcb.Char.Level1SlotsRemaining = curChar.Level1SlotsRemaining
			curChar.CastSuccess = true
		}

		// Secret handshake to cast a spell
		if pcb.Char.Level1SlotsRemaining == -3 {
			spellLevel := pcb.Char.Level2SlotsRemaining
			ok := spellLevel == 0
			if spellLevel >= 0 && spellLevel <= 9 {
				if spellLevel == 0 {
					ok = true
				} else if spellLevel <= 1 && curChar.Level1SlotsRemaining > 0 {
					curChar.Level1SlotsRemaining--
					ok = true
				} else if spellLevel <= 2 && curChar.Level2SlotsRemaining > 0 {
					curChar.Level2SlotsRemaining--
					ok = true
				} else if spellLevel <= 3 && curChar.Level3SlotsRemaining > 0 {
					curChar.Level3SlotsRemaining--
					ok = true
				} else if spellLevel <= 4 && curChar.Level4SlotsRemaining > 0 {
					curChar.Level4SlotsRemaining--
					ok = true
				} else if spellLevel <= 5 && curChar.Level5SlotsRemaining > 0 {
					curChar.Level5SlotsRemaining--
					ok = true
				} else if spellLevel <= 6 && curChar.Level6SlotsRemaining > 0 {
					curChar.Level6SlotsRemaining--
					ok = true
				} else if spellLevel <= 7 && curChar.Level7SlotsRemaining > 0 {
					curChar.Level7SlotsRemaining--
					ok = true
				} else if spellLevel <= 8 && curChar.Level8SlotsRemaining > 0 {
					curChar.Level8SlotsRemaining--
					ok = true
				} else if spellLevel <= 9 && curChar.Level9SlotsRemaining > 0 {
					curChar.Level9SlotsRemaining--
					ok = true
				}
			}

			curChar.CastSuccess = ok
			pcb.Char.Level1SlotsRemaining = curChar.Level1SlotsRemaining
			pcb.Char.Level2SlotsRemaining = curChar.Level1SlotsRemaining
			if ok {
				curChar.Concentrating = pcb.Char.Concentrating
			}
		}

		if pcb.Char.Concentrating == 0 {
			curChar.Concentrating = 0
		}

		if len(pcb.Char.Name) > 0 {
			curChar.Name = pcb.Char.Name
		}
		if pcb.Char.AbilityScore > 0 {
			curChar.AbilityScore = pcb.Char.AbilityScore
		}
		curChar.SetLevel(pcb.Char.Level)

		err := curChar.updateCharacter(tx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("couldn't rollback patchCharacter changes: %v", err)
			}
			log.Printf("couldn't modify character info: %v", err)
			return errorutil.New(http.StatusInternalServerError, "couldn't modify character")
		}
	}

	// if portrait included in patch, save it
	if pcb.Portrait != "" {
		log.Printf("patch includes portrait...")
		if err := c.savePortrait(tx, pcb.Portrait, pcb.Char); err != nil {
			log.Printf("couldn't save portrait err: %v", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error committing patchCharacter transaction: %v", err)
		return errorutil.New(http.StatusInternalServerError, "internal error")
	}

	log.Printf("patched character %v", pcb.Char)

	return nil
}

// deleteCharacter deletes the character in the DB with the received ID
func (c *client) deleteCharacter(cID int64, uID string) error {
	authorized, err := c.isOwner(cID, uID)
	if err != nil {
		return err
	}
	if !authorized {
		return errorutil.New(http.StatusForbidden, "can't delete someone else's character...")
	}

	if _, err := c.db.Exec(
		`DELETE FROM Characters
		 WHERE CharacterKey = ?`,
		cID); err != nil {
		log.Printf("\ncouldn't exec row delete for deleteCharacter: %v", err)
		return errorutil.New(http.StatusInternalServerError, "internal error")
	}
	return nil
}

func (c *client) isOwner(cID int64, uID string) (bool, error) {
	rows, err := c.db.Query(
		`SELECT OwnerID
		 FROM Characters
		 WHERE CharacterKey = ?`,
		cID)
	if err != nil {
		log.Printf("\ncouldn't get chararacter data for deleteCharacter: %v", err)
		return false, errorutil.New(http.StatusInternalServerError, "internal error")
	}
	if ok := rows.Next(); !ok {
		return false, errorutil.Newf(http.StatusNotFound, "character %d not found", cID)
	}
	var ownerID string
	if err := rows.Scan(&ownerID); err != nil {
		log.Printf("\ncouldn't scan query row for deleteCharacter: %v", err)
		return false, errorutil.New(http.StatusInternalServerError, "internal error")
	}
	rows.Close()
	return ownerID == uID, nil
}

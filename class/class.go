// Package class provides an interface for class and subclass data for tinfoilwizard.
package class

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/svipul/tinfoil-wizard/errorutil"
	"github.com/svipul/tinfoil-wizard/httputil"
)

var spellSlotsByLevel = map[int][]int{
	/* c  1  2  3  4  5  6  7  8  9 */
	1:  []int{3, 2, 0, 0, 0, 0, 0, 0, 0, 0},
	2:  []int{3, 3, 0, 0, 0, 0, 0, 0, 0, 0},
	3:  []int{4, 4, 2, 0, 0, 0, 0, 0, 0, 0},
	4:  []int{4, 4, 3, 0, 0, 0, 0, 0, 0, 0},
	5:  []int{4, 4, 3, 2, 0, 0, 0, 0, 0, 0},
	6:  []int{4, 4, 3, 3, 0, 0, 0, 0, 0, 0},
	7:  []int{4, 4, 3, 3, 1, 0, 0, 0, 0, 0},
	8:  []int{4, 4, 3, 3, 2, 0, 0, 0, 0, 0},
	9:  []int{4, 4, 3, 3, 3, 1, 0, 0, 0, 0},
	10: []int{5, 4, 3, 3, 3, 2, 0, 0, 0, 0},
	11: []int{5, 4, 3, 3, 3, 2, 1, 0, 0, 0},
	12: []int{5, 4, 3, 3, 3, 2, 1, 0, 0, 0},
	13: []int{5, 4, 3, 3, 3, 2, 1, 1, 0, 0},
	14: []int{5, 4, 3, 3, 3, 2, 1, 1, 0, 0},
	15: []int{5, 4, 3, 3, 3, 2, 1, 1, 1, 0},
	16: []int{5, 4, 3, 3, 3, 2, 1, 1, 1, 0},
	17: []int{5, 4, 3, 3, 3, 2, 1, 1, 1, 1},
	18: []int{5, 4, 3, 3, 3, 3, 1, 1, 1, 1},
	19: []int{5, 4, 3, 3, 3, 3, 2, 1, 1, 1},
	20: []int{5, 4, 3, 3, 3, 3, 2, 2, 1, 1},
}

// GetSpellSlots looks up in the spell slots map
func GetSpellSlots(level int, spellLevel int) int64 {
	if (level < 1 || level > 20 || spellLevel < 0 || spellLevel > 9) {
		return int64(0)
	}
	return int64(spellSlotsByLevel[level][spellLevel])
}

// Class represents a single Class entity's properties
type Class struct {
	ID                  int         `json:"id"`
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	SpellcastingAbility string      `json:"spellcastingAbility"`
	SubclassingLevel    int         `json:"subclassingLevel"`
	Subclasses          []*Subclass `json:"subclasses,omitempty"`
}

// Subclass represents a single Subclass entity's properties
type Subclass struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Client provides an interface for interacting with
// the Class and Subclass tables in tinfoilwizard's DB
type Client interface {
	HandleClasses(w http.ResponseWriter, r *http.Request)
	HandleClass(w http.ResponseWriter, r *http.Request)
	HandleSubclasses(w http.ResponseWriter, r *http.Request)
	HandleSubclass(w http.ResponseWriter, r *http.Request)
	AllClasses(includeSubs bool) ([]*Class, error)
	ClassByID(id int, includeSubs bool) (*Class, error)
	SubclassesForClass(id int) ([]*Subclass, error)
	SubclassByID(id int) (*Subclass, error)
}

// The class package's client struct encapsulates the package's actual operations,
// so consumers of the client don't have to be concerned with its implementation.
type client struct {
	db *sql.DB
}

// NewClient Creates a new class client.
func NewClient(db *sql.DB) Client {
	return &client{db: db}
}

/******************************************************************************
** HANDLERS:
	** HandleClass: ENDPOINT "/api/classes/{cID}"
	** HandleClasses: ENDPOINT "/api/classes"
	** HandleSubclass: ENDPOINT "/api/subclasses/{scID}"
	** HandleSubclasses: ENDPOINT "/api/classes/{cID}/subclasses"
******************************************************************************/

// HandleClass: GET entity from Class table by class ID, and, opt. class's subclasses
// ENDPOINT: "/api/classes/{cID}"
func (c *client) HandleClass(w http.ResponseWriter, r *http.Request) {
	log.Print("Class endpoint")
	switch r.Method {
	case http.MethodGet:
		cl, err := c.getClass(r)
		if err != nil {
			httputil.ErrorResponse(w, err)
			return
		}
		httputil.JSONResponse(w, 200, cl)
	default:
		httputil.ErrorResponse(w, errorutil.New(http.StatusForbidden, "invalid request"))
	}
}

// HandleClasses GET all entities in Class table, and, optionally, their subclasses
// ENDPOINT: "/api/classes"
func (c *client) HandleClasses(w http.ResponseWriter, r *http.Request) {
	log.Print("Classes endpoint")
	switch r.Method {
	case http.MethodGet:
		cl, err := c.getClasses(r)
		if err != nil {
			httputil.ErrorResponse(w, err)
			return
		}
		httputil.JSONResponse(w, 200, cl)
	default:
		httputil.ErrorResponse(w, errorutil.New(http.StatusForbidden, "invalid request"))
	}
}

// HandleSubclass: GET entity from Subclass table by subclass ID
// ENDPOINT: "/api/subclasses/{scID}"
func (c *client) HandleSubclass(w http.ResponseWriter, r *http.Request) {
	log.Print("Subclass endpoint")
	switch r.Method {
	case http.MethodGet:
		sc, err := c.getSubclass(r)
		if err != nil {
			httputil.ErrorResponse(w, err)
			return
		}
		httputil.JSONResponse(w, 200, sc)
	default:
		httputil.ErrorResponse(w, errorutil.New(http.StatusForbidden, "invalid request"))
	}
}

// HandleSubclasses: GET all entities in Subclass table by class ID
// ENDPOINT: "/api/classes/{cID}/subclasses"
func (c *client) HandleSubclasses(w http.ResponseWriter, r *http.Request) {
	log.Print("Subclasses endpoint")
	switch r.Method {
	case http.MethodGet:
		scs, err := c.getSubclasses(r)
		if err != nil {
			httputil.ErrorResponse(w, err)
			return
		}
		httputil.JSONResponse(w, 200, scs)
	default:
		httputil.ErrorResponse(w, errorutil.New(http.StatusForbidden, "invalid request"))
	}
}

/******************************************************************************
** HELPERS... and their helpers
	** getClass					WORKS
		** ClassByID			WORKS
	** getClasses				WORKS
		** AllClasses			WORKS
	** getSubclass				WORKS
		** SubclassByID			WORKS
	** getSubclasses			WORKS
		** SubclassesForClass	WORKS
******************************************************************************/

func (c *client) getClass(r *http.Request) (*Class, error) {
	vars := mux.Vars(r)
	cID := vars["cID"]
	id, err := strconv.Atoi(cID)
	if err != nil {
		return nil, errorutil.New(400, "request must contain a valid, integer class ID")
	}
	var includeSubs bool
	if scQuery, ok := r.URL.Query()["subclasses"]; ok {
		includeSubs = scQuery[0] == "true"
	}

	return c.ClassByID(id, includeSubs)
}

// ClassByID returns class by id with subclasses per bool param
func (c *client) ClassByID(id int, includeSubs bool) (*Class, error) {
	rows, err := c.db.Query(
		`SELECT cl.ClassKey, cl.Name, cl.Description, cl.SpellcastingAbility, cl.SubclassingLevel
		FROM Class cl
		WHERE ClassKey = ?`,
		id)
	defer rows.Close()
	if err != nil {
		log.Printf("couldn't get class (ClassKey %d) from database: %v", id, err)
		return nil, errorutil.New(http.StatusInternalServerError, "internal error")
	}
	cl := &Class{}
	if ok := rows.Next(); !ok {
		log.Printf("couldn't find row for Class ID %d", id)
		return nil, errorutil.New(404, "class not found")
	}
	if err := rows.Scan(&cl.ID, &cl.Name, &cl.Description, &cl.SpellcastingAbility, &cl.SubclassingLevel); err != nil {
		log.Printf("couldn't read row: %v", err)
		return nil, errorutil.New(500, "internal error")
	}
	if includeSubs {
		scs, err := c.SubclassesForClass(cl.ID)
		if err != nil {
			log.Printf("couldn't find subclasses for class %d: %v", cl.ID, err)
			return nil, err
		}
		cl.Subclasses = scs
	}
	return cl, nil
}

func (c *client) getClasses(r *http.Request) ([]*Class, error) {
	var includeSubs bool
	if scQuery, ok := r.URL.Query()["subclasses"]; ok {
		includeSubs = scQuery[0] == "true"
	}
	return c.AllClasses(includeSubs)
}

func (c *client) AllClasses(includeSubs bool) ([]*Class, error) {
	rows, err := c.db.Query(
		`SELECT *
		 FROM Class`)
	defer rows.Close() // close the connection when done!

	if err != nil {
		log.Printf("couldn't get classeses from database: %v", err)
		return nil, errorutil.New(http.StatusInternalServerError, "internal error")
	}
	var ac []*Class
	for rows.Next() {
		cl := &Class{}
		if err := rows.Scan(&cl.ID, &cl.Name, &cl.Description, &cl.SpellcastingAbility, &cl.SubclassingLevel); err != nil {
			log.Printf("couldn't read row: %v", err)
			return nil, errorutil.New(500, "internal error")
		}
		if includeSubs {
			scs, err := c.SubclassesForClass(cl.ID)
			if err != nil {
				log.Printf("couldn't find subclasses for class %d: %v", cl.ID, err)
				return nil, err
			}
			cl.Subclasses = scs
		}
		ac = append(ac, cl)
	}
	return ac, nil
}

func (c *client) getSubclass(r *http.Request) (*Subclass, error) {
	vars := mux.Vars(r)
	scID := vars["scID"]
	id, err := strconv.Atoi(scID)
	if err != nil {
		return nil, errorutil.New(400, "request must contain a valid, integer subclass ID")
	}
	return c.SubclassByID(id)
}

func (c *client) SubclassByID(id int) (*Subclass, error) {
	rows, err := c.db.Query(
		`SELECT sc.SubclassKey, sc.Name, sc.Description
		 FROM Subclass sc
		 WHERE SubclassKey = ?`,
		id)
	defer rows.Close()
	if err != nil {
		log.Printf("couldn't get subclass (SubclassKey %d) from database: %v", id, err)
		return nil, errorutil.New(http.StatusInternalServerError, "internal error")
	}
	if ok := rows.Next(); !ok {
		log.Printf("couldn't find row for Subclass ID %d", id)
		return nil, errorutil.New(404, "subclass not found")
	}
	sc := &Subclass{}
	if err := rows.Scan(&sc.ID, &sc.Name, &sc.Description); err != nil {
		log.Printf("couldn't read row: %v", err)
		return nil, errorutil.New(500, "internal error")
	}
	return sc, nil
}

func (c *client) getSubclasses(r *http.Request) ([]*Subclass, error) {
	vars := mux.Vars(r)
	cID := vars["cID"]
	id, err := strconv.Atoi(cID)
	if err != nil {
		return nil, errorutil.New(400, "request must contain a valid, integer class ID")
	}
	return c.SubclassesForClass(id)
}

func (c *client) SubclassesForClass(id int) ([]*Subclass, error) {
	rows, err := c.db.Query(
		`SELECT sc.SubclassKey, sc.Name, sc.Description
		 FROM Subclass sc
		 WHERE ParentClass = ?`,
		id)
	defer rows.Close()
	if err != nil {
		log.Printf("couldn't get subclasses for class (ClassKey %d) from database: %v", id, err)
		return nil, errorutil.New(http.StatusInternalServerError, "internal error")
	}
	var scs []*Subclass
	for rows.Next() {
		sc := &Subclass{}
		if err := rows.Scan(&sc.ID, &sc.Name, &sc.Description); err != nil {
			log.Printf("couldn't read row: %v", err)
			return nil, errorutil.New(500, "internal error")
		}
		scs = append(scs, sc)
	}
	return scs, nil
}

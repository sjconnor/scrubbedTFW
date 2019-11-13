package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql" //"github.com/jmoiron/sqlx" // switch to this driver if we need to filter on db results
	"github.com/gorilla/mux"
	"github.com/svipul/tinfoil-wizard/character"
	"github.com/svipul/tinfoil-wizard/class"
	"github.com/svipul/tinfoil-wizard/spell"
)

func main() {

	// configure and set up database
	db, err := configDB()
	if err != nil {
		log.Fatalf("couldn't configure DB: %v", err)
	}

	ctx := context.Background()

	clc := class.NewClient(db)
	cc, err := character.NewClient(ctx, db, clc)
	if err != nil {
		log.Fatalf("couldn't config CharacterClient: %v", err)
	}
	sc := spell.NewSpellClient(db)

	// set up routes and handling
	r := mux.NewRouter() // create router instance
	http.Handle("/", r)  // gorilla mux to catch all incoming

	// Handlers in spells.go
	r.HandleFunc("/api/spells", sc.AllSpells)
	r.HandleFunc("/api/spells/{sID}", sc.HandleSpell)
	r.HandleFunc("/api/characters/{cID}/spells", sc.CharacterSpells)
	r.HandleFunc("/api/characters/{cID}/spells/{sID}", sc.CharacterSpells)

	// Handlers in characters.go
	r.HandleFunc("/api/characters", cc.AllCharacters)
	r.HandleFunc("/api/characters/{cID}", cc.HandleCharacter)
	//r.HandleFunc("/api/characters/{cID}/level", cc.HandleCharacterSlots)
	//r.HandleFunc("/api/characters/{cID}/level/{lID}}", cc.HandleCharacterSlots)

	// Handlers in class.go
	r.HandleFunc("/api/classes", clc.HandleClasses)
	r.HandleFunc("/api/classes/{cID}", clc.HandleClass)
	r.HandleFunc("/api/classes/{cID}/subclasses", clc.HandleSubclasses)
	r.HandleFunc("/api/subclasses/{scID}", clc.HandleSubclass)

	// Serve the portraits directory
	r.PathPrefix("/portraits/").Handler(http.StripPrefix("/portraits/", http.FileServer(http.Dir("./portraits"))))

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// enforce timeouts
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Launching server on 127.0.0.1:8080")
	log.Fatal(srv.ListenAndServe())
}

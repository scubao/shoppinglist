// Package main provides a Shopping List
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

type ShoppingEntry struct {
	Id      int64  `db:"id" json:"id"`
	Created int64  `db:"created" json:"created"`
	User    string `json:"user"`
	Amount  int64  `json:"amount"`
	Name    string `json:"name"`
	Market  string `json:"market"`
	Done    bool   `json:"done"`
}

func newShoppingEntry(user, name, market string, amount int64) ShoppingEntry {
	return ShoppingEntry{
		Created: time.Now().UnixNano(),
		User:    user,
		Amount:  amount,
		Name:    name,
		Market:  market,
		Done:    false,
	}
}

func jsonShoppingEntry(p ShoppingEntry) {
	m, err := json.Marshal(p)
	checkErr(err, "Marshal failed")
	log.Printf("m = %+v\n", string(m))
}

func parseCommandline() ShoppingEntry {
	userflag := flag.String("user", "Oliver", "Username")
	produktflag := flag.String("produkt", "Eier", "Produkt")
	anzahlflag := flag.Int64("anzahl", 1, "Anzahl")
	ortflag := flag.String("ort", "Rewe", "Ort")
	//erledigtflag := flag.Bool("erledigt", false, "Erledigt")
	flag.Parse()
	return newShoppingEntry(*userflag, *produktflag, *ortflag, *anzahlflag)
	//	return shopObject{1, *userflag, *produktflag, *anzahlflag, *ortflag, *erledigtflag}
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("sqlite3", "shopping_sqlite3.bin")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(ShoppingEntry{}, "shoppingentry").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func handleEntries(w http.ResponseWriter, r *http.Request) {
	var shoppinglist []ShoppingEntry
	dbmap := initDb()
	defer dbmap.Db.Close()
	// fetch all rows
	_, err := dbmap.Select(&shoppinglist, "select * from shoppingentry order by id")
	checkErr(err, "Select failed")
	b, err := json.Marshal(shoppinglist)
	checkErr(err, "Marshal failed")
	fmt.Fprint(w, string(b))
}

func addEntry(w http.ResponseWriter, r *http.Request) {
	var shoppinglist []ShoppingEntry
	dbmap := initDb()
	defer dbmap.Db.Close()
	// fetch all rows
	_, err := dbmap.Select(&shoppinglist, "select * from shoppingentry order by id")
	checkErr(err, "Select failed")
	b, err := json.Marshal(shoppinglist)
	checkErr(err, "Marshal failed")
	fmt.Fprint(w, string(b))
}

func doneEntry(w http.ResponseWriter, r *http.Request) {
	var shoppinglist []ShoppingEntry
	dbmap := initDb()
	defer dbmap.Db.Close()
	// fetch all rows
	_, err := dbmap.Select(&shoppinglist, "select * from shoppingentry order by id")
	checkErr(err, "Select failed")
	b, err := json.Marshal(shoppinglist)
	checkErr(err, "Marshal failed")
	fmt.Fprint(w, string(b))
}

func handleEntry(w http.ResponseWriter, r *http.Request) {
	var entry ShoppingEntry
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println(id)
	dbmap := initDb()
	defer dbmap.Db.Close()

	switch r.Method {
	case "GET":
		err := dbmap.SelectOne(&entry, "select * from shoppingentry where id=?", id)
		checkErr(err, "SelectOne failed")
		b, err := json.Marshal(entry)
		checkErr(err, "Marshal failed")
		fmt.Fprint(w, string(b))
	case "DELETE":
		_, err := dbmap.Exec("delete from shoppingentry where id=?", id)
		checkErr(err, "SelectOne failed")
		w.WriteHeader(http.StatusNoContent)
	}

}

func main() {

	// initialize the DbMap
	//dbmap := initDb()
	//defer dbmap.Db.Close()

	// delete any existing rows
	// err := dbmap.TruncateTables()
	// checkErr(err, "TruncateTables failed")

	// insert rows - auto increment PKs will be set properly after the insert
	//e1 := newShoppingEntry("Oliver", "Bananen", "Rewe", 6)
	//e2 := newShoppingEntry("Oliver", "Birnen", "Aldi", 3)

	// insert two entries
	//err := dbmap.Insert(&e1, &e2)
	//checkErr(err, "Insert Failed!")

	// fetch all rows
	//	_, err = dbmap.Select(&shoppinglist, "select * from shoppingentry order by id")
	//	checkErr(err, "Select failed")
	//	for _, v := range shoppinglist {
	//		log.Println(v)
	//	}

	router := mux.NewRouter()
	router.HandleFunc("/entries", handleEntries).Methods("GET")
	router.HandleFunc("/entry/{id}", handleEntry).Methods("GET", "DELETE")
	http.ListenAndServe(":8080", router)

	//log.Println("All rows:")
	//for x, p := range shoppinglist {
	//	log.Printf("    %d: %v\n", x, p)
	//}

	// use convenience SelectInt
	//count, err := dbmap.SelectInt("select count(*) from shoppingentry")
	//checkErr(err, "select count(*) failed")
	//log.Println("Rows after inserting:", count)

	// update a row
	//e2.Name = "Äpfel"
	//count, err = dbmap.Update(&e2)
	//checkErr(err, "Update failed")
	//log.Println("Rows updated:", count)

	// fetch one row - note use of "post_id" instead of "Id" since column is aliased
	//
	// Postgres users should use $1 instead of ? placeholders
	// See 'Known Issues' below
	//
	//err = dbmap.SelectOne(&e2, "select * from shoppingentry where shopping_id=?", e2.Id)
	//checkErr(err, "SelectOne failed")
	//log.Println("e2 row:", e2)

	// delete row by PK
	//count, err = dbmap.Delete(&e1)
	//checkErr(err, "Delete failed")
	//log.Println("Rows deleted:", count)

	// delete row manually via Exec
	//_, err = dbmap.Exec("delete from shoppingentry where shopping_id=?", e2.Id)
	//checkErr(err, "Exec failed")

	// confirm count is zero
	//count, err = dbmap.SelectInt("select count(*) from shoppingentry")
	//checkErr(err, "select count(*) failed")
	//log.Println("Row count - should be zero:", count)

}

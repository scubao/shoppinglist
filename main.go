// Package main provides a Shopping List
package main

import (
	"database/sql"
	"flag"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

type shopObject struct {
	id       int64
	user     string
	produkt  string
	anzahl   int64
	ort      string
	erledigt bool
}

type ShoppingEntry struct {
	Id      int64 `db:"shopping_id"`
	Created int64
	User    string
	Amount  int64
	Name    string
	Market  string
	Done    bool
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

type Post struct {
	// db tag lets you specify the column name if it differs from the struct field
	Id      int64 `db:"post_id"`
	Created int64
	Title   string `db:",size:50"`               // Column size set to 50
	Body    string `db:"article_body,size:1024"` // Set both column name and size
}

func newPost(title, body string) Post {
	return Post{
		Created: time.Now().UnixNano(),
		Title:   title,
		Body:    body,
	}
}

func main() {
	var todos map[int64]ShoppingEntry
	todos = make(map[int64]ShoppingEntry)
	testobject := parseCommandline()
	todos[1] = testobject
	log.Println(todos)

	// initialize the DbMap
	dbmap := initDb()
	defer dbmap.Db.Close()

	// delete any existing rows
	//err := dbmap.TruncateTables()
	//checkErr(err, "TruncateTables failed")

	// create two posts
	e1 := newShoppingEntry("Oliver", "Bananen", "Rewe", 6)
	e2 := newShoppingEntry("Oliver", "Birnen", "Aldi", 3)

	// insert rows - auto increment PKs will be set properly after the insert
	err := dbmap.Insert(&e1, &e2)
	checkErr(err, "Insert failed")

	// use convenience SelectInt
	count, err := dbmap.SelectInt("select count(*) from shoppingentry")
	checkErr(err, "select count(*) failed")
	log.Println("Rows after inserting:", count)

	// update a row
	e2.Name = "Äpfel"
	count, err = dbmap.Update(&e2)
	checkErr(err, "Update failed")
	log.Println("Rows updated:", count)

	// fetch one row - note use of "post_id" instead of "Id" since column is aliased
	//
	// Postgres users should use $1 instead of ? placeholders
	// See 'Known Issues' below
	//
	err = dbmap.SelectOne(&e2, "select * from shoppingentry where shopping_id=?", e2.Id)
	checkErr(err, "SelectOne failed")
	log.Println("e2 row:", e2)

	// fetch all rows
	var shoppingentries []ShoppingEntry
	_, err = dbmap.Select(&shoppingentries, "select * from shoppingentry order by shopping_id")
	checkErr(err, "Select failed")
	log.Println("All rows:")
	for x, p := range shoppingentries {
		log.Printf("    %d: %v\n", x, p)
	}

	// delete row by PK
	count, err = dbmap.Delete(&e1)
	checkErr(err, "Delete failed")
	log.Println("Rows deleted:", count)

	// delete row manually via Exec
	_, err = dbmap.Exec("delete from shoppingentry where shopping_id=?", e2.Id)
	checkErr(err, "Exec failed")

	// confirm count is zero
	count, err = dbmap.SelectInt("select count(*) from shoppingentry")
	checkErr(err, "select count(*) failed")
	log.Println("Row count - should be zero:", count)

	log.Println("Done!")

}

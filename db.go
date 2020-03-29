package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"
)

func initDB() {
	var err error
	db, err = sql.Open("postgres", config.DB)
	if err != nil {
		return
	}
	Pongerr := db.Ping()
	if Pongerr != nil {
		log.Fatal("Database is down\n")
		return
	}
	ProcFolders(config.Folders)

}

// ExistsWithinDB returns whether the sha1 exists or not.
func ExistsWithinDB(Sha1 string) bool {
	var exists bool
	err := db.QueryRow(`SELECT * from items where item @> '{"sha1":"` + Sha1 + `"}'`).Scan(&exists)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	return true
}

// UpdateLocation is used to update the location of a file in the database with the same sha1 as the one given as an argument.
func UpdateLocation(Value Item) {
	db.Query(`UPDATE items SET item = jsonb_set(item, '{location}', '["` + strings.Join(Value.File, `","`) + `"]') where item->>'sha1' = '` + Value.Sha1 + `';`)
}

// UpdateTags is used to update the existing tags of an item with the same sha1 as Value
// Func UpdateTags(Value Item){}

// UpdateRow removes the item with the same sha1 as Value and adds Value into the database.
func UpdateRow(Value Item) {
	db.Query(`DELETE from items where item @> '{"sha1":"` + Value.Sha1 + `"}'`)
	AddItem(Value)
}

// AddItem adds the given item into the database.
func AddItem(Value Item) {
	marsh, err := json.Marshal(Value)
	if err != nil {
		panic(err)
	}
	db.Query(`insert into items (item) values ('` + string(marsh) + `')`)
}

// QuerySha is used to query the database to find an item with a given hash.
func QuerySha(Hash string) *sql.Row {
	return db.QueryRow(`select * from items where item @> '{"sha1":"` + Hash + `"}'`)
}

// QueryLikeItem returns rows that are like Value.
func QueryLikeItem(Value Item) (*sql.Rows, error) {
	var where []string
	s := "select * from items "
	if len(Value.Tags) != 0 {
		if Value.Tags[0] == "not null" { // Special case to return items that have tags.
			return db.Query("select * from items where item->>'tags' is not null")
		}
		s := `item->'tags' @> '["` + strings.Join(Value.Tags, `","`) + `"]'`
		where = append(where, s)
	}
	if len(Value.File) != 0 {
		s := `item->'location' @> '["` + strings.Join(Value.File, `","`) + `"]'`
		where = append(where, s)
	}
	if Value.Thumbnail != "" {
		where = append(where, `item @> '{"thumbnail":"`+Value.Thumbnail+`"}'`)
	}
	if Value.Sha1 != "" {
		where = append(where, `item @> '{"sha1":"`+Value.Sha1+`"}'`)
	}
	if len(Value.Tags) != 0 || len(Value.File) != 0 || Value.Sha1 != "" || Value.Thumbnail != "" {
		s = `select * from items where ` + strings.Join(where, " OR ")
	}

	return db.Query(s)
}

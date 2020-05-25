package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"
)

// Item is the main json used to read from and write to the database, and to communicate with the HTML front-end.
type Item struct {
	File      []string `json:"location"`
	Thumbnail string   `json:"thumbnail"`
	Tags      []string `json:"tags"`
	Sha1      string   `json:"sha1"`
	Mode      bool     `json:"strict"`
	Size      int64    `json:"size"`
}

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
func QuerySha(Hash string) []Item {
	return QueryToItemArray(db.Query(`select * from items where item @> '{"sha1":"` + Hash + `"}'`))
}

// StrictQuery is used to query using exact match.
func StrictQuery(Value Item) []Item {
	var where []string
	s := "select * from items "
	if len(Value.Tags) != 0 {
		if Value.Tags[0] == "not null" { // Special case to return items that have tags.
			return QueryToItemArray(db.Query("select * from items where item->>'tags' is not null"))
		} else if Value.Tags[0] != "" {
			s := `item->'tags' @> '["` + strings.Join(Value.Tags, `","`) + `"]'`
			where = append(where, s)
		}
	}
	if len(Value.File) != 0 && Value.File[0] != "" {
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
		s = `select * from items where ` + strings.Join(where, " AND ")
	}

	return QueryToItemArray(db.Query(s))
}

//SubstrQuery queries using *any* matchup.
func SubstrQuery(Value Item) []Item {
	var where []string
	s := "select * from items "
	if len(Value.Tags) != 0 {
		if Value.Tags[0] == "not null" { // Special case to return items that have tags.
			return QueryToItemArray(db.Query("select * from items where item->>'tags' is not null"))
		} else if Value.Tags[0] != "" {
			s := `item->>'tags' ~ '` + strings.Join(Value.Tags, `' AND item->>'tags' ~ '`) + `'`
			where = append(where, s)
		}
	}
	if len(Value.File) != 0 && Value.File[0] != "" {
		s := `item->>'location' ~ '` + Value.File[0] + `'`
		where = append(where, s)
	}
	if Value.Thumbnail != "" {
		where = append(where, `item @> '{"thumbnail":"`+Value.Thumbnail+`"}'`)
	}
	if Value.Sha1 != "" {
		where = append(where, `item->>'sha1' ~ '`+Value.Sha1+`'`)
	}
	if len(Value.Tags) != 0 || len(Value.File) != 0 || Value.Sha1 != "" || Value.Thumbnail != "" {
		s = `select * from items where ` + strings.Join(where, " AND ")
	}
	return QueryToItemArray(db.Query(s))
}

// QueryLikeItem returns rows that are like Value.
func QueryLikeItem(Value Item) []Item {
	if !Value.Mode {
		return StrictQuery(Value)
	}
	return SubstrQuery(Value)

}

//QueryToItemArray takes a query and turns turns them into Item objects then returns them all in an array.
func QueryToItemArray(rows *sql.Rows, err error) []Item {
	if rows == nil {
		return []Item{}
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var received Item
		var testing string
		err := rows.Scan(&testing)
		if err != nil {
			log.Println(err)
		}
		json.Unmarshal([]byte(testing), &received)
		items = append(items, received)
	}
	return items
}

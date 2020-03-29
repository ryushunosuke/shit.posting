package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// File is used for writing into the database from a folder.
// type File struct {
// 	Location string
// 	Sha1     string
// }

// Sha1 returns the sha1 hashed string value.
func Sha1(data []byte) string {
	Hasher := sha1.New()
	Hasher.Write(data)
	bs := Hasher.Sum(nil)
	return fmt.Sprintf("%x", bs)

}

// ProcFolders processes folders and puts the files within the folders to the database.
func ProcFolders(Folders []string) {
	for _, folder := range Folders {
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			ftype, err := GetFileType(folder + "/" + file.Name())
			if err != nil {
				panic(err)
			}
			ftype = Convert(ftype)
			if !config.TypeMap[ftype] {
				continue
			}
			if file.Size() > 3*1024*1024*1024 {
				continue
			}
			f, err := ioutil.ReadFile(folder + "/" + file.Name())
			Hashed := Sha1(f)
			if err != nil {
				panic(err)
			}
			if !ThumbnailExists(Hashed + ".jpg") {
				ThumbnailFile(folder+"/"+file.Name(), Hashed)
			}
			if !ExistsWithinDB(Hashed) {
				AddItem(Item{File: []string{folder + "/" + file.Name()}, Thumbnail: "thumbnail/" + Hashed + ".jpg", Sha1: Hashed})
				if !ThumbnailExists(Hashed + ".jpg") {
					ThumbnailFile(folder+"/"+file.Name(), Hashed)
				}
				fmt.Printf("Added " + Hashed + " to the database.\n")

			} else {
				row := QuerySha(Hashed)
				var Value Item
				var between string
				row.Scan(&between)
				json.Unmarshal([]byte(between), &Value)
				Dupe := false
				for _, a := range Value.File {
					if a == folder+"/"+file.Name() {
						Dupe = true
					}
				}
				if Dupe {
					continue
				}
				Value.Sha1 = Hashed
				Value.File = append(Value.File, folder+"/"+file.Name())
				UpdateLocation(Value)
				fmt.Printf("Updated " + Value.Sha1 + " in the database.\n")

			}
			if err != nil {
				panic(err)
			}
			fmt.Printf("File: %v, type: %v\n", file.Name(), ftype)
		}

	}
}

// ThumbnailExists checks if the given file has a thumbnail in the path ./thumbnail/{sha1}.{ftype}
func ThumbnailExists(Path string) bool {
	_, err := os.Stat("./thumbnail/" + Path)
	if os.IsNotExist(err) {
		return false
	}
	return true

}

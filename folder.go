package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
)

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
			debug.FreeOSMemory() // Used due to windoms problems.
			ftype, err := GetFileType(folder + "/" + file.Name())
			if err != nil {
				panic(err)
			}
			ftype = Convert(ftype)
			if !config.TypeMap[ftype] {
				continue
			}
			if file.Size() > config.Filesize { // Max file size.
				continue
			}
			Test := QueryLikeItem(Item{
				File: []string{folder + "/" + file.Name()},
			})
			if len(Test) != 0 && Test[0].Size == file.Size() { // An item with the same location and size is already within the database.
				continue
			}
			if len(Test) != 0 && Test[0].Size == 0 {
				Test[0].Size = file.Size()
				UpdateRow(Test[0])
				fmt.Printf("Updated size of item with hash %s to size %d\n", Test[0].Sha1, Test[0].Size)
				continue
			}
			f, err := ioutil.ReadFile(folder + "/" + file.Name())
			Hashed := Sha1(f)
			if err != nil {
				panic(err)
			}
			if !ThumbnailExists(Hashed) {
				ThumbnailFile(folder+"/"+file.Name(), Hashed)
			}
			if !ExistsWithinDB(Hashed) {
				AddItem(Item{File: []string{folder + "/" + file.Name()}, Thumbnail: "thumbnail/" + Hashed + ".jpg", Sha1: Hashed})
				if !ThumbnailExists(Hashed) {
					ThumbnailFile(folder+"/"+file.Name(), Hashed)
				}
				fmt.Printf("Added %s to the database.\n", Hashed)

			} else {
				Value := QuerySha(Hashed)[0]
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
				UpdateRow(Value)
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
	_, err := os.Stat("./thumbnail/" + Path + ".jpg")
	if os.IsNotExist(err) {
		return false
	}
	return true

}

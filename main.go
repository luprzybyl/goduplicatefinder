package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type fileInfo struct {
	hash     string
	fullname string
}

type fileInfos []fileInfo

func (slice fileInfos) Len() int {
	return len(slice)
}

func (slice fileInfos) Less(i, j int) bool {
	return slice[i].hash < slice[j].hash
}

func (slice fileInfos) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func writeHash(filepath string, ch chan fileInfo) {
	ch <- fileInfo{getHash(filepath), filepath}
}

func getHash(filepath string) string {
	hasher := sha256.New()
	s, err := ioutil.ReadFile(filepath)
	hasher.Write(s)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func getFolderContents(folderpath string) ([]string, []string) {
	files := []string{}
	folders := []string{}

	fls, err := ioutil.ReadDir(folderpath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fls {
		if file.IsDir() {
			folders = append(folders, folderpath+"/"+file.Name())
		} else {
			files = append(files, folderpath+"/"+file.Name())
		}
	}
	return files, folders
}

func processDir(folderpath string) {
	fmt.Println("Processing " + folderpath)
	files, folders := getFolderContents(folderpath)

	filechan := make(chan fileInfo, len(files))
	defer close(filechan)
	for _, file := range files {
		go writeHash(file, filechan)
	}
	for range files {
		filesTable = append(filesTable, <-filechan)
	}
	folderchan := make(chan bool, len(folders))
	defer close(folderchan)
	for _, folder := range folders {
		processDir(folder)
	}
}

var filesTable fileInfos

func main() {

	dirs := os.Args[1:]

	if len(dirs) == 0 {
		processDir(".")
	}

	var prev fileInfo

	for _, dir := range dirs {
		processDir(dir)
	}

	fmt.Println("sorting...")
	sort.Sort(filesTable)
	fmt.Println("done")

	for _, f := range filesTable {
		if f.hash == prev.hash {
			fmt.Println("----------------------------------------------------------------------------------------------------------------------------")
			fmt.Println(f.hash, f.fullname)
			fmt.Println(prev.hash, prev.fullname)
			fmt.Println("----------------------------------------------------------------------------------------------------------------------------")
		}
		prev = f
	}

}

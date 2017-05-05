package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
)

type fileInfo struct {
	hash     string
	name     string
	fullname string
	size     string
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

func getHash(filepath string) string {
	hasher := sha256.New()
	s, err := ioutil.ReadFile(filepath)
	hasher.Write(s)
	if err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func getFolderContents(folderpath string, files []fileInfo) []fileInfo {
	fmt.Println(len(files), folderpath)
	d, err := os.Open(folderpath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, fi := range fi {
		if fi.Mode().IsRegular() {
			//fmt.Println(getHash(folderpath+"/"+fi.Name()), fi.Name(), fi.Size(), "bytes")
			files = append(files, fileInfo{getHash(folderpath + "/" + fi.Name()), fi.Name(), folderpath + "/" + fi.Name(), strconv.FormatInt(fi.Size(), 10)})
		} else if fi.Mode().IsDir() {
			files = getFolderContents(folderpath+"/"+fi.Name(), files)
		}
	}
	return files
}

func main() {
	var files fileInfos
	var prev fileInfo

	files = getFolderContents("/data", files)
	fmt.Println("sorting...")
	sort.Sort(files)
	fmt.Println("done")

	for _, f := range files {
		if f.hash == prev.hash {
			fmt.Println("----------------------------------------------------------------------------------------------------------------------------")
			fmt.Println(f.hash, f.name, f.size, f.fullname)
			fmt.Println(prev.hash, prev.name, prev.size, prev.fullname)
			fmt.Println("----------------------------------------------------------------------------------------------------------------------------")
		}
		//fmt.Println(f.hash, f.name, f.size, f.fullname)
		prev = f
	}
}

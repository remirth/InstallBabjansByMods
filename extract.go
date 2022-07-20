package main

import (
	"log"
	"os"
	"strings"
)

type ExtractAddress struct {
	zipfile     string
	destination string
}

type ExtractAddresses []ExtractAddress

func (ea ExtractAddresses) includesZipfile(zipfile string) bool {
	for _, v := range ea {
		if v.zipfile == zipfile {
			return true
		}
	}
	return false
}

func (ea ExtractAddresses) extract(src string, workdir string, dest string) {

	for _, v := range ea {
		err := Unzip(src+"\\"+v.zipfile, workdir)

		if err != nil {
			log.Fatal(err)
		}

		unzippedDir := workdir + "\\" + strings.Split(v.zipfile, ".zip")[0]
		files, err := os.ReadDir(unzippedDir)

		if err != nil {
			log.Fatal(err)
		}

		os.MkdirAll(dest+"\\"+v.destination, 0755)

		for _, file := range files {
			if !pathExists(dest + "\\" + v.destination + "\\" + file.Name()) {
				err := CopyFile(unzippedDir+"\\"+file.Name(), dest+"\\"+v.destination+"\\"+file.Name())
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}

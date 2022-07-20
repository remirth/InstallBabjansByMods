package main

import (
	"os"
	"strings"

	"github.com/fatih/color"
)

type ExtractAddress struct {
	zipfile     string
	destination string
}

type ExtractAddresses []ExtractAddress

type WorkerInfo struct {
	dir         string
	destination string
}

func (ea ExtractAddresses) includesZipfile(zipfile string) bool {
	for _, v := range ea {
		if v.zipfile == zipfile {
			return true
		}
	}
	return false
}

func (ea ExtractAddresses) extract(src, workdir, dest string, Quit chan<- bool) {
	wi := make(chan WorkerInfo)
	go ExtractionWorker(wi, Quit, src, workdir, dest)
	for _, v := range ea {
		err := Unzip(src+"\\"+v.zipfile, workdir)

		checkPanic(err)

		os.MkdirAll(dest+"\\"+v.destination, 0755)

		unzippedDir := workdir + "\\" + strings.Split(v.zipfile, ".zip")[0]

		wi <- WorkerInfo{unzippedDir, v.destination}
	}

	close(wi)

}

func ExtractionWorker(infos <-chan WorkerInfo, Quit chan<- bool, src, workdir, dest string) {
	for info := range infos {
		files, err := os.ReadDir(info.dir)
		checkPanic(err)

		for i := 0; i < len(files); i++ {
			if !pathExists(dest + "\\" + info.destination + "\\" + files[i].Name()) {
				color.Cyan("assssss")
				err := CopyFile(info.dir+"\\"+files[i].Name(), dest+"\\"+info.destination+"\\"+files[i].Name())
				checkPanic(err)
				color.Cyan("Successfully added " + files[i].Name() + " to " + info.destination + ".")
			} else {
				color.Cyan(files[i].Name() + " already exists in " + info.destination + ".")
			}
		}
	}
	Quit <- true

}

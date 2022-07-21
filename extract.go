package main

import (
	"os"
	"strings"

	"github.com/fatih/color"
)

// The information needed to extract a zipfile to it's destination in .minecraft.
type ExtractBundle struct {
	zipfile     string
	destination string
}

// A slice of ExtractBundle structs.
type ExtractRecord []ExtractBundle

// ExtractionWorkerOptions contains the options needed for the ExtractionWorker to move the files.
type ExtractionWorkerOptions struct {
	workdir string
	dest    string
}

// includesZipFile returns true if the record contains the given zipfile.
func (ea ExtractRecord) includesZipfile(zipfile string) bool {
	for _, v := range ea {
		if v.zipfile == zipfile {
			return true
		}
	}
	return false
}

// Extract extracts the record's zipfiles to their destination.
// src corresponds to the source directory, workdir corresponds to the temporary directory, and dest corresponds to the destination directory.
func (ea ExtractRecord) extract(src, workdir, dest string, Quit chan<- bool) {
	workerChannel := make(chan ExtractionWorkerOptions)
	go moveExtractionsWorker(workerChannel, Quit, src, workdir, dest)
	for _, v := range ea {
		err := Unzip(src+"\\"+v.zipfile, workdir)

		CheckPanic(err)

		os.MkdirAll(dest+"\\"+v.destination, 0755)

		unzippedDir := workdir + "\\" + strings.Split(v.zipfile, ".zip")[0]

		workerChannel <- ExtractionWorkerOptions{unzippedDir, v.destination}
	}

	close(workerChannel)

}

// ExtractionWorker copies the files indicated by the given options-channel to the destination indicated by the options.
// src corresponds to the source directory, workdir corresponds to the temporary directory, and dest corresponds to the destination directory.
func moveExtractionsWorker(options <-chan ExtractionWorkerOptions, Quit chan<- bool, src, workdir, dest string) {
	for option := range options {
		files, err := os.ReadDir(option.workdir)
		CheckPanic(err)

		for _, file := range files {
			if !PathExists(dest + "\\" + option.dest + "\\" + file.Name()) {
				err := CopyFile(option.workdir+"\\"+file.Name(), dest+"\\"+option.dest+"\\"+file.Name())
				CheckPanic(err)
				color.Cyan("Successfully added " + file.Name() + " to " + option.dest + ".")
			} else {
				color.Cyan(file.Name() + " already exists in " + option.dest + ".")
			}
		}
	}
	Quit <- true

}

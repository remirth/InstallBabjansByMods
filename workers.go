package main

import (
	"os"
	"strings"
)

type SourceFile struct {
	Filename string
	Path     string
}

func extractModWorker(zipfiles <-chan string, results chan<- string, src string, workdir string, dest string) {
	mods := make(chan SourceFile)

	go moveModWorker(mods, results, dest)
	for file := range zipfiles {
		err := Unzip(src+"\\"+file, workdir)

		checkPanic(err)

		fullDir := workdir + "\\" + strings.Split(file, ".zip")[0]

		files, err := os.ReadDir(fullDir)

		checkPanic(err)

		for _, file := range files {
			mods <- SourceFile{file.Name(), fullDir}
		}
	}

	close(mods)
}

func moveModWorker(mods <-chan SourceFile, results chan<- string, dest string) {
	for source := range mods {
		if !pathExists(dest + "\\mods\\" + source.Filename) {
			err := CopyFile(source.Path+"\\"+source.Filename, dest+"\\mods\\"+source.Filename)
			checkPanic(err)
			results <- "Successfully added " + source.Filename + " to " + "mods folder."
		} else {
			results <- source.Filename + " already exists in " + "mods folder."
		}
	}

	close(results)
}

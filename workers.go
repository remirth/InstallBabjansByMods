package main

import (
	"os"
	"strings"
)

type SourceFile struct {
	Filename string
	Path     string
}

// extractModWorker extracts the mod files from the given zip file to the mods directory.
// src corresponds to the source directory, workdir corresponds to the temporary directory, and dest corresponds to the destination directory.
func extractModWorker(zipfiles <-chan string, results chan<- string, src, workdir, dest string) {
	mods := make(chan SourceFile)

	go moveModWorker(mods, results, dest)
	for file := range zipfiles {
		err := Unzip(src+"\\"+file, workdir)

		CheckPanic(err)

		fullDir := workdir + "\\" + strings.Split(file, ".zip")[0]

		files, err := os.ReadDir(fullDir)

		CheckPanic(err)

		for _, file := range files {
			mods <- SourceFile{file.Name(), fullDir}
		}
	}

	close(mods)
}

// moveModWorker moves the mod files from the temporary directory to the mods directory.
func moveModWorker(mods <-chan SourceFile, results chan<- string, dest string) {
	for source := range mods {
		if !PathExists(dest + "\\mods\\" + source.Filename) {
			err := CopyFile(source.Path+"\\"+source.Filename, dest+"\\mods\\"+source.Filename)
			CheckPanic(err)
			results <- "Successfully added " + source.Filename + " to " + "mods folder."
		} else {
			results <- source.Filename + " already exists in " + "mods folder."
		}
	}

	close(results)
}

package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
)

// extractMods extracts the mod files from the given zip file to the mods directory.
// src corresponds to the source directory, workdir corresponds to the temporary directory, and dest corresponds to the destination directory.
func extractMods(fname, src, workdir, dest string) {

	wg := sync.WaitGroup{}
	err := Unzip(src+"\\"+fname, workdir)

	CheckPanic(err)

	fullDir := workdir + "\\" + strings.Split(fname, ".zip")[0]

	files, err := os.ReadDir(fullDir)

	CheckPanic(err)

	for _, file := range files {
		go func(file string) {
			wg.Add(1)
			moveMod(file, fullDir, dest)
			wg.Done()
		}(file.Name())
	}

	wg.Wait()
}

// moveMod moves a mod files from the temporary directory to the mods directory.
func moveMod(filename, src, dest string) {
	if !PathExists(dest + "\\mods\\" + filename) {
		err := CopyFile(src+"\\"+filename, dest+"\\mods\\"+filename)
		CheckPanic(err)
		color.Cyan(fmt.Sprint("Successfully added " + filename + " to " + "mods folder."))
	} else {
		color.Cyan(fmt.Sprint(filename + " already exists in " + "mods folder."))
	}

}

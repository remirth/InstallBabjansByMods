package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/zs5460/art"
)

func main() {
	mcDir := os.Getenv("APPDATA") + "\\.minecraft"
	appDataDir := os.Getenv("APPDATA") + "\\BabjansByTemp"
	zipsDownload := "https://github.com/exsjabe/BabjansByZips/archive/refs/heads/master.zip"
	wg := sync.WaitGroup{}
	os.RemoveAll(appDataDir)

	toExtract := ExtractRecord{
		ExtractBundle{".ResourcePacks.zip", "resourcepacks"},
		ExtractBundle{".Shaders.zip", "shaderpacks"},
	}

	if !CommandExists("java") {
		log.Fatal("Please install Java before initializing mods.")
	}

	if !PathExists(mcDir) {
		log.Fatal("Please install Minecraft before initializing mods.")
	}

	print(art.String("BabjansBy"))
	fmt.Println("")

	if !PathExists(mcDir + "\\versions\\1.16.5-forge-36.2.23") {
		wg.Add(1)
		go func() {
			InstallForge(appDataDir, mcDir)
			wg.Done()
		}()
	}
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = ("Downloading mods...")
	s.Start()
	fileName := DownloadFile(zipsDownload, appDataDir)
	s.Suffix = ("Extracting mods...")

	Unzip(appDataDir+"\\"+fileName, appDataDir)

	folderName := appDataDir + "\\" + strings.Split(fileName, ".zip")[0]

	wg.Add(1)
	go func() {
		defer wg.Done()
		HandleExtraction(toExtract, folderName, appDataDir, mcDir)
	}()

	zips, err := os.ReadDir(folderName)

	CheckError(err)

	os.MkdirAll(mcDir+"\\mods", 0755)

	res := make(chan string)
	go func() {
		modsWg := sync.WaitGroup{}
		defer close(res)
		for _, zip := range zips {
			if toExtract.includesZipfile(zip.Name()) {
				continue
			}
			modsWg.Add(1)
			go func(fname string) {
				defer modsWg.Done()
				extractMods(fname, folderName, appDataDir, mcDir, res)
			}(zip.Name())

		}
		modsWg.Wait()
	}()

	for str := range res {
		color.Yellow(str)
	}

	wg.Wait()
	s.Suffix = ("Cleaning up...")
	os.RemoveAll(appDataDir)
	color.Magenta("Removed all temp files.")
	s.Stop()
	color.Green("Successfully installed all BabjansBy mods!")
	color.Magenta("Press any key to close this window...")
	fmt.Scanln()

}

func HandleExtraction(er ExtractRecord, src, tempdir, targetDir string) {
	res := make(chan string)
	go func() {
		defer close(res)
		er.extract(src, tempdir, targetDir, res)
	}()
	for r := range res {
		color.Cyan(r)
	}
}
func InstallForge(workdir, dest string) {
	color.Magenta("Installing Forge...")
	forgeDownload := "https://github.com/exsjabe/1.16.5-forge-36.2.23/archive/refs/heads/master.zip"
	forgeZip := DownloadFile(forgeDownload, workdir)
	err := Unzip(workdir+"\\"+forgeZip, workdir)
	CheckError(err)
	forgeFolder := workdir + "\\" + strings.Split(forgeZip, ".zip")[0]

	fullDir := forgeFolder + "\\forge-1.16.5-36.2.23-installer.jar"

	err = exec.Command("java", "-jar", fullDir).Run()
	os.RemoveAll(".//forge-1.16.5-36.2.23-installer.jar.log")
	CheckError(err)
	color.Green("Successfully installed Forge!")

}

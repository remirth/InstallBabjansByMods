package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	spinner "github.com/alecrabbit/go-cli-spinner"
	"github.com/fatih/color"
	"github.com/zs5460/art"
)

func main() {
	mcDir := os.Getenv("APPDATA") + "\\.minecraft"
	appDataDir := os.Getenv("APPDATA") + "\\BabjansByTemp"
	zipsDownload := "https://github.com/exsjabe/BabjansByZips/archive/refs/heads/master.zip"
	closeExtraction := make(chan bool)
	closeInstall := make(chan bool)
	os.RemoveAll(appDataDir)

	toExtract := ExtractAddresses{
		ExtractAddress{".ResourcePacks.zip", "resourcepacks"},
		ExtractAddress{".Shaders.zip", "shaderpacks"},
	}

	if !commandExists("java") {
		log.Fatal("Please install Java before initializing mods.")
	}

	if !pathExists(mcDir) {
		log.Fatal("Please install Minecraft before initializing mods.")
	}

	print(art.String("BabjansBy"))
	fmt.Println("")

	if !pathExists(mcDir + "\\versions\\1.16.5-forge-36.2.23") {
		go installForge(appDataDir, mcDir, closeInstall)
	} else {
		close(closeInstall)
	}

	s, _ := spinner.New()
	s.Message("Downloading mods...")
	s.Start()
	fileName := downloadFile(zipsDownload, appDataDir)
	s.Message("Extracting mods...")

	Unzip(appDataDir+"\\"+fileName, appDataDir)

	folderName := appDataDir + "\\" + strings.Split(fileName, ".zip")[0]

	go toExtract.extract(folderName, appDataDir, mcDir, closeExtraction)

	zips, err := os.ReadDir(folderName)

	checkError(err)

	os.MkdirAll(mcDir+"\\mods", 0755)

	files := make(chan string, len(zips))

	workerResults := [3]chan string{make(chan string), make(chan string), make(chan string)}
	go extractModWorker(files, workerResults[0], folderName, folderName, mcDir)
	go extractModWorker(files, workerResults[1], folderName, folderName, mcDir)
	go extractModWorker(files, workerResults[2], folderName, folderName, mcDir)

	for _, zip := range zips {
		if toExtract.includesZipfile(zip.Name()) {
			continue
		}

		color.Magenta(fmt.Sprint("Adding ", zip.Name(), " to worker channel."))
		files <- zip.Name()
	}

	close(files)

	for i, result := range workerResults {
		for r := range result {
			color.Yellow(fmt.Sprint("Worker ", i+1, ": ", r))
		}
	}

	if <-closeExtraction {
		close(closeExtraction)
	}

	if <-closeInstall {
		close(closeInstall)
	}

	s.Message("Cleaning up...")
	os.RemoveAll(appDataDir)
	s.Stop()
	color.Magenta("Removed all temp files.")
	color.Green("Successfully installed all BabjansBy mods!.")

}

func installForge(workdir string, dest string, Quit chan<- bool) {
	color.Magenta("Installing Forge...")
	forgeDownload := "https://github.com/exsjabe/1.16.5-forge-36.2.23/archive/refs/heads/master.zip"
	forgeZip := downloadFile(forgeDownload, workdir)
	err := Unzip(workdir+"\\"+forgeZip, workdir)
	checkError(err)
	forgeFolder := workdir + "\\" + strings.Split(forgeZip, ".zip")[0]

	fullDir := forgeFolder + "\\forge-1.16.5-36.2.23-installer.jar"

	err = exec.Command("java", "-jar", fullDir).Run()
	os.RemoveAll(".//forge-1.16.5-36.2.23-installer.jar.log")
	checkError(err)
	color.Green("Successfully installed Forge!")

	Quit <- true
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	spinner "github.com/alecrabbit/go-cli-spinner"
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
	s, _ := spinner.New()
	s.Message("Downloading mods...")
	s.Start()
	fileName := DownloadFile(zipsDownload, appDataDir)
	s.Message("Extracting mods...")

	Unzip(appDataDir+"\\"+fileName, appDataDir)

	folderName := appDataDir + "\\" + strings.Split(fileName, ".zip")[0]

	wg.Add(1)
	go func() {
		toExtract.extract(folderName, appDataDir, mcDir)
		wg.Done()
	}()

	zips, err := os.ReadDir(folderName)

	CheckError(err)

	os.MkdirAll(mcDir+"\\mods", 0755)

	for _, zip := range zips {
		if toExtract.includesZipfile(zip.Name()) {
			continue
		}
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()
			extractMods(fname, folderName, appDataDir, mcDir)
		}(zip.Name())

	}

	wg.Wait()
	s.Message("Cleaning up...")
	os.RemoveAll(appDataDir)
	color.Magenta("Removed all temp files.")
	s.Stop()
	color.Green("Successfully installed all BabjansBy mods!")
	color.Magenta("Press any key to close this window...")
	fmt.Scanln()

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

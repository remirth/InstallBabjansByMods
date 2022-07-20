package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	mcDir := os.Getenv("APPDATA") + "\\.minecraft"
	appDataDir := os.Getenv("APPDATA") + "\\BabjansByMods"
	zipsDownload := "https://github.com/exsjabe/BabjansByZips/archive/refs/heads/master.zip"
	forgeDownload := "https://github.com/exsjabe/1.16.5-forge-36.2.23/archive/refs/heads/master.zip"

	if !commandExists("java") {
		log.Fatal("Please install Java before before initalizing mods.")
	}

	if !pathExists(mcDir) {
		log.Fatal("Please install Minecraft before initalizing mods.")
	}

	if pathExists(mcDir + "\\versions\\1.16.5-forge-36.2.23") {
		forgeZip := downloadFile(forgeDownload, appDataDir)
		err := Unzip(appDataDir+"\\"+forgeZip, appDataDir)
		if err != nil {
			log.Fatal(err)
		}
		forgeFolder := appDataDir + "\\" + strings.Split(forgeZip, ".zip")[0]

		fullDir := forgeFolder + "\\forge-1.16.5-36.2.23-installer.jar"

		cmd := exec.Command("java", "-jar", fullDir)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		log.Fatal("ass")
	}

	fmt.Println(mcDir)
	fmt.Println(forgeDownload)

	fileName := downloadFile(zipsDownload, appDataDir)
	fmt.Println(fileName)

	Unzip(appDataDir+"\\"+fileName, appDataDir)

	folderName := strings.Split(fileName, ".zip")[0]

	fmt.Println(folderName)

}

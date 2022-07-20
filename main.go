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
	appDataDir := os.Getenv("APPDATA") + "\\BabjansByTemp"
	zipsDownload := "https://github.com/exsjabe/BabjansByZips/archive/refs/heads/master.zip"
	forgeDownload := "https://github.com/exsjabe/1.16.5-forge-36.2.23/archive/refs/heads/master.zip"

	if !commandExists("java") {
		log.Fatal("Please install Java before before initializing mods.")
	}

	if !pathExists(mcDir) {
		log.Fatal("Please install Minecraft before initializing mods.")
	}

	if !pathExists(mcDir + "\\versions\\1.16.5-forge-36.2.23") {
		forgeZip := downloadFile(forgeDownload, appDataDir)
		err := Unzip(appDataDir+"\\"+forgeZip, appDataDir)
		if err != nil {
			log.Fatal(err)
		}
		forgeFolder := appDataDir + "\\" + strings.Split(forgeZip, ".zip")[0]

		fullDir := forgeFolder + "\\forge-1.16.5-36.2.23-installer.jar"

		err = exec.Command("java", "-jar", fullDir).Run()
		if err != nil {
			log.Fatal(err)
		}

		os.RemoveAll(forgeFolder)
		os.RemoveAll(appDataDir + "\\" + forgeZip)
	}

	fmt.Println(mcDir)
	fmt.Println(forgeDownload)

	fileName := downloadFile(zipsDownload, appDataDir)
	fmt.Println(fileName)

	Unzip(appDataDir+"\\"+fileName, appDataDir)

	folderName := appDataDir + "\\" + strings.Split(fileName, ".zip")[0]

	fmt.Println(folderName)

	err := Unzip(folderName+"\\.ResourcePacks.zip", appDataDir)

	if err != nil {
		log.Fatal(err)
	}

	packsDir := appDataDir + "\\.ResourcePacks"
	resourcePacks, err := os.ReadDir(packsDir)

	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(mcDir+"\\resourcepacks", 0755)

	for _, pack := range resourcePacks {
		if !pathExists(mcDir + "\\resourcepacks\\" + pack.Name()) {
			err := CopyFile(packsDir+"\\"+pack.Name(), mcDir+"\\resourcepacks\\"+pack.Name())
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	zips, err := os.ReadDir(folderName)

	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(mcDir+"\\mods", 0755)

	for _, zip := range zips {
		if zip.Name() == ".ResourcePacks.zip" || zip.Name() == ".Shaders.zip" {
			continue
		}

		err := Unzip(folderName+"\\"+zip.Name(), appDataDir)

		if err != nil {
			log.Fatal(err)
		}

		fullDir := appDataDir + "\\" + strings.Split(zip.Name(), ".zip")[0]

		files, err := os.ReadDir(fullDir)

		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if !pathExists(mcDir + "\\mods\\" + file.Name()) {
				err := CopyFile(fullDir+"\\"+file.Name(), mcDir+"\\mods\\"+file.Name())
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	os.RemoveAll(appDataDir)

}

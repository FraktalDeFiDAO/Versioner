package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	Versioner "github.com/FraktalDeFiDAO/versioner"
)

func main() {
	versioner := *Versioner.NewVersioner()

	var (
		updateVersion string
		updateType    string
	)
	flag.StringVar(&versioner.(*Versioner.Versioner).BasePath, "base-path", "./", "Set base path")
	flag.StringVar(&updateType, "update", "", "Update to new version...")
	flag.StringVar(&updateVersion, "version", "", "Update to new version [manually set version]...")
	flag.Parse()

	versioner.(*Versioner.Versioner).VersionFile = versioner.(*Versioner.Versioner).BasePath + "version"
	versioner.SetCurrentVersion()
	if _, err := os.Stat(versioner.(*Versioner.Versioner).VersionFile); err != nil {
		if err := os.WriteFile(versioner.(*Versioner.Versioner).VersionFile, []byte("0.0.1"), 0644); err != nil {
			log.Fatal("Error creating version file...")
		}
	}
	currentVersion, err := versioner.GetCurrentVersion()

	if updateType != "" {
		fmt.Printf("Updating version => %+v \t=> \t%+v\n", updateType, currentVersion)
		if err = versioner.UpdateVersion(updateType, updateVersion); err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		log.Fatal("Error:", err)
	}

}

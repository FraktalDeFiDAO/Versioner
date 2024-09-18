package versioner

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type IVersioner interface {
	GetCurrentVersion() (string, error)
	SetCurrentVersion() error
	UpdateVersion(updateType string, updateVersion string) error
	IncMinor(version string) (string, error)
	IncMajor(version string) (string, error)
	IncRelease(version string) (string, error)
	ShowVersion()
	WriteVersionFIle(versionFile, version string) error
}

type Versioner struct {
	Version     string
	BasePath    string
	VersionFile string
}

func NewVersioner() *IVersioner {
	var v IVersioner = &Versioner{}

	return &v
}
func (v *Versioner) GetCurrentVersion() (string, error) {
	var (
		result = ""
		err    error
	)

	if _, err := os.Stat(v.VersionFile); err != nil {
		return result, errors.New("Version file not found...")
	}
	if data, err := os.ReadFile(v.VersionFile); err != nil {
		return result, errors.New("Version file not found")
	} else {
		result = string(data)
	}

	return result, err
}

func (v *Versioner) SetCurrentVersion() error {
	var err error
	v.Version, err = v.GetCurrentVersion()
	if err != nil {
		return err
	}
	return nil
}
func (v *Versioner) UpdateVersion(updateType string, updateVersion string) error {
	var (
		newVersion = ""
		err        error
	)

	if updateType == "" {
		updateType = "minor"
	}

	currVersion, err := v.GetCurrentVersion()

	if err != nil {
		return err
	}
	switch updateType {
	case "minor":
		fmt.Println("IncMinor:")
		newVersion, err = v.IncMinor(currVersion)
		if err != nil {
			return err
		}
		fmt.Println("newVersion =>", newVersion)
		if err = v.WriteVersionFIle(v.VersionFile, newVersion); err != nil {
			return err
		}

		return nil

	case "major":
		fmt.Println("IncMajor:")

		newVersion, err = v.IncMajor(currVersion)
		if err != nil {
			return err
		}

		fmt.Println("newVersion =>", newVersion)
		if err = v.WriteVersionFIle(v.VersionFile, newVersion); err != nil {
			return err
		}

		return nil
	case "release":
		fmt.Println("IncRelease:")

		newVersion, err = v.IncRelease(currVersion)
		if err != nil {
			return err
		}
		fmt.Println("newVersion:", newVersion)
		fmt.Println("newVersion =>", newVersion)
		if err = v.WriteVersionFIle(v.VersionFile, newVersion); err != nil {
			return err
		}

		return nil

	case "manual":
		fmt.Println("Set Manually:")

		fmt.Println("updateVersion =>", updateVersion)
		v.GetCurrentVersion()
		isVersionGreater := v.Cmp(v.Version, updateVersion)
		if isVersionGreater < 1 {
			return errors.New("error: can not update to an equal or lesser version number [ oldVersion => " + v.Version + " ] [ newVersion => " + updateVersion + " ] " + strconv.Itoa(isVersionGreater))
		}
		if err = v.WriteVersionFIle(v.VersionFile, updateVersion); err != nil {
			return err
		}

		return nil

	default:
		return errors.New("Version Update Type Not Found")
	}

	return nil
}

func (v *Versioner) IncMinor(version string) (string, error) {
	newVersion := ""

	versionChunks := strings.Split(version, ".")

	u, err := strconv.Atoi(versionChunks[2])
	if err != nil {
		return newVersion, err
	}
	versionChunks[2] = strconv.Itoa(u + 1)

	newVersion = strings.Join(versionChunks, ".")

	return newVersion, nil
}

func (v *Versioner) IncMajor(version string) (string, error) {
	newVersion := ""

	versionChunks := strings.Split(version, ".")

	u, err := strconv.Atoi(versionChunks[1])
	versionChunks[1] = strconv.Itoa(u + 1)

	if err != nil {
		log.Fatal(err)
	}
	versionChunks[2] = strconv.Itoa(0)

	newVersion = strings.Join(versionChunks, ".")

	return newVersion, nil

}

func (v *Versioner) IncRelease(version string) (string, error) {
	newVersion := ""

	versionChunks := strings.Split(version, ".")

	u, err := strconv.Atoi(versionChunks[0])
	versionChunks[0] = strconv.Itoa(u + 1)

	if err != nil {
		log.Fatal(err)
	}
	versionChunks[1] = strconv.Itoa(0)
	versionChunks[2] = strconv.Itoa(0)

	newVersion = strings.Join(versionChunks, ".")

	return newVersion, nil

}

func (v *Versioner) ShowVersion() {
	file, err := os.ReadFile(v.VersionFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Version =>", string(file))
}

func (v *Versioner) WriteVersionFIle(versionFile, version string) error {
	fmt.Println("Writing version file...", versionFile, version)
	if _, err := os.Stat(versionFile); err != nil {
		return errors.New("version file not found")
	}
	if err := os.WriteFile(versionFile, []byte(version), 0644); err != nil {
		return errors.New("Error writing verison file")
	}
	return nil
}

func (v *Versioner) Cmp(oldVersion, newVersion string) int {

	oldChunks := strings.Split(oldVersion, ".")
	newChunks := strings.Split(newVersion, ".")

	if len(oldChunks) < 3 {
		log.Fatal("invalid version number")
	}
	if len(newChunks) < 3 {
		log.Fatal("invalid version number")
	}
	if oldVersion == newVersion {
		return 0
	}
	chunks := ([3][2]int{
		[2]int{},
		[2]int{},
		[2]int{},
	})
	chunks[0][0], _ = strconv.Atoi(oldChunks[0])
	chunks[0][1], _ = strconv.Atoi(newChunks[0])

	chunks[1][0], _ = strconv.Atoi(oldChunks[1])
	chunks[1][1], _ = strconv.Atoi(newChunks[1])

	chunks[2][0], _ = strconv.Atoi(oldChunks[2])
	chunks[2][1], _ = strconv.Atoi(newChunks[2])

	if chunks[0][0] < chunks[0][1] {
		return 1
	} else if chunks[1][0] < chunks[1][1] {
		return 1
	} else if chunks[2][0] < chunks[2][1] {
		return 1
	}

	return -1

}

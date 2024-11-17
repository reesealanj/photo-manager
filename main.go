package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type BackupBuddyConfig struct {
	Src          string `json:"source"`
	Dest         string `json:"destination"`
	DryRun       bool   `json:"is_dry_run"`
	Debug        bool   `json:"is_debug_mode"`
	CreateDest   bool   `json:"should_create_dest"`
	SkipExisting bool   `json:"should_skip_existing_files"`
}

var globalConfig BackupBuddyConfig = BackupBuddyConfig{
	Src:          "./source",
	Dest:         "./dest",
	DryRun:       false,
	Debug:        false,
	CreateDest:   false,
	SkipExisting: false,
}

func main() {
	args := os.Args[1:]

	log.Print("[info] Parsing Arguments")
	parseArgs(args, &globalConfig)
	log.Print("[info] Arg Parsing Complete")

	configBytes, _ := json.Marshal(globalConfig)
	debugLog(" config loaded: " + string(configBytes))
	// Notify flags as appropriate
	if globalConfig.DryRun {
		log.Print("[dry-run] Running in Dry Run Mode, no changes will be made")
	}

	if globalConfig.Debug {
		debugLog("[debug] Running in Debug mode")
	}

	if globalConfig.CreateDest {
		debugLog("[debug] Destination file will be created if it does not exist!")
	}

	log.Print("[info] Validating Config Directories")
	// Validate Source
	validateSrc := validateDirectory(globalConfig.Src)
	if !validateSrc {
		log.Fatal("[error] Invalid source arg. Exiting...")
	} else {
		debugLog("[init] Source directory valid!")
	}

	// Validate or Create Destination (based on flags)
	validateDest := validateDirectory(globalConfig.Dest)
	if !validateDest && !globalConfig.CreateDest {
		log.Fatal("[error] Invalid destination arg. Exiting...")
	}
	if !validateDest && globalConfig.CreateDest {
		log.Print("[info] Destination folder does not exist and will be created")
	}

	debugLog("[init] Destination directory valid!")
	log.Print("[info] Config Directories Validated")

	var files []string

	log.Print("[info] gathering files")
	fileWalkErr := filepath.Walk(globalConfig.Src, func(path string, itemInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if shouldCheckFile(itemInfo) {
			files = append(files, path)
		}
		return nil
	})

	log.Print("[info] " + fmt.Sprint(len(files)) + " files gathered.")
	if fileWalkErr != nil {
		log.Fatal("[error] Issue reading src directory: " + fileWalkErr.Error())
	}

	log.Print("[info] Preparing destination directories for all files in " + globalConfig.Src)
	for index, file := range files {
		fileData, fileErr := os.Stat(file)
		if fileErr != nil {
			log.Fatal("[error] Unable to stat file " + file)
		}
		year, month, date := parseTime(fileData.ModTime())
		dateString := month + "/" + date + "/" + year

		debugLog("[found-file] File " + file + " created at " + dateString)

		destPath, dirCreateErr := prepareDirectoriesForDate(year, month, date)
		if dirCreateErr != nil {
			log.Fatal("[error][dest-create] Unable to prepare file path " + year + "/" + month + "/" + date + ". Exiting...")
		} else {
			debugLog("[dest-create] File path " + year + "/" + month + "/" + date + " prepared!")
		}
		debugLog("[dest-create] File path " + dateString + " created or validated successfully!")

		destPath = globalConfig.Dest + "/" + destPath + "/" + fileData.Name()
		srcPath := file

		if !globalConfig.DryRun {
			copyFile(srcPath, destPath)
		} else {
			log.Print("[dryRun] cp " + srcPath + " " + destPath)
		}

		log.Print("[info] Processed " + fmt.Sprint(index) + "/" + fmt.Sprint(len(files)) + " (" + fmt.Sprintf("%.2f", (float64(index)/float64(len(files))*100)) + "%)")
	}

	log.Print("[info] Processed " + fmt.Sprint(len(files)) + "/" + fmt.Sprint(len(files)) + " (100.00%)")
	log.Print("[info] Destination directories prepared at " + globalConfig.Dest + " for all files in " + globalConfig.Src)
}

func debugLog(logString string) {
	if globalConfig.Debug {
		log.Print("[debug]" + logString)
	}
}

func shouldCheckFile(info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	if info.Name() == ".DS_Store" {
		return false
	}

	return true
}

func validateDirectory(path string) bool {
	pathData, err := os.Stat(path)

	if os.IsNotExist(err) {
		log.Print("[warn][validateDirectory] Directory " + path + "does not exist.")
		return false
	}

	if !pathData.IsDir() {
		log.Print("[warn][validateDirectory]" + path + " is not a directory.")
		return false
	}

	return true
}

func validateOrCreateDirectory(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		var configPath string = globalConfig.Dest + "/" + path + "/"

		if !globalConfig.DryRun {
			mkErr := os.MkdirAll(configPath, 0755)
			if mkErr != nil {
				log.Print("[error][create-dir] Directory " + configPath + " could not be created.")
				return mkErr
			}
		} else {
			log.Print("[dry-run] MkdirAll " + configPath)
		}
	}

	return nil
}

func parseTime(time time.Time) (year, month, day string) {
	year = time.Format("2006")
	month = time.Format("01")
	day = time.Format("02")
	return
}

func parseArgs(args []string, config *BackupBuddyConfig) {
	quoteReg := regexp.MustCompile("\"")
	for _, arg := range args {
		if strings.Contains(arg, "--source=") {
			reg := regexp.MustCompile(`--source=`)
			parsedArg := reg.ReplaceAllString(arg, "")
			parsedArg = quoteReg.ReplaceAllString(parsedArg, "")

			config.Src = parsedArg
		}

		if strings.Contains(arg, "--dest=") {
			reg := regexp.MustCompile(`--dest=`)
			parsedArg := reg.ReplaceAllString(arg, "")
			parsedArg = quoteReg.ReplaceAllString(parsedArg, "")

			config.Dest = parsedArg
		}

		if strings.Contains(arg, "--dryRun") {
			config.DryRun = true
		}

		if strings.Contains(arg, "--debug") || arg == "-d" {
			config.Debug = true
		}

		if strings.Contains(arg, "--create-dest") || arg == "-c" {
			config.CreateDest = true
		}

		if strings.Contains(arg, "--allow-skip") || strings.Contains(arg, "--skip") {
			config.SkipExisting = true
		}
	}
}

func prepareDirectoriesForDate(year, month, date string) (string, error) {
	baseDir := year
	monthDir := month + "-" + year
	dateDir := month + "-" + date + "-" + year
	fullPath := baseDir + "/" + monthDir + "/" + dateDir

	err := validateOrCreateDirectory(fullPath)
	if err != nil {
		return fullPath, err
	}

	return fullPath, nil
}

func copyFile(src, dst string) {
	fin, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	_, checkErr := os.Stat(dst)

	if !os.IsNotExist(checkErr) && !globalConfig.SkipExisting {
		log.Fatal("[error] Unable to copy " + src + " to " + dst + ". Reason: File Already Exists")
	}

	if !os.IsNotExist(checkErr) && globalConfig.SkipExisting {
		debugLog("[skip] File " + dst + " exists. Copy failed, skipping.")
		return
	}

	fout, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)

	if err != nil {
		log.Fatal(err)
	}
}

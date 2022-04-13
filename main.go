package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	path := flag.String("path", "", "Path to the folder to archive")
	backupFolder := flag.String("backup", "./archive", "Path to the backup folder")
	count := flag.Int("count", 3, "Number of backup to keep in the archive")
	flag.Parse()
	if len(*path) == 0 {
		fmt.Println("missing --path argument")
		os.Exit(2)
	}
	if *count < 1 {
		fmt.Println("--count must be greater than 0")
		os.Exit(3)
	}
	if _, err := os.Stat(*path); err != nil {
		err = os.MkdirAll(*path, 0750)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
	}
	absSourcePath, err := filepath.Abs(*path)
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}
	absDstPath, err := filepath.Abs(*backupFolder)
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}
	if err := makeBackup(absSourcePath, absDstPath, *count); err != nil {
		fmt.Println(err)
		os.Exit(7)
	}
}

// makeBackup make a backup
func makeBackup(source, dst string, count int) error {
	err := cleanOldBackup(dst, count-1)
	if err != nil {
		return err
	}
	name := time.Now().Format("2006_01_02__15_04_05")
	output, err := exec.Command("cp", "-r", source, filepath.Join(dst, name)).CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}

// cleanOldBackup clean up old archives so as not to exceed the quota
func cleanOldBackup(backupFolder string, count int) error {
	files, err := ioutil.ReadDir(backupFolder)
	if err != nil {
		return err
	}
	for len(files) > count {
		oldestTime := time.Now()
		var oldestFile os.FileInfo
		for _, file := range files {

			if file.ModTime().Before(oldestTime) {
				oldestFile = file
				oldestTime = file.ModTime()
			}
		}
		err = os.RemoveAll(filepath.Join(backupFolder, oldestFile.Name()))
		if err != nil {
			return err
		}
		files, err = ioutil.ReadDir(backupFolder)
		if err != nil {
			return err
		}
	}
	return nil
}

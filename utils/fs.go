package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func CopyFile(src string, dest string) {
	srcFile, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		panic(err)
	}

	err = destFile.Sync()
	if err != nil {
		panic(err)
	}
}

func HasFile(glob string) (matches []string, ok bool) {
	matches, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}
	return matches, len(matches) > 0
}

func FindFiles(startDir, targetFile string) (targetFiles []string, err error) {
	err = filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(info.Name(), targetFile) {
			targetFiles = append(targetFiles, path)
		}
		return nil
	})
	return targetFiles, err
}

func FindInFile(filepath string, rgx *regexp.Regexp, excludes *regexp.Regexp, buf *[]byte) (ok bool, matches string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	foundMatch := false
	var matchingLines strings.Builder

	scanner := bufio.NewScanner(file)
	if buf != nil {
		scanner.Buffer(*buf, len(*buf))
	}
	lineNumber := 1
	for scanner.Scan() {
		lineText := scanner.Text()
		matches := rgx.FindStringSubmatch(lineText)
		excl := []string{}
		if excludes != nil {
			excl = excludes.FindStringSubmatch(lineText)
		}
		if len(matches) > 0 && len(excl) == 0 {
			matchingLines.WriteString(fmt.Sprintf("Line %d: %s\n", lineNumber, lineText))
			// fmt.Printf("Line %d: %s\n", lineNumber, lineText)
			foundMatch = true
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		if err != bufio.ErrTooLong {
			log.Fatal(err)
		}
	}

	return foundMatch, matchingLines.String()
}

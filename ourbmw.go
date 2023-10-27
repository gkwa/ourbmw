package ourbmw

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Main(gitRootDir, oldStr, newStr string) int {
	err := doit(gitRootDir, oldStr, newStr)
	if err != nil {
		slog.Error("doit", "error", err)
		return 1
	}

	return 0
}

func replaceInFile(filePath string, oldStr, newStr string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tempFilePath := filePath + ".temp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Replace(scanner.Text(), oldStr, newStr, -1)
		_, err := tempFile.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	err = os.Rename(tempFilePath, filePath)
	if err != nil {
		return err
	}

	return nil
}

func processFiles(gitRootDir, oldStr, newStr string) error {
	return filepath.Walk(gitRootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.Contains(info.Name(), oldStr) {
			newFileName := strings.Replace(info.Name(), oldStr, newStr, -1)
			newFilePath := filepath.Join(filepath.Dir(path), newFileName)
			err := os.Rename(path, newFilePath)
			if err != nil {
				return err
			}
			path = newFilePath
		}
		return replaceInFile(path, oldStr, newStr)
	})
}

func processFiles2(gitRootDir, oldStr, newStr string) error {
	cookiecutterPath := filepath.Join(gitRootDir, "cookiecutter.json")
	if _, err := os.Stat(cookiecutterPath); os.IsNotExist(err) {
		if err := createCookiecutterFile(cookiecutterPath); err != nil {
			slog.Error("error creating cookiecutter.json", "error", err)
			return err
		}
	}

	return nil
}

func doit(gitRootDir, oldStr, newStr string) error {
	err := processFiles(gitRootDir, oldStr, newStr)
	if err != nil {
		return err
	}

	err = processFiles2(gitRootDir, oldStr, newStr)
	if err != nil {
		return err
	}
	return nil
}

func createCookiecutterFile(filePath string) error {
	cookiecutterData := map[string]interface{}{
		"project_name": "Cookiecutter Website Simple",
		"project_slug": "{{ cookiecutter.project_name.lower().replace(' ', '_') }}",
		"author":       "Anonymous",
	}

	data, err := json.MarshalIndent(cookiecutterData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

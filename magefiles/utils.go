package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/sirgwain/craig-stars/test"
)

// Test both the backend and frontend in succession.
func Test() error {
	if err := sh.RunV("go", "test", "./..."); err != nil {
		return err
	}

	cmd := exec.Command("npm", "run-script", "test")
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Download frontend image files, replacing existent ones if present.
func Images() error {

	// switch dir
	originalDir, err := os.Getwd()
	if err != nil {
		return mg.Fatalf(1, "could not get working directory to revert to: \n%w", err)
	}

	if err := os.Chdir("./frontend/static"); err != nil {
		return mg.Fatalf(1, "error during os.Chdir: \n%w", err)
	}

	// revert the current dir after this call
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			panic(fmt.Errorf("error reverting to original directory: \n%v", err))
		}
	}()

	tmpName, err := downloadImagesZip()
	if err != nil {
		return err
	}
	defer func() {
		if err := sh.Rm(tmpName); err != nil {
			panic(err)
		}
	}()

	if err := unzipTempFile(tmpName); err != nil {
		return err
	}

	return nil
}

func downloadImagesZip() (string, error) {
	// create temp file to store zip file from http request
	tmpFile, err := os.CreateTemp("", "images.zip")
	if err != nil {
		return "", mg.Fatalf(1, "error during os.CreateTemp: \n%w", err)
	}
	defer func() {
		// close temp file after we're done
		tmpFile.Close()
	}()

	var tmpName = tmpFile.Name()

	// download the images zip from the web using an http request
	request, err := http.Get("https://craig-stars.net/images/images.zip")
	if err != nil {
		return "", mg.Fatalf(1, "error during http.Get: \n%w", err)
	}
	defer func() {
		// don't forget to close it!
		if err := request.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if request.StatusCode != 200 {
		statusText := http.StatusText(request.StatusCode)
		if statusText == "" {
			statusText = "unknown status code"
		}
		return "", mg.Fatalf(1, "http web request returned status code %d (%s)", request.StatusCode, statusText)
	}

	// Copy the response body to the temp file
	_, err = io.Copy(tmpFile, request.Body)
	if err != nil {
		return "", mg.Fatalf(1, "error during io.Copy: \n%w", err)
	}

	fmt.Printf("downloaded images.zip to %s\n", tmpName)
	return tmpName, nil
}

// unzip the temp file with the given path; used during image download
func unzipTempFile(tmpName string) error {
	// create zip reader to unzip temp file contents
	reader, err := zip.OpenReader(tmpName)
	if err != nil {
		return mg.Fatalf(1, "error during zip.OpenReader: \n%w", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			panic(err)
		}
	}()

	// Delete previous images folder
	if err := sh.Rm("images"); err != nil {
		return err
	}

	// copy all the files one by one to the images folder
	for _, file := range reader.File {
		filePath := file.Name

		// Create directories as needed
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return mg.Fatalf(1, "error during os.MkdirAll: \n%w", err)
			}
			continue
		}

		// Create a file
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return mg.Fatalf(1, "error during os.MkdirAll: \n%w", err)
		}
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return mg.Fatalf(1, "error during os.OpenFile: \n%w", err)
		}

		// Extract the file
		srcFile, err := file.Open()
		if err != nil {
			dstFile.Close()
			return mg.Fatalf(1, "error during zipFile.Open: \n%w", err)
		}
		_, err = io.Copy(dstFile, srcFile)

		// Close the open files
		dstFile.Close()
		srcFile.Close()
		if err != nil {
			return mg.Fatalf(1, "error during io.Copy: \n%w", err)
		}
	}
	fmt.Println("Unzipped images to frontend/static/images")

	return nil
}

// Merge multiple temp json files from tmp folder together into 1, delimiting them by package
func CI_Merge_JSON() error {
	tmp, err := os.Open("./tmp")
	if err != nil {
		return mg.Fatalf(1, "error while opening temp folder: \n%w", err)
	}
	files, err := tmp.ReadDir(-1)
	if err != nil {
		return mg.Fatalf(1, "error while reading temp folder files: \n%w", err)
	}

	if len(files) == 0 {
		fmt.Println("No files in temp folder; exiting")
		return nil
	}

	count := 0
	for _, fileEntry := range files {
		fileName := fileEntry.Name()
		if !strings.HasPrefix(fileName, "diff_") ||
			!strings.HasSuffix(fileName, ".jsonl") {
			// file doesn't start with correct prefix; probably not a json file
			continue
		}

		// extract name of package from file name
		pkgName, _ := strings.CutPrefix(fileName, "diff_")
		pkgName, _ = strings.CutSuffix(fileName, ".jsonl")

		// grab file data
		file, _ := os.Open("./tmp/" + fileName)
		defer file.Close()
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return mg.Fatalf(1, "error during io.ReadAll: \n%w", err)
		}

		// Add a short comment mentioning which package we're in to the start of the file
		header := "// " + strings.ToUpper(pkgName) + "\n"
		fileContents := header + string(fileBytes)
		if err = test.AppendFile("./tmp/diff.jsonl", fileContents); err != nil {
			return mg.Fatalf(1, "error during test.AppendFile: \n%w", err)
		}

		os.Remove(file.Name())
		count++
	}

	fmt.Printf("Successfully merged %d json files into ./tmp/diff.jsonl\n", count)
	return nil
}

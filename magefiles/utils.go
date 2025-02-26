package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Test both the backend and frontend in succession.
func Test() error {
	fmt.Println("go test ./...")
	if err := sh.RunV("go", "test", "./..."); err != nil {
		return err
	}

	fmt.Println("npm run test")
	cmd := exec.Command("npm", "run-script", "test")
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("npm run lint")
	cmd = exec.Command("npm", "run-script", "lint")
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

	if err := unzipTempFile(tmpName); err != nil {
		return err
	}

	if err := sh.Rm(tmpName); err != nil {
		panic(err)
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
		// close and remove temp file after we're done
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

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

// Run all frontend/backend tests and lint checks.
func Test() error {
	if err := Lint(); err != nil {
		return err
	}

	err := Test_Golang("./...")
	if err != nil {
		return err
	}

	if err := Test_Vitest(""); err != nil {
		return err
	}

	return Test_Playwright("")
}

// Run ESLint lint checks on frontend code.
func Lint() error {
	fmt.Println("Running ESLint linting checks...")
	cmd := exec.Command("npm", "run-script", "lint")
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run backend golang tests via gotestsum, passing the passed in args to "go test".
// Gotestsum args are dependent on the value of $CI and $GITHUB_REPOSITORY/$GH_REPO.
func Test_Golang(goTestArgs string) error {
	fmt.Println("Running backend tests...")

	if goTestArgs = strings.TrimSpace(goTestArgs); goTestArgs == "" {
		goTestArgs = "./..."
	}

	// read gotestsum config args from text file
	// use CI config if on CI; else regular config
	var filePath string
	if is_CI() {
		fmt.Println("CI run detected; using CI config")
		filePath = "gotestsum/gotestsum_ci.config.txt"
	} else {
		fmt.Println("Non-CI run detected; using default config")
		filePath = "gotestsum/gotestsum.config.txt"
	}
	configBytes, err := os.ReadFile(filePath)
	if err != nil {
		return mg.Fatalf(1, "error reading gotestsum config file: \n%w", err)
	}

	// extract config values delimited by commas and whitespace
	configVals := strings.FieldsFunc(string(configBytes), func(r rune) bool {
		return (r == ',' || r == ' ' || r == '\n' || r == '\r')
	})
	fmt.Printf("Config file at %s successfully read.\nContents: %s", filePath, strings.Join(configVals, ", "))

	// if $GITHUB_REPOSITORY is set and nonempty, use that as package name for JUnit report.
	// Otherwise, check for $GH_REPO before falling back to a default string.
	var repoName string = "craig-stars"
	if r := strings.TrimSpace(os.Getenv("GITHUB_REPOSITORY")); r != "" {
		repoName = r
	} else if r = strings.TrimSpace(os.Getenv("GH_REPO")); r != "" {
		repoName = r
	}

	// merge any produced json files together once we're done testing
	// we only do this after all the setup to save time
	defer func() {
		if err := Merge_Temp_JSON(); err != nil {
			fmt.Printf("error merging temp JSON diffs after test run:\n%v\n", err)
		}
	}()

	return sh.RunWithV(map[string]string{"GITHUB_REPOSITORY": repoName},
		configVals[0], configVals[1:]...) // "go", "tool", "gotest.tools/gotestsum"...
}

// Remove all temp json files inside tmp and merge them into 1 large file.
// This takes all files matching the format "diff_**.jsonl"
// and merges them together into 1 large file for easy parsing & CI uploading.
// Comments are added between failing tests from different packages.
func Merge_Temp_JSON() error {
	tmp, err := os.Open("tmp")
	if err != nil {
		return mg.Fatalf(1, "error while opening temp folder: \n%w", err)
	}
	fileNames, err := tmp.Readdirnames(-1)
	if err != nil {
		return mg.Fatalf(1, "error while reading temp folder files: \n%w", err)
	}

	if len(fileNames) == 0 {
		fmt.Println("No JSON diffs were found inside tmp to merge; exiting")
		return nil
	}

	count := 0
	for _, fileName := range fileNames {
		if !strings.HasPrefix(fileName, "diff_") ||
			!strings.HasSuffix(fileName, ".jsonl") {
			// file doesn't start with correct prefix; probably not a json file
			continue
		}

		// extract name of package from file name
		pkgName, _ := strings.CutPrefix(fileName, "diff_")
		pkgName, _ = strings.CutSuffix(pkgName, ".jsonl")

		// grab file data
		fileBytes, err := os.ReadFile("tmp/" + fileName)
		if err != nil {
			return mg.Fatalf(1, "error during os.ReadFile: \n%w", err)
		}

		// Add a header mentioning which package we're in to the start of the file
		contents := "//*" +
			strings.ToUpper(pkgName) + "\n" +
			string(fileBytes)
		if count == 0 {
			// truncate file if it already exists
			if err := os.WriteFile("tmp/diff.jsonl", []byte(contents), 0644); err != nil {
				return mg.Fatalf(1, "error during os.WriteFile: \n%w", err)
			}
		} else {
			if err := test.AppendFile("tmp/diff.jsonl", "\n"+contents); err != nil {
				return mg.Fatalf(1, "error during test.AppendFile: \n%w", err)
			}
		}

		count++
		// remove test file after being merged
		if err := sh.Rm(fileName); err != nil {
			return err
		}
	}

	var message string
	if count > 0 {
		message = fmt.Sprintf("Successfully merged %d temp json files into tmp/diff.jsonl.", count)
	} else {
		message = "No JSON files to merge were found."
	}
	fmt.Println(message, "\nHave a nice day.")
	return nil
}

// Run frontend tests using Vitest.
func Test_Vitest(vitestArgs string) error {
	fmt.Println("Running vitest tests...")
	cmd := exec.Command("npm", "run-script", "test:unit", "--", vitestArgs)
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run end-to-end tests using Playwright.
func Test_Playwright(playwrightArgs string) error {
	fmt.Println("Running playwright tests...")
	cmd := exec.Command("npm", "run-script", "test:e2e", "--", playwrightArgs)
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
		fmt.Printf("removed temp file at %s", tmpName)
	}()

	if err := unzipTempFile(tmpName); err != nil {
		return err
	}

	return nil
}

// Download and store images zip to a temp file
func downloadImagesZip() (tmpFileName string, err error) {
	// create temp file to store zip file from http request
	tmpFile, err := os.CreateTemp("", "images_*.zip")
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
			statusText = "unknown"
		}
		return "", mg.Fatalf(1, "http web request returned status code %d (%s)", request.StatusCode, statusText)
	}

	// Copy the response body to the temp file
	_, err = io.Copy(tmpFile, request.Body)
	if err != nil {
		return "", mg.Fatalf(1, "error during io.Copy: \n%w", err)
	}

	fmt.Println("downloaded images.zip to", tmpName)
	return tmpName, nil
}

// unzip the temp file at the given path
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
	fmt.Println("unzipped images to frontend/static/images")

	return nil
}

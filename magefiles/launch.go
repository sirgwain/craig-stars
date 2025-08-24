package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// is_CI reports whether the current process is running in CI (continuous integration)
// by checking the "CI" environment variable.
func is_CI() bool {
	return os.Getenv("CI") != ""
}

// Build and launch the server for local development.
// This calls both Build and Launch consecutively.
func Run() error {
	Build()
	return Launch()
}

// Build the frontend and backend consecutively, alongside some setup work.
func Build() {
	// Returns no errors since mg.Deps panics
	mg.SerialDeps(Clean,
		Tidy,
		Copy_Wasm_Exec,
		Generate,
		Build_Frontend,
		Build_Backend,
	)
}

// Clean up various temporary directories.
// This runs "go clean" and removes everything in dist and frontend/build.
func Clean() error {
	if err := sh.RunV("go", "clean"); err != nil {
		return err
	}
	if err := sh.Rm("dist"); err != nil {
		return err
	}

	return sh.Rm("frontend/build")
}

// Copy wasm executable from GOROOT to frontend folder.
// This copies the "wasm_exec.js" file from GOROOT/lib/wasm into
// frontend/src/lib/wasm, creating the folder if not already present.
func Copy_Wasm_Exec() error {
	if err := os.MkdirAll("frontend/src/lib/wasm", 0755); err != nil {
		return mg.Fatalf(1, "error during os.MkdirAll: \n%w", err)
	}

	// Find GOROOT
	goroot, err := sh.Output("go", "env", "GOROOT")
	if err != nil {
		return err
	}
	goroot = strings.ReplaceAll(goroot, "\\", "/") // replace backslashes on windows

	// Check if wasm executable exists or not.
	// Go 1.24 moved wasm_exec.js from misc/wasm to lib/wasm,
	// but we require go 1.24 anyways to run our tool deps so it shouldn't matter.
	if _, err := os.Stat(goroot + "/lib/wasm/wasm_exec.js"); errors.Is(err, os.ErrNotExist) {
		// file doesn't exist
		return mg.Fatalf(1, "executable was not found inside GOROOT: %v", goroot)
	} else if err != nil {
		// some other random error
		return mg.Fatalf(1, "error during os.Stat(): \n%w", err)
	}

	// file exists
	path := goroot + "/lib/wasm/wasm_exec.js"
	if err := sh.Copy("frontend/src/lib/wasm/wasm_exec.js", path); err != nil {
		return mg.Fatalf(1, "error while copying wasm exec: \n%w", err)
	}
	return nil
}

// Copy wasm_exec.js from tinygo to the frontend wasm folder
// This copies the "wasm_exec.js" file from GOROOT/lib/wasm into
// frontend/src/lib/wasm, creating the folder if not already present.
// cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js
func Copy_Wasm_Exec_TinyGo() error {
	if err := os.MkdirAll("frontend/src/lib/wasm", 0755); err != nil {
		return mg.Fatalf(1, "error during os.MkdirAll: \n%w", err)
	}

	// Find TINYGOROOT
	goroot, err := sh.Output("tinygo", "env", "TINYGOROOT")
	if err != nil {
		return err
	}
	goroot = strings.ReplaceAll(goroot, "\\", "/") // replace backslashes on windows

	// Check if wasm executable exists or not.
	// Go 1.24 moved wasm_exec.js from misc/wasm to lib/wasm,
	// but we require go 1.24 anyways to run our tool deps so it shouldn't matter.
	if _, err := os.Stat(goroot + "/targets/wasm_exec.js"); errors.Is(err, os.ErrNotExist) {
		// file doesn't exist
		return mg.Fatalf(1, "executable was not found inside TINYGOROOT: %v", goroot)
	} else if err != nil {
		// some other random error
		return mg.Fatalf(1, "error during os.Stat(): \n%w", err)
	}

	// file exists
	path := goroot + "/targets/wasm_exec.js"
	if err := sh.Copy("frontend/src/lib/wasm/wasm_exec.js", path); err != nil {
		return mg.Fatalf(1, "error while copying wasm exec: \n%w", err)
	}
	return nil
}

// Tidy up go.mod (equivalent to "go mod tidy -v")
func Tidy() error {
	return sh.RunV("go", "mod", "tidy", "-v")
}

// Generate go code and techs.JSON files.
func Generate() error {
	fmt.Println("running sqlc generate ./...")
	if err := sh.RunV("go", "tool", "sqlc", "generate"); err != nil {
		return err
	}

	fmt.Println("running buf gen ./...")
	if err := sh.RunV("go", "tool", "buf", "generate"); err != nil {
		return err
	}

	fmt.Println("running go generate ./...")
	if err := sh.RunV("go", "generate", "./..."); err != nil {
		return err
	}

	fmt.Println("running npm generate ./...")
	cmd := exec.Command("npm", "run", "generate")
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return mg.Fatalf(1, "error during npm run generate: %w", err)
	}

	fmt.Println("generating techs.json")
	techs2json, err := sh.Output("go", "run", "main.go", "generate", "techsjson")
	if err != nil {
		return err
	}
	if err := os.WriteFile("frontend/src/lib/ssr/techs.json", []byte(techs2json), 0644); err != nil {
		return mg.Fatalf(1, "error during os.WriteFile for techs.json: \n%w", err)
	}

	fmt.Println("generating rules.json")
	rules2json, err := sh.Output("go", "run", "main.go", "generate", "rulesjson")
	if err != nil {
		return err
	}
	if err := os.WriteFile("frontend/src/lib/ssr/rules.json", []byte(rules2json), 0644); err != nil {
		return mg.Fatalf(1, "error during os.WriteFile for rules.json: \n%w", err)
	}

	if !is_CI() {
		fmt.Println("formatting generated json")
		cmd := exec.Command("npx", "prettier", "--ignore-unknown", "--write", "src/lib/ssr")
		cmd.Dir = "./frontend"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return mg.Fatalf(1, "error during npm run generate: %w", err)
		}
	}

	return nil
}

// Build the frontend using SvelteKit.
func Build_Frontend() error {
	mg.Deps(Generate)

	cmd := exec.Command("npm", "install")
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("npm", "run-script", "build")
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Build the backend Golang executable for local dev, as well as the WASM binary.
// This builds the binary for main.go without any version control info;
// Air runs this whenever changes are detected in backend files.
func Build_Backend() error {
	return build_backend(ldflags, "-buildvcs=false")
}

// Variant of Build_Backend used during release containing embedded version control info.
// This takes arguments for the version number, commit hash and build time and passes them
// to go build's ldflags.
func Build_Backend_CI(version, hash, releaseTime string) error {
	mg.SerialDeps(Tidy, Generate, Build_WASM)

	// Benchmarks might suggest otherwise, but these string literals get concatenated during compile time
	args := ldflags + fmt.Sprintf(" -X 'github.com/sirgwain/craig-stars/cmd.semver=%s'"+
		" -X 'github.com/sirgwain/craig-stars/cmd.commit=%s'"+
		" -X 'github.com/sirgwain/craig-stars/cmd.buildTime=%s'", version, hash, releaseTime)
	return build_backend(args)
}

// Internal implementation for building backend with custom go build args
func build_backend(buildArgs ...string) error {
	if err := os.MkdirAll("dist", 0755); err != nil { // re-create dist if folder no exist
		return mg.Fatalf(1, "error during os.MkdirAll: \n%w", err)
	}

	// wrap buildArgs with extra stuff before & after command
	args := make([]string, 1, len(buildArgs)+4)
	args[0] = "build"
	args = append(args, buildArgs...)
	args = append(args, "-o", "dist/"+binary_name, "main.go")

	println("building backend")
	if err := sh.RunV("go", args...); err != nil {
		return err
	}

	// use tinygo for real builds
	if is_CI() {
		mg.Deps(Build_WASM_TinyGo)
	} else {
		mg.Deps(Build_WASM)
	}
	return nil
}

// Build Web-Assembly binary into frontend.
func Build_WASM() error {
	var err error
	if err = os.MkdirAll("frontend/src/lib/wasm", 0755); err != nil {
		return mg.Fatalf(1, "error during os.MkdirAll: \n%w", err)
	}
	println("building wasm")
	if is_CI() {
		err = sh.RunWithV(map[string]string{"GOOS": "js", "GOARCH": "wasm"},
			"go", "build", "-o", "frontend/src/lib/wasm/cs.wasm", "-ldflags", "-s -w", "wasm/main.go")
	} else {
		err = sh.RunWithV(map[string]string{"GOOS": "js", "GOARCH": "wasm"},
			"go", "build", "-o", "frontend/src/lib/wasm/cs.wasm", "wasm/main.go")
	}
	if err != nil {
		return err
	}
	return Copy_Wasm_Exec()
}

// Build tinygo Web-Assembly binary into frontend.
func Build_WASM_TinyGo() error {
	var err error
	if err = os.MkdirAll("frontend/src/lib/wasm", 0755); err != nil {
		return mg.Fatalf(1, "error during os.MkdirAll: \n%w", err)
	}
	println("building tinygo wasm")
	if is_CI() {
		err = sh.RunWithV(map[string]string{"GOOS": "js", "GOARCH": "wasm"},
			"tinygo", "build", "-o", "frontend/src/lib/wasm/cs.wasm", "-no-debug", "wasm/main.go")
	} else {
		err = sh.RunWithV(map[string]string{"GOOS": "js", "GOARCH": "wasm"},
			"tinygo", "build", "-o", "frontend/src/lib/wasm/cs.wasm", "wasm/main.go")
	}
	if err != nil {
		return err
	}

	return Copy_Wasm_Exec_TinyGo()
}

// Launch both backend and frontend servers simultaneously.
func Launch() error {
	// use a WaitGroup to wait until a single goroutine finishes
	wg := sync.WaitGroup{}
	wg.Add(1)

	// run both commands in separate goroutines, using a channel to recieve any errors
	var c chan error
	go func() {
		err := Launch_Backend()
		c <- err
		wg.Done()
		close(c)
	}()
	go func() {
		err := Launch_Frontend()
		c <- err
		wg.Done()
		close(c)
	}()

	// Block until either goroutine finishes and then return the error
	wg.Wait()
	return <-c
}

// Launch the backend go server using air for hot reloads.
func Launch_Backend() error {
	return sh.RunV("go", "tool", "github.com/air-verse/air")
}

// Launch the frontend svelte server.
func Launch_Frontend() error {
	cmd := exec.Command("npm", "run-script", "dev")
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

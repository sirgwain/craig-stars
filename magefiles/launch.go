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

var Aliases = map[string]interface{}{
	"dev":          Launch,
	"dev_frontend": Launch_Frontend,
	"dev_backend":  Launch_Backend,
	"copy_wasm":    Copy_Wasm_Exec,
}

// Build and launch the server for local development.
// This calls both Build and Launch consecutively.
func Run() error {
	if err := Build(); err != nil {
		return err
	}
	return Launch()

}

// Build the frontend and backend consecutively, alongside some setup work.
func Build() error {
	mg.Deps(Clean)
	mg.Deps(Copy_Wasm_Exec)
	mg.Deps(Tidy)
	mg.Deps(Generate)
	mg.Deps(Build_Frontend)
	mg.Deps(Build_Backend)

	return nil
}

// Clean up various temporary directories.
// This runs go clean and removes everything in dist, tmp and frontend/build.
func Clean() error {
	if err := sh.RunV("go", "clean"); err != nil {
		return err
	}
	if err := sh.Rm("dist"); err != nil {
		return err
	}
	if err := sh.Rm("tmp"); err != nil {
		return err
	}
	if err := os.MkdirAll("tmp", 0755); err != nil {
		return mg.Fatalf(1, "error re-creating tmp dir: \n%w", err)
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

	// check if wasm executable exists or not.
	// Go 1.24 moved wasm_exec.js from misc/wasm to lib/wasm,
	// but we require go 1.24 anyways so it shouldn't matter.
	if _, err := os.Stat(goroot + "/lib/wasm/wasm_exec.js"); errors.Is(err, os.ErrNotExist) {
		// file doesn't exist
		return mg.Fatalf(1, "executable was not found inside GOROOT %v", goroot)
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

// Tidy up go.mod (equivalent to "go mod tidy -v")
func Tidy() error {
	return sh.RunV("go", "mod", "tidy", "-v")
}

// Generate go code and techs.JSON files.
func Generate() error {
	fmt.Println("running go generate ./...")
	if err := sh.RunV("go", "generate", "./..."); err != nil {
		return err
	}

	fmt.Println("running tygo generate")
	if err := sh.RunV("go", "tool", "github.com/gzuidhof/tygo", "generate"); err != nil {
		return err
	}

	// format generated tygo file on non-CI runs
	if _, ok := os.LookupEnv("CI"); !ok {
		fmt.Println("running prettier on tygo generated file")
		cmd := exec.Command("npx", "prettier", "--write", "./src/lib/types/cs.ts")
		cmd.Dir = "./frontend"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return mg.Fatalf(1, "error during prettier formatting after generation: \n%w", err)
		}
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
// This builds the binary for main.go without any version control info.
// Air runs this whenever changes are detected.
func Build_Backend() error {
	return build_backend(ldflags, "-buildvcs=false")
}

// Variant of Build_Backend used during release containing embedded version control info.
// This takes arguments for the version number, commit hash and build time and passes them
// to go build's ldflags if not empty.
func Build_Backend_CI(version, hash, releaseTime string) error {
	mg.Deps(Tidy)
	mg.Deps(Generate)
	mg.Deps(Build_WASM)
	// Go passes these arguments directly to build without any quoting or escaping (hence why no surrounding quotes)
	args := ldflags
	// TODO: Change these strings if/when mage updates to support default arguments
	if version != "" {
		args += fmt.Sprintf(" -X 'github.com/sirgwain/craig-stars/cmd.semver=%s'", version)
	}
	if hash != "" {
		args += fmt.Sprintf(" -X 'github.com/sirgwain/craig-stars/cmd.commit=%s'", hash)
	}
	if releaseTime != "" {
		args += fmt.Sprintf(" -X 'github.com/sirgwain/craig-stars/cmd.buildTime=%s'", releaseTime)
	}
	return build_backend(args)
}

// Internal implementation for building backend with custom go build args
func build_backend(buildArgs ...string) error {
	if err := os.MkdirAll("dist", 0755); err != nil { // MkdirAll used due to no-oping if folder already exists
		return mg.Fatalf(1, "error during os.MkdirAll: \n%w", err)
	}

	flags := append(append([]string{"build"}, buildArgs...), "-o",
		fmt.Sprintf("dist/%s", binary_name), "main.go")
	if err := sh.RunV("go", flags...); err != nil {
		return err
	}

	mg.Deps(Build_WASM)
	return nil
}

// Build Web-Assembly binary into frontend.
func Build_WASM() error {
	if err := os.MkdirAll("frontend/src/lib/wasm", 0755); err != nil {
		return mg.Fatalf(1, "error during os.MkdirAll: \n%w", err)
	}
	return sh.RunWithV(map[string]string{"GOOS": "js", "GOARCH": "wasm"},
		"go", "build", "-o", "frontend/src/lib/wasm/cs.wasm", "wasm/main.go")
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

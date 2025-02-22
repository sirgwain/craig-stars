package main

import (
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
// This removes everything in dist and frontend/build.
func Clean() error {
	if err := sh.RunV("go", "clean"); err != nil {
		return err
	}
	if err := sh.Rm("dist"); err != nil {
		return err
	}
	if err := sh.Rm("frontend/build"); err != nil {
		return err
	}
	return nil
}

// Copy wasm executable from GOROOT to frontend folder.
// This copies the "wasm_exec.js" file from your GOROOT into
// frontend/src/lib/wasm, creating the folder if not already present.
func Copy_Wasm_Exec() error {
	if err := os.MkdirAll("frontend/src/lib/wasm", 0755); err != nil {
		return mg.Fatalf(1, "error during os.MkdirAll: \n%w", err)
	}

	goRoot, err := sh.Output("go", "env", "GOROOT")
	if err != nil {
		return mg.Fatalf(1, "error finding GOROOT: \n%w", err)
	}

	if err := sh.Copy("frontend/src/lib/wasm/wasm_exec.js",
		strings.ReplaceAll(goRoot, "\\", "/")+
			"/lib/wasm/wasm_exec.js"); err != nil {
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
	// If/when mage supports default arguments, these should probably be changed to account for it
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

// Launch the backend go server using air.
func Launch_Backend() error {
	return sh.RunV("air")
}

// Launch the frontend svelte server.
func Launch_Frontend() error {
	mg.Deps(Copy_Wasm_Exec)

	cmd := exec.Command("npm", "run-script", "dev")
	cmd.Dir = "./frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return sh.RunV("npm", "run-script", "dev")
}

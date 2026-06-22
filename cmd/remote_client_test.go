package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
)

func TestResolveCLIClientConfigPrecedence(t *testing.T) {
	oldServer, oldToken, oldConfig := cliServer, cliToken, cliConfigPath
	defer func() {
		cliServer, cliToken, cliConfigPath = oldServer, oldToken, oldConfig
	}()

	t.Setenv(cliEnvServer, "https://env.example")
	t.Setenv(cliEnvToken, "env-token")

	path := filepath.Join(t.TempDir(), "cli.json")
	if err := writeCLIConfig(path, cliConfig{
		Server:      "https://file.example",
		AccessToken: "file-token",
	}); err != nil {
		t.Fatal(err)
	}

	cliConfigPath = path
	cliServer = "https://flag.example"
	cliToken = "flag-token"

	cfg, err := resolveCLIClientConfig(true)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server != "https://flag.example" {
		t.Fatalf("server = %q", cfg.Server)
	}
	if cfg.AccessToken != "flag-token" {
		t.Fatalf("token = %q", cfg.AccessToken)
	}

	cliServer, cliToken = "", ""
	cfg, err = resolveCLIClientConfig(true)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server != "https://env.example" || cfg.AccessToken != "env-token" {
		t.Fatalf("env config = %#v", cfg)
	}

	os.Unsetenv(cliEnvServer)
	os.Unsetenv(cliEnvToken)
	cfg, err = resolveCLIClientConfig(true)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server != "https://file.example" || cfg.AccessToken != "file-token" {
		t.Fatalf("file config = %#v", cfg)
	}
}

func TestResolveCLIClientConfigMultipleServers(t *testing.T) {
	oldServer, oldToken, oldConfig := cliServer, cliToken, cliConfigPath
	defer func() {
		cliServer, cliToken, cliConfigPath = oldServer, oldToken, oldConfig
	}()

	path := filepath.Join(t.TempDir(), "cli.json")
	if err := writeCLIConfig(path, cliConfig{
		DefaultServer: defaultCLIServer,
		Servers: map[string]cliServerConfig{
			"craig-stars.net": {
				Server:      defaultCLIServer,
				AccessToken: "prod-token",
			},
			"localhost:8080": {
				Server:      "http://localhost:8080",
				AccessToken: "local-token",
			},
		},
	}); err != nil {
		t.Fatal(err)
	}

	cliConfigPath = path
	cliServer, cliToken = "", ""
	t.Setenv(cliEnvServer, "")
	t.Setenv(cliEnvToken, "")

	cfg, err := resolveCLIClientConfig(true)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server != defaultCLIServer || cfg.AccessToken != "prod-token" {
		t.Fatalf("default config = %#v", cfg)
	}

	cliServer = "http://localhost:8080"
	cfg, err = resolveCLIClientConfig(true)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server != "http://localhost:8080" || cfg.AccessToken != "local-token" {
		t.Fatalf("local config = %#v", cfg)
	}
}

func TestSaveServerTokenPreservesExistingServers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cli.json")
	if err := writeCLIConfig(path, cliConfig{
		DefaultServer: defaultCLIServer,
		Servers: map[string]cliServerConfig{
			"craig-stars.net": {
				Server:      defaultCLIServer,
				AccessToken: "prod-token",
			},
		},
	}); err != nil {
		t.Fatal(err)
	}

	err := saveServerToken(path, "http://localhost:8080", tokenResponse{
		AccessToken: "local-token",
		Scope:       "api:player",
	}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := readCLIConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Servers["craig-stars.net"].AccessToken != "prod-token" {
		t.Fatalf("prod token was not preserved: %#v", cfg.Servers)
	}
	if cfg.Servers["localhost:8080"].AccessToken != "local-token" {
		t.Fatalf("local token was not saved: %#v", cfg.Servers)
	}
}

func TestResolveCLIClientConfigDefaultServer(t *testing.T) {
	oldServer, oldToken, oldConfig := cliServer, cliToken, cliConfigPath
	defer func() {
		cliServer, cliToken, cliConfigPath = oldServer, oldToken, oldConfig
	}()

	cliServer, cliToken = "", ""
	cliConfigPath = filepath.Join(t.TempDir(), "missing.json")
	t.Setenv(cliEnvServer, "")
	t.Setenv(cliEnvToken, "")

	cfg, err := resolveCLIClientConfig(false)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server != defaultCLIServer {
		t.Fatalf("default server = %q", cfg.Server)
	}
}

func TestConnectBaseURL(t *testing.T) {
	if got := connectBaseURL("craig-stars.net/"); got != "https://craig-stars.net/api/grpc" {
		t.Fatalf("connectBaseURL = %q", got)
	}
	if got := connectBaseURL("http://localhost:8080"); got != "http://localhost:8080/api/grpc" {
		t.Fatalf("connectBaseURL = %q", got)
	}
}

func TestMakePKCEPair(t *testing.T) {
	verifier, challenge, err := makePKCEPair()
	if err != nil {
		t.Fatal(err)
	}
	if verifier == "" || challenge == "" || verifier == challenge {
		t.Fatalf("bad PKCE pair verifier=%q challenge=%q", verifier, challenge)
	}
}

func TestFindRaceByName(t *testing.T) {
	races := []*craig_starsv1.Race{{Name: "Humanoid"}}
	if got := findRaceByName(races, "humanoid"); got == nil || got.GetName() != "Humanoid" {
		t.Fatalf("findRaceByName = %#v", got)
	}
	if got := findRaceByName(races, "Robots"); got != nil {
		t.Fatalf("findRaceByName unexpected = %#v", got)
	}
}

func TestPrintGamesTable(t *testing.T) {
	var out bytes.Buffer
	err := printGamesTable(&out, []*craig_starsv1.GameWithPlayers{{
		Game:    &craig_starsv1.Game{Id: 7, Name: "Tiny", Year: 2400, OpenPlayerSlots: 1, Public: true},
		Players: []*craig_starsv1.PlayerStatus{{Name: "A"}, {Name: "B"}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	text := out.String()
	for _, want := range []string{"ID", "Tiny", "2.4K", "Yes"} {
		if !strings.Contains(text, want) {
			t.Fatalf("table missing %q: %s", want, text)
		}
	}
}

func TestPrintRacesTable(t *testing.T) {
	var out bytes.Buffer
	err := printRacesTable(&out, []*craig_starsv1.Race{{
		Id:            3,
		Name:          "Humanoid",
		PluralName:    "Humanoids",
		GrowthRate:    15,
		FactoryOutput: 10,
		FactoryCost:   10,
		NumFactories:  10,
		MineOutput:    10,
		MineCost:      5,
		NumMines:      10,
	}})
	if err != nil {
		t.Fatal(err)
	}
	text := out.String()
	for _, want := range []string{"Humanoid", "Humanoids", "15", "10/10/10"} {
		if !strings.Contains(text, want) {
			t.Fatalf("table missing %q: %s", want, text)
		}
	}
}

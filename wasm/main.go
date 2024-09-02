//go:build wasi || wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/wasm/wasm"
)

// wasm wrapper for calculating race points
func calculateRacePoints(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "ERROR: number of arguments doesn't match"
	}

	raceJson := args[0].String()
	race := &cs.Race{}
	json.Unmarshal([]byte(raceJson), &race)

	rules := cs.NewRules()
	return race.ComputeRacePoints(rules.RaceStartingPoints)
}

func main() {
	wasm.ExposeFunction("calculateRacePoints", js.FuncOf(calculateRacePoints))
	wasm.Ready()

	// fmt.Println("wasm initialized")
	<-make(chan bool) // To use anything from Go WASM, the program may not exit.
}

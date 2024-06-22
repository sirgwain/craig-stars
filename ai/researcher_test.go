package ai

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
)

func Test_aiPlayer_research(t *testing.T) {

	tests := []struct {
		name      string
		techLevel cs.TechLevel
		want      cs.TechField
	}{
		{"No tech", cs.TechLevel{}, cs.Propulsion},
		{"One prop", cs.TechLevel{Propulsion: 1}, cs.Propulsion},
		{"After 2 prop", cs.TechLevel{Propulsion: 2}, cs.Biotechnology},
		{"After 2 prop, 1 bio", cs.TechLevel{Propulsion: 2, Biotechnology: 1}, cs.Energy},
		{"After 2 prop, 1 bio, 1 energy", cs.TechLevel{Energy: 1, Propulsion: 2, Biotechnology: 1}, cs.Weapons},
		{"After 2 prop, 1 bio, 1 energy, 1 weapons", cs.TechLevel{Energy: 1, Weapons: 1, Propulsion: 2, Biotechnology: 1}, cs.Construction},
		{"After 2 prop, 1 bio, 1 energy, 1 weapons 1 con", cs.TechLevel{Energy: 1, Weapons: 1, Propulsion: 2, Construction: 1, Biotechnology: 1}, cs.Construction},
		{"After 2 prop, 1 bio, 1 energy, 1 weapons 2 con", cs.TechLevel{Energy: 1, Weapons: 1, Propulsion: 2, Construction: 2, Biotechnology: 1}, cs.Construction},
		{"After 2 prop, 1 bio, 1 energy, 1 weapons 4 con", cs.TechLevel{Energy: 1, Weapons: 1, Propulsion: 2, Construction: 4, Biotechnology: 1}, cs.Electronics},
		{"After plan done", cs.TechLevel{Energy: 18, Weapons: 24, Propulsion: 16, Construction: 16, Electronics: 19, Biotechnology: 7}, cs.Energy},
		{"Almost max", cs.TechLevel{Energy: 26, Weapons: 26, Propulsion: 26, Construction: 26, Electronics: 25, Biotechnology: 25}, cs.Electronics},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ai := testAIPlayer()
			ai.TechLevels = tt.techLevel

			// trigger research
			ai.research()

			if ai.Researching != tt.want {
				t.Errorf("research() = %v, want %v", ai.Researching, tt.want)
			}
		})
	}
}

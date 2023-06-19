package cs

import "testing"

func TestSalvage_decay(t *testing.T) {
	tests := []struct {
		name  string
		Cargo Cargo
		want  Cargo
	}{
		{"no op, no cargo", Cargo{}, Cargo{}},
		{"decay each 10%", Cargo{Ironium: 100, Boranium: 100, Germanium: 100}, Cargo{Ironium: 90, Boranium: 90, Germanium: 90}},
		{"decay to 0", Cargo{Ironium: 2, Boranium: 2, Germanium: 2}, Cargo{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			salvage := &Salvage{
				Cargo: tt.Cargo,
			}
			salvage.decay(&rules)
			if salvage.Cargo != tt.want {
				t.Errorf("Salvage.decay() = %v, want %v", salvage.Cargo, tt.want)
			}

		})
	}
}

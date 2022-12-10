package cs

type Bomb struct {
	Quantity             int     `json:"quantity,omitempty"`
	KillRate             float64 `json:"killRate,omitempty"`
	MinKillRate          int     `json:"minKillRate,omitempty"`
	StructureDestroyRate float64 `json:"structureDestroyRate,omitempty"`
	UnterraformRate      int     `json:"unterraformRate,omitempty"`
}

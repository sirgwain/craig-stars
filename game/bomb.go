package game

type Bomb struct {
	Quantity             int     `json:"quantity,omitempty"`
	KillRate             float64 `json:"kill_rate,omitempty"`
	MinKillRate          int     `json:"min_kill_rate,omitempty"`
	StructureDestroyRate float64 `json:"structure_destroy_rate,omitempty"`
	UnterraformRate      int     `json:"unterraform_rate,omitempty"`
}

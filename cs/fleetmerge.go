package cs

// A SplitFleetOrder is for splitting tokens off of single fleet to form a second new fleet
// TODO: not yet implemented
type SplitFleetOrder struct {
	Source      *Fleet      `json:"source,omitempty"`
	SplitTokens []ShipToken `json:"splitTokens,omitempty"`
}

// A MergeFleetOrder is for moving tokens between two fleets
// TODO: not yet implemented
type MergeFleetOrder struct {
	Source       *Fleet      `json:"source,omitempty"`
	Dest         *Fleet      `json:"dest,omitempty"`
	SourceTokens []ShipToken `json:"splitTokens,omitempty"`
	DestTokens   []ShipToken `json:"destTokens,omitempty"`
}

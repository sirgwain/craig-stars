package cs

type SplitFleetOrder struct {
	// The tokens to split out of this fleet into a new one
	Source      *Fleet
	SplitTokens []ShipToken `json:"splitTokens,omitempty"`
}

type MergeFleetOrder struct {
	// The tokens to split out of this fleet into a new one
	SplitTokens []ShipToken `json:"splitTokens,omitempty"`
}

package cs

type Wormhole struct {
	MapObject
	DestinationNum   int               `json:"destinationNum,omitempty"`
	Stability        WormholeStability `json:"stability,omitempty"`
	YearsAtStability int               `json:"yearsAtStability,omitempty"`
}

type WormholeStats struct {
	YearsToDegrade int     `json:"yearsToDegrade"`
	ChanceToJump   float64 `json:"chanceToJump"`
	JiggleDistance int     `json:"jiggleDistance"`
}

type WormholeStability string

const (
	WormholeStabilityRockSolid         WormholeStability = "RockSolid"
	WormholeStabilityStable            WormholeStability = "Stable"
	WormholeStabilityMostlyStable      WormholeStability = "MostlyStable"
	WormholeStabilityAverage           WormholeStability = "Average"
	WormholeStabilitySlightlyVolatile  WormholeStability = "SlightlyVolatile"
	WormholeStabilityVolatile          WormholeStability = "Volatile"
	WormholeStabilityExtremelyVolatile WormholeStability = "ExtremelyVolatile"
)

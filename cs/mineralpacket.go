package cs

type MineralPacket struct {
	MapObject
	TargetPlanetNum   uint    `json:"targetPlanetNum,omitempty"`
	Cargo             Cargo   `json:"cargo,omitempty"`
	SafeWarpSpeed     int     `json:"safeWarpSpeed,omitempty"`
	WarpFactor        int     `json:"warpFactor,omitempty"`
	DistanceTravelled float64 `json:"distanceTravelled,omitempty"`
	Heading           Vector  `json:"position"`
}

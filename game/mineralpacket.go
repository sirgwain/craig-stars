package game

type MineralPacket struct {
	MapObject
	TargetPlanetNum   uint    `json:"targetPlanetNum,omitempty"`
	Cargo             Cargo   `json:"cargo,omitempty" gorm:"embedded;embeddedPrefix:cargo_"`
	SafeWarpSpeed     int     `json:"safeWarpSpeed,omitempty"`
	WarpFactor        int     `json:"warpFactor,omitempty"`
	DistanceTravelled float64 `json:"distanceTravelled,omitempty"`
	Heading           Vector  `json:"position" gorm:"embedded;embeddedPrefix:heading_"`
}

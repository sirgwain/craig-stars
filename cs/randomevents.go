package cs

type RandomCometSize int

const (
	RandomCometSizeSmall RandomCometSize = iota
	RandomCometSizeMedium
	RandomCometSizeLarge
	RandomCometSizeHuge
)

type RandomEventType int

const (
	RandomEventTypeComet RandomEventType = iota
	RandomEventTypeMineralDeposit
	RandomEventTypePlanetaryChange
	RandomEventTypeAncientArtifact
	RandomEventTypeMysteryTrader
)

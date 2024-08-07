package cs

// The mystery trader travels through space and gives a boon to any player that gives it a fleet
// full of minerals
// TODO: not yet implemented
type MysteryTrader struct {
	MapObject
	Heading   Vector            `json:"heading,omitempty"`
	WarpSpeed int               `json:"warpSpeed,omitempty"`
	Spec      MysteryTraderSpec `json:"spec,omitempty"`
}

type MysteryTraderSpec struct {
}

// create a new mysterytrader object
func newMysteryTrader(position Vector, num int) *MysteryTrader {
	return &MysteryTrader{
		MapObject: MapObject{
			Type:     MapObjectTypeMysteryTrader,
			Position: position,
			Num:      num,
			Dirty:    true,
		},
	}
}

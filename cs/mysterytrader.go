package cs

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

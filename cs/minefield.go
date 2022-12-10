package cs

import "math"

type MineFieldType string

const (
	MineFieldTypeStandard  MineFieldType = "Standard"
	MineFieldTypeHeavy     MineFieldType = "Heavy"
	MineFieldTypeSpeedBump MineFieldType = "SpeedBump"
)

type MineFieldStats struct {
	MinDamagePerFleetRS int     `json:"minDamagePerFleetRS"`
	DamagePerEngineRS   int     `json:"damagePerEngineRS"`
	MaxSpeed            int     `json:"maxSpeed"`
	ChanceOfHit         float64 `json:"chanceOfHit"`
	MinDamagePerFleet   int     `json:"minDamagePerFleet"`
	DamagePerEngine     int     `json:"damagePerEngine"`
	SweepFactor         float64 `json:"sweepFactor"`
	MinDecay            int     `json:"minDecay"`
	CanDetonate         bool    `json:"canDetonate"`
}

type MineField struct {
	MapObject
	MineFieldOrders
	Type     MineFieldType `json:"type,omitempty"`
	NumMines int           `json:"numMines,omitempty"`
}

type MineFieldOrders struct {
	Detonate bool `json:"detonate,omitempty"`
}

// The radius of a minefield is the sqrt of its mines
func (mf *MineField) Radius() float64 {
	return math.Sqrt(float64(mf.NumMines))
}

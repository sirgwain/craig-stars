package cs

import (
	"fmt"
	"math"
)

type MineFieldType string

const (
	MineFieldTypeStandard  MineFieldType = "Standard"
	MineFieldTypeHeavy     MineFieldType = "Heavy"
	MineFieldTypeSpeedBump MineFieldType = "SpeedBump"
)

type MineField struct {
	MapObject
	MineFieldOrders
	MineFieldType MineFieldType `json:"mineFieldType"`
	NumMines      int           `json:"numMines"`
	Spec          MineFieldSpec `json:"spec"`
}

type MineFieldOrders struct {
	Detonate bool `json:"detonate,omitempty"`
}

type MineFieldSpec struct {
	Radius float64 `json:"radius"`
}

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

// The radius of a minefield is the sqrt of its mines
func (mf *MineField) Radius() float64 {
	return math.Sqrt(float64(mf.NumMines))
}

func computeMinefieldSpec(mineField *MineField) MineFieldSpec {
	spec := MineFieldSpec{}
	spec.Radius = mineField.Radius()
	return spec
}

func newMineField(player *Player, mineFieldType MineFieldType, numMines int, num int, position Vector) *MineField {
	return &MineField{
		MapObject: MapObject{
			Type:      MapObjectTypeFleet,
			PlayerNum: player.Num,
			Dirty:     true,
			Num:       num,
			Name:      fmt.Sprintf("%s MineField #%d", player.Race.PluralName, num),
			Position:  position,
		},
		MineFieldType: mineFieldType,
		NumMines:      numMines,
	}
}

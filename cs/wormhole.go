package cs

import (
	"fmt"

	"slices"

	"github.com/rs/zerolog/log"
)

type Wormhole struct {
	MapObject
	DestinationNum   int               `json:"destinationNum,omitempty"`
	Stability        WormholeStability `json:"stability,omitempty"`
	YearsAtStability int               `json:"yearsAtStability,omitempty"`
	Spec             WormholeSpec      `json:"spec,omitempty"`
}

type WormholeSpec struct {
	Stats WormholeStats
}

type WormholeStats struct {
	YearsToDegrade int     `json:"yearsToDegrade"`
	ChanceToJump   float64 `json:"chanceToJump"`
	JiggleDistance int     `json:"jiggleDistance"`
}

type WormholeStability string

const (
	WormholeStabilityNone              WormholeStability = ""
	WormholeStabilityRockSolid         WormholeStability = "RockSolid"
	WormholeStabilityStable            WormholeStability = "Stable"
	WormholeStabilityMostlyStable      WormholeStability = "MostlyStable"
	WormholeStabilityAverage           WormholeStability = "Average"
	WormholeStabilitySlightlyVolatile  WormholeStability = "SlightlyVolatile"
	WormholeStabilityVolatile          WormholeStability = "Volatile"
	WormholeStabilityExtremelyVolatile WormholeStability = "ExtremelyVolatile"
)

var WormholeStabilities []WormholeStability = []WormholeStability{
	WormholeStabilityRockSolid,
	WormholeStabilityStable,
	WormholeStabilityMostlyStable,
	WormholeStabilityAverage,
	WormholeStabilitySlightlyVolatile,
	WormholeStabilityVolatile,
	WormholeStabilityExtremelyVolatile,
}

func (w *Wormhole) String() string {
	return fmt.Sprintf("Wormhole #%d %v", w.Num, w.Position)
}

func newWormhole(position Vector, num int, stability WormholeStability) *Wormhole {
	return &Wormhole{
		MapObject: MapObject{
			Type:     MapObjectTypeWormhole,
			Position: position,
			Num:      num,
		},
		Stability: stability,
	}
}

func computeWormholeSpec(w *Wormhole, rules *Rules) WormholeSpec {
	return WormholeSpec{
		Stats: rules.WormholeStatsByStability[w.Stability],
	}
}

func generateWormhole(mapObjectGetter mapObjectGetter, area Vector, random rng, planetPositions []Vector, wormholePositions []Vector, minDistanceFromPlanets int) (position Vector, stability WormholeStability, err error) {
	width, height := int(area.X), int(area.Y)

	position = Vector{X: float64(random.Intn(width)), Y: float64(random.Intn(height))}

	minWormholeDistance := (height + width) / 2 / 4
	posCheckCount := 0
	for !mapObjectGetter.isPositionValid(position, &planetPositions, float64(minDistanceFromPlanets)) ||
		!mapObjectGetter.isPositionValid(position, &wormholePositions, float64(minWormholeDistance)) {
		position = Vector{X: float64(random.Intn(width)), Y: float64(random.Intn(height))}
		posCheckCount++
		if posCheckCount > 1000 {
			return Vector{}, WormholeStabilityNone, fmt.Errorf("find a valid position for a wormhole in 1000 tries, min: %d, numPlanets: %d, numWormholes: %d, area: %v", minDistanceFromPlanets, len(planetPositions), len(wormholePositions), area)
		}
	}

	stability = WormholeStabilities[random.Intn(len(WormholeStabilities))]

	return position, stability, nil
}

func (w *Wormhole) jiggle(area Vector, mapObjectGetter mapObjectGetter, random rng) {
	stats := w.Spec.Stats

	// don't infinite jiggle
	jiggleCount := 0
	var newPosition Vector
	for {
		newPosition = Vector{
			ClampFloat64(w.Position.X+float64(random.Intn(stats.JiggleDistance/2)-stats.JiggleDistance/2), 0, area.X),
			ClampFloat64(w.Position.Y+float64(random.Intn(stats.JiggleDistance/2)-stats.JiggleDistance/2), 0, area.Y),
		}
		log.Debug().Msgf("%v jiggled to %v", w, newPosition)
		jiggleCount++

		if mapObjectGetter.getMapObjectsAtPosition(newPosition) == nil || jiggleCount > 100 {
			break
		}
	}
	w.Position = newPosition
	w.MarkDirty()
}

func (w *Wormhole) degrade() {
	stats := w.Spec.Stats
	w.YearsAtStability++
	w.MarkDirty()
	if w.YearsAtStability > stats.YearsToDegrade && stats.YearsToDegrade != Infinite {
		if w.Stability != WormholeStabilityExtremelyVolatile {
			// go to the next stability
			w.Stability = WormholeStabilities[slices.Index(WormholeStabilities, w.Stability)+1]
			w.YearsAtStability = 0
			w.MarkDirty()
			log.Debug().Msgf("%v degraded to %s", w, w.Stability)
		}
	}
}

func (w *Wormhole) shouldJump(random rng) bool {
	randVal := random.Float64()
	return w.Spec.Stats.ChanceToJump > randVal
}

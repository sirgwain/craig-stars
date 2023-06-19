package cs

import (
	"fmt"
	"math/rand"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type Wormhole struct {
	MapObject        `json:"mapObject,omitempty"`
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

func generateWormhole(mapObjectGetter mapObjectGetter, num int, area Vector, random *rand.Rand, planetPositions []Vector, wormholePositions []Vector, minDistanceFromPlanets int) (position Vector, stability WormholeStability, err error) {
	width, height := int(area.X), int(area.Y)

	position = Vector{X: float64(random.Intn(width)), Y: float64(random.Intn(height))}

	minWormholeDistance := (height + width) / 2 / 4
	posCheckCount := 0
	for !mapObjectGetter.isPositionValid(position, &planetPositions, float64(minDistanceFromPlanets)) ||
		!mapObjectGetter.isPositionValid(position, &wormholePositions, float64(minWormholeDistance)) {
		position = Vector{X: float64(random.Intn(width)), Y: float64(random.Intn(height))}
		posCheckCount++
		if posCheckCount > 1000 {
			return Vector{}, WormholeStabilityNone, fmt.Errorf("find a valid position for a planet in 1000 tries, min: %d, numPlanets: %d, numWormholes: %d, area: %v", minDistanceFromPlanets, len(planetPositions), len(wormholePositions), area)
		}
	}

	stability = WormholeStabilities[random.Intn(len(WormholeStabilities))]

	return position, stability, nil
}

func (w *Wormhole) jiggle(mapObjectGetter mapObjectGetter, random *rand.Rand) {
	stats := w.Spec.Stats

	// don't infinite jiggle
	jiggleCount := 0
	var newPosition Vector
	for {
		newPosition = Vector{
			w.Position.X + float64(random.Intn(stats.JiggleDistance/2)-stats.JiggleDistance/2),
			w.Position.Y + float64(random.Intn(stats.JiggleDistance/2)-stats.JiggleDistance/2),
		}
		log.Debug().Msgf("%v jiggled to %v", w, newPosition)
		jiggleCount++

		if mapObjectGetter.getMapObjectsAtPosition(newPosition) == nil || jiggleCount > 100 {
			break
		}
	}
	w.Position = newPosition
	w.Dirty = true
}

func (w *Wormhole) degrade() {
	stats := w.Spec.Stats
	w.YearsAtStability++
	w.Dirty = true
	if w.YearsAtStability > stats.YearsToDegrade {
		if w.Stability != WormholeStabilityExtremelyVolatile {
			// go to the next stability
			w.Stability = WormholeStabilities[slices.Index(WormholeStabilities, w.Stability)+1]
			w.YearsAtStability = 0
			w.Dirty = true
			log.Debug().Msgf("%v degraded to %s", w, w.Stability)
		}
	}
}

func (w *Wormhole) shouldJump(random *rand.Rand) bool {
	randVal := random.Float64()
	return w.Spec.Stats.ChanceToJump > randVal
}

//go:build wasi || wasm

package wasm

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	craig_starsv1wasm "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1wasm"
)

// Each wasm instannce is unique to a browser session, so keep track of state so we don't have to
// send it and serialize it for each call
type state struct {
	rules    *cs.Rules
	player   *cs.Player
	universe *cs.Universe
}

var debug = false

var serviceHandler = craig_starsv1wasm.NewWasmServiceHandler(&wasmService{
	state: state{rules: &cs.StandardRules},
})

type wasmService struct {
	state
}

func (s *wasmService) AddWaypoint(ctx context.Context, req *craig_starsv1.AddWaypointRequest) (*craig_starsv1.AddWaypointResponse, error) {
	fleet := converter.C.ConvertFleet(req.Fleet)
	dest := converter.C.ConvertWaypointDest(req.Dest)

	fleet.InjectDesigns(s.player.Designs)

	index := fleet.AddWaypoint(s.player, dest, int(req.CurrentSelectedWaypointIndex), req.FastestWaypoint)
	return &craig_starsv1.AddWaypointResponse{Index: int32(index), Fleet: converter.C.ConvertCSFleet(fleet)}, nil
}

func (s *wasmService) CalculateRacePoints(ctx context.Context, req *craig_starsv1.CalculateRacePointsRequest) (*craig_starsv1.CalculateRacePointsResponse, error) {
	race := converter.C.ConvertRaceP(req.Race)
	race.Spec = cs.ComputeRaceSpec(race, s.rules)
	points := race.ComputeRacePoints(s.rules.RaceStartingPoints)
	return &craig_starsv1.CalculateRacePointsResponse{Points: int32(points)}, nil
}

func (s *wasmService) ComputeMinefieldSpec(ctx context.Context, req *craig_starsv1.ComputeMinefieldSpecRequest) (*craig_starsv1.ComputeMinefieldSpecResponse, error) {
	minefield := converter.C.ConvertMinefield(req.Minefield)
	spec := cs.ComputeMinefieldSpec(s.rules, s.player, minefield, cs.NumMapObjectsWithin(s.player.Intels.PlanetIntels, minefield.Position, minefield.Radius()))
	return &craig_starsv1.ComputeMinefieldSpecResponse{Spec: converter.C.ConvertCSMinefieldSpec(spec)}, nil
}

func (s *wasmService) ComputePlayerResearchSpec(ctx context.Context, req *craig_starsv1.ComputePlayerResearchSpecRequest) (*craig_starsv1.ComputePlayerResearchSpecResponse, error) {
	planets := []*cs.Planet{}
	for _, p := range s.player.PlanetIntels {
		if p.PlayerNum != s.player.Num {
			continue
		}
		planets = append(planets, p)
	}

	spec := cs.ComputePlayerResearchSpec(s.player, s.rules, planets)
	return &craig_starsv1.ComputePlayerResearchSpecResponse{Spec: converter.C.ConvertCSPlayerResearchSpec(spec)}, nil

}

func (s *wasmService) ComputeRaceSpec(ctx context.Context, req *craig_starsv1.ComputeRaceSpecRequest) (*craig_starsv1.ComputeRaceSpecResponse, error) {
	race := converter.C.ConvertRaceP(req.Race)
	spec := cs.ComputeRaceSpec(race, s.rules)

	return &craig_starsv1.ComputeRaceSpecResponse{
		Spec: converter.C.ConvertCSRaceSpec(spec),
	}, nil
}

func (s *wasmService) ComputeShipDesignSpec(ctx context.Context, req *craig_starsv1.ComputeShipDesignSpecRequest) (*craig_starsv1.ComputeShipDesignSpecResponse, error) {
	design := converter.C.ConvertShipDesignP(req.Design)
	spec, err := cs.ComputeShipDesignSpec(s.rules, s.player.TechLevels, cs.ComputeRaceSpec(&s.player.Race, s.rules), design)
	if err != nil {
		return nil, err
	}
	return &craig_starsv1.ComputeShipDesignSpecResponse{
		Spec: converter.C.ConvertCSShipDesignSpec(spec),
	}, nil
}

func (s *wasmService) EnableDebug(ctx context.Context, req *craig_starsv1.EnableDebugRequest) (*craig_starsv1.EnableDebugResponse, error) {
	debug = req.Debug
	EnableDebug()
	return &craig_starsv1.EnableDebugResponse{}, nil
}

func (s *wasmService) EstimateProduction(ctx context.Context, req *craig_starsv1.EstimateProductionRequest) (*craig_starsv1.EstimateProductionResponse, error) {
	planet := converter.C.ConvertPlanetP(req.Planet)
	// setup the starbase
	if planet.Spec.HasStarbase {
		planet.Starbase = &cs.Fleet{
			Tokens: []cs.ShipToken{
				{Quantity: 1, DesignNum: planet.Spec.StarbaseDesignNum},
			},
		}
	}

	// make sure if we have a starbase, it has a design so we can compute
	// upgrade costs
	if err := planet.PopulateStarbaseDesign(s.player); err != nil {
		return nil, fmt.Errorf("failed to populate %s starbase with player design.: %v", planet.Name, err)
	}

	if err := planet.PopulateProductionQueueDesigns(s.player); err != nil {
		return nil, fmt.Errorf("failed to populate %s production queue designs.: %v", planet.Name, err)
	}

	planet.PopulateProductionQueueEstimates(s.rules, s.player)

	return &craig_starsv1.EstimateProductionResponse{
		Planet: converter.C.ConvertCSPlanet(planet),
	}, nil
}

func (s *wasmService) GetMaxBuildable(ctx context.Context, req *craig_starsv1.GetMaxBuildableRequest) (*craig_starsv1.GetMaxBuildableResponse, error) {
	planet := converter.C.ConvertPlanetP(req.Planet)
	itemType := cs.QueueItemType(req.ItemType)

	maxBuild := 5000
	if !itemType.IsAuto() /* || itemType == cs.QueueItemTypeAutoMineralAlchemy */ {
		maxBuild = planet.MaxBuildable(s.player, itemType)
		// maxBuildable is set to infinite for most things, but we want big number
		if maxBuild == cs.Infinite {
			maxBuild = 5000
		}
	}

	return &craig_starsv1.GetMaxBuildableResponse{
		Result: int32(maxBuild),
	}, nil
}

// GetPlanetHabitability implements craig_starsv1wasm.WasmServiceHandler.
func (s *wasmService) GetPlanetHabitability(ctx context.Context, req *craig_starsv1.GetPlanetHabitabilityRequest) (*craig_starsv1.GetPlanetHabitabilityResponse, error) {
	race := converter.C.ConvertRaceP(req.Race)
	hab := converter.C.ConvertHab(req.Hab)

	return &craig_starsv1.GetPlanetHabitabilityResponse{
		Result: int32(race.GetPlanetHabitability(hab)),
	}, nil
}

func (s *wasmService) GetResearchCost(ctx context.Context, req *craig_starsv1.GetResearchCostRequest) (*craig_starsv1.GetResearchCostResponse, error) {
	resources := s.player.GetResearchCost(s.rules, converter.C.ConvertTechLevel(req.TechLevel))
	return &craig_starsv1.GetResearchCostResponse{Resources: int32(resources)}, nil
}

func (s *wasmService) GetStarbaseUpgradeCost(ctx context.Context, req *craig_starsv1.GetStarbaseUpgradeCostRequest) (*craig_starsv1.GetStarbaseUpgradeCostResponse, error) {
	costCalculatoor := cs.NewCostCalculator()
	design := converter.C.ConvertShipDesignP(req.Design)
	newDesign := converter.C.ConvertShipDesignP(req.NewDesign)
	cost, err := costCalculatoor.StarbaseUpgradeCost(s.rules, s.player.TechLevels, s.player.Race.Spec, design, newDesign)
	if err != nil {
		return nil, fmt.Errorf("unable to calculate starbase upgrade cost: %v", err)
	}
	return &craig_starsv1.GetStarbaseUpgradeCostResponse{
		Cost: converter.C.ConvertCSCost(cost),
	}, nil
}

func (s *wasmService) GetTechCost(ctx context.Context, req *craig_starsv1.GetTechCostRequest) (*craig_starsv1.GetTechCostResponse, error) {
	if s.player == nil {
		return nil, errors.New("no player, can't calculate tech cost")
	}
	costCalculatoor := cs.NewCostCalculator()

	cost := costCalculatoor.GetTechCost(s.rules, s.player.TechLevels, s.player.Race.Spec, converter.C.ConvertTech(req.Tech))
	return &craig_starsv1.GetTechCostResponse{Cost: converter.C.ConvertCSCost(cost)}, nil
}

func (s *wasmService) SetDesigns(ctx context.Context, req *craig_starsv1.SetDesignsRequest) (*craig_starsv1.SetDesignsResponse, error) {
	s.player.Designs = converter.C.ConvertShipDesigns(req.Designs)
	return &craig_starsv1.SetDesignsResponse{}, nil
}

func (s *wasmService) SetPlayer(ctx context.Context, req *craig_starsv1.SetPlayerRequest) (*craig_starsv1.SetPlayerResponse, error) {
	player := converter.C.ConvertPlayer(req.Player)
	player.Race.Spec = cs.ComputeRaceSpec(&player.Race, s.rules)

	// TODO: remove this code after removing designs from base player object
	var designs []*cs.ShipDesign
	var intels cs.Intels
	if s.player != nil {
		designs = s.player.Designs
		intels = s.player.Intels
	}
	s.player = player
	s.player.Designs = designs
	s.player.Intels = intels
	return &craig_starsv1.SetPlayerResponse{}, nil
}

func (s *wasmService) SetIntels(ctx context.Context, req *craig_starsv1.SetIntelsRequest) (*craig_starsv1.SetIntelsResponse, error) {
	s.player.Intels = converter.C.ConvertIntels(req.Intels)
	return &craig_starsv1.SetIntelsResponse{}, nil
}

func (s *wasmService) UpdatePlanet(ctx context.Context, req *craig_starsv1.UpdatePlanetRequest) (*craig_starsv1.UpdatePlanetResponse, error) {
	planet := converter.C.ConvertPlanetP(req.Planet)

	if planet.Num <= 0 || planet.Num > len(s.player.PlanetIntels) {
		return nil, errors.New("planet num out of range")
	}

	// save this updated planet back to our intel
	s.player.PlanetIntels[planet.Num-1] = planet
	return &craig_starsv1.UpdatePlanetResponse{}, nil
}

func (s *wasmService) UpdatePlanets(ctx context.Context, req *craig_starsv1.UpdatePlanetsRequest) (*craig_starsv1.UpdatePlanetsResponse, error) {
	planets := converter.C.ConvertPlanets(req.Planets)

	for _, planet := range planets {
		if planet.Num <= 0 || planet.Num > len(s.player.PlanetIntels) {
			return nil, errors.New("planet num out of range")
		}

		// save this updated planet back to our intel
		s.player.PlanetIntels[planet.Num-1] = planet

	}
	return &craig_starsv1.UpdatePlanetsResponse{}, nil
}

func (s *wasmService) UpdateWaypoint(ctx context.Context, req *craig_starsv1.UpdateWaypointRequest) (*craig_starsv1.UpdateWaypointResponse, error) {
	fleet := converter.C.ConvertFleet(req.Fleet)
	dest := converter.C.ConvertWaypointDest(req.Dest)

	fleet.InjectDesigns(s.player.Designs)

	result := fleet.UpdateWaypoint(s.player, dest, int(req.CurrentSelectedWaypointIndex), req.FastestWaypoint)
	return &craig_starsv1.UpdateWaypointResponse{
		Result: converter.CSUpdateWaypointResultToUpdateWaypointResult(result),
		Fleet:  converter.C.ConvertCSFleet(fleet),
	}, nil
}

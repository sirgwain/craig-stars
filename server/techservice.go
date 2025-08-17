//go:build !wasi && !wasm

package server

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/proto/converter"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1/craig_starsv1connect"
)

// NewTechServiceHandler constructs the TechService handler.
func NewTechServiceHandler() craig_starsv1connect.TechServiceHandler {
	return &techService{}
}

type techService struct{}

// GetTech returns a single Tech by name, or NOT_FOUND if missing.
func (s *techService) GetTech(ctx context.Context, req *connect.Request[craig_starsv1.GetTechRequest]) (*connect.Response[craig_starsv1.GetTechResponse], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}
	tech := cs.StaticTechStore.GetTech(req.Msg.GetName())
	if tech == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("tech not found"))
	}

	resp := &craig_starsv1.GetTechResponse{}
	switch t := tech.(type) {
	case *cs.TechPlanetaryScanner:
		resp.Tech = &craig_starsv1.GetTechResponse_TechPlanetaryScanner{
			TechPlanetaryScanner: converter.C.ConvertCSTechPlanetaryScanner(t),
		}
	case *cs.TechTerraform:
		resp.Tech = &craig_starsv1.GetTechResponse_TechTerraform{
			TechTerraform: converter.C.ConvertCSTechTerraform(t),
		}
	case *cs.TechDefense:
		resp.Tech = &craig_starsv1.GetTechResponse_TechDefense{
			TechDefense: converter.C.ConvertCSTechDefense(t),
		}
	case *cs.TechPlanetary:
		resp.Tech = &craig_starsv1.GetTechResponse_TechPlanetary{
			TechPlanetary: converter.C.ConvertCSTechPlanetary(t),
		}
	case *cs.TechHullComponent:
		resp.Tech = &craig_starsv1.GetTechResponse_TechHullComponent{
			TechHullComponent: converter.C.ConvertCSTechHullComponentP(t),
		}
	case *cs.TechHull:
		resp.Tech = &craig_starsv1.GetTechResponse_TechHull{
			TechHull: converter.C.ConvertCSTechHull(t),
		}
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unsupported tech type"))
	}

	return connect.NewResponse(resp), nil
}

// GetTechs returns the full TechStore (StaticTechStore).
func (s *techService) GetTechs(ctx context.Context, req *connect.Request[craig_starsv1.GetTechsRequest]) (*connect.Response[craig_starsv1.GetTechsResponse], error) {
	techs := cs.StaticTechStore
	return connect.NewResponse(&craig_starsv1.GetTechsResponse{
		PlanetaryScanners: converter.C.ConvertCSTechPlanetaryScanners(techs.PlanetaryScanners),
		Terraforms:        converter.C.ConvertCSTechTerraforms(techs.Terraforms),
		Defenses:          converter.C.ConvertCSTechDefenses(techs.Defenses),
		Planetaries:       converter.C.ConvertCSTechPlanetaries(techs.Planetaries),
		HullComponents:    converter.C.ConvertCSTechHullComponents(techs.HullComponents),
		Hulls:             converter.C.ConvertCSTechHulls(techs.Hulls),
	}), nil
}

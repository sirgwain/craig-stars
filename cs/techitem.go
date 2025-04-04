package cs

// The "type" of a given TechItem.
type TechItemType string

const (
	TechItemTypeNone             TechItemType = ""
	TechItemTypeDefense          TechItemType = "Defense"
	TechItemTypeEngine           TechItemType = "Engine"
	TechItemTypeHull             TechItemType = "Hull"
	TechItemTypeHullComponent    TechItemType = "HullComponent"
	TechItemTypePlanetary        TechItemType = "Planetary"
	TechItemTypePlanetaryScanner TechItemType = "PlanetaryScanner"
	TechItemTypeTerraform        TechItemType = "Terraform"
)

var TechItemTypes = []TechItemType{
	TechItemTypeNone,
	TechItemTypeDefense,
	TechItemTypeEngine,
	TechItemTypeHull,
	TechItemTypeHullComponent,
	TechItemTypePlanetary,
	TechItemTypeTerraform,
	TechItemTypePlanetaryScanner,
}

// A tech item of unknown type.
type TechItem interface {
	getType() TechItemType
}

func (defense *TechDefense) getType() TechItemType          { return TechItemTypeDefense }
func (engine *TechEngine) getType() TechItemType            { return TechItemTypeEngine }
func (hull *TechHull) getType() TechItemType                { return TechItemTypeHull }
func (hc *TechHullComponent) getType() TechItemType         { return TechItemTypeHullComponent }
func (planetary *TechPlanetary) getType() TechItemType      { return TechItemTypePlanetary }
func (scanner *TechPlanetaryScanner) getType() TechItemType { return TechItemTypePlanetaryScanner }
func (terraform *TechTerraform) getType() TechItemType      { return TechItemTypeTerraform }

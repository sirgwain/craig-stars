package cs

import (
	"fmt"
	"math"

	"github.com/rs/zerolog/log"
)

// handle player orders
type Orderer interface {
	UpdatePlayerOrders(player *Player, playerPlanets []*Planet, order PlayerOrders, rules *Rules)
	UpdatePlanetOrders(rules *Rules, player *Player, planet *Planet, orders PlanetOrders) error
	UpdateFleetOrders(player *Player, fleet *Fleet, orders FleetOrders)
	UpdateMineFieldOrders(player *Player, minefield *MineField, orders MineFieldOrders)
	TransferFleetCargo(rules *Rules, player, destPlayer *Player, source, dest *Fleet, transferAmount Cargo) error
	TransferPlanetCargo(rules *Rules, player *Player, source *Fleet, dest *Planet, transferAmount Cargo) error
	TransferSalvageCargo(rules *Rules, player *Player, source *Fleet, dest *Salvage, nextSalvageNum int, transferAmount Cargo) (*Salvage, error)
	SplitFleetTokens(rules *Rules, player *Player, playerFleets []*Fleet, source *Fleet, tokens []ShipToken) (*Fleet, error)
	SplitAll(rules *Rules, player *Player, playerFleets []*Fleet, source *Fleet) ([]*Fleet, error)
	Merge(rules *Rules, player *Player, fleets []*Fleet) (*Fleet, error)
}

type orders struct {
}

func NewOrderer() Orderer {
	return &orders{}
}

// update a player's orders
func (o *orders) UpdatePlayerOrders(player *Player, playerPlanets []*Planet, orders PlayerOrders, rules *Rules) {
	researchAmountUpdated := orders.ResearchAmount != player.ResearchAmount

	// save the new orders
	player.PlayerOrders = orders

	// if we updated the research amount, update the planet specs
	if researchAmountUpdated {
		for _, planet := range playerPlanets {
			if planet.PlayerNum == player.Num {
				spec := &planet.Spec
				spec.computeResourcesPerYearAvailable(player, planet)
				planet.MarkDirty()
			}
		}
	}

	player.Spec = computePlayerSpec(player, rules, playerPlanets)

}

// update a planet orders
func (o *orders) UpdatePlanetOrders(rules *Rules, player *Player, planet *Planet, orders PlanetOrders) error {
	planet.PlanetOrders = orders

	// make sure if we have a starbase, it has a design so we can compute
	// upgrade costs
	if err := planet.PopulateStarbaseDesign(player); err != nil {
		return err
	}

	spec := &planet.Spec

	// update the player spec with new values from this planet
	oldResourcesPerYearResearch := spec.ResourcesPerYearResearch
	oldResourcesPerYearResearchEstimatedLeftover := spec.ResourcesPerYearResearchEstimatedLeftover

	spec.computeResourcesPerYearAvailable(player, planet)
	if err := planet.PopulateProductionQueueDesigns(player); err != nil {
		return err
	}

	// make sure we can actually build this stuff
	for _, item := range planet.ProductionQueue {
		if item.design != nil && item.design.OriginalPlayerNum != None {
			return fmt.Errorf("cannot build %s, it was transferred from another player", item.design.Name)
		}
	}

	if err := planet.PopulateProductionQueueCosts(player); err != nil {
		return err
	}

	planet.PopulateProductionQueueEstimates(rules, player)

	// update the player spec with the change in resources for this planet
	player.Spec.ResourcesPerYearResearch = player.Spec.ResourcesPerYearResearch - oldResourcesPerYearResearch + spec.ResourcesPerYearResearch
	player.Spec.ResourcesPerYearResearchEstimated = player.Spec.ResourcesPerYearResearchEstimated - oldResourcesPerYearResearchEstimatedLeftover + spec.ResourcesPerYearResearchEstimatedLeftover

	planet.MarkDirty()
	return nil
}

// update the orders to a fleet
func (o *orders) UpdateFleetOrders(player *Player, fleet *Fleet, orders FleetOrders) {
	// copy user modifiable things to the fleet fleet
	fleet.RepeatOrders = orders.RepeatOrders
	wp0 := &fleet.Waypoints[0]
	newWP0 := orders.Waypoints[0]

	// TODO: do we want to lookup the target?
	wp0.WarpSpeed = newWP0.WarpSpeed
	wp0.Task = newWP0.Task
	wp0.TransportTasks = newWP0.TransportTasks
	wp0.LayMineFieldDuration = newWP0.LayMineFieldDuration
	wp0.WaitAtWaypoint = newWP0.WaitAtWaypoint
	wp0.TargetName = newWP0.TargetName
	wp0.TargetType = newWP0.TargetType
	wp0.TargetNum = newWP0.TargetNum
	wp0.TargetPlayerNum = newWP0.TargetPlayerNum
	wp0.TransferToPlayer = newWP0.TransferToPlayer

	fleet.Waypoints = append(fleet.Waypoints[:1], orders.Waypoints[1:]...)
	fleet.computeFuelUsage(player)
	fleet.MarkDirty()
}

func (o *orders) UpdateMineFieldOrders(player *Player, minefield *MineField, orders MineFieldOrders) {
	minefield.MineFieldOrders = orders
}

// transfer cargo from a fleet to/from a fleet
func (o *orders) TransferFleetCargo(rules *Rules, player, destPlayer *Player, source, dest *Fleet, transferAmount Cargo) error {

	if source.availableCargoSpace() < transferAmount.Total() {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", source.Name, source.availableCargoSpace(), transferAmount.Total(), dest.Name)
	}

	if dest.availableCargoSpace() < -transferAmount.Total() {
		return fmt.Errorf("dest %s has %d cargo space available, cannot transfer %dkT from %s", dest.Name, dest.availableCargoSpace(), transferAmount.Total(), dest.Name)
	}

	if !dest.Cargo.CanTransfer(transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %v from %s, there is not enough to transfer", source.Name, transferAmount, dest.Name)
	}

	if !source.Cargo.CanTransfer(transferAmount.Negative()) {
		return fmt.Errorf("fleet %s cannot transfer %v to %s, the fleet does not have enough the required cargo", source.Name, transferAmount.Negative(), dest.Name)
	}

	// transfer the cargo
	source.Cargo = source.Cargo.Add(transferAmount)
	dest.Cargo = dest.Cargo.Subtract(transferAmount)

	source.Spec = ComputeFleetSpec(rules, player, source)
	dest.Spec = ComputeFleetSpec(rules, destPlayer, dest)
	source.MarkDirty()
	dest.MarkDirty()

	return nil
}

// transfer cargo from a planet to/from a fleet
func (o *orders) TransferPlanetCargo(rules *Rules, player *Player, source *Fleet, dest *Planet, transferAmount Cargo) error {

	if source.availableCargoSpace() < transferAmount.Total() {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", source.Name, source.availableCargoSpace(), transferAmount.Total(), dest.Name)
	}

	if !dest.Cargo.CanTransfer(transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %v from %s, the planet does not have the required cargo", source.Name, transferAmount, dest.Name)
	}

	if !source.Cargo.CanTransfer(transferAmount.Negative()) {
		return fmt.Errorf("fleet %s cannot transfer %v to %s, the fleet does not have enough the required cargo", source.Name, transferAmount.Negative(), dest.Name)
	}

	// transfer the cargo
	source.Cargo = source.Cargo.Add(transferAmount)
	dest.Cargo = dest.Cargo.Subtract(transferAmount)
	source.Spec = ComputeFleetSpec(rules, player, source)

	if dest.Starbase != nil {
		dest.Starbase.Tokens[0].design = player.GetDesign(dest.Starbase.Tokens[0].DesignNum)
	}
	starbaseSpec := dest.Spec.PlanetStarbaseSpec
	dest.Spec = computePlanetSpec(rules, player, dest)
	dest.Spec.PlanetStarbaseSpec = starbaseSpec

	source.MarkDirty()
	dest.MarkDirty()
	return nil
}

// transfer cargo from a planet to/from a fleet
func (o *orders) TransferSalvageCargo(rules *Rules, player *Player, source *Fleet, dest *Salvage, nextSalvageNum int, transferAmount Cargo) (*Salvage, error) {

	if source.availableCargoSpace() < transferAmount.Total() {
		return nil, fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", source.Name, source.availableCargoSpace(), transferAmount.Total(), dest.Name)
	}

	if dest != nil && !dest.Cargo.CanTransfer(transferAmount) {
		return nil, fmt.Errorf("fleet %s cannot transfer %v from %s, the salvage does not have the required cargo", source.Name, transferAmount, dest.Name)
	}

	if !source.Cargo.CanTransfer(transferAmount.Negative()) {
		return nil, fmt.Errorf("fleet %s cannot transfer %v to %s, the fleet does not have enough the required cargo", source.Name, transferAmount.Negative(), dest.Name)
	}

	if dest == nil {
		dest = newSalvage(source.Position, nextSalvageNum, source.PlayerNum, transferAmount.Negative())
	} else {
		dest.Cargo = dest.Cargo.Subtract(transferAmount)
	}

	// transfer the cargo
	source.Cargo = source.Cargo.Add(transferAmount)
	source.Spec = ComputeFleetSpec(rules, player, source)

	// make our player aware of this salvage
	discover := newDiscoverer(player)
	discover.discoverSalvage(dest)

	source.MarkDirty()
	dest.MarkDirty()
	return dest, nil
}

// split a fleet's tokens into a new fleet
func (o *orders) SplitFleetTokens(rules *Rules, player *Player, playerFleets []*Fleet, source *Fleet, tokens []ShipToken) (*Fleet, error) {
	if source == nil {
		return nil, fmt.Errorf("no source fleet to split")
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to split")
	}

	// build a map of tokens by their design
	splitTokensByDesign := map[int]*ShipToken{}
	for i := range tokens {
		token := &tokens[i]
		splitTokensByDesign[token.DesignNum] = token
	}

	remainingTokens := 0
	for _, token := range source.Tokens {
		if split, found := splitTokensByDesign[token.DesignNum]; found {
			if split.Quantity > token.Quantity {
				return nil, fmt.Errorf("tried to split token (design %d) from %d to %d tokens", token.DesignNum, token.Quantity, split.Quantity)
			}
			if split.Quantity < token.Quantity {
				// we have some leftover
				remainingTokens++
			}
		} else {
			remainingTokens++
		}
	}

	// build a map of source tokens by their design
	sourceTokensByDesign := map[int]*ShipToken{}
	for i := range source.Tokens {
		token := &source.Tokens[i]
		sourceTokensByDesign[token.DesignNum] = token
	}

	// make sure we don't have split tokens that aren't in our source
	for _, token := range tokens {
		if _, found := sourceTokensByDesign[token.DesignNum]; !found {
			return nil, fmt.Errorf("tried to create an entirely new token for design %d", token.DesignNum)
		}
	}

	if remainingTokens == 0 {
		return nil, fmt.Errorf("no ships left in source fleet")
	}

	// finally, time to split!
	// make sure our fleets all have designs
	if err := player.InjectDesigns(append(playerFleets, source)); err != nil {
		return nil, err
	}

	// now create the new fleet
	fleetNum := player.getNextFleetNum(playerFleets)

	// build a map of designs so we can fill in the tokens
	designsByNum := make(map[int]*ShipDesign, len(player.Designs))
	for i := range player.Designs {
		design := player.Designs[i]
		designsByNum[design.Num] = design
	}

	// fill in designs
	for i := range tokens {
		token := &tokens[i]
		token.design = designsByNum[token.DesignNum]

		if token.design == nil {
			return nil, fmt.Errorf("unable to find design %d", token.DesignNum)
		}
	}

	for i := range source.Tokens {
		token := &source.Tokens[i]
		token.design = designsByNum[token.DesignNum]

		if token.design == nil {
			return nil, fmt.Errorf("unable to find design %d for fleet %s", token.DesignNum, source.Name)
		}
	}

	var baseDesign = tokens[0].design
	fleet := newFleet(player, baseDesign, fleetNum, baseDesign.Name, source.Waypoints)
	fleet.OrbitingPlanetNum = source.OrbitingPlanetNum
	fleet.Heading = source.Heading
	fleet.WarpSpeed = source.WarpSpeed
	fleet.PreviousPosition = source.PreviousPosition
	fleet.battlePlan = source.battlePlan
	fleet.Tokens = tokens

	// the fleet has some percentage of fuel fullness
	fleetFuelFullness := float64(source.Fuel) / float64(source.Spec.FuelCapacity)

	totalCargo := source.Cargo.Total()
	totalCargoCapacity := source.Spec.CargoCapacity

	// now remove all the tokens from the old fleet
	for i := range tokens {
		splitToken := &tokens[i]
		sourceToken := sourceTokensByDesign[splitToken.DesignNum]
		// quantity := sourceToken.Quantity
		// quantityDamaged := sourceToken.QuantityDamaged

		// each ship in a token has some amount of fuel, based on the total fuel on the fleet
		shipFuel := fleetFuelFullness * float64(sourceToken.design.Spec.FuelCapacity)

		var shipCargoPercent float64 = 0
		if totalCargo > 0 && splitToken.design.Spec.CargoCapacity > 0 {
			// see how much this ship's cargo capacity is compared to the fleet total
			shipCargoPercent = float64(splitToken.design.Spec.CargoCapacity) / float64(totalCargoCapacity)
		}

		// leave any remainder fuel on the source
		fuelToMove := int(math.Floor(shipFuel * float64(splitToken.Quantity)))
		fleet.Fuel += fuelToMove
		source.Fuel -= fuelToMove

		if totalCargo > 0 && shipCargoPercent > 0 {
			fleet.Cargo = Cargo{
				Ironium:   int(math.Floor(shipCargoPercent * float64(splitToken.Quantity) * float64(source.Cargo.Ironium))),
				Boranium:  int(math.Floor(shipCargoPercent * float64(splitToken.Quantity) * float64(source.Cargo.Boranium))),
				Germanium: int(math.Floor(shipCargoPercent * float64(splitToken.Quantity) * float64(source.Cargo.Germanium))),
				Colonists: int(math.Floor(shipCargoPercent * float64(splitToken.Quantity) * float64(source.Cargo.Colonists))),
			}
			source.Cargo = source.Cargo.Subtract(fleet.Cargo)
		}

		// split damage
		if sourceToken.QuantityDamaged > 0 {
			// split tokens always take the damaged tokens
			splitToken.Damage = sourceToken.Damage

			// move all damaged tokens to the split
			// if all damaged tokens are moved, leave none on th emain stack
			if splitToken.Quantity >= sourceToken.QuantityDamaged {
				splitToken.QuantityDamaged = sourceToken.QuantityDamaged
				sourceToken.QuantityDamaged = 0
				sourceToken.Damage = 0
			} else {
				splitToken.QuantityDamaged = MinInt(sourceToken.QuantityDamaged, splitToken.Quantity)
				sourceToken.QuantityDamaged -= splitToken.QuantityDamaged
			}
		}

		sourceToken.Quantity -= splitToken.Quantity
	}

	// remove any empty tokens
	updatedTokens := make([]ShipToken, 0, 1)
	for _, token := range source.Tokens {
		if token.Quantity > 0 {
			updatedTokens = append(updatedTokens, token)
		}
	}
	source.Tokens = updatedTokens

	// update fleet specs
	fleet.Spec = ComputeFleetSpec(rules, player, &fleet)
	source.Spec = ComputeFleetSpec(rules, player, source)
	return &fleet, nil
}

// split all fleets
func (o *orders) SplitAll(rules *Rules, player *Player, playerFleets []*Fleet, source *Fleet) ([]*Fleet, error) {
	// call split on each of these tokens
	newFleets := []*Fleet{}
	for index := 0; index < len(source.Tokens); index++ {
		token := &source.Tokens[index]
		startingTokenIndex := 0

		// don't split the first token
		if index == 0 {
			if token.Quantity == 1 {
				continue
			}
			// we have more than one ship in the first token, skip the first one
			// when splitting, it will become our new source fleet
			startingTokenIndex = 1
		}

		tokenQuantity := token.Quantity
		for tokenIndex := startingTokenIndex; tokenIndex < tokenQuantity; tokenIndex++ {
			splitToken := *token
			splitToken.Quantity = 1
			fleet, err := o.SplitFleetTokens(rules, player, playerFleets, source, []ShipToken{splitToken})
			if err != nil {
				return nil, err
			}
			newFleets = append(newFleets, fleet)
			// add this new fleet to our player fleets so the fleet num increment goes up
			playerFleets = append(playerFleets, fleet)
		}
		index--
	}

	return newFleets, nil
}

// merge fleets into a single fleet
func (o *orders) Merge(rules *Rules, player *Player, fleets []*Fleet) (*Fleet, error) {
	if len(fleets) <= 1 {
		return nil, fmt.Errorf("no fleets to merge")
	}

	// make sure all our fleets have designs for each token
	player.InjectDesigns(fleets)

	fleet := fleets[0]

	tokensByDesign := map[int]*ShipToken{}
	for i := range fleet.Tokens {
		token := &fleet.Tokens[i]
		tokensByDesign[token.DesignNum] = token
	}

	src := fleets[0].Name
	dest := make([]string, 0, len(fleets)-1)

	for i := 1; i < len(fleets); i++ {
		mergingFleet := fleets[i]
		dest = append(dest, mergingFleet.Name)

		for _, token := range mergingFleet.Tokens {
			if t, found := tokensByDesign[token.DesignNum]; found {
				t.Quantity += token.Quantity
			} else {
				fleet.Tokens = append(fleet.Tokens, token)
				tokensByDesign[token.DesignNum] = &fleet.Tokens[len(fleet.Tokens)-1]
			}

			if token.QuantityDamaged > 0 {
				mergingTokenTotalDamage := float64(token.QuantityDamaged) * token.Damage
				// the token we're merging in has damage
				// figure out the total and add it to the fleet we're merging into
				t := tokensByDesign[token.DesignNum]
				destTokenTotalDamage := float64(t.QuantityDamaged) * t.Damage

				// if we merge 2 damaged tokens in 1 damaged token, split the damage
				// between the 3 damaged tokens
				t.QuantityDamaged += token.QuantityDamaged
				t.Damage = (mergingTokenTotalDamage + destTokenTotalDamage) / float64(t.QuantityDamaged)
			}

		}

		// add cargo and fuel to the dest fleet
		fleet.Cargo = fleet.Cargo.Add(mergingFleet.Cargo)
		fleet.Fuel += mergingFleet.Fuel

		// mark the merging fleet for deletion
		mergingFleet.Delete = true
	}

	log.Info().
		Int64("GameID", player.GameID).
		Int("PlayerNum", player.Num).
		Str("Source", src).
		Strs("Dest", dest).
		Msg("merged fleet")

	fleet.Spec = ComputeFleetSpec(rules, player, fleet)
	fleet.MarkDirty()

	return fleet, nil
}

package cs

import (
	"fmt"
	"math"

	"github.com/rs/zerolog/log"
)

type CargoTransferRequest struct {
	Cargo `tstype:",extends"`
	Fuel  int `json:"fuel,omitempty"`
}

type SplitFleetRequest struct {
	// The source fleet to split tokens from
	Source *Fleet `json:"sourcefleet,omitempty"`
	// an optional destination fleet to give tokens to. If nil a new fleet will be crated
	Dest *Fleet `json:"destfleet,omitempty"`

	// a matching slice of source and dest tokens that only differ in token.Quantity
	SourceTokens []ShipToken `json:"sourceTokens,omitempty"`
	DestTokens   []ShipToken `json:"destTokens,omitempty"`

	// a name for the dest fleet, if it is newly created
	DestBaseName string `json:"destBaseName,omitempty"`

	// the amount of cargo to transfer from the source fleet to the dest when splitting
	TransferAmount CargoTransferRequest `json:"transferAmount,omitempty"`
}

// The Orderer interface is used to handle any game logic with updating orders. This is used for
// updating planet and fleet psecs after cargo transfer, splitting and merging fleets, updating research, etc.
type Orderer interface {
	UpdatePlayerOrders(player *Player, playerPlanets []*Planet, order PlayerOrders, rules *Rules)
	UpdatePlanetOrders(rules *Rules, player *Player, planet *Planet, orders PlanetOrders, playerPlanets []*Planet) error
	UpdateFleetOrders(player *Player, fleet *Fleet, orders FleetOrders)
	UpdateMineFieldOrders(player *Player, minefield *MineField, orders MineFieldOrders) error
	TransferByHand(rules *Rules, player *Player, fleet *Fleet, dest CargoHolder, transferAmount CargoTransferRequest) error
	SplitFleet(rules *Rules, player *Player, playerFleets []*Fleet, request SplitFleetRequest) (source, dest *Fleet, err error)
	SplitAll(rules *Rules, player *Player, playerFleets []*Fleet, source *Fleet) ([]*Fleet, error)
	Merge(rules *Rules, player *Player, fleets []*Fleet) (*Fleet, error)
}

type orders struct {
}

func NewOrderer() Orderer {
	return &orders{}
}

func (ctr CargoTransferRequest) Negative() CargoTransferRequest {
	return CargoTransferRequest{Cargo: ctr.Cargo.Negative(), Fuel: -ctr.Fuel}
}

func (ctr CargoTransferRequest) HasNegative() bool {
	return ctr.Cargo.HasNegative() || ctr.Fuel < 0
}

func (ctr CargoTransferRequest) HasPositive() bool {
	return ctr.Ironium > 0 || ctr.Boranium > 0 || ctr.Germanium > 0 || ctr.Fuel > 0 || ctr.Colonists > 0
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

	log.Info().
		Int64("GameID", player.GameID).
		Int("PlayerNum", player.Num).
		Str("Player", player.Name).
		Msg("update player orders")

}

// update a planet orders
func (o *orders) UpdatePlanetOrders(rules *Rules, player *Player, planet *Planet, orders PlanetOrders, playerPlanets []*Planet) error {
	planet.PlanetOrders = orders
	o.updatePlanetSpec(rules, player, planet)

	// update the current planet in our list of planets
	for i, playerPlanet := range playerPlanets {
		if playerPlanet.Num == planet.Num {
			playerPlanets[i] = planet
			break
		}
	}

	// update the player spec with the change in resources for this planet
	// if we turned on/off Contribute Only Leftover Resources to Research, the amount this planet contributes to research goes up
	player.Spec.PlayerResearchSpec = ComputePlayerResearchSpec(player, rules, playerPlanets)
	planet.MarkDirty()

	log.Info().
		Int64("GameID", player.GameID).
		Int("PlayerNum", player.Num).
		Str("Planet", planet.Name).
		Interface("Orders", orders).
		Msg("update planet orders")

	return nil
}

// update a planet spec after some change
func (o *orders) updatePlanetSpec(rules *Rules, player *Player, planet *Planet) error {

	// make sure if we have a starbase, it has a design so we can compute
	// upgrade costs
	if err := planet.PopulateStarbaseDesign(player); err != nil {
		return err
	}

	planet.Spec = computePlanetSpec(rules, player, planet)
	if err := planet.PopulateProductionQueueDesigns(player); err != nil {
		return err
	}

	// make sure we can actually build this stuff
	for _, item := range planet.ProductionQueue {
		if item.design != nil && item.design.OriginalPlayerNum != None {
			return fmt.Errorf("cannot build %s, it was transferred from another player", item.design.Name)
		}
	}

	if err := planet.PopulateProductionQueueEstimates(rules, player); err != nil {
		return fmt.Errorf("planet %s unable to populate queue estimates: %w", planet.Name, err)
	}

	return nil
}

// update the orders to a fleet
func (o *orders) UpdateFleetOrders(player *Player, fleet *Fleet, orders FleetOrders) {
	// copy user modifiable things to the fleet fleet
	fleet.RepeatOrders = orders.RepeatOrders
	fleet.BattlePlanNum = orders.BattlePlanNum

	wp0 := &fleet.Waypoints[0]
	newWP0 := orders.Waypoints[0]

	// TODO: do we want to lookup the target?
	wp0.WarpSpeed = newWP0.WarpSpeed
	wp0.Task = newWP0.Task
	wp0.TransportTasks = newWP0.TransportTasks
	wp0.LayMineFieldDuration = newWP0.LayMineFieldDuration
	wp0.PatrolRange = newWP0.PatrolRange
	wp0.PatrolWarpSpeed = newWP0.PatrolWarpSpeed
	wp0.WaitAtWaypoint = newWP0.WaitAtWaypoint
	wp0.TargetName = newWP0.TargetName
	wp0.TargetType = newWP0.TargetType
	wp0.TargetNum = newWP0.TargetNum
	wp0.TargetPlayerNum = newWP0.TargetPlayerNum
	wp0.TransferToPlayer = newWP0.TransferToPlayer

	fleet.Waypoints = append(fleet.Waypoints[:1], orders.Waypoints[1:]...)
	if len(fleet.Waypoints) > 1 {
		fleet.Heading = (fleet.Waypoints[1].Position.Subtract(fleet.Position)).Normalized()
	}

	fleet.computeFuelUsage(player)

	log.Info().
		Int64("GameID", player.GameID).
		Int("PlayerNum", player.Num).
		Str("Fleet", fleet.Name).
		Interface("Orders", orders).
		Msg("update fleet orders")

}

func (o *orders) UpdateMineFieldOrders(player *Player, minefield *MineField, orders MineFieldOrders) error {
	if !player.Race.Spec.CanDetonateMineFields {
		return fmt.Errorf("%s cannot detonate minefields", player.Race.PluralName)
	}
	if !minefield.MineFieldType.CanDetonate() {
		return fmt.Errorf("%s minefields cannot detonate", minefield.MineFieldType)
	}

	minefield.MineFieldOrders = orders
	return nil
}

// TransferByHand does a by hand transfer of cargo to/from a dest
// Note: this will allow "illegal" orders like trying to steal cargo without tech. The client
// is expected to prevent this where possible for by hand transfers and if not, the player will
// get a message saying their by hand transfer was unsuccessful when it is resolved at turn generation side.
func (o *orders) TransferByHand(rules *Rules, player *Player, fleet *Fleet, dest CargoHolder, transferAmount CargoTransferRequest) error {

	var destName string
	if dest == nil {
		// create a new jettison from any existing jettison
		dest = newJettison(fleet.Position, player.getByHandTransfer(MapObjectTarget{TargetPosition: fleet.Position}))
		destName = "jettison"
	} else {
		destName = dest.GetMapObject().Name
	}

	if fleet.availableCargoSpace() < transferAmount.Total() {
		return fmt.Errorf("fleet %s has %d cargo space available, cannot transfer %dkT from %s", fleet.Name, fleet.availableCargoSpace(), transferAmount.Total(), destName)
	}

	if !dest.CanTransfer(fleet, transferAmount) {
		return fmt.Errorf("fleet %s cannot transfer %v from %s, the dest does not have the required cargo or does not allow this transfer", fleet.Name, transferAmount, destName)
	}

	if !fleet.CanTransfer(fleet, transferAmount.Negative()) {
		return fmt.Errorf("fleet %s cannot transfer %v to %s, the fleet does not have enough the required cargo", fleet.Name, transferAmount.Negative(), destName)
	}

	if dest.GetCargoCapacity() != Infinite && Clamp(dest.GetCargoCapacity()-dest.GetCargo().Total(), 0, dest.GetCargoCapacity()) < -transferAmount.Total() {
		return fmt.Errorf("dest %s has %d cargo space available, cannot transfer %dkT from %s", destName, dest.GetCargoCapacity(), transferAmount.Total(), destName)
	}

	if fleet.availableFuelSpace() < transferAmount.Fuel {
		return fmt.Errorf("fleet %s has %d fuel space available, cannot transfer %dmg from %s", fleet.Name, fleet.availableFuelSpace(), transferAmount.Fuel, destName)
	}

	if dest.GetFuelCapacity() != Infinite && Clamp(dest.GetFuelCapacity()-dest.GetFuel(), 0, dest.GetFuelCapacity()) < -transferAmount.Fuel {
		return fmt.Errorf("dest %s has %d fuel space available, cannot transfer %dmg from %s", destName, dest.GetFuelCapacity(), transferAmount.Fuel, destName)
	}

	// transfer the cargo
	initialFleetCargo := fleet.Cargo
	initialDestCargo := dest.GetCargo()
	fleet.Cargo = fleet.Cargo.Add(transferAmount.Cargo)
	fleet.Fuel += transferAmount.Fuel
	dest.SetCargo(dest.GetCargo().Subtract(transferAmount.Cargo))
	switch t := dest.(type) {
	case *Fleet:
		t.Fuel -= transferAmount.Fuel
		if t.PlayerNum == player.Num {
			t.Spec = ComputeFleetSpec(rules, player, t)
		}
	case *FleetIntel:
		t.Fuel -= transferAmount.Fuel
	}

	// record this call with the player
	target := dest.GetMapObject().ToTarget()
	player.transferByHand(fleet, target, transferAmount.Cargo.Negative())

	log.Info().
		Int64("GameID", player.GameID).
		Int("PlayerNum", player.Num).
		Str("Fleet", fleet.Name).
		Str("InitialFleetCargo", initialFleetCargo.PrettyString()).
		Str("InitialDestCargo", initialDestCargo.PrettyString()).
		Str("Cargo", fleet.Cargo.PrettyString()).
		Str("Transfer", transferAmount.Cargo.PrettyString()).
		Str("DestCargo", dest.GetCargo().PrettyString()).
		Str("Dest", destName).
		Msg("by hand transfer")

	fleet.Spec = ComputeFleetSpec(rules, player, fleet)
	fleet.computeFuelUsage(player)

	// update the spec of the dest if we own it
	switch t := dest.(type) {
	case *Planet:
		if t.OwnedBy(player.Num) {
			o.updatePlanetSpec(rules, player, t)
		}
	case *Fleet:
		if t.OwnedBy(player.Num) {
			t.Spec = ComputeFleetSpec(rules, player, t)
			t.computeFuelUsage(player)
		}
	}

	return nil
}

// split a fleet into two fleets based on a request
func (o *orders) SplitFleet(rules *Rules, player *Player, playerFleets []*Fleet, request SplitFleetRequest) (source, dest *Fleet, err error) {
	source = request.Source
	dest = request.Dest

	// validate the request
	if source == nil {
		return nil, nil, fmt.Errorf("no source fleet to split")
	}

	// build a map of tokens by design
	tokensByDesign := map[int]*ShipToken{}
	stackDamageByDesign := map[int]int{}
	for _, token := range source.Tokens {
		tokensByDesign[token.DesignNum] = &token
		stackDamageByDesign[token.DesignNum] += int(token.Damage) * token.QuantityDamaged
	}
	if dest != nil {
		for _, token := range dest.Tokens {
			stackDamageByDesign[token.DesignNum] += int(token.Damage) * token.QuantityDamaged
			if t, found := tokensByDesign[token.DesignNum]; found {
				t.Quantity += token.Quantity
				t.QuantityDamaged += token.QuantityDamaged
			} else {
				tokensByDesign[token.DesignNum] = &token
			}
		}
	}

	// build a map of the split request tokens by design
	splitTokensByDesign := map[int]*ShipToken{}
	splitStackDamageByDesign := map[int]int{}
	for _, token := range request.SourceTokens {
		splitTokensByDesign[token.DesignNum] = &token
		splitStackDamageByDesign[token.DesignNum] += int(token.Damage) * token.QuantityDamaged
	}
	for _, token := range request.DestTokens {
		splitStackDamageByDesign[token.DesignNum] += int(token.Damage) * token.QuantityDamaged
		if t, found := splitTokensByDesign[token.DesignNum]; found {
			t.Quantity += token.Quantity
			t.QuantityDamaged += token.QuantityDamaged
		} else {
			splitTokensByDesign[token.DesignNum] = &token
		}
	}

	if len(tokensByDesign) != len(splitTokensByDesign) {
		return nil, nil, fmt.Errorf("source fleet tokens and split request tokens don't match")
	}

	for designNum, token := range tokensByDesign {
		splitToken := splitTokensByDesign[designNum]
		if splitToken == nil {
			return nil, nil, fmt.Errorf("found token in original fleets but not in split request")
		}

		// we might lose a point of damage in a split/merge, that's ok
		stackDamage := stackDamageByDesign[token.DesignNum]
		splitStackDamage := splitStackDamageByDesign[token.DesignNum]
		if splitToken.Quantity != token.Quantity || splitToken.QuantityDamaged != token.QuantityDamaged || Abs(stackDamage-splitStackDamage) > 1 {
			return nil, nil, fmt.Errorf("token in original fleet has different quantity/damage that token in split request")
		}
	}

	// create a new dest fleet if dest is nil
	if dest == nil {
		// create a new fleet
		// now create the new fleet
		fleetNum := player.GetNextFleetNum(playerFleets)
		baseName := source.BaseName
		if request.DestBaseName != "" {
			baseName = request.DestBaseName
		}
		fleet := NewFleet(player, fleetNum, baseName, source.Waypoints)
		fleet.OrbitingPlanetNum = source.OrbitingPlanetNum
		fleet.Heading = source.Heading
		fleet.WarpSpeed = source.WarpSpeed
		fleet.PreviousPosition = source.PreviousPosition
		fleet.FleetOrders = source.FleetOrders

		// create a slice of empty tokens we will populate
		fleet.Tokens = make([]ShipToken, len(source.Tokens))

		dest = fleet
	}

	if !source.CanTransfer(dest, request.TransferAmount.Negative()) {
		return nil, nil, fmt.Errorf("source cannot transfer %v to new fleet, the fleet does not have enough of the required cargo", request.TransferAmount.Negative())
	}

	// update the tokens for each fleet
	source.Tokens = request.SourceTokens
	dest.Tokens = request.DestTokens

	// remove any empty tokens
	tokens := []ShipToken{}
	for _, token := range source.Tokens {
		if token.Quantity > 0 {
			tokens = append(tokens, token)
		}
	}
	source.Tokens = tokens

	tokens = []ShipToken{}
	for _, token := range dest.Tokens {
		if token.Quantity > 0 {
			tokens = append(tokens, token)
		}
	}
	dest.Tokens = tokens

	// update fleet specs
	player.InjectDesigns([]*Fleet{source, dest})
	source.Spec = ComputeFleetSpec(rules, player, source)
	dest.Spec = ComputeFleetSpec(rules, player, dest)

	// transfer the cargo as per the player's request
	if err = o.TransferByHand(rules, player, source, dest, request.TransferAmount); err != nil {
		return nil, nil, err
	}

	// split any immediate cargo transfers we did before based on capacity
	if err := player.CargoTransfers.splitByHandTransfers(source, dest); err != nil {
		return nil, nil, err
	}

	if len(source.Tokens) == 0 {
		source.Delete = true
	}
	if len(dest.Tokens) == 0 {
		dest.Delete = true
	}

	log.Info().
		Int64("GameID", player.GameID).
		Int("PlayerNum", player.Num).
		Str("Source", source.Name).
		Str("Dest", dest.Name).
		Msg("split fleet")

	return source, dest, nil
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
			fleet, err := o.splitFleetTokens(rules, player, playerFleets, source, []ShipToken{splitToken})
			if err != nil {
				return nil, err
			}
			newFleets = append(newFleets, fleet)
			// add this new fleet to our player fleets so the fleet num increment goes up
			playerFleets = append(playerFleets, fleet)
		}
		index--
	}

	log.Info().
		Int64("GameID", player.GameID).
		Int("PlayerNum", player.Num).
		Str("Fleet", source.Name).
		Msg("split all fleet")

	return newFleets, nil
}

// split a fleet's tokens into a new fleet
func (o *orders) splitFleetTokens(rules *Rules, player *Player, playerFleets []*Fleet, source *Fleet, tokens []ShipToken) (*Fleet, error) {
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
	fleetNum := player.GetNextFleetNum(playerFleets)

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
	fleet := newFleetForDesign(player, baseDesign, 1, fleetNum, baseDesign.Name, source.Waypoints)
	fleet.OrbitingPlanetNum = source.OrbitingPlanetNum
	fleet.Heading = source.Heading
	fleet.WarpSpeed = source.WarpSpeed
	fleet.PreviousPosition = source.PreviousPosition
	fleet.BattlePlanNum = source.BattlePlanNum
	fleet.Tokens = tokens
	fleet.FleetOrders = source.FleetOrders

	// for splitting fuel and cargo
	originalCargoCapacity := source.Spec.CargoCapacity
	originalFuelCapacity := source.Spec.FuelCapacity
	destCargoCapacity := 0
	destFuelCapacity := 0

	// now remove all the tokens from the old fleet
	for i := range tokens {
		splitToken := &tokens[i]
		sourceToken := sourceTokensByDesign[splitToken.DesignNum]

		// track how much cargo capacity our destination will get
		destCargoCapacity += splitToken.design.Spec.CargoCapacity * splitToken.Quantity
		destFuelCapacity += splitToken.design.Spec.FuelCapacity * splitToken.Quantity

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
				splitToken.QuantityDamaged = min(sourceToken.QuantityDamaged, splitToken.Quantity)
				sourceToken.QuantityDamaged -= splitToken.QuantityDamaged
			}
		}

		sourceToken.Quantity -= splitToken.Quantity
	}

	fuel1, fuel2, err := splitValues(originalFuelCapacity, originalFuelCapacity-destFuelCapacity, destFuelCapacity, source.Fuel)
	if err != nil {
		return nil, fmt.Errorf("unable to split fuel %w", err)
	}
	source.Fuel = fuel1[0]
	fleet.Fuel = fuel2[0]

	// split cargo
	if originalCargoCapacity > 0 {
		source.Cargo, fleet.Cargo, err = source.Cargo.Split(originalCargoCapacity, originalCargoCapacity-destCargoCapacity, destCargoCapacity)
		if err != nil {
			return nil, fmt.Errorf("unable to split cargo %w", err)
		}
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

	// split any immediate cargo transfers as well
	if err := player.CargoTransfers.splitByHandTransfers(source, &fleet); err != nil {
		return nil, fmt.Errorf("unable to split immediate cargo transfers %w", err)
	}

	// update fuel usage estimates
	source.computeFuelUsage(player)
	fleet.computeFuelUsage(player)

	return &fleet, nil
}

// merge fleets into a single fleet
func (o *orders) Merge(rules *Rules, player *Player, fleets []*Fleet) (*Fleet, error) {
	if len(fleets) <= 1 {
		return nil, fmt.Errorf("no fleets to merge")
	}

	// make sure all our fleets have designs for each token
	player.InjectDesigns(fleets)

	fleet := fleets[0]

	src := fleets[0].Name
	dest := make([]string, 0, len(fleets)-1)

	for i := 1; i < len(fleets); i++ {
		mergingFleet := fleets[i]
		dest = append(dest, mergingFleet.Name)

		for _, token := range mergingFleet.Tokens {
			existingToken := fleet.getTokenByDesign(token.DesignNum)
			if existingToken == nil {
				// we don't have a token for this design yet, so add it
				fleet.Tokens = append(fleet.Tokens, token)
				existingToken = &fleet.Tokens[len(fleet.Tokens)-1]
			} else {
				existingToken.Quantity += token.Quantity
			}

			if token.QuantityDamaged > 0 {
				mergingTokenTotalDamage := token.QuantityDamaged * int(token.Damage)
				// the token we're merging in has damage
				// figure out the total and add it to the fleet we're merging into
				destTokenTotalDamage := existingToken.QuantityDamaged * int(existingToken.Damage)

				// if we merge 2 damaged tokens in 1 damaged token, split the damage
				// between the 3 damaged tokens
				existingToken.QuantityDamaged += token.QuantityDamaged
				existingToken.Damage = math.Floor(float64(mergingTokenTotalDamage+destTokenTotalDamage) / float64(existingToken.QuantityDamaged))
			}
		}

		// add cargo and fuel to the dest fleet
		fleet.Cargo = fleet.Cargo.Add(mergingFleet.Cargo)
		fleet.Fuel += mergingFleet.Fuel

		// mark the merging fleet for deletion
		mergingFleet.Delete = true
	}

	// merge cargo transfers
	player.CargoTransfers.mergeByHandTransfers(fleet, fleets)

	log.Info().
		Int64("GameID", player.GameID).
		Int("PlayerNum", player.Num).
		Str("Source", src).
		Strs("Dest", dest).
		Msg("merged fleet")

	fleet.Spec = ComputeFleetSpec(rules, player, fleet)
	fleet.computeFuelUsage(player)

	return fleet, nil
}

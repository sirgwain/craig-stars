package cs

// The techTrader interface handles checks for tech level increases from trading
type techTrader interface {
	// Perform an invasion tech trading check, if the player hasn't traded for tech yet this year.
	//
	// Wrapper function for techLevelGained
	checkInvasionTechTrade(rules *Rules, player *Player, targetLevel TechLevel) (field TechField)

	// Check fleet-based tech trading for both parts & levels at once.
	//
	// Returns the field or part gained from this trade instance, if any
	checkFleetTechTrade(rules *Rules, player *Player, tokens []ShipToken) (field TechField, acquiredPart *Tech)

	// Check for a tech level boost for a player's tech level and some target we scrapped, invaded, etc.
	//
	// Returns the field gained from this trade instance, if any
	techLevelGained(rules *Rules, current, target TechLevel) TechField

	// Checks for acquirable part gain from this tech trade event
	//
	// Returns the part gained from this trade instance, if any
	acquirablePartGained(rules *Rules, player *Player, tokens []ShipToken) *Tech
}

type techTrade struct{}

func newTechTrader() techTrader {
	return &techTrade{}
}

// checkInvasionTechTrade performs an invasion tech trading check, if the player hasn't traded for tech yet this year.
//
// Wrapper function for techLevelGained
func (t *techTrade) checkInvasionTechTrade(rules *Rules, player *Player, targetLevel TechLevel) TechField {
	if player.techLevelGained || player.acquirablePartGained {
		return TechFieldNone
	}
	return t.techLevelGained(rules, player.TechLevels, targetLevel)
}

// checkFleetTechTrade checks a scrapped fleet for a tech trade. After one
// shared eligibility roll, an acquirable part is checked first, then a tech level
// from the fleet's designs.
//
// Returns the field or part gained from this trade instance, if any
func (t *techTrade) checkFleetTechTrade(rules *Rules, player *Player, tokens []ShipToken) (TechField, *Tech) {
	if len(tokens) == 0 || player.techLevelGained || player.acquirablePartGained || rules.random.Float64() >= rules.TechTradeChance {
		return TechFieldNone, nil
	}
	if part := t.acquirablePartGained(rules, player, tokens); part != nil {
		return TechFieldNone, part
	}
	target := TechLevel{}
	for _, token := range tokens {
		target = target.Max(token.design.Spec.TechLevel)
	}
	return t.chooseField(rules, player.TechLevels, target), nil
}

// techLevelGained checks for a tech level boost for a player's tech level and some target we scrapped, invaded, etc.
//
// Returns the field gained from this trade instance, if any
func (t *techTrade) techLevelGained(rules *Rules, current, target TechLevel) TechField {
	if rules.random.Float64() >= rules.TechTradeChance {
		return TechFieldNone
	}
	return t.chooseField(rules, current, target)
}

// chooseField makes six random field draws, with replacement, returning the first
// field where the target is ahead of the player. The size of the level gap
// doesn't change the odds, only which fields qualify.
func (t *techTrade) chooseField(rules *Rules, current, target TechLevel) TechField {
	for range len(TechFields) {
		field := TechFields[rules.random.Intn(len(TechFields))]
		if current.Get(field) < min(target.Get(field), rules.MaxTechLevel) {
			return field
		}
	}
	return TechFieldNone
}

// acquirablePartGained makes 13 random draws from the trader part slots,
// checking any part these tokens carry that the player hasn't acquired yet.
//
// Returns the part gained from this trade instance, if any
func (t *techTrade) acquirablePartGained(rules *Rules, player *Player, tokens []ShipToken) *Tech {
	if player.techLevelGained || player.acquirablePartGained {
		return nil
	}
	// Keep the original 13 trader slots, including empty/non-component slots.
	parts := []*Tech{&MultiCargoPod.Tech, &MultiFunctionPod.Tech, &LangstonShell.Tech, &MegaPolyShell.Tech, &AlienMiner.Tech, &HushABoom.Tech, &AntiMatterTorpedo.Tech, &MultiContainedMunition.Tech, &MiniMorph.Tech, &EnigmaPulsar.Tech, &GenesisDevice.Tech, &JumpGate.Tech, nil}
	counts := make(map[string]int)
	for _, token := range tokens {
		for _, slot := range token.design.Slots {
			counts[slot.HullComponent] += slot.Quantity
		}
	}
	for range len(parts) {
		part := parts[rules.random.Intn(len(parts))]
		if part == nil || player.AcquiredTechs[part.Name] || counts[part.Name] == 0 {
			continue
		}
		if checkAcquirablePartChance(rules, counts[part.Name]) {
			return part
		}
	}
	return nil
}

// checkAcquirablePartChance checks if a part trade succeeds, with a chance per installed part, up to a max quantity
func checkAcquirablePartChance(rules *Rules, qty int) bool {
	return qty > 0 && rules.random.Float64() < rules.AcquirablePartTradeChanceBase*float64(min(qty, rules.AcquirablePartTradeItemMax))
}

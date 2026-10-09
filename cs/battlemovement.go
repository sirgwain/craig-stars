package cs

import "math"

// battleDamageScore ranks predicted damage according to the tactic; lower scores are better.
func battleDamageScore(given, taken float64, tactic BattleTactic) float64 {
	switch tactic {
	case BattleTacticDisengage, BattleTacticMinimizeDamageToSelf:
		return taken
	case BattleTacticMaximizeNetDamage, BattleTacticMaximizeDamageRatio:
		if given == 0 {
			return taken
		}
		return -given / (taken + 1)
	default:
		return -given
	}
}

// movementDamage estimates the source's damage at a distance for movement planning.
// Proximity scoring includes reduced threats from weapons beyond their firing range.
func (b *battle) movementDamage(source, target *battleToken, distance int, proximity bool) float64 {
	if !source.isStillInBattle() {
		return 0
	}
	damage := 0.0
	for _, weapon := range source.weaponSlots {
		outOfRange := distance > weapon.weaponRange
		if outOfRange && !proximity {
			continue
		}
		power := float64(weapon.power * weapon.slotQuantity)
		if weapon.weaponType == battleWeaponTypeBeam {
			power *= weapon.beamBonus
			if weapon.weaponRange > 0 {
				power *= max(0, 1-float64(distance)/float64(weapon.weaponRange)*b.rules.BeamRangeDropoff)
			}
			power *= target.beamDamageMultiplier()
			if weapon.damagesShieldsOnly {
				power = min(power, float64(target.stackShields)*float64(source.Quantity)/float64(max(1, target.Quantity)))
			}
			if outOfRange {
				power = max(float64(weapon.slotQuantity), power/float64(distance+10-weapon.weaponRange))
			}
			damage += power * float64(source.Quantity)
		} else {
			accuracy := weapon.getAccuracy(target.torpedoJamming)
			power *= float64(source.Quantity) * accuracy
			if target.stackShields > 0 {
				power += float64(weapon.power*weapon.slotQuantity*source.Quantity) * (1 - accuracy) / 8
			}
			if outOfRange {
				power = max(float64(weapon.slotQuantity), power/float64(distance+10-weapon.weaponRange))
			}
			damage += power
		}
	}
	if !proximity {
		defense := float64(target.stackShields+target.armor*target.Quantity) - target.Damage*float64(target.QuantityDamaged)
		damage = min(damage, max(1, defense))
	}
	return damage
}

// movementTargetType chooses the primary target class while any matching hostile
// token remains active, otherwise choosing the secondary class.
func (b *battle) movementTargetType(token *battleToken) BattleTarget {
	for _, enemy := range b.tokens {
		if enemy.isStillInBattle() && token.willAttack(enemy.PlayerNum) && enemy.isTargetOf(token.PrimaryTarget) {
			return token.PrimaryTarget
		}
	}
	return token.SecondaryTarget
}

// movementSearch returns the planning radius and, when targets are beyond reach,
// the nearest matching enemy that the token can damage.
func (b *battle) movementSearch(token *battleToken, targetType BattleTarget) (int, *battleToken) {
	moves := max(1, token.movesLeft)
	minRange, maxRange, nonSapperRange := math.MaxInt, 0, -1
	for _, weapon := range token.weaponSlots {
		minRange = min(minRange, weapon.weaponRange)
		maxRange = max(maxRange, weapon.weaponRange)
		if weapon.weaponType == battleWeaponTypeBeam && !weapon.damagesShieldsOnly {
			nonSapperRange = max(nonSapperRange, weapon.weaponRange)
		}
	}
	rangeToConsider := maxRange
	if token.Tactic == BattleTacticMaximizeDamage || token.Tactic == BattleTacticMaximizeNetDamage {
		if minRange != math.MaxInt {
			rangeToConsider = minRange
		}
	}
	var closest *battleToken
	closestDistance := math.MaxInt
	for _, enemy := range b.tokens {
		if !enemy.isStillInBattle() || !token.willAttack(enemy.PlayerNum) || !enemy.isTargetOf(targetType) {
			continue
		}
		distance := token.getDistanceAway(enemy.Position)
		if enemy.movesLeft >= moves {
			distance++
		}
		canReach := distance <= rangeToConsider+moves
		if maxRange == 3 && enemy.stackShields == 0 && nonSapperRange < maxRange {
			canReach = canReach && distance <= nonSapperRange+moves
		}
		if canReach && len(token.weaponSlots) > 0 {
			return moves, nil
		}
		if distance < closestDistance && b.movementDamage(token, enemy, 0, false) > 0 {
			closest, closestDistance = enemy, distance
		}
	}
	return 1, closest
}

// movementScore evaluates a position against predicted enemy responses and the
// token's tactic, including proximity threats and friendly occupancy during retreat.
func (b *battle) movementScore(token *battleToken, position Vector, targetType BattleTarget) float64 {
	flee := token.Tactic == BattleTacticDisengage
	givenBest, takenTotal := 0.0, 0.0
	for _, enemy := range b.tokens {
		if !enemy.isStillInBattle() || !token.willAttack(enemy.PlayerNum) {
			continue
		}
		distance := position.chebyshevDistance(enemy.Position)
		minDistance, maxDistance := distance, distance
		if enemy.movesLeft >= max(1, token.movesLeft) {
			minDistance = max(0, distance-1)
			for _, dx := range []int{-1, 1} {
				for _, dy := range []int{-1, 1} {
					corner := Vector{Clamp(enemy.Position.X+dx, 0, battleWidth-1), Clamp(enemy.Position.Y+dy, 0, battleHeight-1)}
					maxDistance = max(maxDistance, position.chebyshevDistance(corner))
				}
			}
		}
		bestEnemyScore := math.Inf(1)
		bestGiven, bestTaken := 0.0, 0.0
		for distance := minDistance; distance <= maxDistance; distance++ {
			given := 0.0
			if enemy.isTargetOf(targetType) {
				given = b.movementDamage(token, enemy, distance, false)
			}
			taken := b.movementDamage(enemy, token, distance, flee)
			score := battleDamageScore(taken, given, enemy.Tactic)
			if score <= bestEnemyScore {
				bestEnemyScore, bestGiven, bestTaken = score, given, taken
			}
		}
		givenBest = max(givenBest, bestGiven)
		takenTotal += bestTaken
	}
	score := battleDamageScore(givenBest, takenTotal, token.Tactic)
	if flee {
		for _, friend := range b.tokens {
			if friend.PlayerNum == token.PlayerNum && friend.Position == position {
				score += 2
			}
		}
		if position == token.Position {
			score--
		}
	}
	return score
}

// getBestMoves plans destinations across remaining movement and returns the
// best next steps, breaking damage score ties by position.
func (b *battle) getBestMoves(token *battleToken) []Vector {
	// Attackers break ties toward the center so they don't get pinned in a corner.
	// Retreating tokens break ties toward their current square instead, so the
	// center doesn't pull them back toward the enemy.
	distance := battleCenterDistance
	if token.Tactic == BattleTacticDisengage {
		distance = token.getDistanceAway
	}
	targetType := b.movementTargetType(token)
	radius, closest := b.movementSearch(token, targetType)
	best := newBattleMoveChoice(distance)
	nearScores := make(map[Vector]float64)
	for x := max(0, token.Position.X-radius); x <= min(battleWidth-1, token.Position.X+radius); x++ {
		for y := max(0, token.Position.Y-radius); y <= min(battleHeight-1, token.Position.Y+radius); y++ {
			position := Vector{x, y}
			score := b.movementScore(token, position, targetType)
			if position.chebyshevDistance(token.Position) <= 1 {
				nearScores[position] = score
			}
			best.add(position, score)
		}
	}
	destinations := best.moves
	if closest != nil {
		destinations = []Vector{closest.Position}
	}
	// Plan destinations using the remaining moves, but take only one step now.
	steps := newBattleMoveChoice(distance)
	seen := make(map[Vector]bool)
	for _, destination := range destinations {
		for _, step := range battleStepsToward(token.Position, destination) {
			if !seen[step] {
				seen[step] = true
				steps.add(step, nearScores[step])
			}
		}
	}
	return steps.moves
}

// battleMoveChoice collects the lowest scoring moves, breaking ties with updateBestMoves.
type battleMoveChoice struct {
	moves    []Vector
	score    float64
	distance func(Vector) int
}

func newBattleMoveChoice(distance func(Vector) int) *battleMoveChoice {
	return &battleMoveChoice{moves: []Vector{}, score: math.Inf(1), distance: distance}
}

// add considers a move, ignoring float rounding when comparing scores.
func (c *battleMoveChoice) add(position Vector, score float64) {
	const epsilon = 1e-9
	if score < c.score-epsilon {
		c.score = score
		c.moves = updateBestMoves(true, position, c.moves, c.distance)
	} else if math.Abs(score-c.score) <= epsilon {
		c.moves = updateBestMoves(false, position, c.moves, c.distance)
	}
}

// updateBestMoves adds a position that scored equal to or better than the current
// best moves, breaking score ties by the smallest preference distance:
// if better, take the new position and discard the others
// if equal but nearer, take the new move and discard the others
// if equal and farther, keep the bestMoves as is
// if equal and the same distance, add it to bestMoves
func updateBestMoves(better bool, newPosition Vector, bestMoves []Vector, distance func(Vector) int) []Vector {
	if len(bestMoves) == 0 || better {
		return []Vector{newPosition}
	}
	newDistance := distance(newPosition)
	oldDistance := distance(bestMoves[0])
	if newDistance == oldDistance {
		return append(bestMoves, newPosition)
	}
	if newDistance < oldDistance {
		return []Vector{newPosition}
	}
	return bestMoves
}

// battleCenterDistance returns the squared straight-line distance from the center
// of the board. Coordinates are doubled so the center at (4.5, 4.5) stays an
// integer, and corners are farther from the center than the middle of an edge.
func battleCenterDistance(position Vector) int {
	dx, dy := 2*position.X-(battleWidth-1), 2*position.Y-(battleHeight-1)
	return dx*dx + dy*dy
}

// battleStepsToward returns legal single-step positions leading toward the destination.
func battleStepsToward(position, destination Vector) []Vector {
	dx, dy := destination.X-position.X, destination.Y-position.Y
	if max(Abs(dx), Abs(dy)) <= 1 {
		return []Vector{destination}
	}
	xStep, yStep := 0, 0
	if dx > 0 {
		xStep = 1
	} else if dx < 0 {
		xStep = -1
	}
	if dy > 0 {
		yStep = 1
	} else if dy < 0 {
		yStep = -1
	}
	steps := []Vector{{position.X + xStep, position.Y + yStep}}
	if dx == 0 {
		steps = append(steps, Vector{position.X - 1, position.Y + yStep}, Vector{position.X + 1, position.Y + yStep})
	} else if dy == 0 {
		steps = append(steps, Vector{position.X + xStep, position.Y - 1}, Vector{position.X + xStep, position.Y + 1})
	} else if Abs(dx) > Abs(dy) {
		steps = append(steps, Vector{position.X + xStep, position.Y})
	} else if Abs(dy) > Abs(dx) {
		steps = append(steps, Vector{position.X, position.Y + yStep})
	}
	valid := steps[:0]
	for _, step := range steps {
		if step.X >= 0 && step.X < battleWidth && step.Y >= 0 && step.Y < battleHeight {
			valid = append(valid, step)
		}
	}
	return valid
}

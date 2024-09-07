package cs

import (
	"fmt"
	"math"
	"sort"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// From: https://wiki.starsautohost.org/wiki/Guts_of_the_Battle_Engine
// ===================================================================
// Here are the guts of the battle engine as I understand it from both experience, observation and the help file
// (please pull me up on any points I get wrong)
//
// For a battle to take place 2 or more fleets (or a fleet and a starbase) must be at the same location and at
// least one of the fleets must be armed and have orders to attack ships of the others race (the type of ships
// involved doesn't matter). If are race has a fleet present at a location where there is a battle, but doesn't
// have orders to attack any of the other races there and none of the other races present has orders to attack it
// then it will not take part in the battle (and can not benefit from potential tech gain -- actually you can benefit
// from tech gain, a fact I learned from trying not to get the tech gain in a wolf/lamb tech exchange - LEit).
// Each ship present at the battle will form part of a token (AKA a stack), it is possible to have a token comprised
// of just a single ship. Tokens are always of ships of the same design. Each ship design in each fleet will create
// a token, splitting a few ships off to form a second fleet before the battle will create a second token on the
// battle board.
//
// The battle grid is made up of 10 squares by 10 squares. Each token is in a single square, there can be more
// than one token in the same square.There is an limit of 256 tokens per battle event for all players involved,
// if this limit is exceeded, then excess tokens will be left out (those created from fleets with the highest
// fleet numbers), in such a case each player will have an equal number of tokens, each player will be guaranteed
// to get their "share" of the available token slots (ie in a 4 race battle 256 / 4 = 64 token slots), if a race
// doesn't use up all their "slots" then they are shared equally between the other players.
//
// Each battle is made up of rounds. There are a maximum of 16 rounds in each battle. Each round has two parts,
// movement and shooting. Each token has a speed rating, and will be able to move between 0 and 3 squares in a
// single turn. If a token has a fractional speed rating then they will get a bonus square of movement every set
// number of turns. a 1/4 bonus means an extra square of movement on the first round and then on every fourth round
// after that starting with the fifth. A 1/2 speed bonus gets a bonus square of movement every other turn starting
// with the first, and a 3/4 speed bonus gets a bonus square of movement for the first three rounds of every 4 round
// cycle. The order of movement is this, each token with 3 movement squares moves a single square, then each token
// with 2+ movement moves a single square (if it had speed 3 then it would move for its second square) and then all
// ships with at least one square of movement move again. At each stage the ships with the most weight will move first
// though there is less than a 15% difference in weight then there is a chance that the lighter ship will go first.
// The smaller the weight % difference the greater the chance of the lighter ship going first.
//
// Each token has an attractiveness rating. This is used in both working out where ships move to and which ships are
// shot at first. The essence of the formula is cost / defence. A ship will have different attractiveness ratings
// verses different types of weapons (beams, sappers, torpedoes and capital missiles). Cost is calculated by summing
// the resource and boranium costs of the ship design used (iron and germ costs don't affect the attractiveness rating).
// Defence is calculated by the shield and armour dp modified by the enemies torpedo accuracy (after base accuracy,
// comps and jammers are worked out) if defending vs torps or capital missiles, the effects of double damage for unshielded
// targets vs capital missiles and the effects of deflectors against beam weapons. The attractiveness rating can be
// change during the course of the battle as shields and armour deplete. Attractiveness doesn't take into account the
// one missile one kill rule, thus chaff has become a fairly effective tactic.
//
// battle orders are comprised of 4 parts. A primary and secondary target type, legitimate races to attack and the tactic
// to use in battle. Ships will only attack tokens belonging to legitimate target races, however if another race present
// has any ships (including unarmed ships) with battles orders to attack your race then that race will also be considered
// a legitimate target. When attacking ships will try and shoot the most attractive ship of a type listed as a primary
// target and if no ships are available which are primary targets then the most attractive ship of a type listed as a
// secondary target will be targeted. Ships which are not listed as primary or secondary targets will not get shot at,
// even if they are shooting back.
//
// There are 6 different battle orders which determine the movement AI of the ships in battle, the movement AI is applied
// each time a ship wants to move a square on the battle board.:
//   - Disengage - If there is any enemy ship in firing range then move to any square further away than your current square.
//     If you are in range of an enemy weapon but cannot move further away then try move to a square that is of
//     the same distance away. If you are in range of the enemies weapons and cannot move away or maintain distance
//     then move to a random square. If you are not in range of the enemies weapons then move randomly. Also you
//     will try and disengage which will require 7 squares of movement to be clocked up before you can leave from
//     the battle board.
//   - Disengage if Challenged - Behaves like Maximise Damage until token takes damage and then behaves like Disengage.
//   - Minimise Damage to Self - (Not 100% sure on this one) If within range of an enemy weapon then move away from the
//     enemy (just like Disengage). If out of range of the enemies weapons or cannot move away from
//     the enemy then try and get in range of the best available target without moving towards the enemy.
//   - Maximise Net Damage - Locate most attractive primary target (or secondary if no primary targets are left). If out
//     of range with ANY weapon then move towards target. If in range with all weapons them move as to
//     maximise damage_done/damage_taken. The effect of this is if your weapons are longer range then
//     try to stay at maximum range. If your weapons range is the same then do random movement while
//     staying in range. If your weapons are shorter range and also beam weapons then attempt to close
//     in to zero range.
//   - Maximise Damage Ratio - As Maximise Net Damage but only considers the longest range weapon.
//   - Maximise Damage - Locate most attractive primary target (or secondary if no primary targets are left). If any of
//     your weapons are out of range of that token then keep moving to squares that are closer to it until
//     in range with all weapons. If using any beam weapons (as they have range dissipation) then attempt
//     to close to 0 range. If just using missiles or torps and in range then move randomly to a squares
//     still in range.
//
// Note that there is a bug when fighting starbases, the battle AI doesn't count the +1 range bonus when calculating movement.
// This mainly applies when your ships are attempting to get out of range of the enemy, so vs starbase with range 6 missiles,
// your ships will move to distance 7, the movement AI won't calculate that they are still in range even when they keep getting
// shot at.
// After the movement phase all ships will shoot their weapons, a token will fire all weapons from the same slot in a single
// shot. The weapon slot with the highest initiative will fire first. If there are two ships with slots of the same init,
// then the ships will be randomly given a priority over who can fire first (which will stick for the entire battle). The
// rest of the weapon slots are then fired in init order. Damage is worked out in between each shot and applied to the ships.
// If ships or tokens are destroyed before their turn to shoot then they won't be able to fire back. The movement AI will
// go after the most attractive primary target on the board, but if this token is not in range, then the ship will fire on
// the most attractive primary target within range (or secondary if none available). Starbases have a +1 range bonus to all
// their weapons (this also gets applied to minefield sweeping rates), though cannot move. The movement AI doesn't take this
// bonus into account when moving ships to close in on an enemy starbase.
//
// Damage for each shot is calculated by multiplying the number of weapons in the slot by the number of ships in the token
// by the amount of dp the weapon does. For beam weapons, this damage will dissipate by 10% over the range of the beam
// (for a range 2 beam - no dissipation at range 0, 5% dissipation at range 1 and 10% dissipation at range 2). Also
// capacitors and deflectors will modify the damage actually done to the enemy ship. Damage will be applied first to the
// tokens shield stack and then to armour only when the entire shield stack of the token is down. For missile ships,
// each missile fired will be tested to see if it will hit, the chance to hit is based on the base accuracy, the computers
// on the ship and the enemy jammers. Missiles that miss will do 1/8 of their damage to the shields and won't affect armour.
// For missiles that hit, upto half will be taken by the shields, the rest will go to the armour. For capital missiles
// any damage done after the shields are taken down will do double damage to the armour. Whole ship kills are worked out
// by adding up all the damage done to the armour by a single salvo (from a token's slot) and dividing this by the amount
// of armour each single ship in the token has left (total armour x token damage %). The number of complete ships the shot
// could kill will be removed from the enemy token, the rest of the damage will divided equally among the rest of the ships
// in the token and applied as damage. As token armour is stored in 1/512ths (about 0.2%s) of total armour and not as an
// exact dp figure (shields are stored as an exact figure), there may be some rounding of the damage after each salvo
// (AFAIK its always rounds up). This fact can be abused by creating lots of small fleet tokens with weak missiles
// and many slots, where each slot that hits will do 0.2% damage to the enemy token even if each individual missile
// would do less damage normally (especially the case with a beta torp shooting a large nub stack).
//
// After all the weapons that are in range have fired, the next round begins, starting with ship movement.
// The battle is ended when either the 16 round timer runs out, there is only one race left present on the battle board
// or if there are two or more races which have no hostile intentions towards each other.
//
// After a battle, salvage is created. This is equal to 1/3 of the current mineral costs of all the ships that where
// destroyed during the battle. This is left at the location of the battle and will decay over time, or if the battle
// happened over a planet, then the minerals will get deposited there.
//
// Any races that took part in a battle and had at least one ship that managed to survive (either through surviving
// till the end or retreating beforehand) has a potential to gain tech levels from ships that where destroyed during
// the battle. For the exact details of the formulas and chances involved see the Guts of Tech Trading.
//
// Movement speed and moves per round.
// 3/4 is      1110
// 1 is        1111
// 1 1/4 is    2111
// 1 1/2 is    2121
// 1 3/4 is    2212
// 2 is        2222
// 2 1/4 is    3222
// 2 1/2 is    3232
//
// ===================================================================

const battleWidth, battleHeight = 10, 10

// The battler interface is the main entrypoint into running battles
// A batter is created for fleets in the same location in the universe and it is used
// to check if a battle will occur, and run the battle.
// Running a battle returns a BattleRecord that is passed along to each player.
type battler interface {
	findTargets() bool
	runBattle() *BattleRecord
}

// battle defines the state of a battle as it progresses
type battle struct {
	planet   *Planet
	position Vector
	tokens   []*battleToken
	round    int
	players  map[int]*Player
	board    [battleWidth][battleHeight]int // the number of tokens in each square
	record   *BattleRecord
	rules    *Rules
	log      zerolog.Logger
}

var positionsByPlayer = []BattleVector{
	{1, 4},
	{8, 5},
	{4, 1},
	{5, 8},
	{1, 1},
	{8, 1},
	{8, 8},
	{1, 8},
}

var movementByRound = [9][4]int{
	{1, 0, 1, 0},
	{1, 1, 0, 1},
	{1, 1, 1, 1},
	{2, 1, 1, 1},
	{2, 1, 2, 1},
	{2, 2, 1, 2},
	{2, 2, 2, 2},
	{3, 2, 2, 2},
	{3, 2, 3, 2},
}

// get the movement of this design with additional cargo
func getBattleMovement(idealEngineSpeed, movementBonus, mass, numEngines int) int {
	if numEngines == 0 {
		return 0
	}
	return Clamp(((idealEngineSpeed+movementBonus)-2)-((mass)/numEngines/70), 2, 10)
}

// BuildBattle builds a battle recording with all the battle tokens for a list of fleets that contains more than one player.
// We'll use this to determine if a battle should take place at this location.
// Also, any players that have a potential battle will discover each other's designs.
func newBattler(log zerolog.Logger, rules *Rules, techFinder TechFinder, battleNum int, players map[int]*Player, fleets []*Fleet, planet *Planet) battler {
	battleLogger := log.With().Int("Battle", battleNum).Logger()
	if len(fleets) == 0 {
		battleLogger.Error().Msg("Can't build battle with no fleets.")
		return nil
	}

	// track
	sortedPlayerNums := make([]int, 0, len(players))
	for _, player := range players {
		sortedPlayerNums = append(sortedPlayerNums, player.Num)
	}
	sort.Ints(sortedPlayerNums)

	playerStartingPositions := make(map[int]BattleVector)
	for i, num := range sortedPlayerNums {
		if i >= len(positionsByPlayer) {
			battleLogger.Warn().Msg("Oh noes! We have a battle with more players than we have positions for...")
		}
		playerStartingPositions[num] = positionsByPlayer[i%len(positionsByPlayer)]
	}

	board := [battleWidth][battleHeight]int{}
	// add each fleet's token to the battle
	tokens := []*battleToken{}
	tokenRecords := []BattleRecordToken{}
	num := 0
	dampening := 0
	for _, fleet := range fleets {
		totalCargo := fleet.Cargo.Total()
		totalCargoCapacity := fleet.Spec.CargoCapacity

		player := players[fleet.PlayerNum]

		for i := range fleet.Tokens {
			token := &fleet.Tokens[i]
			num++

			// add cargo from the fleet to each token
			cargoMass := 0
			if totalCargo > 0 && token.design.Spec.CargoCapacity > 0 {
				// see how much this ship's cargo capacity is compared to the fleet total
				shipCargoPercent := float64(token.design.Spec.CargoCapacity) / float64(totalCargoCapacity)
				cargoMass = int(float64(totalCargo) * shipCargoPercent)
			}

			position := playerStartingPositions[player.Num]
			battleToken := newBattleToken(num, position, cargoMass, token, *fleet.battlePlan, player, techFinder)
			tokens = append(tokens, battleToken)
			tokenRecords = append(tokenRecords, battleToken.BattleRecordToken)

			// put this token on the board
			board[position.X][position.Y] += battleToken.StartingQuantity

			// find the highest dampener we have
			dampening = MaxInt(dampening, token.design.Spec.ReduceMovement)
		}
	}

	// apply dampening
	if dampening > 0 {
		for _, token := range tokens {
			// we only dampen movement of ships that move, not starbases (obviously)
			// and we can't go below 2
			if token.Movement > 0 {
				token.Movement = Clamp(token.Movement-dampening, 2, 10)
			}
		}
	}

	planetNum := 0
	if planet != nil {
		planetNum = planet.Num
	}

	battle := &battle{
		planet:   planet,
		position: fleets[0].Position,
		tokens:   tokens,
		record:   newBattleRecord(battleNum, planetNum, fleets[0].Position, tokenRecords),
		players:  players,
		rules:    rules,
		log:      battleLogger,
	}

	return battle
}

// findTargets allocates targets for each token in a battle
// this returns false if no targets are found
func (b *battle) findTargets() bool {
	hasTargets := false
	for _, token := range b.tokens {
		if !token.hasWeapons() || !token.isStillInBattle() {
			continue
		}

		// first determine all the targets for this token's weapons
		token.findWeaponsTargets(b.tokens)

		// if we move, find a move target
		if token.Movement == 0 {
			continue
		}
		token.moveTarget = token.findMoveTarget()
		hasTargets = token.moveTarget != nil || hasTargets
	}

	return hasTargets
}

// runBattle runs a battle!
func (b *battle) runBattle() *BattleRecord {
	if b.planet != nil {
		b.log.Info().Msgf("Running a battle at %s involving %d players and %d tokens.", b.planet.Name, len(b.players), len(b.tokens))
	} else {
		b.log.Info().Msgf("Running a battle at (%.2f, %.2f) involving %d players and %d tokens.", b.position.X, b.position.Y, len(b.players), len(b.tokens))
	}

	// movement order is set at the start of battle and doesn't change
	moveOrder := b.buildMovementOrder(b.tokens)
	for b.round = 1; b.round <= b.rules.NumBattleRounds; b.round++ {

		// each round we build the SortedWeaponSlots list
		// anew to account for ships that were destroyed
		weaponSlots := b.getSortedWeaponSlots(b.tokens)

		// find new targets
		if !b.findTargets() {
			// out of targets, battle is done
			break
		}

		b.record.recordNewRound()

		// movement is a repeating pattern of 4 movement blocks
		// which we figured out in buildMovementOrder
		roundBlock := ((b.round - 1) % 4)
		for _, token := range moveOrder[roundBlock] {
			b.moveToken(token, weaponSlots)
		}

		// iterate over each weapon and fire if they have a target in range
		for _, weaponSlot := range weaponSlots {
			// find all available targets for this weapon
			targets := weaponSlot.getTargetsInRange()
			b.fireWeaponSlot(weaponSlot, targets)
		}

		// at the end of this round, regenerate shields
		for _, token := range b.tokens {
			if token.stackShields > 0 && token.isStillInBattle() {
				token.regenerateShields()
			}
		}
	}

	// record destroyed tokens
	for _, token := range b.tokens {
		if token.quantityDestroyed > 0 {
			b.record.recordDestroyedToken(token, token.quantityDestroyed)
		}
	}

	return b.record
}

// buildMovementOrder builds a list of Movers in this battle.
// Each ship moves in order of mass with heavier ships moving first.
// Ships that can move 3 times in a round move first, then ships that move 2 times, then 1.
func (b *battle) buildMovementOrder(tokens []*battleToken) (moveOrder [4][]*battleToken) {
	// our tokens are moved by mass
	tokensByMass := make([]*battleToken, 0)
	for _, token := range tokens {
		if token.Movement > 0 { // starbases don't move
			tokensByMass = append(tokensByMass, token)
		}
	}
	sort.Slice(tokensByMass, func(i, j int) bool {
		return tokensByMass[i].Mass > tokensByMass[j].Mass
	})

	// each token can move up to 3 times in a round
	// ships that can move 3 times go first, so we loop through the moveNum backwards
	// so that our Movers list has ships that move 3 times first
	for moveNum := 2; moveNum >= 0; moveNum-- {
		// for each block of 4 rounds, add each ship to the movement list if it's supposed to move that round
		for roundBlock := 0; roundBlock < 4; roundBlock++ {
			// add each battle token to the movement for this roundBlock
			for _, token := range tokensByMass {
				// movement is between 2 and 10, so we offset it to fit in our MovementByRound table
				movement := token.Movement

				// see if this token can move on this moveNum (i.e. move 1, 2, or 3)
				if movementByRound[movement-2][roundBlock] > moveNum {
					moveOrder[roundBlock] = append(moveOrder[roundBlock], token)
				}
			}
		}
	}
	return moveOrder
}

// get all weapon slots on the board, sorted by initiative
func (b *battle) getSortedWeaponSlots(tokens []*battleToken) []*battleWeaponSlot {
	slots := []*battleWeaponSlot{}
	for _, token := range tokens {
		if token.isStillInBattle() {
			slots = append(slots, token.weaponSlots...)
		}
	}
	sort.Slice(slots, func(i, j int) bool {
		return slots[i].initiative > slots[j].initiative
	})
	return slots
}

// moveToken moves a token towards or away from its target
func (b *battle) moveToken(token *battleToken, weaponSlots []*battleWeaponSlot) {
	if !token.isStillInBattle() {
		return
	}
	// count this token's moves
	token.movesMade++

	// always disengage if we have no weapons or targets
	if token.moveTarget == nil || !token.hasWeapons() {
		token.Tactic = BattleTacticDisengage
	}

	// disengage if we are a scared little token
	if token.Tactic == BattleTacticDisengageIfChallenged && token.damaged {
		token.Tactic = BattleTacticDisengage
	}

	// assume we move nowhere
	oldPosition := token.Position
	var bestMoves []BattleVector

	if token.Tactic == BattleTacticDisengage {
		if token.movesMade >= b.rules.MovesToRunAway {
			// we've moved enough to leave the board
			token.ranAway = true
			b.board[token.Position.Y][token.Position.X] -= token.Quantity
			b.record.recordRunAway(b.round, token)
			return
		}

		// move out of harms way
		bestMoves = b.getBestFleeMoves(token, weaponSlots)

	} else {
		// move to best do some damage
		bestMoves = b.getBestAttackMoves(token, weaponSlots)
	}

	// update the board after a token moves
	bestMove := bestMoves[b.rules.random.Intn(len(bestMoves))]
	b.record.recordMove(b.round, token, token.Position, bestMove)
	token.Position = bestMove
	b.board[oldPosition.Y][oldPosition.X] -= token.Quantity
	b.board[token.Position.Y][token.Position.X] += token.Quantity
}

// getEstimatedDamageForWeapons estimates the damage for a group of weapons against a target
// if it were to move to a position
func (b *battle) getEstimatedDamageForWeapons(weapons []*battleWeaponSlot, target *battleToken, position BattleVector) int {
	damageDone := 0
	for _, weapon := range weapons {
		if !weapon.token.willTarget(target) {
			// this weapon wouldn't target the attacker, so don't add its damage
			continue
		}
		distanceToWeapon := position.distance(weapon.token.Position)
		damageDone += b.getEstimatedDamageForWeapon(weapon, target, distanceToWeapon)
	}

	return damageDone
}

// getEstimatedDamageForWeapon estimates the damage for a single weapon to a target at a distance
func (b *battle) getEstimatedDamageForWeapon(weapon *battleWeaponSlot, target *battleToken, distance int) int {
	if !weapon.isInRangeValue(distance) {
		// no damage, skip this weapon
		return 0
	}

	var bwd battleWeaponDamage
	if weapon.weaponType == battleWeaponTypeBeam {
		bwd = weapon.getBeamDamageToTargetAtDistance(weapon.power*weapon.slotQuantity*weapon.token.Quantity, target, distance, b.rules.BeamRangeDropoff)
	} else {
		bwd = weapon.getEstimatedTorpedoDamageToTarget(target)
	}
	// add up actual shield/armor damage plus any tokens that would be destroyed by this shot
	return bwd.shieldDamage + bwd.armorDamage
}

// for a new position on the board, and an existing slice of bestMoves, return the new bestMoves
// Note: this function is only called if a move is equal or better, never for worse moves
// if better, take the new position and discard the others
// if not better, but closer to center, take the new move and discard the others
// if not better and farther away, keep the bestMoves as is
// if equivalent and the same distance from center, add it to bestMoves
func updateMovesWithCenterPreference(better bool, newPosition BattleVector, bestMoves []BattleVector) []BattleVector {
	if len(bestMoves) == 0 {
		// we have no moves at all, so this is newPosition is our best move
		return append(bestMoves, newPosition)
	}

	if better {
		// this move is better, start over with it
		return []BattleVector{newPosition}
	}

	// center of battle board is at (4.5, 4.5) so scale to 100x100 for 45,45 scaled center
	newDistanceFromCenter := newPosition.scale(10).distance(BattleVector{45, 45})
	oldDistanceFromCenter := bestMoves[0].scale(10).distance(BattleVector{45, 45})

	if oldDistanceFromCenter == newDistanceFromCenter {
		// move is equivalent in damage and closeness to center, add new move
		return append(bestMoves, newPosition)
	}
	if oldDistanceFromCenter > newDistanceFromCenter {
		// damage is the same, but new position is closer to center, use only new move
		return []BattleVector{newPosition}
	}

	// damage is the same but farther from center, keep what we have but don't add this move
	return bestMoves
}

// get the best move to attack a target
// this checks all the nearby squares and moves in the direction that gets the token
// closer or does the best damage, according to the battle plan
// i.e. max damage, max damage ratio, max net damage, etc
func (b *battle) getBestAttackMoves(token *battleToken, weapons []*battleWeaponSlot) []BattleVector {

	bestDamageDone := 0
	bestDamageTaken := math.MaxInt
	bestDamageRatio := 0.0
	bestNetDamage := -math.MaxInt
	bestDamageMoves := []BattleVector{}
	bestMoveCloserMoves := []BattleVector{}

	// if no other options are presented, we move towards the moveTarget
	target := token.moveTarget
	bestDistanceToTarget := token.getDistanceAway(target.Position)

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			newPosition := BattleVector{token.Position.X + dx, token.Position.Y + dy}
			if newPosition.X < 0 || newPosition.X >= battleWidth || newPosition.Y < 0 || newPosition.Y >= battleHeight {
				// skip invalid squares
				continue
			}

			// see if this move puts us closer to the target
			distanceToTarget := newPosition.distance(target.Position)
			if distanceToTarget == bestDistanceToTarget {
				bestMoveCloserMoves = append(bestMoveCloserMoves, newPosition)
			} else if distanceToTarget < bestDistanceToTarget {
				bestDistanceToTarget = distanceToTarget
				bestMoveCloserMoves = []BattleVector{newPosition}
			}

			// figure out how much damage we will do if we move to this spot
			damageDone := 0
			for _, weapon := range token.weaponSlots {
				// each weapon will fire a volley at the most attractive target
				for _, target := range weapon.targets {
					distance := newPosition.distance(target.Position)
					weaponDamage := b.getEstimatedDamageForWeapon(weapon, target, distance)
					if weaponDamage > 0 {
						// targets are sorted by attractiveness and
						// we would be able to damage this target, so add it to
						// total damage and move to the next weapon
						damageDone += weaponDamage
						break
					}
				}
			}

			// figure out damage taken from all enemy weapons if we move to this square
			damageTaken := b.getEstimatedDamageForWeapons(weapons, token, newPosition)

			var damageRatio float64
			if damageTaken > 0 {
				damageRatio = float64(damageDone) / float64(damageTaken)
			}

			// b.log.Debug().Msgf("moving to %#v, damageDone: %d, damageTaken: %d, netDamage: %d, damageRatio: %f", newPosition, damageDone, damageTaken, damageDone-damageTaken, damageRatio)

			switch token.Tactic {
			case BattleTacticMaximizeDamageRatio:
				// first see if this move is doing damage but taking no damage
				// if our previous "best" move involved taking damage, this one is better
				// reset our bestDamageDone to 0 to "start over"
				if damageTaken == 0 && bestDamageTaken > 0 && damageDone > 0 {
					bestDamageDone = 0
					bestDamageTaken = 0
				}
				if damageTaken == 0 && damageDone > 0 {
					// we took no damage, so sort by best damage
					if damageDone >= bestDamageDone {
						bestDamageMoves = updateMovesWithCenterPreference(damageDone > bestDamageDone, newPosition, bestDamageMoves)
						bestDamageDone = damageDone
					}
				} else if bestDamageDone == 0 {
					// we don't have a best damage yet, and this move
					if damageRatio >= bestDamageRatio {
						bestDamageMoves = updateMovesWithCenterPreference(damageRatio > bestDamageRatio, newPosition, bestDamageMoves)
						bestDamageRatio = damageRatio
					}
				}
			case BattleTacticMaximizeNetDamage:
				if damageDone > 0 && damageDone-damageTaken >= bestNetDamage {
					bestDamageMoves = updateMovesWithCenterPreference(damageDone-damageTaken > bestNetDamage, newPosition, bestDamageMoves)
					bestDamageDone = damageDone
					bestNetDamage = damageDone - damageTaken
				}

			case BattleTacticMinimizeDamageToSelf:
				if damageTaken <= bestDamageTaken {
					bestDamageMoves = updateMovesWithCenterPreference(damageTaken < bestDamageTaken, newPosition, bestDamageMoves)
					bestDamageTaken = damageTaken
					if damageDone > bestDamageDone {
						bestDamageDone = damageDone
					}
				}

			case BattleTacticMaximizeDamage:
				if damageDone >= bestDamageDone {
					bestDamageMoves = updateMovesWithCenterPreference(damageDone > bestDamageDone, newPosition, bestDamageMoves)
					bestDamageDone = damageDone
				}
			}
		}
	}

	// if none of our moves lead to damage, pick the move that moves us towards our target
	if bestDamageDone == 0 && bestDamageRatio == 0.0 {
		return bestMoveCloserMoves
	}

	return bestDamageMoves
}

// get the best move this token should fleet to based on weapons on the board
func (b *battle) getBestFleeMoves(token *battleToken, weapons []*battleWeaponSlot) []BattleVector {
	// find the best move for running away
	lowestDamageMoves := make([]BattleVector, 0, 9)

	// if we stayed still, figure out our damage
	damageTaken := b.getEstimatedDamageForWeapons(weapons, token, token.Position)
	lowestDamage := damageTaken

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			newPosition := BattleVector{token.Position.X + dx, token.Position.Y + dy}
			if newPosition.X < 0 || newPosition.X >= battleWidth || newPosition.Y < 0 || newPosition.Y >= battleHeight {
				// skip invalid squares
				continue
			}

			// figure out damage taken from all weapons targeting us if we moved to this square
			damageTaken := b.getEstimatedDamageForWeapons(weapons, token, newPosition)
			// b.log.Debug().Msgf("moving to %#v causes %d damage", newPosition, damageTaken)

			// if this move is the same as our previous one, add it to our possible list
			if damageTaken == lowestDamage {
				lowestDamageMoves = append(lowestDamageMoves, newPosition)
			}

			// if this move is better, replace the list
			if damageTaken < lowestDamage {
				lowestDamage = damageTaken
				lowestDamageMoves = []BattleVector{newPosition}
			}
		}
	}

	// pick a random best move and move there
	return lowestDamageMoves
}

// Fire the weapon slot towards its target
func (b *battle) fireWeaponSlot(weapon *battleWeaponSlot, targets []*battleToken) {
	if len(targets) == 0 || !weapon.token.isStillInBattle() {
		// no targets, nothing to do
		return
	}

	switch weapon.weaponType {
	case battleWeaponTypeBeam:
		b.fireBeamWeapon(weapon, targets)
	case battleWeaponTypeTorpedo:
		b.fireTorpedo(weapon, targets)
	}
}

// fire a beam weapon slot at a slice of targets
// weapons fire in volleys. A single slot fires a volley at a target
// if you have 10 frigates with a 3 laser slot, they fire 30 lasers as a "volley"
// if you have 10 destroyers with three 1-laser slots, they fire three volleys of 10 lasers each
func (b *battle) fireBeamWeapon(weapon *battleWeaponSlot, targets []*battleToken) {

	// get the damage for this volley
	attacker := weapon.token
	damage := weapon.power * weapon.slotQuantity * attacker.Quantity
	b.log.Debug().Msgf("%v is to firing %vx%d at %d targets for a total of %v damage", weapon.token, weapon.slot.HullComponent, weapon.slotQuantity*attacker.Quantity, len(targets), damage)

	for _, target := range targets {
		if !target.isStillInBattle() {
			continue
		}
		// skip targets that are out of shields for sappers
		if weapon.damagesShieldsOnly && target.stackShields <= 0 {
			continue
		}
		// reset the damage to base damage if this weapon hits all targets
		if weapon.hitsAllTargets {
			damage = weapon.power * weapon.slotQuantity * attacker.Quantity
		}

		if damage == 0 {
			// no more damage to do
			break
		}

		// check the damage against this target
		bwd := weapon.getBeamDamageToTarget(damage, target, b.rules.BeamRangeDropoff)
		b.log.Debug().Msgf("%v fired a %vx%d at %v (shields: %v, armor: %v, %v@%v damage) for %v damage", weapon.token, weapon.slot.HullComponent, weapon.slotQuantity*weapon.token.Quantity, target, target.totalStackShields, target.armor, target.Quantity, target.Damage, bwd.armorDamage)

		// update stack shields
		target.stackShields -= bwd.shieldDamage

		// update damage for the next target
		damage = bwd.leftover

		if bwd.numDestroyed >= target.Quantity {
			target.Quantity = 0
			target.QuantityDamaged = 0
			target.Damage = 0
			target.quantityDestroyed += bwd.numDestroyed
			b.board[target.Position.Y][target.Position.X] -= bwd.numDestroyed
			target.destroyed = true
			b.log.Debug().Msgf("%v %v did %v shield damage and %v armor damage (leftoverDamage %d) and completely destroyed %v", weapon.token, weapon.slot.HullComponent, bwd.shieldDamage, bwd.armorDamage, bwd.leftover, target)

			// record one round of beam fire per target
			b.record.recordBeamFire(b.round, weapon.token, weapon.token.Position, target.Position, weapon.slot.HullSlotIndex, *target, bwd.shieldDamage, bwd.armorDamage, bwd.numDestroyed)

			// next target
			continue
		}

		// handle any destroyed tokens
		target.Quantity -= bwd.numDestroyed
		target.quantityDestroyed += bwd.numDestroyed
		b.board[target.Position.Y][target.Position.X] -= bwd.numDestroyed

		// apply damage to this ship
		target.Damage = bwd.damage
		target.QuantityDamaged = bwd.quantityDamaged

		b.log.Debug().Msgf("%v destroyed %v ships (leftoverDamage %d), leaving %v damaged %v@%v damage", weapon.token, bwd.numDestroyed, bwd.leftover, target, target.Quantity, target.Damage)
		target.damaged = true

		// record one round of beam fire per target
		b.record.recordBeamFire(b.round, weapon.token, weapon.token.Position, target.Position, weapon.slot.HullSlotIndex, *target, bwd.shieldDamage, bwd.armorDamage, bwd.numDestroyed)
	}
}

// Fire a torpedo slot from a ship. Torpedos are different than beam weapons
// A ship will fire each torpedo at its target until the target is destroyed, then
// fire remaining torpedos at the next target.
// Each torpedo has an accuracy rating. That determines if it hits. A torpedo that
// misses still explodes and does 1/8th damage to shields
func (b *battle) fireTorpedo(weapon *battleWeaponSlot, targets []*battleToken) {
	attacker := weapon.token
	damage := weapon.power
	numTorpedos := weapon.slotQuantity * attacker.Quantity

	b.log.Debug().Msgf("%s is attempting to fire at %d targets with %d torpedos at %.2f%% accuracy for %d damage each",
		weapon.token, len(targets), numTorpedos, (weapon.getAccuracy(0))*100.0, damage)

	// fire each torpedo at each target until it's destroyed or we're out of torpedos
	remainingTorpedos := numTorpedos
	torpedoNum := 0
	for _, target := range targets {
		if !target.isStillInBattle() {
			// this token isn't valid anymore, skip it
			continue
		}

		// no more damage to spread, break out
		if remainingTorpedos == 0 {
			break
		}

		// shields are shared among all tokens
		armor := target.armor
		shipDamage := target.Damage

		totalShieldDamage := 0
		totalArmorDamage := 0
		hits := 0
		misses := 0
		shipsDestroyed := 0

		for remainingTorpedos > 0 && !target.destroyed {
			// fire a torpedo
			torpedoNum++
			remainingTorpedos--
			hit := b.rules.random.Float64() <= weapon.getAccuracy(target.torpedoJamming)

			if hit {
				hits++

				// torpedos do half damage to shields, half to armor (until shields are gone, when they do full armor damage)
				shieldDamage := float64(0.5) * float64(damage)
				armorDamage := float64(0.5) * float64(damage)

				// apply up to half our damage to shields
				// anything leftover goes to armor
				afterShieldsDamaged := float64(target.stackShields) - shieldDamage
				var actualShieldDamage float64
				if afterShieldsDamaged < 0 {
					// We did more damage to shields than they had remaining
					// apply the difference to armor
					actualShieldDamage = shieldDamage + afterShieldsDamaged
					armorDamage += float64(-afterShieldsDamaged)

				} else {
					actualShieldDamage = shieldDamage
				}
				target.stackShields -= int(actualShieldDamage)

				if target.stackShields <= 0 && weapon.capitalShipMissile {
					// capital ship missiles double damage after shields are gone
					armorDamage *= 2
				}

				totalShieldDamage += int(actualShieldDamage)
				totalArmorDamage += int(armorDamage)
				shipDamage += armorDamage

				// this torpedo blew up a ship, hooray!
				if shipDamage >= float64(armor) {
					// remove a ship from this stack
					target.Quantity--
					target.quantityDestroyed++
					b.board[target.Position.Y][target.Position.X] -= 1
					target.QuantityDamaged = MaxInt(target.QuantityDamaged-1, 0)

					if target.QuantityDamaged > 0 {
						// we destroyed a token, but we still have damaged tokens in the stack
						// so reset our shipDamage counter to the damage + any leftover. We apply that
						// to the rest of the tokens
						// i.e. if we fire 2 omega torpedos for 300 damage each at 3 damaged 1700dp@1300 ships
						// the first shot damages the top ship, the second one kills it but we have 200 leftover
						// this will carry over to damage the remaining ships
						leftoverDamage := shipDamage - float64(armor)
						shipDamage = target.Damage + leftoverDamage
					} else {
						// we have no more damaged tokens, so remove the stack's damage
						// and reset our ship damage to 0
						// this could happen if we are firing on a stack with 3 ships but 2@10 damage or something
						shipDamage = 0
						target.Damage = 0
					}
					if target.Quantity <= 0 {
						// record that we destroyed this token
						target.destroyed = true
						b.log.Debug().Msgf("%v torpedo number %v hit %v, did %v shield damage and %v armor damage and completely destroyed %v", weapon.token, torpedoNum, target, actualShieldDamage, armorDamage, target)
						shipsDestroyed++
					}
				} else {
					b.log.Debug().Msgf("%v torpedo number %v hit %v, did %v shield damage and %v armor damage (%v accumulated damage so far)", weapon.token, torpedoNum, target, actualShieldDamage, armorDamage, shipDamage)
				}
			} else {
				misses++
				// damage shields by 1/8th
				// round up, do a minimum of 1 damage
				shieldDamage := int(math.Min(1, math.Round(b.rules.TorpedoSplashDamage*float64(damage))))
				actualShieldDamage := shieldDamage
				if shieldDamage > target.stackShields {
					actualShieldDamage = target.stackShields
				}
				target.stackShields = int(math.Max(0, float64(target.stackShields-shieldDamage)))
				b.log.Debug().Msgf("%s torpedo number %d missed %s, did %d damage to shields leaving %d shields", weapon.token, torpedoNum, target, shieldDamage, target.stackShields)

				totalShieldDamage += actualShieldDamage
			}
		}

		// we have leftover damage, apply it to all remaining tokens evenly
		if shipDamage > 0 && target.Quantity > 0 {
			target.damaged = true // target lived, but is damaged
			var previousDamage float64
			if target.QuantityDamaged > 0 {
				// we had some tokens damaged previously that we didn't touch
				previousDamage = target.Damage * float64(target.QuantityDamaged)
				shipDamage -= target.Damage // we already include this in our ship damage
			}
			target.Damage = (shipDamage + previousDamage) / float64(target.Quantity)
			target.QuantityDamaged = target.Quantity
			b.log.Debug().Msgf("%s had %d hits and %d misses to %v for %v total damage leaving %d@%v", weapon.token, hits, misses, target, totalArmorDamage+totalShieldDamage, target.QuantityDamaged, target.Damage)

		}
		b.record.recordTorpedoFire(b.round, weapon.token, weapon.token.Position, target.Position, weapon.slot.HullSlotIndex, target, totalShieldDamage, totalArmorDamage, shipsDestroyed, hits, misses)
	}
}

// Function to allow server to run test battles
func RunTestBattle(players []*Player, fleets []*Fleet) (*BattleRecord, error) {
	rules := NewRules()

	playersByNum := map[int]*Player{}
	designsByNum := make(map[playerObject]*ShipDesign)
	battlePlansByNum := make(map[playerBattlePlanNum]*BattlePlan)

	for _, player := range players {
		playersByNum[player.Num] = player
		player.Race.Spec = computeRaceSpec(&player.Race, &rules)
		player.Spec = computePlayerSpec(player, &rules, []*Planet{})

		for _, design := range player.Designs {
			var err error
			design.Spec, err = ComputeShipDesignSpec(&rules, player.TechLevels, player.Race.Spec, design)
			if err != nil {
				return nil, fmt.Errorf("ComputeShipDesignSpec returned error %w", err)
			}
			designsByNum[playerObjectKey(design.PlayerNum, design.Num)] = design
		}

		for i := range player.BattlePlans {
			plan := &player.BattlePlans[i]
			battlePlansByNum[playerBattlePlanNum{PlayerNum: player.Num, Num: plan.Num}] = plan
		}

	}

	for _, fleet := range fleets {
		for i := range fleet.Tokens {
			token := &fleet.Tokens[i]
			token.design = designsByNum[playerObjectKey(fleet.PlayerNum, token.DesignNum)]
		}
		fleet.Spec = ComputeFleetSpec(&rules, playersByNum[fleet.PlayerNum], fleet)
		fleet.battlePlan = battlePlansByNum[playerBattlePlanNum{fleet.PlayerNum, fleet.BattlePlanNum}]
	}

	battler := newBattler(log.Logger, &rules, &StaticTechStore, 1, playersByNum, fleets, nil)
	record := battler.runBattle()
	for _, player := range players {

		for _, otherplayer := range players {
			player.discoverer.discoverPlayer(otherplayer)
		}
		for _, fleet := range fleets {
			if fleet.PlayerNum != player.Num {
				for _, token := range fleet.Tokens {
					player.discoverer.discoverDesign(token.design, true)
				}
			}
		}
	}

	return record, nil
}

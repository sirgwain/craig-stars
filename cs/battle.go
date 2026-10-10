package cs

import (
	"log/slog"
	"math"
	"sort"
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
//
// craig-stars differences from the description above: there is no limit on the number
// of tokens in a battle, and armor damage is stored in 1/500ths of a ship's armor.

const battleWidth, battleHeight = 10, 10

// The battler interface is the main entrypoint into running battles
// A batter is created for fleets in the same location in the universe and it is used
// to check if a battle will occur, and run the battle.
// Running a battle returns a BattleRecord that is passed along to each player.
type battler interface {
	hasHostility() bool
	runBattle() *BattleRecord
}

// battle defines the state of a battle as it progresses
type battle struct {
	num      int
	planet   *Planet
	position Vector
	tokens   []*battleToken
	round    int
	players  map[int]*Player
	record   *BattleRecord
	rules    *Rules
	log      *slog.Logger
	fleets   []*Fleet
}

// battleStartingPositions groups board coordinates by participant count.
var battleStartingPositions = [...][]Vector{
	1:  {{4, 4}},
	2:  {{1, 4}, {8, 5}},
	3:  {{4, 1}, {8, 8}, {1, 8}},
	4:  {{1, 1}, {8, 8}, {1, 8}, {8, 1}},
	5:  {{4, 1}, {6, 8}, {1, 4}, {8, 4}, {2, 8}},
	6:  {{1, 4}, {8, 5}, {2, 8}, {7, 1}, {6, 8}, {3, 1}},
	7:  {{1, 1}, {1, 5}, {2, 8}, {6, 8}, {8, 6}, {8, 2}, {5, 1}},
	8:  {{1, 3}, {1, 6}, {3, 8}, {6, 8}, {8, 6}, {8, 3}, {6, 1}, {3, 1}},
	9:  {{1, 3}, {8, 6}, {3, 8}, {6, 1}, {1, 6}, {8, 3}, {6, 8}, {3, 1}, {4, 4}},
	10: {{2, 1}, {5, 1}, {8, 1}, {1, 4}, {8, 4}, {4, 5}, {1, 7}, {8, 7}, {3, 8}, {6, 8}},
	11: {{1, 3}, {8, 6}, {3, 8}, {6, 1}, {1, 6}, {8, 3}, {6, 8}, {3, 1}, {3, 4}, {6, 3}, {6, 6}},
	12: {{1, 4}, {8, 5}, {2, 8}, {7, 1}, {6, 8}, {3, 1}, {1, 6}, {8, 3}, {1, 2}, {4, 8}, {5, 1}, {8, 7}},
	13: {{1, 1}, {1, 3}, {1, 5}, {1, 7}, {3, 1}, {5, 1}, {7, 1}, {8, 3}, {8, 5}, {3, 8}, {5, 8}, {7, 8}, {4, 4}},
	14: {{1, 1}, {1, 3}, {1, 5}, {1, 7}, {2, 8}, {4, 8}, {6, 8}, {8, 8}, {8, 6}, {8, 4}, {8, 2}, {7, 1}, {5, 1}, {3, 1}},
	15: {{1, 1}, {1, 3}, {1, 5}, {1, 7}, {2, 8}, {4, 8}, {6, 8}, {8, 8}, {8, 6}, {8, 4}, {8, 2}, {7, 1}, {5, 1}, {3, 1}, {4, 4}},
	16: {{1, 1}, {1, 3}, {1, 5}, {1, 7}, {2, 8}, {4, 8}, {6, 8}, {8, 8}, {8, 6}, {8, 4}, {8, 2}, {7, 1}, {5, 1}, {3, 1}, {3, 3}, {6, 6}},
}

// getMovesForRound returns how many squares a token moves in a round.
// - speed: token movement minus 2, so 0 (3/4 moves per round) to 8 (2 1/2).
// - round: zero-based round number; the pattern repeats every 4 rounds.
func getMovesForRound(speed, round int) int {
	base := (speed + 2) >> 2 // whole moves per round

	switch speed & 0x3 {
	case 0:
		// +1 on even rounds
		if round&0x1 == 0 {
			return base + 1
		}
		return base
	case 1:
		// +1 except when (round % 4) == 2
		if round&0x3 != 2 {
			return base + 1
		}
		return base
	case 2:
		// never +1
		return base
	case 3:
		// +1 when (round % 4) == 0
		if round&0x3 == 0 {
			return base + 1
		}
		return base
	default:
		// unreachable, but return base
		return base
	}
}

// get the movement of this design with additional cargo
func getBattleSpeed(movementMin, movementMax, idealEngineSpeed int, movementBonus float64, mass, numEngines int) int {
	if numEngines == 0 {
		return 0
	}
	mb := int(math.Ceil(movementBonus)) // round up fractional movement bonus
	speed := idealEngineSpeed + mb - 2
	massPenalty := ((mass / 70) / numEngines)
	return Clamp(speed-massPenalty, movementMin, movementMax)
}

// newBattler creates the battle tokens for the fleets of players that take part in a
// battle at this location. Use hasHostility to check whether a battle takes place.
func newBattler(log *slog.Logger, rules *Rules, battleNum int, players map[int]*Player, fleets []*Fleet, planet *Planet) battler {
	if len(fleets) == 0 {
		return nil
	}
	hostility := getBattleHostility(players, fleets)
	participants := make(map[int]*Player)
	for number, targets := range hostility {
		if len(targets) > 0 {
			participants[number] = players[number]
		}
	}
	fleets = selectBattleFleets(fleets, participants)
	playerNums := make([]int, 0, len(participants))
	for number := range participants {
		playerNums = append(playerNums, number)
	}
	sort.Ints(playerNums)
	startingPositions := make(map[int]Vector)
	for i, number := range playerNums {
		startingPositions[number] = getBattleStartingPosition(len(playerNums), i)
	}
	tokens := make([]*battleToken, 0)
	for _, fleet := range fleets {
		for i := range fleet.Tokens {
			token := &fleet.Tokens[i]
			if token.Quantity == 0 {
				continue
			}
			bt := newBattleToken(rules, len(tokens)+1, startingPositions[fleet.PlayerNum], token, *fleet.battlePlan, players[fleet.PlayerNum])
			bt.fleet = fleet
			bt.attackPlayers = hostility[fleet.PlayerNum]
			tokens = append(tokens, bt)
		}
	}
	position := Vector{}
	if len(fleets) > 0 {
		position = fleets[0].Position
	}
	return &battle{num: battleNum, planet: planet, position: position, tokens: tokens, rules: rules, log: log, players: participants, fleets: fleets}
}

// getBattleHostility combines attack orders, retaliation, and allied support into
// each player's hostile intentions when an armed fleet can initiate combat.
func getBattleHostility(players map[int]*Player, fleets []*Fleet) map[int]map[int]bool {
	hostility, canStart := getBattleAttackOrders(players, fleets)
	if !canStart {
		for number := range hostility {
			hostility[number] = make(map[int]bool)
		}
		return hostility
	}
	addBattleRetaliation(hostility)
	addBattleSupport(players, hostility)
	return hostility
}

// getBattleAttackOrders returns the players each player's armed fleets will attack,
// and whether any armed fleet other than a starbase has orders that start a battle.
func getBattleAttackOrders(players map[int]*Player, fleets []*Fleet) (map[int]map[int]bool, bool) {
	hostility := make(map[int]map[int]bool, len(players))
	for number := range players {
		hostility[number] = make(map[int]bool)
	}
	canStart := false
	for _, fleet := range fleets {
		armed := false
		for _, token := range fleet.Tokens {
			armed = armed || token.Quantity > 0 && token.design.Spec.HasWeapons
		}
		if !armed || fleet.battlePlan == nil || !fleet.Starbase && fleet.battlePlan.PrimaryTarget == BattleTargetNone {
			continue
		}
		if !fleet.Starbase {
			switch fleet.battlePlan.AttackWho {
			case BattleAttackWhoEnemies, BattleAttackWhoEnemiesAndNeutrals, BattleAttackWhoEveryone:
				canStart = true
			}
		}
		for number := range players {
			if number != fleet.PlayerNum && fleet.willAttack(players[fleet.PlayerNum], number) {
				hostility[fleet.PlayerNum][number] = true
			}
		}
	}
	return hostility, canStart
}

// addBattleRetaliation makes every attacked player hostile toward its attackers.
func addBattleRetaliation(hostility map[int]map[int]bool) {
	for number, targets := range hostility {
		for target := range targets {
			hostility[target][number] = true
		}
	}
}

// addBattleSupport lets players with no hostilities of their own join their friends'
// fights, repeating until no one else joins so support can chain through friends.
// A player whose friends are fighting each other stays out. Support is one-sided:
// the opponents do not add the supporter to their own targets.
func addBattleSupport(players map[int]*Player, hostility map[int]map[int]bool) {
	playerNums := make([]int, 0, len(players))
	for number := range players {
		playerNums = append(playerNums, number)
	}
	sort.Ints(playerNums)
	for changed := true; changed; {
		changed = false
		for _, number := range playerNums {
			if len(hostility[number]) > 0 {
				continue
			}
			player := players[number]
			inherited := make(map[int]bool)
			for _, friend := range playerNums {
				if !player.IsFriend(friend) {
					continue
				}
				if inherited[friend] {
					// one friend is fighting another
					inherited = make(map[int]bool)
					break
				}
				for target := range hostility[friend] {
					if target != number {
						inherited[target] = true
					}
				}
			}
			hostility[number] = inherited
			if len(inherited) > 0 {
				changed = true
			}
		}
	}
}

// selectBattleFleets returns the fleets belonging to participating players.
func selectBattleFleets(fleets []*Fleet, players map[int]*Player) []*Fleet {
	selected := make([]*Fleet, 0, len(fleets))
	for _, fleet := range fleets {
		if players[fleet.PlayerNum] != nil {
			selected = append(selected, fleet)
		}
	}
	return selected
}

// getBattleStartingPosition returns a player's board position in the formation
// for the given participant count and player index. Battles with more players
// than the largest formation share its positions.
func getBattleStartingPosition(players, index int) Vector {
	formation := battleStartingPositions[Clamp(players, 1, len(battleStartingPositions)-1)]
	return formation[index%len(formation)]
}

// prepareBattle dumps ordered mineral cargo, updates mass and speed, applies movement
// dampening, and shuffles tokens once to establish initiative tie order before
// recording the starting state.
func (b *battle) prepareBattle() {
	dumped := Mineral{}
	dumpedFleets := make(map[*Fleet]bool)
	for _, fleet := range b.fleets {
		if fleet.battlePlan.DumpCargo && fleet.Cargo.HasMinerals() {
			dumped = dumped.Add(fleet.Cargo.ToMineral())
			fleet.Cargo = Cargo{Colonists: fleet.Cargo.Colonists}
			dumpedFleets[fleet] = true
		}
	}
	dampening := 0
	for _, token := range b.tokens {
		dampening = max(dampening, token.design.Spec.ReduceMovement)
	}
	for _, token := range b.tokens {
		cargo := getCargoPerShip(token.fleet.Cargo.Total(), token.fleet.Spec.CargoCapacity, token.design.Spec.CargoCapacity)
		token.Mass = token.design.Spec.Mass + cargo
		token.Movement = token.design.getMovement(b.rules, cargo)
		if dumpedFleets[token.fleet] && token.design.Spec.CargoCapacity > 0 {
			token.Movement = max(b.rules.MovementMin, token.Movement-1)
		}
		if token.Movement > 0 {
			token.Movement = Clamp(token.Movement-dampening, b.rules.MovementMin, b.rules.MovementMax)
		}
	}
	b.rules.random.Shuffle(len(b.tokens), func(i, j int) { b.tokens[i], b.tokens[j] = b.tokens[j], b.tokens[i] })
	records := make([]BattleRecordToken, 0, len(b.tokens))
	for _, token := range b.tokens {
		records = append(records, token.BattleRecordToken)
	}
	planetNum := None
	if b.planet != nil {
		planetNum = b.planet.Num
	}
	b.record = newBattleRecord(b.num, planetNum, b.position, records)
	b.record.fleets = b.fleets
	b.record.dumpedMinerals = dumped
}

// hasHostility reports whether any active token has hostile intentions toward another active player.
func (b *battle) hasHostility() bool {
	for _, token := range b.tokens {
		if !token.isStillInBattle() {
			continue
		}
		for _, other := range b.tokens {
			if other.isStillInBattle() && token.willAttack(other.PlayerNum) {
				return true
			}
		}
	}
	return false
}

// hasMultiplePlayers reports whether tokens from more than one player remain in the battle.
func (b *battle) hasMultiplePlayers() bool {
	playerNum := None
	for _, token := range b.tokens {
		if !token.isStillInBattle() {
			continue
		}
		if playerNum == None {
			playerNum = token.PlayerNum
		} else if token.PlayerNum != playerNum {
			return true
		}
	}
	return false
}

// runBattle runs a battle!
func (b *battle) runBattle() *BattleRecord {
	for _, fleet := range b.fleets {
		fleet.noHeal = true
	}
	b.prepareBattle()
	for b.round = 1; b.round <= b.rules.NumBattleRounds; b.round++ {
		if b.round > 1 {
			for _, token := range b.tokens {
				if token.isStillInBattle() {
					token.regenerateShields()
				}
			}
		}
		if !b.hasHostility() {
			// battle done, no more hostility
			break
		}

		// start a new round
		b.record.recordNewRound()

		// 1. movement
		for _, token := range b.tokens {
			token.movesLeft = 0
			if token.Movement > 0 {
				token.movesLeft = getMovesForRound(token.Movement-2, b.round-1)
			}
			token.movementMass = float64(token.Mass) * (1 + (2*b.rules.random.Float64()-1)*b.rules.MovementMassVariance)
		}
		weapons := b.getSortedWeaponSlots(b.tokens)
		for _, token := range b.buildMovementOrder(b.tokens, b.round-1) {
			b.moveToken(token)
			token.movesLeft--
		}

		if !b.hasHostility() {
			// if runners ran and we're done, exit
			break
		}

		// 2. fire weapons
		for _, weapon := range weapons {
			if !b.hasMultiplePlayers() {
				break
			}
			b.fireWeaponSlot(weapon, weapon.findTargets(b.tokens))
		}
	}

	// Damage is tracked as fractions during battle, but ships keep only whole points.
	for _, token := range b.tokens {
		token.Damage = math.Floor(token.Damage)
		if token.Damage == 0 {
			token.QuantityDamaged = 0
		}
		if token.quantityDestroyed > 0 {
			b.record.recordDestroyedToken(token, token.quantityDestroyed)
		}
	}
	return b.record
}

// buildMovementOrder returns the order tokens move in a round, one entry per square moved.
// Each ship moves in order of mass with heavier ships moving first.
// Ships that can move 3 times in a round move first, then ships that move 2 times, then 1.
func (b *battle) buildMovementOrder(tokens []*battleToken, round int) (moveOrder []*battleToken) {
	// our tokens are moved by mass
	tokensByMass := make([]*battleToken, 0)
	for _, token := range tokens {
		if token.Movement > 0 { // starbases don't move
			tokensByMass = append(tokensByMass, token)
		}
	}
	sort.SliceStable(tokensByMass, func(i, j int) bool {
		return tokensByMass[i].movementMass > tokensByMass[j].movementMass
	})

	// each token can move up to 3 times in a round
	// ships that can move 3 times go first, so we loop through the moveNum backwards
	// so that the order has ships that move 3 times first
	for moveNum := 2; moveNum >= 0; moveNum-- {
		for _, token := range tokensByMass {
			if getMovesForRound(token.Movement-2, round) > moveNum {
				moveOrder = append(moveOrder, token)
			}
		}
	}
	return moveOrder
}

// get all weapon slots on the board, sorted by initiative
func (b *battle) getSortedWeaponSlots(tokens []*battleToken) []*battleWeaponSlot {
	slots := []*battleWeaponSlot{}
	for i := len(tokens) - 1; i >= 0; i-- {
		token := tokens[i]
		if token.isStillInBattle() {
			slots = append(slots, token.weaponSlots...)
		}
	}
	sort.SliceStable(slots, func(i, j int) bool {
		return slots[i].initiative > slots[j].initiative
	})
	return slots
}

// moveToken moves a token towards or away from its target
func (b *battle) moveToken(token *battleToken) {
	if !token.isStillInBattle() {
		return
	}
	if !token.hasWeapons() {
		token.Tactic = BattleTacticDisengage
	}
	oldPosition := token.Position
	if token.Tactic == BattleTacticDisengage {
		if token.movesMade >= b.rules.MovesToRunAway {
			token.ranAway = true
			b.record.recordRunAway(b.round, token)
			return
		}
		token.movesMade++
	}
	moves := b.getBestMoves(token)
	if len(moves) == 0 {
		return
	}
	token.Position = moves[b.rules.random.Intn(len(moves))]
	b.record.recordMove(b.round, token, oldPosition, token.Position)
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
	damage := weapon.beamPower()
	for _, target := range targets {
		if !target.isStillInBattle() || !weapon.willDamage(target) {
			continue
		}
		if damage <= 0 {
			break
		}
		if weapon.hitsAllTargets {
			damage = weapon.beamPower()
		}
		result := weapon.getBeamDamageToTarget(damage, target, b.rules.BeamRangeDropoff)
		adjustedDamage := result.shieldDamage + result.armorDamage + result.leftover
		b.applyWeaponDamage(target, result)
		b.record.recordBeamFire(b.round, weapon.token, weapon.token.Position, target.Position, weapon.slot.HullSlotIndex, *target, result.shieldDamage, result.armorDamage, result.numDestroyed)
		if weapon.hitsAllTargets {
			continue
		}
		if result.leftover > 0 && adjustedDamage > 0 {
			// Carry unused base power forward so each target applies its own defenses.
			damage = min(damage-1, int(math.Round(float64(damage)*float64(result.leftover)/float64(adjustedDamage))))
		} else {
			damage = 0
		}
	}
}

// Fire a torpedo slot from a ship.
// A ship will fire each torpedo at its target until the target is destroyed, then
// fire any remaining torpedoes at the next target.
//
// Each torpedo has an accuracy rating that determines how often it hits the target.
// A torpedo that misses still explodes and does 1/8th damage to shields (if any).
func (b *battle) fireTorpedo(weapon *battleWeaponSlot, targets []*battleToken) {
	remaining := weapon.slotQuantity * weapon.token.Quantity
	for _, target := range targets {
		if !target.isStillInBattle() {
			continue
		}
		if remaining == 0 {
			break
		}
		hits := b.torpedoHits(remaining, weapon.getAccuracy(target.torpedoJamming))
		fired := remaining
		armorLeft := float64(target.armor*target.Quantity) - target.Damage*float64(target.QuantityDamaged)
		power := weapon.power
		if weapon.capitalShipMissile && target.stackShields == 0 {
			power *= 2
		}
		if target.Quantity < remaining && float64(hits*power) > armorLeft {
			// The volley can destroy this stack, so find the smallest number of torpedoes
			// that does it and leave the rest for the next target. Each torpedo destroys
			// at most one ship, so start at the ship count. Hits are scaled in proportion,
			// and missed torpedoes splash shields before the hits are applied.
			for count := target.Quantity; count <= remaining; count++ {
				hitCount := int(math.Ceil(float64(count*hits) / float64(remaining)))
				missCount := count - hitCount
				shields := max(0, float64(target.stackShields)-float64(missCount*power)*b.rules.TorpedoSplashDamage)
				armorDamage := float64(hitCount*power)/2 + max(0, float64(hitCount*power)/2-shields)
				if armorDamage >= armorLeft {
					fired = count
					break
				}
			}
		}
		firedHits := hits
		if fired < remaining {
			firedHits = int(math.Ceil(float64(fired*hits) / float64(remaining)))
		}
		misses := fired - firedHits
		result := weapon.getTorpedoVolleyDamage(target, float64(firedHits), float64(misses), fired, b.rules.TorpedoSplashDamage)
		b.applyWeaponDamage(target, result)
		b.record.recordTorpedoFire(b.round, weapon.token, weapon.token.Position, target.Position, weapon.slot.HullSlotIndex, target, result.shieldDamage, result.armorDamage, result.numDestroyed, firedHits, misses)
		remaining -= fired
	}
}

// torpedoHits rolls individual hits for small volleys and uses the expected hit
// count for volleys containing more than 200 torpedoes.
func (b *battle) torpedoHits(count int, accuracy float64) int {
	if count > 200 {
		return int(float64(count) * accuracy)
	}
	hits := 0
	for i := 0; i < count; i++ {
		if b.rules.random.Float64() < accuracy {
			hits++
		}
	}
	return hits
}

// applyWeaponDamage updates a target's defenses and casualties, accounts for cargo
// loss and salvage, and starts challenged retreat when armor is damaged.
func (b *battle) applyWeaponDamage(target *battleToken, result battleWeaponDamage) {
	oldQuantity := target.Quantity
	target.stackShields -= result.shieldDamage
	target.Quantity -= result.numDestroyed
	target.quantityDestroyed += result.numDestroyed
	target.Damage = result.damage
	target.QuantityDamaged = result.quantityDamaged
	if oldQuantity > 0 && result.numDestroyed > 0 {
		target.stackShields = int(math.Round(float64(target.stackShields) * float64(target.Quantity) / float64(oldQuantity)))
		target.totalStackShields = int(math.Round(float64(target.totalStackShields) * float64(target.Quantity) / float64(oldQuantity)))
		lost := b.removeDestroyedCargo(target, result.numDestroyed)
		if target.attributes&battleTokenAttributeStarbase == 0 && target.design != nil {
			minerals := MultiplyCost(target.design.Spec.Cost, float64(result.numDestroyed)*b.rules.SalvageFromBattleFactor).ToMineral().Add(lost.ToMineral())
			b.record.salvageMinerals = b.record.salvageMinerals.Add(minerals.MultiplyFloat64(b.salvageRecovery(), math.Round))
		}
	}
	if (result.armorDamage > 0 || result.numDestroyed > 0) && target.Tactic == BattleTacticDisengageIfChallenged {
		target.Tactic = BattleTacticDisengage
		target.movesMade = 0
	}
	target.destroyed = target.Quantity == 0
}

// removeDestroyedCargo removes the destroyed ships' share of fleet cargo and fuel,
// records cargo losses, and returns the lost cargo.
func (b *battle) removeDestroyedCargo(target *battleToken, quantity int) Cargo {
	fleet := target.fleet
	if fleet == nil {
		return Cargo{}
	}
	capacity := 0
	fuelCapacity := 0
	ships := 0
	for _, token := range fleet.Tokens {
		capacity += token.Quantity * token.design.Spec.CargoCapacity
		fuelCapacity += token.Quantity * token.design.Spec.FuelCapacity
		ships += token.Quantity
	}
	lostCapacity := quantity * target.design.Spec.CargoCapacity
	lostFuelCapacity := quantity * target.design.Spec.FuelCapacity
	lost := Cargo{}
	if ships == 0 {
		lost = fleet.Cargo
		fleet.Fuel = 0
	} else if capacity+lostCapacity > 0 {
		lost = fleet.Cargo.Multiply(float64(lostCapacity) / float64(capacity+lostCapacity))
	}
	if ships > 0 && fuelCapacity+lostFuelCapacity > 0 {
		fleet.Fuel -= int(float64(fleet.Fuel) * float64(lostFuelCapacity) / float64(fuelCapacity+lostFuelCapacity))
	}
	fleet.Cargo = fleet.Cargo.Subtract(lost)
	b.record.Stats.CargoLostByPlayer[target.PlayerNum] = b.record.Stats.CargoLostByPlayer[target.PlayerNum].Add(lost)
	return lost
}

// salvageRecovery returns the share of wreckage recovered in deep space, at a
// planet, or at a planet whose starbase is still in the battle.
func (b *battle) salvageRecovery() float64 {
	if b.planet == nil {
		return 0.75
	}
	for _, token := range b.tokens {
		if token.attributes&battleTokenAttributeStarbase != 0 && token.Quantity > 0 {
			return 0.8
		}
	}
	return 0.5
}

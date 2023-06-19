package cs

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

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
// TODO:
// * Accuracy/Beam Defense modifications to attractiveness
// * Beam Bonus
// * Capital Ship missiles doing double damage to capital ships
// * Unit test FireTorpedo and FireBeamWeapon
type battler interface {
	HasTargets() bool
	RunBattle() *BattleRecord
}

// battle defines the state of a battle as it progresses
type battle struct {
	Num               int
	Planet            *Planet
	Position          Vector
	Fleets            []*Fleet
	Tokens            []*battleToken
	Round             int
	hasTargets        bool
	SortedWeaponSlots []*battleWeaponSlot
	MoveOrder         [4][]*battleToken
	Record            *BattleRecord
	players           map[int]*Player
	rules             *Rules
	techFinder        TechFinder
}

type battleTarget string

const (
	BattleTargetNone              battleTarget = ""
	BattleTargetAny               battleTarget = "Any"
	BattleTargetStarbase          battleTarget = "Starbase"
	BattleTargetArmedShips        battleTarget = "ArmedShips"
	BattleTargetBombersFreighters battleTarget = "BombersFreighters"
	BattleTargetUnarmedShips      battleTarget = "UnarmedShips"
	BattleTargetFuelTransports    battleTarget = "FuelTransports"
	BattleTargetFreighters        battleTarget = "Freighters"
)

type battleTokenAttribute int

const (
	BattleTokenAttributeUnarmed       battleTokenAttribute = 0
	BattleTokenAttributeArmed         battleTokenAttribute = 1 << 0
	BattleTokenAttributeBomber        battleTokenAttribute = 1 << 1
	BattleTokenAttributeFreighter     battleTokenAttribute = 1 << 2
	BattleTokenAttributeStarbase      battleTokenAttribute = 1 << 3
	BattleTokenAttributeFuelTransport battleTokenAttribute = 1 << 4
)

// a token for a battle
type battleToken struct {
	BattleRecordToken
	Fleet        *Fleet
	Attributes   battleTokenAttribute
	MoveTarget   *battleToken
	TargetedBy   []*battleToken
	Position     Vector
	Destroyed    bool
	Damaged      bool
	RanAway      bool
	MovesMade    int
	Shields      int
	TotalShields int
}

func (bt *battleToken) GetDistanceAway(position Vector) int {
	return int(math.Max(math.Abs(bt.Position.X-position.X), math.Abs(bt.Position.Y-position.Y)))
}

func (bt *battleToken) GetDistanceAwayWeaponSlot(weaponSlot *battleWeaponSlot) int {
	return bt.GetDistanceAway(weaponSlot.Token.Position)
}

func (bt *battleToken) String() string {
	return fmt.Sprintf("%d %s (%d)", bt.Fleet.PlayerNum, bt.Token.design.Name, bt.Token.Quantity)
}

type battleWeaponType int

const (
	BattleWeaponTypeBeam battleWeaponType = iota
	BattleWeaponTypeTorpedo
)

// A token firing weapons
type battleWeaponSlot struct {
	// The token with the weapon
	Token *battleToken

	// The weapon slot
	Slot ShipDesignSlot

	// The type of weapon this weapon slot is
	WeaponType battleWeaponType

	// Each weapon has a potential list of targets that is updated each turn
	// This list is sorted by attractiveness
	Targets []*battleToken

	// The range of this weapon
	Range int

	// the power of the weapon
	Power int

	// the accuracy of the weapon, if it's a torpedo
	Accuracy int

	// the initiative of the weapon
	Initiative int
}

// Return true if this weapon slot is in range of the token target
func (slot *battleWeaponSlot) isInRange(target *battleToken) bool {
	if target == nil {
		return false
	}
	return slot.isInRangePosition(target.Position)
}

func (slot *battleWeaponSlot) isInRangePosition(position Vector) bool {
	// diagonal shots count as one move, so we take the max distance on the x or y as our actual distance away
	// i.e. 4 over, 1 up is 4 range away, 3 over 2 up is 3 range away, etc.
	return slot.Token.GetDistanceAway(position) <= slot.Range
}

func (slot *battleWeaponSlot) isInRangeValue(rangeValue int) bool {
	return rangeValue <= slot.Range
}

var positionsByPlayer = []Vector{
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

// BuildBattle builds a battle recording with all the battle tokens for a list of fleets that contains more than one player.
// We'll use this to determine if a battle should take place at this location.
// Also, any players that have a potential battle will discover each other's designs.
func newBattler(rules *Rules, techFinder TechFinder, players map[int]*Player, fleets []*Fleet, planet *Planet) battler {
	if len(fleets) == 0 {
		log.Error().Msg("Can't build battle with no fleets.")
		return nil
	}

	// add each fleet's token to the battle
	tokens := []*battleToken{}
	tokenRecords := []BattleRecordToken{}
	num := 0
	for _, fleet := range fleets {
		for _, token := range fleet.Tokens {
			num++
			battleToken := &battleToken{
				BattleRecordToken: BattleRecordToken{
					Num:       num,
					PlayerNum: fleet.PlayerNum,
					Token:     token,
				},
				Fleet:        fleet,
				Shields:      token.Quantity * token.design.Spec.Shield,
				TotalShields: token.Quantity * token.design.Spec.Shield,
				Attributes:   getBattleTokenAttributes(token),
			}
			tokens = append(tokens, battleToken)
			tokenRecords = append(tokenRecords, battleToken.BattleRecordToken)
		}
	}

	return &battle{
		Planet:     planet,
		Position:   fleets[0].Position,
		Tokens:     tokens,
		Record:     newBattleRecord(tokenRecords),
		players:    players,
		rules:      rules,
		techFinder: techFinder,
	}
}

// newBattleWeaponSlot creates a new BattleWeaponSlot object
func newBattleWeaponSlot(token *battleToken, slot ShipDesignSlot, hc TechHullComponent, rangeBonus int) *battleWeaponSlot {
	weaponSlot := &battleWeaponSlot{
		Token:      token,
		Slot:       slot,
		Targets:    []*battleToken{},
		Range:      hc.Range + rangeBonus,
		Power:      hc.Power,
		Accuracy:   hc.Accuracy,
		Initiative: hc.Initiative,
	}

	if hc.Category == TechCategoryBeamWeapon {
		weaponSlot.WeaponType = BattleWeaponTypeBeam
	} else if hc.Category == TechCategoryTorpedo {
		weaponSlot.WeaponType = BattleWeaponTypeTorpedo
	}

	return weaponSlot
}

func (b *battle) HasTargets() bool {
	for _, token := range b.getRemainingTokens() {
		if !token.Token.design.Spec.HasWeapons {
			continue
		}
		if b.getTarget(token, b.getRemainingTokens()) != nil {
			return true
		}
	}
	return false
}

// RunBattle runs a battle!
func (b *battle) RunBattle() *BattleRecord {
	planet := b.Planet
	if planet != nil {
		log.Info().Msgf("Running a battle at %s involving %d players and %d tokens.", b.Planet.Name, len(b.players), len(b.Tokens))
	} else {
		log.Info().Msgf("Running a battle at (%.2f, %.2f) involving %d players and %d tokens.", b.Position.X, b.Position.Y, len(b.players), len(b.Tokens))
	}

	b.placeTokensOnBoard()
	b.buildMovementOrder()
	for b.Round = 1; b.Round <= b.rules.NumBattleRounds; b.Round++ {
		b.Record.RecordNewRound()

		// each round we build the SortedWeaponSlots list
		// anew to account for ships that were destroyed
		b.buildSortedWeaponSlots()

		// find new targets
		b.findMoveTargets()
		if b.hasTargets {
			// if we still have targets, process the round

			// movement is a repeating pattern of 4 movement blocks
			// which we figured out in BuildMovement
			roundBlock := b.Round % 4
			for _, token := range b.MoveOrder[roundBlock] {
				if !token.Destroyed && !token.RanAway {
					b.moveToken(token)
				}
			}

			// iterate over each weapon and fire if they have a target
			for _, weaponSlot := range b.SortedWeaponSlots {
				// find all available targets for this weapon
				b.findTargets(weaponSlot)
				b.fireWeaponSlot(weaponSlot)
			}

			// at the end of this round, regenerate shields
			for _, token := range b.Tokens {
				if !token.Destroyed && !token.RanAway && token.Shields > 0 {
					b.regenerateShields(token)
				}
			}
		} else {
			// no one has targets, we are done
			break
		}
	}

	// tell players about the battle
	for _, player := range b.players {
		_ = player
		// Message.battle(game.Players[playerID], b.Planet, b.Position, playerRecord)
		// messager.battle(player, b.Planet, b.Position, b.Record)
	}

	return b.Record
}

func (battle *battle) placeTokensOnBoard() {
	tokensByPlayer := make(map[int][]*battleToken)
	for _, token := range battle.Tokens {
		tokensByPlayer[token.PlayerNum] = append(tokensByPlayer[token.PlayerNum], token)
	}
	playerIndex := 0
	for _, tokens := range tokensByPlayer {
		for _, token := range tokens {
			token.Position = positionsByPlayer[playerIndex]
			battle.Record.RecordStartingPosition(token.Num, token.Position)
		}
		playerIndex++
		if playerIndex >= len(positionsByPlayer) {
			log.Warn().Int("BattleNum", battle.Num).Msg("Oh noes! We have a battle with more players than we have positions for...")
			playerIndex = 0
		}
	}
}

// getRemainingTokens returns an enumerable of remaining tokens
func (b *battle) getRemainingTokens() []*battleToken {
	remaining := []*battleToken{}
	for _, token := range b.Tokens {
		if !(token.Destroyed || token.RanAway) {
			remaining = append(remaining, token)
		}
	}
	return remaining
}

// convert hulltype to BattleTokenAttributes
func getBattleTokenAttributes(token ShipToken) battleTokenAttribute {
	attributes := BattleTokenAttributeUnarmed
	hullType := token.design.Spec.HullType

	if hullType == TechHullTypeStarbase {
		attributes |= BattleTokenAttributeStarbase
	}

	if token.design.Spec.HasWeapons {
		attributes |= BattleTokenAttributeArmed
	}

	if hullType == TechHullTypeFreighter {
		attributes |= BattleTokenAttributeFreighter
	}

	if hullType == TechHullTypeFuelTransport {
		attributes |= BattleTokenAttributeFuelTransport
	}

	if hullType == TechHullTypeBomber {
		attributes |= BattleTokenAttributeBomber
	}

	return attributes
}

// Find all the targets for a weapon
func (b *battle) findTargets(weapon *battleWeaponSlot) {
	weapon.Targets = weapon.Targets[:0]
	attacker := weapon.Token
	primaryTarget := attacker.Fleet.battlePlan.PrimaryTarget
	secondaryTarget := attacker.Fleet.battlePlan.SecondaryTarget

	var primaryTargets []*battleToken
	var secondaryTargets []*battleToken

	// Find all enemy tokens
	for _, token := range b.getRemainingTokens() {
		if !b.willAttack(attacker.Fleet, b.players[attacker.PlayerNum], token.PlayerNum) {
			continue
		}

		// if we will target this
		if b.willTarget(primaryTarget, token) && weapon.isInRange(token) {
			primaryTargets = append(primaryTargets, token)
		} else if b.willTarget(secondaryTarget, token) && weapon.isInRange(token) {
			secondaryTargets = append(secondaryTargets, token)
		}
	}

	// our list of available targets is all primary and all secondary targets in range
	sort.Slice(primaryTargets, func(i, j int) bool {
		return weapon.getAttractiveness(primaryTargets[i]) < weapon.getAttractiveness(primaryTargets[j])
	})
	sort.Slice(secondaryTargets, func(i, j int) bool {
		return weapon.getAttractiveness(secondaryTargets[i]) < weapon.getAttractiveness(secondaryTargets[j])
	})

	weapon.Targets = append(weapon.Targets, primaryTargets...)
	weapon.Targets = append(weapon.Targets, secondaryTargets...)
}

// findMoveTargets allocates targets for each token in a battle.
func (b *battle) findMoveTargets() {
	b.hasTargets = false
	for _, token := range b.getRemainingTokens() {
		if !token.Token.design.Spec.HasWeapons {
			continue
		}
		token.MoveTarget = b.getTarget(token, b.getRemainingTokens())
		if token.MoveTarget != nil {
			token.MoveTarget.TargetedBy = append(token.MoveTarget.TargetedBy, token)
		}
		b.hasTargets = token.MoveTarget != nil || b.hasTargets
	}
}

// buildMovementOrder builds a list of Movers in this battle.
// Each ship moves in order of mass with heavier ships moving first.
// Ships that can move 3 times in a round move first, then ships that move 2 times, then 1.
func (b *battle) buildMovementOrder() {
	// our tokens are moved by mass
	tokensByMass := make([]*battleToken, 0)
	for _, token := range b.Tokens {
		if token.Token.design.Spec.Movement > 0 { // starbases don't move
			tokensByMass = append(tokensByMass, token)
		}
	}
	sort.Slice(tokensByMass, func(i, j int) bool {
		return tokensByMass[i].Token.design.Spec.Mass > tokensByMass[j].Token.design.Spec.Mass
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
				movement := token.Token.design.Spec.Movement

				// see if this token can move on this moveNum (i.e. move 1, 2, or 3)
				if movementByRound[movement-2][roundBlock] > moveNum {
					b.MoveOrder[roundBlock] = append(b.MoveOrder[roundBlock], token)
				}
			}
		}
	}
}

// Fire the weapon slot towards its target
func (b *battle) fireWeaponSlot(weapon *battleWeaponSlot) {
	if len(weapon.Targets) == 0 || weapon.Token.Destroyed || weapon.Token.RanAway {
		// no targets, nothing to do
		return
	}

	switch weapon.WeaponType {
	case BattleWeaponTypeBeam:
		b.fireBeamWeapon(weapon)
	case BattleWeaponTypeTorpedo:
		b.fireTorpedo(weapon)
	}
}

func (b *battle) fireBeamWeapon(weapon *battleWeaponSlot) {
	attackerShipToken := weapon.Token.Token
	damage := weapon.Power * weapon.Slot.Quantity * attackerShipToken.Quantity
	remainingDamage := float64(damage)

	log.Debug().Msgf("%v is firing at %v targets for a total of %v damage", weapon.Token, len(weapon.Targets), damage)

	for _, target := range weapon.Targets {
		if target.Destroyed || target.RanAway {
			continue
		}
		if remainingDamage == 0 {
			break
		}

		targetShipToken := target.Token
		shields := target.Shields
		armor := targetShipToken.design.Spec.Armor

		distance := weapon.Token.GetDistanceAway(target.Position)
		rangedDamage := int(float64(remainingDamage) * (1 - b.rules.BeamRangeDropoff*float64(distance)/float64(weapon.Range)))

		log.Debug().Msgf("%v fired %v %v(s) at %v (shields: %v, armor: %v, distance: %v, %v@%v damage) for %v (range adjusted to %v)", weapon.Token, weapon.Slot.Quantity, weapon.Slot.HullComponent, target, shields, armor, distance, targetShipToken.Quantity, targetShipToken.Damage, remainingDamage, rangedDamage)

		if rangedDamage > shields {
			remainingDamage = float64(rangedDamage - shields)
			target.Shields = 0

			existingDamage := targetShipToken.Damage * float64(targetShipToken.QuantityDamaged)
			remainingDamage += existingDamage

			numDestroyed := int(remainingDamage / float64(armor))
			if numDestroyed >= targetShipToken.Quantity {
				numDestroyed = targetShipToken.Quantity
				targetShipToken.Quantity = 0
				remainingDamage -= float64(armor * numDestroyed)

				target.Destroyed = true
				log.Debug().Msgf("%v %v %v(s) hit %v, did %v shield damage and %v armor damage and completely destroyed %v", weapon.Token, weapon.Slot.Quantity, weapon.Slot.HullComponent, target, shields, rangedDamage-shields, target)

				b.Record.RecordBeamFire(b.Round, weapon.Token, weapon.Token.Position, target.Position, weapon.Slot.HullSlotIndex, target, shields, rangedDamage-shields, numDestroyed)
			} else {
				if numDestroyed > 0 {
					targetShipToken.Quantity -= numDestroyed
				}

				remainingDamage -= float64(armor * numDestroyed)

				if remainingDamage > 0 {
					targetShipToken.Damage = remainingDamage / float64(targetShipToken.Quantity)
					targetShipToken.QuantityDamaged = targetShipToken.Quantity
					remainingDamage = 0
					log.Debug().Msgf("%v destroyed %v ships, leaving %v damaged %v@%v damage", weapon.Token, numDestroyed, target, targetShipToken.Quantity, targetShipToken.Damage)
				}

				b.Record.RecordBeamFire(b.Round, weapon.Token, weapon.Token.Position, target.Position, weapon.Slot.HullSlotIndex, target, shields, rangedDamage-shields, numDestroyed)
			}

		} else {
			target.Shields -= rangedDamage
			log.Debug().Msgf("%v firing %v %v(s) did %v damage to %v shields, leaving %v shields still operational.", weapon.Token, weapon.Slot.Quantity, weapon.Slot.HullComponent, rangedDamage, target, target.Shields)
			b.Record.RecordBeamFire(b.Round, weapon.Token, weapon.Token.Position, target.Position, weapon.Slot.HullSlotIndex, target, rangedDamage, 0, 0)
		}
		log.Debug().Msgf("%v %v %v(s) has %v remaining dp to burn through %v additional targets.", weapon.Token, weapon.Slot.Quantity, weapon.Slot.HullComponent, remainingDamage, len(weapon.Targets)-1)

		target.Damaged = true
	}
}

// Fire a torpedo slot from a ship. Torpedos are different than beam weapons
// A ship will fire each torpedo at its target until the target is destroyed, then
// fire remaining torpedos at the next target.
// Each torpedo has an accuracy rating. That determines if it hits. A torpedo that
// misses still explodes and does 1/8th damage to shields
func (b *battle) fireTorpedo(weapon *battleWeaponSlot) {
	attackerShipToken := weapon.Token.Token
	// damage is power * number of weapons * number of attackers.
	damage := weapon.Power
	torpedoInaccuracyFactor := weapon.Token.Token.design.Spec.TorpedoInaccuracyFactor
	accuracy := (100.0 - (100.0-float64(weapon.Accuracy))*torpedoInaccuracyFactor) / 100.0
	numTorpedos := weapon.Slot.Quantity * attackerShipToken.Quantity
	remainingTorpedos := float64(numTorpedos)

	log.Debug().Msgf("%s is firing at %d targets with %d torpedos at %.2f%% accuracy for %d damage each",
		weapon.Token, len(weapon.Targets), numTorpedos, accuracy*100.0, damage)

	torpedoNum := 0
	for _, target := range weapon.Targets {
		if target.Destroyed || target.RanAway {
			// this token isn't valid anymore, skip it
			continue
		}

		// no more damage to spread, break out
		if remainingTorpedos == 0 {
			break
		}

		targetShipToken := target.Token
		// shields are shared among all tokens
		shields := target.Shields
		armor := targetShipToken.design.Spec.Armor

		totalShieldDamage := 0
		totalArmorDamage := 0
		hits := 0
		misses := 0
		shipsDestroyed := 0

		for remainingTorpedos > 0 && !target.Destroyed {
			// fire a torpedo
			torpedoNum++
			remainingTorpedos--
			hit := accuracy > float64(rand.Float64())

			if hit {
				hits++
				shieldDamage := float64(0.5) * float64(damage)
				armorDamage := float64(0.5) * float64(damage)

				// apply up to half our damage to shields
				// anything leftover goes to armor
				afterShieldsDamaged := float64(shields) - shieldDamage
				actualShieldDamage := float64(0)
				if afterShieldsDamaged < 0 {
					// We did more damage to shields than they had remaining
					// apply the difference to armor
					armorDamage += float64(-afterShieldsDamaged)
					actualShieldDamage = shieldDamage + afterShieldsDamaged
				} else {
					actualShieldDamage = shieldDamage
				}
				target.Shields -= int(actualShieldDamage)

				// this torpedo blew up a ship, hooray!
				if armorDamage+targetShipToken.Damage/float64(targetShipToken.QuantityDamaged) >= float64(armor) {
					targetShipToken.Quantity--
					if targetShipToken.Quantity <= 0 {
						// record that we destroyed this token
						target.Destroyed = true
						log.Debug().Msgf("%v torpedo number %v hit %v, did %v shield damage and %v armor damage and completely destroyed %v", weapon.Token, torpedoNum, target, actualShieldDamage, armorDamage, target)

						totalShieldDamage += int(actualShieldDamage)
						totalArmorDamage += int(armorDamage)
						shipsDestroyed++
					}
				} else {
					// damage all remaining ships in this token
					// leave at least 1dp per token. We can't destroy more than one token with a torpedo
					// it's possible a torpedo blast into 100 ships will leave them all severely damaged
					// but will only destroy one
					previousDamage := targetShipToken.Damage
					targetShipToken.Damage = float64(math.Min(float64(armor*targetShipToken.Quantity-targetShipToken.Quantity), float64(armorDamage+targetShipToken.Damage)))
					targetShipToken.QuantityDamaged = targetShipToken.Quantity
					actualArmorDamage := int(targetShipToken.Damage - previousDamage)

					log.Debug().Msgf("%v torpedo number %v hit %v, did %v shield damage and %v armor damage leaving %v@%v damage", weapon.Token, torpedoNum, target, actualShieldDamage, actualArmorDamage, targetShipToken.Quantity, targetShipToken.Damage)

					totalShieldDamage += int(actualShieldDamage)
					totalArmorDamage += int(targetShipToken.Damage - previousDamage)
				}
			} else {
				misses++
				// damage shields by 1/8th
				// round up, do a minimum of 1 damage
				shieldDamage := int(math.Min(1, math.Round(b.rules.TorpedoSplashDamage*float64(damage))))
				actualShieldDamage := shieldDamage
				if shieldDamage > target.Shields {
					actualShieldDamage = target.Shields
				}
				target.Shields = int(math.Max(0, float64(target.Shields-shieldDamage)))
				log.Debug().Msgf("%s torpedo number %d missed %s, did %d damage to shields leaving %d shields", weapon.Token, torpedoNum, target, shieldDamage, target.Shields)

				totalShieldDamage += actualShieldDamage
			}
		}
		b.Record.RecordTorpedoFire(b.Round, weapon.Token, weapon.Token.Position, target.Position, weapon.Slot.HullSlotIndex, target, totalShieldDamage, totalArmorDamage, shipsDestroyed, hits, misses)
	}
}

// moveToken moves a token towards or away from its target
// TODO: figure out moving away/random
func (b *battle) moveToken(token *battleToken) {
	// count this token's moves
	token.MovesMade++
	if token.MoveTarget == nil || !token.Token.design.Spec.HasWeapons {
		// tokens with no weapons always run away
		b.runAway(token)
		return
	}

	// we have weapons, figure out our tactic and targets
	switch token.Fleet.battlePlan.Tactic {
	case BattleTacticDisengage:
		b.runAway(token)
	case BattleTacticDisengageIfChallenged:
		if token.Damaged {
			b.runAway(token)
		} else {
			b.maximizeDamage(token)
		}
	case BattleTacticMinimizeDamageToSelf, BattleTacticMaximizeNetDamage, BattleTacticMaximizeDamageRatio, BattleTacticMaximizeDamage:
		b.maximizeDamage(token)
	}
}

// regenerateShields regenerates the shields of the given token if the player regenerates shields
// and the token has shields.
func (b *battle) regenerateShields(token *battleToken) {
	player := b.players[token.PlayerNum]

	if player.Race.Spec.ShieldRegenerationRate > 0 && token.Shields > 0 {
		regenerationAmount := int(float64(token.TotalShields)*player.Race.Spec.ShieldRegenerationRate + 0.5)
		token.Shields = int(clamp(token.Shields+regenerationAmount, 0, token.TotalShields))
	}
}

func (b *battle) maximizeDamage(token *battleToken) {
	if token.MoveTarget != nil {
		newPosition := token.Position
		if token.Position.Y > token.MoveTarget.Position.Y {
			newPosition.Y--
		} else {
			newPosition.Y++
		}
		if token.Position.X > token.MoveTarget.Position.X {
			newPosition.X--
		} else {
			newPosition.X++
		}

		// we can't move off board
		newPosition.X = float64(clamp(int(newPosition.X), 0, 9))
		newPosition.Y = float64(clamp(int(newPosition.Y), 0, 9))

		// create a move record for the viewer and then move the token
		b.Record.RecordMove(b.Round, token, token.Position, newPosition)
		token.Position = newPosition
	}
}

func (b *battle) runAway(token *battleToken) {

	if token.MovesMade >= b.rules.MovesToRunAway {
		token.RanAway = true
		b.Record.RecordRunAway(b.Round, token)
	}

	// if we are in range of a weapon, move away, otherwise move randomly
	weaponsInRange := make([]*battleWeaponSlot, 0)
	for _, weapon := range b.SortedWeaponSlots {
		if weapon.isInRange(token) {
			weaponsInRange = append(weaponsInRange, weapon)
		}
	}

	possiblePositions := []Vector{
		token.Position.Add(VectorRight),
		token.Position.Add(VectorLeft),
		token.Position.Add(VectorDown),
		token.Position.Add(VectorUp),
		token.Position.Add(VectorUp).Add(VectorRight),
		token.Position.Add(VectorUp).Add(VectorLeft),
		token.Position.Add(VectorDown).Add(VectorRight),
		token.Position.Add(VectorDown).Add(VectorLeft),
	}

	var newPosition Vector
	if len(weaponsInRange) > 0 {
		// default to move to a random position
		newPosition = possiblePositions[b.rules.random.Intn(len(possiblePositions))]

		// move to a position that is out of range, or to the greatest distance away we can get
		maxNumWeaponsInRange := math.MinInt32
		for _, possiblePosition := range possiblePositions {
			// can't move here
			if possiblePosition.X < 0 || possiblePosition.X > 9 || possiblePosition.Y < 0 || possiblePosition.Y > 9 {
				continue
			}
			numWeaponsInRange := 0
			for _, weapon := range weaponsInRange {
				distanceAway := weapon.Token.GetDistanceAway(possiblePosition)
				if weapon.isInRangeValue(distanceAway) {
					numWeaponsInRange++
					if distanceAway > maxNumWeaponsInRange {
						maxNumWeaponsInRange = distanceAway
						newPosition = possiblePosition
					}
				}
			}

			// no weapons in range of this position, move there
			if numWeaponsInRange == 0 {
				newPosition = possiblePosition
				break
			}
		}

		// we can't move off board (this should never be a problem)
		newPosition.X = float64(clamp(int(newPosition.X), 0, 9))
		newPosition.Y = float64(clamp(int(newPosition.Y), 0, 9))

		// move to our new position
		b.Record.RecordMove(b.Round, token, token.Position, newPosition)
		token.Position = newPosition
	} else {
		// move at random
		newPosition = possiblePositions[b.rules.random.Intn(len(possiblePositions))]
		b.Record.RecordMove(b.Round, token, token.Position, newPosition)
		token.Position = newPosition
	}
}

// getTarget returns the primary or secondary target based on the attacker's battle plan and the defenders available
func (b *battle) getTarget(attacker *battleToken, defenders []*battleToken) *battleToken {
	var primaryTarget *battleToken
	var secondaryTarget *battleToken
	var primaryTargetAttractiveness float64
	var secondaryTargetAttractiveness float64
	primaryTargetOrder := attacker.Fleet.battlePlan.PrimaryTarget
	secondaryTargetOrder := attacker.Fleet.battlePlan.SecondaryTarget

	// TODO: We need to account for the fact that if a fleet targets us, we will target them back
	for _, defender := range defenders {
		if !(defender.Destroyed || defender.RanAway) && b.willAttack(attacker.Fleet, b.players[attacker.PlayerNum], defender.PlayerNum) {
			// if we would target this defender with our primary target and it's more attractive than our current primaryTarget, pick it
			if b.willTarget(primaryTargetOrder, defender) {
				attractiveness := defender.getAttractiveness(attacker)
				if attractiveness >= primaryTargetAttractiveness {
					primaryTarget = defender
					primaryTargetAttractiveness = attractiveness
				}
			}

			// if we would target this defender with our secondary target, pick it
			if b.willTarget(secondaryTargetOrder, defender) {
				attractiveness := defender.getAttractiveness(attacker)
				if attractiveness >= secondaryTargetAttractiveness {
					secondaryTarget = defender
					secondaryTargetAttractiveness = attractiveness
				}
			}
		}
	}

	if primaryTarget != nil {
		return primaryTarget
	}

	return secondaryTarget
}

func (defender *battleToken) getAttractiveness(attacker *battleToken) float64 {
	cost := defender.Token.design.Spec.Cost.MultiplyInt(defender.Token.Quantity)
	defense := (defender.Token.design.Spec.Armor + defender.Token.design.Spec.Shield) * defender.Token.Quantity

	// TODO: change defense based on attacker weapons

	return float64(float64(cost.Germanium+cost.Resources) / float64(defense))
}

func (weapon *battleWeaponSlot) getAttractiveness(target *battleToken) float64 {
	cost := target.Token.design.Spec.Cost.MultiplyInt(target.Token.Quantity)
	defense := (target.Token.design.Spec.Armor + target.Token.design.Spec.Shield) * target.Token.Quantity

	// TODO: change defense based on attacker weapons

	return float64(float64(cost.Germanium+cost.Resources) / float64(defense))
}

// return true if this fleet will attack a fleet by another player based on player
// relations and the fleet battle plan
func (b *battle) willAttack(fleet *Fleet, player *Player, otherPlayerNum int) bool {
	willAttack := false
	// if we have weapons and we don't own this other fleet, see if we
	// would target it
	if fleet.Spec.HasWeapons && fleet.battlePlan.Tactic != BattleTacticDisengage && otherPlayerNum != player.Num {
		switch fleet.battlePlan.AttackWho {
		case BattleAttackWhoEnemies:
			willAttack = player.IsEnemy(otherPlayerNum)
		case BattleAttackWhoEnemiesAndNeutrals:
			willAttack = player.IsEnemy(otherPlayerNum) || player.IsNeutral(otherPlayerNum)
		case BattleAttackWhoEveryone:
			willAttack = true
		}
	}
	return willAttack
}

// willTarget returns true if the BattleOrder Target type would target this token
func (b *battle) willTarget(target battleTarget, token *battleToken) bool {
	switch target {
	case BattleTargetAny:
		return true
	case BattleTargetNone:
		return false
	case BattleTargetStarbase:
		return (token.Attributes & BattleTokenAttributeStarbase) > 0
	case BattleTargetArmedShips:
		return (token.Attributes & BattleTokenAttributeArmed) > 0
	case BattleTargetBombersFreighters:
		return (token.Attributes&BattleTokenAttributeBomber) > 0 || (token.Attributes&BattleTokenAttributeFreighter) > 0
	case BattleTargetUnarmedShips:
		return (token.Attributes & BattleTokenAttributeArmed) == 0
	case BattleTargetFuelTransports:
		return (token.Attributes & BattleTokenAttributeFuelTransport) > 0
	case BattleTargetFreighters:
		return (token.Attributes & BattleTokenAttributeFreighter) > 0
	}

	return false
}

func (b *battle) buildSortedWeaponSlots() {
	b.SortedWeaponSlots = make([]*battleWeaponSlot, 0)
	for _, token := range b.Tokens {
		hull := b.techFinder.GetHull(token.Token.design.Hull)
		for _, slot := range token.Token.design.Spec.WeaponSlots {
			bws := newBattleWeaponSlot(token, slot, *b.techFinder.GetHullComponent(slot.HullComponent), hull.RangeBonus)
			b.SortedWeaponSlots = append(b.SortedWeaponSlots, bws)
		}
	}
	sort.Slice(b.SortedWeaponSlots, func(i, j int) bool {
		return b.SortedWeaponSlots[i].Initiative < b.SortedWeaponSlots[j].Initiative
	})
}

package cs

import (
	"fmt"
	"math"
)

// Starbases with Packet Throwers can build mineral packets and fling them at other planets.
type MineralPacket struct {
	GameDBObject
	MapObject
	TargetPlanetNum   int           `json:"targetPlanetNum"`
	Cargo             Cargo         `json:"cargo"`
	WarpSpeed         int           `json:"warpSpeed"`
	SafeWarpSpeed     int           `json:"safeWarpSpeed"`
	Heading           VectorFloat64 `json:"heading"`
	ScanRange         int           `json:"scanRange"`
	ScanRangePen      int           `json:"scanRangePen"`
	distanceTravelled float64
	builtThisTurn     bool
}

type MineralPacketDamage struct {
	Killed            int `json:"killed,omitempty"`
	DefensesDestroyed int `json:"defensesDestroyed,omitempty"`
	Uncaught          int `json:"uncaught,omitempty"`
}

// this mineral packet will decay to nothing before reaching its target
const MineralPacketDecayToNothing = -1

func newMineralPacket(player *Player, num int, warpSpeed int, safeWarpSpeed int, cargo Cargo, position Vector, targetPlanetNum int) *MineralPacket {
	packet := MineralPacket{
		Type:            MapObjectTypeMineralPacket,
		PlayerNum:       player.Num,
		Num:             num,
		Name:            fmt.Sprintf("%s Mineral Packet", player.Race.PluralName),
		Position:        position,
		WarpSpeed:       warpSpeed,
		SafeWarpSpeed:   safeWarpSpeed,
		Cargo:           cargo,
		TargetPlanetNum: targetPlanetNum,
		ScanRange:       NoScanner,
		ScanRangePen:    NoScanner,
	}

	// PP packets have built in scanners
	if player.Race.Spec.PacketBuiltInScanner {
		packet.ScanRangePen = warpSpeed * warpSpeed
	}

	return &packet
}

// get the rate of decay for a packet between 0 and 1
//
// Depending on how fast a packet is thrown compared to its safe speed, it decays
// Source: https://wiki.starsautohost.org/wiki/%22Mass_Packet_FAQ%22_by_Barry_Kearns_1997-02-07_v2.6b
func (packet *MineralPacket) getPacketDecayRate(rules *Rules, race *Race) float64 {

	// we only care about packets thrown up to 3 warps over the limit
	overSafeWarp := min(packet.WarpSpeed-packet.SafeWarpSpeed, 3)

	// IT is always counted as being 1 more over the safe warp
	overSafeWarp = min(race.Spec.PacketOverSafeWarpPenalty+overSafeWarp, 3)

	packetDecayRate := 0.0
	if overSafeWarp > 0 {
		packetDecayRate = rules.PacketDecayRate[overSafeWarp]
	}

	// PP have half the decay rate
	packetDecayRate *= race.Spec.PacketDecayFactor

	return packetDecayRate
}

// decay the packet's minerals for flying some fraction of a year. Packets thrown over
// their safe speed lose a percentage of each mineral, with a minimum loss per mineral.
func (packet *MineralPacket) decay(rules *Rules, race *Race, fraction float64) {
	decayRate := packet.getPacketDecayRate(rules, race) * fraction
	if decayRate == 0 {
		return
	}

	for _, minType := range MineralTypes {
		mineral := packet.Cargo.GetAmount(minType)
		if mineral > 0 {
			decayAmount := max(int(decayRate*float64(mineral)), int(float64(rules.PacketMinDecay)*race.Spec.PacketDecayFactor))
			packet.Cargo = packet.Cargo.SubtractAmount(minType, decayAmount)
		}
	}
	packet.Cargo = packet.Cargo.MinZero()
}

// move this packet through space
func (packet *MineralPacket) movePacket(rules *Rules, player *Player, target *Planet, planetPlayer *Player) {
	dist := float64(packet.WarpSpeed * packet.WarpSpeed)
	totalDist := packet.Position.DistanceTo(target.Position)

	// move at half distance if this packet was created this turn
	if packet.builtThisTurn {
		dist /= 2
	}

	// round up, if we are <1 LY away, i.e. the target is 81.9 ly away, warp 9 (81 ly travel) should be able to make it there
	if dist < totalDist && totalDist-dist < 1 {
		dist = math.Ceil(totalDist)
	}

	vectorTravelled := target.Position.Subtract(packet.Position).Normalized().Scale(dist)
	dist = vectorTravelled.Length()

	// don't overshoot
	dist = min(totalDist, dist)

	if totalDist == dist {
		// packets decay for the part of the year they travelled before they impact
		packet.decay(rules, &player.Race, min(1, dist/float64(packet.WarpSpeed*packet.WarpSpeed)))
		if packet.Cargo.Total() == 0 {
			// decayed to nothing before reaching the planet
			packet.Delete = true
			return
		}
		packet.completeMove(rules, player, target, planetPlayer)
	} else {
		// move this packet closer to the next planet
		packet.distanceTravelled = dist
		packet.Heading = target.Position.Subtract(packet.Position).Normalized()
		packet.Position = packet.Position.ToFloat64().Add(packet.Heading.Scale(dist)).ToInt(true)
	}
}

// Complete movement of an incoming packet about to impact the planet
//
// Example:
// You fling a 1000kT packet at Warp 10 at a planet with a Warp 5 driver, a population of 250,000 and 50 defenses preventing 60% of incoming damage.
// packet power = 10 x 10 = 100
// receiver power = 5 x 5 = 25 (halved for IT)
// caught = 25 / 100 = 250 thousandths
// minerals recovered = 1000kT x 250/1000 + 1000kT x 750/1000 x 1/9 = 250 + 83 = 333kT
// raw damage = (100 - 25) x 1000 / 160 = 468
// damage = 468 x 40% = 187
// #colonists killed = Max. of (187 x 2,500kT / 1000, 187) x 100
// = Max. of (467, 187) x 100 = 46,700 colonists
// #defenses destroyed = Max. of (50 x 187 / 1000, 187 / 20) = Max. of (9, 9) = 9
// (if this rounds to 0, there is a damage/20 chance of destroying one defense)
//
// If, however, the receiving planet had no mass driver or defenses, the damage is far greater:
// minerals recovered = 1000kT x 0 + 1000kT x 1/9 = only 111kT
// raw damage = damage = 100 x 1000 / 160 = 625
// #colonists killed = Max. of (625 x 2,500kT / 1000, 625) x 100
// = Max. of (1,562, 625) x 100 = 156,200 colonists
//
// If the packet increased speed up to Warp 13, then:
// damage = 169 x 1000 / 160 = 1056
// #colonists killed = Max. of (1056 x 2,500kT / 1000, 1056) x 100
// = Max. of (2,640, 1056) x 100 = 264,000, destroying the colony
func (packet *MineralPacket) completeMove(rules *Rules, player *Player, planet *Planet, planetPlayer *Player) {
	// Capture the receiver before casualties can remove its colony/starbase.
	caughtPermille := packet.caughtPermille(planet, planetPlayer)
	damage := packet.getDamage(planet, planetPlayer)
	if planetPlayer != nil && !planetPlayer.Race.Spec.LivesOnStarbases && damage.Killed < planet.GetPopulation() && planet.Defenses > 0 {
		strength := packet.impactStrength(planet, planetPlayer)
		if strength > 0 && planet.Defenses*strength/1000 == 0 && rules.random.Intn(20) < strength {
			damage.DefensesDestroyed = max(1, damage.DefensesDestroyed)
		}
	}

	// Recovery is quantized in thousandths; the original recovers one ninth
	// of the uncaught material, despite the manual's one-third description.
	keepPermille := caughtPermille + (1000-caughtPermille)/9
	for _, mineral := range MineralTypes {
		planet.Cargo = planet.Cargo.AddAmount(mineral, packet.Cargo.GetAmount(mineral)*keepPermille/1000)
	}
	packet.checkTerraform(rules, player, planet, 1000-caughtPermille)

	if damage == (MineralPacketDamage{}) {
		// caught packet successfully, transfer cargo
		messager.planetPacketCaught(planetPlayer, planet, packet)
	} else if planetPlayer != nil {
		// kill off colonists and defenses
		// note, for AR races, this will be 0 colonists killed or structures destroyed
		planet.addPopulation(-roundTo100(damage.Killed, math.Round))
		planet.Defenses = Clamp(planet.Defenses-damage.DefensesDestroyed, 0, planet.Defenses)

		messager.planetPacketDamage(planetPlayer, planet, packet, damage)
		if planet.GetPopulation() <= 0 {
			planet.emptyPlanet()
			messager.planetDiedOff(planetPlayer, planet)
		}
	}

	// if we didn't receive this planet, notify the sender
	if planet.PlayerNum != packet.PlayerNum {
		if player.Race.Spec.DetectPacketDestinationStarbases && planet.Spec.HasStarbase {
			// discover the receiving planet's starbase design
			player.discoverer.discoverDesign(planet.Starbase.Tokens[0].design, true)
		}

		messager.planetPacketArrived(player, planet, packet)
	}

	// delete the packet
	packet.Delete = true
}

// receiverPower is the square of the planet's mass driver speed. IT drivers
// have half the normal catching power, even at their nominal safe speed.
func receiverPower(planet *Planet, receiver *Player) int {
	if !planet.Owned() || !planet.Spec.HasMassDriver {
		return 0
	}
	power := planet.Spec.SafePacketSpeed * planet.Spec.SafePacketSpeed
	return int(float64(power) * receiver.Race.Spec.PacketReceiverFactor)
}

// caughtPermille is the thousandths of the packet the receiving planet catches safely
func (packet *MineralPacket) caughtPermille(planet *Planet, receiver *Player) int {
	if packet.WarpSpeed <= 0 {
		return 0
	}
	return min(1000, receiverPower(planet, receiver)*1000/(packet.WarpSpeed*packet.WarpSpeed))
}

// get the damage a mineral packet will do when it collides with a planet
func (packet *MineralPacket) getDamage(planet *Planet, planetPlayer *Player) MineralPacketDamage {
	if !planet.Owned() {
		// unowned planets aren't damaged, but all cargo is uncaught
		return MineralPacketDamage{Uncaught: packet.Cargo.Total()}
	}

	caughtPermille := packet.caughtPermille(planet, planetPlayer)
	if caughtPermille == 1000 {
		// planet will successfully catch the packet
		return MineralPacketDamage{}
	}

	uncaught := (1000 - caughtPermille) * packet.Cargo.Total() / 1000
	if planetPlayer.Race.Spec.LivesOnStarbases {
		return MineralPacketDamage{Uncaught: uncaught}
	}
	damageWithDefenses := packet.impactStrength(planet, planetPlayer)
	colonistsKilled := max(planet.Cargo.Colonists*damageWithDefenses/1000, damageWithDefenses) * 100
	defensesDestroyed := max(planet.Defenses*damageWithDefenses/1000, damageWithDefenses/20)

	// kill off colonists and destroy defenses, up to however much actually exists
	return MineralPacketDamage{
		Killed:            min(colonistsKilled, planet.GetPopulation()),
		DefensesDestroyed: min(planet.Defenses, defensesDestroyed),
		Uncaught:          uncaught,
	}

}

// impactStrength is the damage of an uncaught packet after planetary defenses
func (packet *MineralPacket) impactStrength(planet *Planet, planetPlayer *Player) int {
	rawDamage := (packet.WarpSpeed*packet.WarpSpeed - receiverPower(planet, planetPlayer)) * packet.Cargo.Total() / 160
	return max(0, int(float64(rawDamage)*(1-planet.Spec.DefenseCoverage)))
}

// Estimate potential damage of an incoming mineral packet.
// Simulates decay each turn until impact
func (packet *MineralPacket) estimateDamage(rules *Rules, player *Player, target *Planet, planetPlayer *Player) MineralPacketDamage {
	if packet.caughtPermille(target, planetPlayer) == 1000 {
		// planet will have no problem catching the packet on arrival
		return MineralPacketDamage{}
	}

	distPerYear := float64(packet.WarpSpeed * packet.WarpSpeed)
	// save copy of packet so we don't alter the original
	packetCopy := *packet

	// decay the packet each year until impact, including the partial year it arrives
	// (packets within 1ly of the target after a year's travel make it there, see movePacket)
	for totalDist := packet.Position.DistanceTo(target.Position); totalDist > 0; totalDist -= distPerYear {
		packetCopy.decay(rules, &player.Race, min(1, totalDist/distPerYear))

		// packet out of minerals; return special exit code
		if packetCopy.Cargo.Total() == 0 {
			return MineralPacketDamage{Uncaught: MineralPacketDecayToNothing}
		}

		if totalDist < distPerYear+1 {
			// arrives this year
			break
		}
	}

	// calculate damage based on however much cargo is left
	damage := packetCopy.getDamage(target, planetPlayer)

	// clear packet uncaught statistic as we don't care about it
	// (this is a **damage** test function after all)
	damage.Uncaught = 0

	return damage
}

// checkTerraform checks if an uncaught PP packet terraforms the target planet.
// PP makes a 50% roll per 100 kT of each uncaught mineral. A successful
// terraform roll has a further 10% chance of changing the underlying habitat.
// Each mineral type affects one habitat axis: ironium -> grav, boranium -> temp,
// germanium -> rad. Changes always move toward the sender's ideal habitat (or
// outward, to the extremes, for an axis the sender is immune to).
func (packet *MineralPacket) checkTerraform(rules *Rules, player *Player, planet *Planet, uncaughtPermille int) {
	unit := player.Race.Spec.PacketPermaTerraformSizeUnit
	chance := player.Race.Spec.PacketTerraformChance
	if chance <= 0 || unit <= 0 || uncaughtPermille <= 0 {
		return
	}
	terraformer := NewTerraformer()
	for i, mineralType := range MineralTypes {
		habType := HabType(i)

		// roll once per 100 kT chunk of uncaught mineral (a partial chunk gets a
		// proportionally smaller chance). Each hit is one click of temporary
		// terraforming, and some hits are also a permanent click
		temporary, permanent := 0, 0
		for mineral := packet.Cargo.GetAmount(mineralType) * uncaughtPermille / 1000; mineral > 0; mineral -= unit {
			rollChance := chance * float64(min(mineral, unit)) / float64(unit)
			if rules.random.Float64() >= rollChance {
				continue
			}
			temporary++
			if rules.random.Float64() < player.Race.Spec.PacketPermaformChance {
				permanent++
			}
		}
		// pick a direction: toward the sender's ideal, or, if the sender is immune
		// to this axis, away from the middle. Immune senders only get half the
		// temporary clicks
		immune := player.Race.IsImmune(habType)
		direction := 1
		ideal := player.Race.HabCenter().Get(habType)
		if immune {
			if planet.BaseHab.Get(habType) < 50 {
				direction = -1
			}
			temporary /= 2
		} else if planet.BaseHab.Get(habType) > ideal {
			direction = -1
		}
		// permanent clicks move the planet's base habitat, stopping at the sender's
		// ideal (immune axes stop only at the min/max hab)
		originalBase := planet.BaseHab.Get(habType)
		newBase := originalBase + direction*permanent
		if !immune {
			newBase = originalBase + direction*min(permanent, Abs(ideal-originalBase))
		}
		newBase = Clamp(newBase, rules.MinHab, rules.MaxHab)
		planet.BaseHab.Set(habType, newBase)
		if change := newBase - originalBase; change != 0 {
			messager.planetPacketPermaform(player, planet, habType, change)
		}
		// temporary clicks move the current habitat, like the sender's own
		// terraforming would. They need the sender to have terraforming tech, can't go
		// past the sender's ideal or the sender's terraforming range from the
		// (possibly new) base, and never undo terraforming that's already further along
		if temporary > 0 {
			current := planet.Hab.Get(habType)
			target := current
			ability := terraformer.GetTerraformAbility(player).Get(habType)
			if immune {
				if terraformer.GetTerraformAbility(player).absSum() > 0 {
					target = Clamp(current+direction*temporary, rules.MinHab, rules.MaxHab)
				}
			} else if current < ideal && ability > 0 {
				limit := min(ideal, min(rules.MaxHab, newBase+ability))
				if current < limit {
					target = min(limit, current+temporary)
				}
			} else if current > ideal && ability > 0 {
				limit := max(ideal, max(rules.MinHab, newBase-ability))
				if current > limit {
					target = max(limit, current-temporary)
				}
			}
			planet.Hab.Set(habType, target)
			if change := target - current; change != 0 {
				messager.planetPacketTerraform(player, planet, habType, change)
			}
		}
	}
	// keep the terraformed amount in sync with any base habitat changes
	planet.TerraformedAmount = planet.Hab.Subtract(planet.BaseHab)
	planet.MarkDirty()
}

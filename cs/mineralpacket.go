package cs

import (
	"fmt"
	"math"
)

// Starbases with Packet Throwers can build mineral packets and fling them at other planets.
type MineralPacket struct {
	MapObject
	TargetPlanetNum   int    `json:"targetPlanetNum,omitempty"`
	Cargo             Cargo  `json:"cargo,omitempty"`
	WarpSpeed         int    `json:"warpSpeed"`
	SafeWarpSpeed     int    `json:"safeWarpSpeed,omitempty"`
	Heading           Vector `json:"heading"`
	ScanRange         int    `json:"scanRange"`
	ScanRangePen      int    `json:"scanRangePen"`
	distanceTravelled float64
	builtThisTurn     bool
}

type MineralPacketDamage struct {
	Killed            int `json:"killed,omitempty"`
	DefensesDestroyed int `json:"defensesDestroyed,omitempty"`
	Uncaught          int `json:"uncaught,omitempty"`
}

func newMineralPacket(player *Player, num int, warpSpeed int, safeWarpSpeed int, cargo Cargo, position Vector, targetPlanetNum int) *MineralPacket {
	packet := MineralPacket{
		MapObject: MapObject{
			Type:      MapObjectTypeMineralPacket,
			PlayerNum: player.Num,
			Num:       num,
			Dirty:     true,
			Name:      fmt.Sprintf("%s Mineral Packet", player.Race.PluralName),
			Position:  position,
		},
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
// https://wiki.starsautohost.org/wiki/%22Mass_Packet_FAQ%22_by_Barry_Kearns_1997-02-07_v2.6b
// Depending on how fast a packet is thrown compared to its safe speed, it decays
func (packet *MineralPacket) getPacketDecayRate(rules *Rules, race *Race) float64 {

	// we only care about packets thrown up to 3 warps over the limit
	overSafeWarp := MinInt(packet.WarpSpeed-packet.SafeWarpSpeed, 3)

	// IT is always counted as being 1 more over the safe warp
	overSafeWarp = MinInt(race.Spec.PacketOverSafeWarpPenalty+overSafeWarp, 3)

	packetDecayRate := 0.0
	if overSafeWarp > 0 {
		packetDecayRate = rules.PacketDecayRate[overSafeWarp]
	}

	// PP have half the decay rate
	packetDecayRate *= race.Spec.PacketDecayFactor

	return packetDecayRate
}

// move this packet through spcae
func (packet *MineralPacket) movePacket(rules *Rules, player *Player, target *Planet, planetPlayer *Player) {
	dist := float64(packet.WarpSpeed * packet.WarpSpeed)
	totalDist := packet.Position.DistanceTo(target.Position)

	// move at half distance if this packet was created this turn
	if packet.builtThisTurn {
		// half move packets...
		dist /= 2
	}
	// round up, if we are <1 away, i.e. the target is 81.9 ly away, warp 9 (81 ly travel) should be able to make it there
	if dist < totalDist && totalDist-dist < 1 {
		dist = math.Ceil(totalDist)
	}

	vectorTravelled := target.Position.Subtract(packet.Position).Normalized().Scale(dist)
	dist = vectorTravelled.Length()

	// don't overshoot
	dist = math.Min(totalDist, dist)

	if totalDist == dist {
		packet.completeMove(rules, player, target, planetPlayer)
	} else {
		// move this packet closer to the next planet
		packet.distanceTravelled = dist
		packet.Heading = target.Position.Subtract(packet.Position).Normalized()
		packet.Position = packet.Position.Add(packet.Heading.Scale(dist))
		packet.Position = packet.Position.Round()

		packet.MarkDirty()
	}
}

// Damage calcs as per the Stars! Manual
//
// Example:
// You fling a 1000kT packet at Warp 10 at a planet with a Warp 5 driver, a population of 250,000 and 50 defenses preventing 60% of incoming damage.
// spdPacket = 100
// spdReceiver = 25
// %CaughtSafely = 25%
// minerals recovered = 1000kT x 25% + 1000kT x 75% x 1/3 = 250 + 250 = 500kT
// dmgRaw = 75 x 1000 / 160 = 469
// dmgRaw2 = 469 x 40% = 188
// #colonists killed = Max. of ( 188 x 250,000 / 1000, 188 x 100)
// = Max. of ( 47,000, 18800) = 47,000 colonists
// #defenses destroyed = 50 * 188 / 1000 = 9 (rounded down)
//
// If, however, the receiving planet had no mass driver or defenses, the damage is far greater:
// minerals recovered = 1000kT x 0% + 1000kT x 100% x 1/3 = only 333kT dmgRaw = 100 x 1000 / 160 = 625
// dmgRaw2 = 625 x 100% = 625
// #colonists killed = Max. of (625 x 250,000 / 1000, 625 x 100)
// = Max.of(156,250, 62500) = 156,250.
// If the packet increased speed up to Warp 13, then:
// dmgRaw2 = dmgRaw = 169 x 1000 / 160 = 1056
// #colonists killed = Max. of (1056 x 250,000 / 1000, 1056 x 100)
// = Max.of( 264,000, 105600) destroying the colony
func (packet *MineralPacket) completeMove(rules *Rules, player *Player, planet *Planet, planetPlayer *Player) {
	damage := packet.getDamage(rules, planet, planetPlayer)

	if damage == (MineralPacketDamage{}) {
		// caught packet successfully, transfer cargo
		messager.planetPacketCaught(planetPlayer, planet, packet)
	} else if planetPlayer != nil {
		// kill off colonists and defenses
		// note, for AR races, this will be 0 colonists killed or structures destroyed
		planet.setPopulation(roundToNearest100(Clamp(planet.population()-damage.Killed, 0, planet.population())))
		planet.Defenses = Clamp(planet.Defenses-damage.DefensesDestroyed, 0, planet.Defenses)

		messager.planetPacketDamage(planetPlayer, planet, packet, damage.Killed, damage.DefensesDestroyed)
		if planet.population() == 0 {
			planet.emptyPlanet()
			messager.planetDiedOff(planetPlayer, planet)
		}
	}

	mineralsRecovered := 1.0
	if damage.Uncaught > 0 {
		packet.checkTerraform(rules, player, planet, damage.Uncaught)
		packet.checkPermaform(rules, player, planet, damage.Uncaught)
		receiverDriverSpeed := 0
		if planet.Spec.HasStarbase {
			receiverDriverSpeed = planet.Spec.SafePacketSpeed
		}
		percentCaughtSafely := float64((packet.WarpSpeed * packet.WarpSpeed) / (receiverDriverSpeed * receiverDriverSpeed))
		mineralsRecovered = percentCaughtSafely + (1-percentCaughtSafely)/3
	}

	// one way or another, these minerals are ending up on the planet
	planet.Cargo = planet.Cargo.Add(packet.Cargo.Multiply(mineralsRecovered))

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

// get the damage a mineral packet will do when it collides with a planet
func (packet *MineralPacket) getDamage(rules *Rules, planet *Planet, planetPlayer *Player) MineralPacketDamage {
	if !planet.Owned() {
		// unowned planets aren't damaged, but all cargo is uncaught
		return MineralPacketDamage{Uncaught: packet.Cargo.Total()}
	}

	if planet.Spec.HasMassDriver && planet.Spec.SafePacketSpeed >= packet.WarpSpeed {
		// planet successfully caught this packet
		return MineralPacketDamage{}
	}

	if planetPlayer != nil && planetPlayer.Race.Spec.LivesOnStarbases {
		// No damage, but all cargo is uncaught and might impact the planet
		return MineralPacketDamage{Uncaught: packet.Cargo.Total()}
	}

	// uh oh, this packet is going too fast and we'll take damage
	receiverDriverSpeed := 0
	if planet.Spec.HasStarbase {
		receiverDriverSpeed = planet.Spec.SafePacketSpeed
	}

	weight := packet.Cargo.Total()
	speedOfPacket := packet.WarpSpeed * packet.WarpSpeed
	speedOfReceiver := receiverDriverSpeed * receiverDriverSpeed
	percentCaughtSafely := float64(speedOfReceiver) / float64(speedOfPacket)
	uncaught := int((1.0 - percentCaughtSafely) * float64(weight))
	// mineralsRecovered := int(float64(weight)*percentCaughtSafely + float64(weight)*(1/3.0)*(1-percentCaughtSafely))
	rawDamage := float64((speedOfPacket-speedOfReceiver)*weight) / 160
	damageWithDefenses := rawDamage * (1 - planet.Spec.DefenseCoverage)
	colonistsKilled := roundToNearest100f(math.Max(damageWithDefenses*float64(planet.population())/1000, damageWithDefenses*100))
	defensesDestroyed := int(math.Max(float64(planet.Defenses)*damageWithDefenses/1000, damageWithDefenses/20))

	// kill off colonists and defenses
	return MineralPacketDamage{
		Killed:            roundToNearest100(MinInt(colonistsKilled, planet.population())),
		DefensesDestroyed: MinInt(planet.Defenses, defensesDestroyed),
		Uncaught:          uncaught,
	}

}

// Estimate potential damage of incoming mineral packet
// Simulates decay each turn until impact
func (packet *MineralPacket) estimateDamage(rules *Rules, player *Player, target *Planet, planetPlayer *Player) MineralPacketDamage {
	if target.Spec.HasMassDriver && target.Spec.SafePacketSpeed >= packet.WarpSpeed {
		// planet successfully caught this packet
		return MineralPacketDamage{}
	}
	spd := float64(packet.WarpSpeed * packet.WarpSpeed)
	decayRate := 0.0
	totalDist := packet.Position.DistanceTo(target.Position)
	eta := int(totalDist / spd)

	//save copy of packet so we don't alter the original
	packetCopy := *packet

	for i := 0; i <= eta; i++ {
		if totalDist <= spd {
			// 1 turn until impact - only travels/decays partially
			distTraveled := totalDist / float64(spd)
			decayRate = (packetCopy.getPacketDecayRate(rules, &player.Race) * distTraveled)
		} else {
			decayRate = packetCopy.getPacketDecayRate(rules, &player.Race)
			totalDist -= spd
		}

		// no decay, so we don't need to bother calculating decay amount
		if decayRate == 1 {
			break
		}

		// loop through all 3 mineral types and reduce each one in turn
		for _, minType := range [3]CargoType{Ironium, Boranium, Germanium} {
			mineral := packetCopy.Cargo.GetAmount(minType)

			// subtract either the normal or minimum decay amounts, whichever is higher (rounded DOWN)
			if mineral > 0 {
				decayAmount := MaxInt(int(decayRate*float64(mineral)), int(float64(rules.PacketMinDecay)*float64(player.Race.Spec.PacketDecayFactor)))
				packetCopy.Cargo.SubtractAmount(minType, decayAmount)
				}
			}
		
		// packet out of minerals
		if packetCopy.Cargo.GetAmount(Ironium) == 0 && packetCopy.Cargo.GetAmount(Boranium) == 0 && packetCopy.Cargo.GetAmount(Boranium) == 0 {
			return MineralPacketDamage{}
		}
		}
	}

	damage := packetCopy.getDamage(rules, target, planetPlayer)

	// clear packet uncaught statistic as we don't care about it (this is a **damage** test function after all)
	damage.Uncaught = 0

	return damage
}

// Check if an uncaught PP packet will terraform the target planet's environment  (50% chance/100kT)
func (packet *MineralPacket) checkTerraform(rules *Rules, player *Player, planet *Planet, uncaught int) {
	if player.Race.Spec.PacketTerraformChance > 0 {
		for uncaughtCheck := player.Race.Spec.PacketPermaTerraformSizeUnit; uncaughtCheck <= uncaught; uncaughtCheck += player.Race.Spec.PacketPermaTerraformSizeUnit {
			if player.Race.Spec.PacketTerraformChance >= rules.random.Float64() {
				terraformer := NewTerraformer()

				result := terraformer.TerraformOneStep(planet, player, nil, false)
				if result.Terraformed() {
					messager.planetPacketTerraform(player, planet, result.Type, result.Direction)
				}
			}
		}
	}
}

// Check if an uncaught PP packet will permanently alter the target planet's environment (0.1% chance/100kT)
func (p *MineralPacket) checkPermaform(rules *Rules, player *Player, planet *Planet, uncaught int) {
	if player.Race.Spec.PacketPermaformChance > 0 && player.Race.Spec.PacketPermaTerraformSizeUnit > 0 {
		for uncaughtCheck := player.Race.Spec.PacketPermaTerraformSizeUnit; uncaughtCheck <= uncaught; uncaughtCheck += player.Race.Spec.PacketPermaTerraformSizeUnit {
			if player.Race.Spec.PacketPermaformChance >= float64(rules.random.Float64()) {
				habType := HabType(rules.random.Intn(3))
				terraformer := NewTerraformer()
				result := terraformer.PermaformOneStep(planet, player, habType)
				if result.Terraformed() {
					messager.planetPacketPermaform(player, planet, result.Type, result.Direction)
				}
			}
		}
	}
}

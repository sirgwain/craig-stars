package cs

import (
	"fmt"
	"math"
)

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
// Depending on how fast a packet is thrown compared to it's safe speed, it decays
func (packet *MineralPacket) getPacketDecayRate(rules *Rules, race *Race) float64 {
	overSafeWarp := packet.WarpSpeed - packet.SafeWarpSpeed

	// IT is always count as being at least 1 over the safe warp
	overSafeWarp += race.Spec.PacketOverSafeWarpPenalty

	// we only care about packets thrown up to 3 warp over the limit
	overSafeWarp = minInt(packet.WarpSpeed-packet.SafeWarpSpeed, 3)

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

// Damage calcs the Stars! Manual
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
	damage := packet.getDamage(rules, planet)

	if damage == (MineralPacketDamage{}) {
		// caught packet successfully, transfer cargo
		messager.mineralPacketCaught(planetPlayer, planet, packet)
	} else if damage.Killed > 0 || damage.DefensesDestroyed > 0 {
		// kill off colonists and defenses
		planet.setPopulation(roundToNearest100(clamp(planet.population()-damage.Killed, 0, planet.population())))
		planet.Defenses = clamp(planet.Defenses-damage.DefensesDestroyed, 0, planet.Defenses)

		messager.mineralPacketDamage(planetPlayer, planet, packet, damage.Killed, damage.DefensesDestroyed)
		if planet.population() == 0 {
			planet.emptyPlanet()
		}
	}

	if damage.Uncaught > 0 {
		packet.checkTerraform(rules, player, planet, damage.Uncaught)
		packet.checkPermaform(rules, player, planet, damage.Uncaught)
	}

	// one way or another, these minerals are ending up on the planet
	planet.Cargo = planet.Cargo.Add(packet.Cargo)

	// if we didn't receive this planet, notify the sender
	if planet.PlayerNum != packet.PlayerNum {
		if player.Race.Spec.DetectPacketDestinationStarbases && planet.Spec.HasStarbase {
			// discover the receiving planet's starbase design
			discoverer := newDiscoverer(player)
			discoverer.discoverDesign(player, planet.Starbase.Tokens[0].design, true)
		}

		messager.mineralPacketArrived(player, planet, packet)
	}

	// delete the packet
	packet.Delete = true
}

// get the damage a mineral packet will do when it collides with a planet
func (packet *MineralPacket) getDamage(rules *Rules, planet *Planet) MineralPacketDamage {

	if !planet.Owned() {
		// unowned planets aren't damaged, but all cargo is uncaught
		return MineralPacketDamage{Uncaught: packet.Cargo.Total()}
	}

	if planet.Spec.HasMassDriver && planet.Spec.SafePacketSpeed >= packet.WarpSpeed {
		// planet successfully caught this packet
		return MineralPacketDamage{}
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
		Killed:            roundToNearest100(minInt(colonistsKilled, planet.population())),
		DefensesDestroyed: minInt(planet.Defenses, defensesDestroyed),
		Uncaught:          uncaught,
	}

}

// From mazda on starsautohost: https://starsautohost.org/sahforum2/index.php?t=msg&th=1294&start=0&rid=0
//
// For each uncaught 100kT of mineral there is 50% chance of performing 1 click of normal terraforming on the target planet (i.e. the same terra with the same limits as if you were sat on the planet spending resources).
// So a large packet of 2000kT on an unoccupied planet should expect to terrafrom it 10 clicks, which is a lot.
//
// Secondly, the same packet can also perform "permanent terraforming" by altering the underlying planet variables.
// There is a 50% chance of performing 1 click of permanent terraforming per 1000kT of uncaught material.
//
// This ability is only available to PP and CA and for CAs it is random, whereas PP's can choose what planet to permanently alter.
//
// Note these figures are for uncaught amounts, so work best on planets with no defences and especially no flingers.
//
// The direction of terraforming in all cases is towards the optimum for the flinging race.
// Whether the target is yours, friends, enemies or empty makes no difference.
//
// For an immunity I believe the terraforming is towards the closest edge.
func (packet *MineralPacket) checkTerraform(rules *Rules, player *Player, planet *Planet, uncaught int) {
	if player.Race.Spec.PacketTerraformChance > 0 {
		for uncaughtCheck := player.Race.Spec.PacketPermaTerraformSizeUnit; uncaughtCheck <= uncaught; uncaughtCheck += player.Race.Spec.PacketPermaTerraformSizeUnit {
			if player.Race.Spec.PacketTerraformChance >= rules.random.Float64() {
				terraformer := NewTerraformer()

				result := terraformer.TerraformOneStep(planet, player, nil, false)
				if result.Terraformed() {
					messager.packetTerraform(player, planet, result.Type, result.Direction)
				}
			}
		}
	}
}

func (p *MineralPacket) checkPermaform(rules *Rules, player *Player, planet *Planet, uncaught int) {
	if player.Race.Spec.PacketPermaformChance > 0 && player.Race.Spec.PacketPermaTerraformSizeUnit > 0 {
		for uncaughtCheck := player.Race.Spec.PacketPermaTerraformSizeUnit; uncaughtCheck <= uncaught; uncaughtCheck += player.Race.Spec.PacketPermaTerraformSizeUnit {
			if player.Race.Spec.PacketPermaformChance >= float64(rules.random.Float64()) {
				habType := HabType(rules.random.Intn(3))
				terraformer := NewTerraformer()
				result := terraformer.PermaformOneStep(planet, player, habType)
				if result.Terraformed() {
					messager.packetPermaform(player, planet, result.Type, result.Direction)
				}
			}
		}
	}
}

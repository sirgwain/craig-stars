//go:build !wasi && !wasm

package cs

import "slices"

// LongRangeScoutTestSlots is the standard scout design used by small test scenarios.
var LongRangeScoutTestSlots = []ShipDesignSlot{
	{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
	{HullComponent: FuelTank.Name, HullSlotIndex: 3, Quantity: 1},
}

// SantaMariaTestSlots is the standard colony-ship design used by test scenarios.
var SantaMariaTestSlots = []ShipDesignSlot{
	{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: ColonizationModule.Name, HullSlotIndex: 2, Quantity: 1},
}

// SantaMariaARTestSlots is the Alternate Reality colony-ship design used by test scenarios.
var SantaMariaARTestSlots = []ShipDesignSlot{
	{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: OrbitalConstructionModule.Name, HullSlotIndex: 2, Quantity: 1},
}

// TeamsterTestSlots is the standard medium-freighter design used by cargo test scenarios.
var TeamsterTestSlots = []ShipDesignSlot{
	{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
	{HullComponent: RhinoScanner.Name, HullSlotIndex: 3, Quantity: 1},
}

// DestroyerDeltaTestSlots is the standard destroyer design used by battle test scenarios.
var DestroyerDeltaTestSlots = []ShipDesignSlot{
	{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
	{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 2, Quantity: 1},
	{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 3, Quantity: 1},
	{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 4, Quantity: 1},
	{HullComponent: Crobmnium.Name, HullSlotIndex: 5, Quantity: 1},
	{HullComponent: ManeuveringJet.Name, HullSlotIndex: 6, Quantity: 1},
	{HullComponent: BattleComputer.Name, HullSlotIndex: 7, Quantity: 1},
}

// Designs returns independent copies, so customizing a scenario cannot change a preset.
func Designs(designs ...ShipDesign) []ShipDesign {
	result := slices.Clone(designs)
	for i := range result {
		result[i].Slots = slices.Clone(result[i].Slots)
	}
	return result
}

var (
	DesignLongRangeScout = ShipDesign{Name: "Long Range Scout", Hull: Scout.Name, Slots: LongRangeScoutTestSlots}
	DesignTeamster       = ShipDesign{Name: "Teamster", Hull: MediumFreighter.Name, Slots: TeamsterTestSlots}
	DesignSantaMaria     = ShipDesign{Name: "Santa Maria", Hull: ColonyShip.Name, Slots: SantaMariaTestSlots}
	DesignSantaMariaAR   = ShipDesign{Name: "Santa Maria", Hull: ColonyShip.Name, Slots: SantaMariaARTestSlots}
	DesignDestroyerDelta = ShipDesign{Name: "Destroyer", Hull: Destroyer.Name, Slots: DestroyerDeltaTestSlots}
	DesignThief          = ShipDesign{Name: "Thief", Hull: MediumFreighter.Name, Slots: []ShipDesignSlot{
		{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: RobberBaronScanner.Name, HullSlotIndex: 3, Quantity: 1},
	}}
	DesignSmallFreighter = ShipDesign{Name: "Small Freighter", Hull: SmallFreighter.Name, Slots: []ShipDesignSlot{
		{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: CargoPod.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: BatScanner.Name, HullSlotIndex: 3, Quantity: 1},
	}}
	DesignStealingFreighter = ShipDesign{Name: "Stealing Freighter", Hull: MediumFreighter.Name, Slots: []ShipDesignSlot{
		{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: CargoPod.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: RobberBaronScanner.Name, HullSlotIndex: 3, Quantity: 1},
	}}
	DesignGalleon = ShipDesign{Name: "Galleon", Hull: Galleon.Name, Slots: []ShipDesignSlot{
		{HullComponent: SubGalacticFuelScoop.Name, HullSlotIndex: 1, Quantity: 4},
	}}
	DesignMiniMineLayer = ShipDesign{Name: "Little Hen", Hull: MiniMineLayer.Name, Slots: []ShipDesignSlot{
		{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: MineDispenser40.Name, HullSlotIndex: 2, Quantity: 2},
		{HullComponent: MineDispenser40.Name, HullSlotIndex: 3, Quantity: 2},
		{HullComponent: BatScanner.Name, HullSlotIndex: 4, Quantity: 1},
	}}
	DesignCloakedScout = ShipDesign{Name: "Cloaked Scout", Hull: Scout.Name, Slots: []ShipDesignSlot{
		{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: StealthCloak.Name, HullSlotIndex: 3, Quantity: 1},
	}}
	DesignRemoteTerraformer = ShipDesign{Name: "Remote Terraformer", Hull: MiniMiner.Name, Slots: []ShipDesignSlot{
		{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: BatScanner.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: OrbitalAdjuster.Name, HullSlotIndex: 3, Quantity: 1},
		{HullComponent: OrbitalAdjuster.Name, HullSlotIndex: 4, Quantity: 1},
	}}
	DesignGatePrivateer = ShipDesign{Name: "Gate Privateer", Hull: Privateer.Name, Slots: []ShipDesignSlot{
		{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: JumpGate.Name, HullSlotIndex: 3, Quantity: 1},
	}}
	DesignPotatoBug = ShipDesign{Name: "Potato Bug", Hull: MidgetMiner.Name, Slots: []ShipDesignSlot{
		{HullComponent: QuickJump5.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: RoboMidgetMiner.Name, HullSlotIndex: 2, Quantity: 1},
	}}
	DesignSantaMariaIFE = ShipDesign{Name: "Santa Maria", Hull: ColonyShip.Name, Slots: []ShipDesignSlot{
		{HullComponent: FuelMizer.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: ColonizationModule.Name, HullSlotIndex: 2, Quantity: 1},
	}}
	DesignStalwartDefender = ShipDesign{Name: "Stalwart Defender", Hull: Destroyer.Name, Slots: []ShipDesignSlot{
		{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: BetaTorpedo.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: XRayLaser.Name, HullSlotIndex: 3, Quantity: 1},
		{HullComponent: RhinoScanner.Name, HullSlotIndex: 4, Quantity: 1},
		{HullComponent: Crobmnium.Name, HullSlotIndex: 5, Quantity: 1},
		{HullComponent: Overthruster.Name, HullSlotIndex: 6, Quantity: 1},
		{HullComponent: BattleComputer.Name, HullSlotIndex: 7, Quantity: 1},
	}}
	DesignJihadCruiser = ShipDesign{Name: "Jihad Cruiser", Hull: Cruiser.Name, Slots: []ShipDesignSlot{
		{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 2},
		{HullComponent: Overthruster.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: BattleNexus.Name, HullSlotIndex: 3, Quantity: 1},
		{HullComponent: JihadMissile.Name, HullSlotIndex: 4, Quantity: 2},
		{HullComponent: JihadMissile.Name, HullSlotIndex: 5, Quantity: 2},
		{HullComponent: ElephantScanner.Name, HullSlotIndex: 6, Quantity: 2},
		{HullComponent: Kelarium.Name, HullSlotIndex: 7, Quantity: 2},
	}}
	DesignPrivateer = ShipDesign{Name: "Privateer", Hull: Privateer.Name, Slots: []ShipDesignSlot{
		{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: CargoPod.Name, HullSlotIndex: 3, Quantity: 1},
		{HullComponent: CargoPod.Name, HullSlotIndex: 4, Quantity: 1},
		{HullComponent: CargoPod.Name, HullSlotIndex: 5, Quantity: 1},
	}}
	DesignSpaceStation = ShipDesign{Name: "Starbase", Hull: SpaceStation.Name, Slots: []ShipDesignSlot{
		{HullComponent: Laser.Name, HullSlotIndex: 2, Quantity: 8},
		{HullComponent: MoleSkinShield.Name, HullSlotIndex: 3, Quantity: 8},
		{HullComponent: Laser.Name, HullSlotIndex: 4, Quantity: 8},
		{HullComponent: MoleSkinShield.Name, HullSlotIndex: 6, Quantity: 8},
		{HullComponent: Laser.Name, HullSlotIndex: 8, Quantity: 8},
		{HullComponent: Laser.Name, HullSlotIndex: 10, Quantity: 8},
	}}
	DesignDeathStar  = ShipDesign{Name: "Starbase", Hull: DeathStar.Name, Slots: []ShipDesignSlot{}}
	DesignBattleBeam = ShipDesign{Name: "Beam Destroyer", Hull: Destroyer.Name, Slots: []ShipDesignSlot{
		{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 3, Quantity: 1},
		{HullComponent: GorillaDelagator.Name, HullSlotIndex: 4, Quantity: 1},
		{HullComponent: Crobmnium.Name, HullSlotIndex: 5, Quantity: 1},
		{HullComponent: Overthruster.Name, HullSlotIndex: 6, Quantity: 1},
		{HullComponent: EnergyCapacitor.Name, HullSlotIndex: 7, Quantity: 1},
	}}
	DesignBattleTorpedo = ShipDesign{Name: "Torpedo Destroyer", Hull: Destroyer.Name, Slots: []ShipDesignSlot{
		{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 3, Quantity: 1},
		{HullComponent: GorillaDelagator.Name, HullSlotIndex: 4, Quantity: 1},
		{HullComponent: Crobmnium.Name, HullSlotIndex: 5, Quantity: 1},
		{HullComponent: Overthruster.Name, HullSlotIndex: 6, Quantity: 1},
		{HullComponent: BattleSuperComputer.Name, HullSlotIndex: 7, Quantity: 1},
	}}
	DesignBattleStarbase = ShipDesign{Name: "Battle Starbase", Hull: SpaceStation.Name, Slots: []ShipDesignSlot{
		{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 2, Quantity: 8},
		{HullComponent: GorillaDelagator.Name, HullSlotIndex: 3, Quantity: 8},
		{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 4, Quantity: 8},
		{HullComponent: GorillaDelagator.Name, HullSlotIndex: 6, Quantity: 8},
		{HullComponent: Jammer30.Name, HullSlotIndex: 7, Quantity: 3},
		{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 8, Quantity: 8},
		{HullComponent: BattleSuperComputer.Name, HullSlotIndex: 9, Quantity: 3},
		{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 10, Quantity: 8},
	}}
	DesignBattleChaff = ShipDesign{Name: "Chaff", Hull: Scout.Name, Slots: []ShipDesignSlot{
		{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: XRayLaser.Name, HullSlotIndex: 2, Quantity: 1},
	}}
	DesignBattleCruiser = ShipDesign{Name: "Battle Cruiser", Hull: BattleCruiser.Name, Slots: []ShipDesignSlot{
		{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 2},
		{HullComponent: Overthruster.Name, HullSlotIndex: 2, Quantity: 2},
		{HullComponent: BattleSuperComputer.Name, HullSlotIndex: 3, Quantity: 2},
		{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 4, Quantity: 3},
		{HullComponent: DeltaTorpedo.Name, HullSlotIndex: 5, Quantity: 3},
		{HullComponent: Overthruster.Name, HullSlotIndex: 6, Quantity: 3},
		{HullComponent: GorillaDelagator.Name, HullSlotIndex: 7, Quantity: 4},
	}}
	DesignArmoredFreighter = ShipDesign{Name: "Armored Freighter", Hull: SmallFreighter.Name, Slots: []ShipDesignSlot{
		{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: Crobmnium.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: RhinoScanner.Name, HullSlotIndex: 3, Quantity: 1},
	}}
	DesignShieldedScout = ShipDesign{Name: "Shielded Scout", Hull: Scout.Name, Slots: []ShipDesignSlot{
		{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: RhinoScanner.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: CompletePhaseShield.Name, HullSlotIndex: 3, Quantity: 1},
	}}
	DesignJammedDefender = ShipDesign{Name: "Jammed Defender", Hull: Destroyer.Name, Slots: []ShipDesignSlot{
		{HullComponent: TransStar10.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: ColloidalPhaser.Name, HullSlotIndex: 3, Quantity: 1},
		{HullComponent: RhinoScanner.Name, HullSlotIndex: 4, Quantity: 1},
		{HullComponent: Superlatanium.Name, HullSlotIndex: 5, Quantity: 1},
		{HullComponent: Jammer30.Name, HullSlotIndex: 6, Quantity: 1},
		{HullComponent: FluxCapacitor.Name, HullSlotIndex: 7, Quantity: 1},
	}}
	DesignStalwartSapper = ShipDesign{Name: "Stalwart Sapper", Hull: Destroyer.Name, Slots: []ShipDesignSlot{
		{HullComponent: LongHump6.Name, HullSlotIndex: 1, Quantity: 1},
		{HullComponent: PulsedSapper.Name, HullSlotIndex: 2, Quantity: 1},
		{HullComponent: PulsedSapper.Name, HullSlotIndex: 3, Quantity: 1},
		{HullComponent: RhinoScanner.Name, HullSlotIndex: 4, Quantity: 1},
		{HullComponent: Superlatanium.Name, HullSlotIndex: 5, Quantity: 1},
		{HullComponent: Overthruster.Name, HullSlotIndex: 6, Quantity: 1},
		{HullComponent: Overthruster.Name, HullSlotIndex: 7, Quantity: 1},
	}}
)

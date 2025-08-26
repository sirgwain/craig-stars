/**
 * warpspeed for using a stargate vs moving with warp drive
 */
export const StargateWarpSpeed = 11;

/**
 * time period to perform a task, like patrol
 */
export const Indefinite = 0;

/**
 * use automatic warp speed for patrols
 */
export const PatrolWarpSpeedAutomatic = 0;

/**
 * target fleets in any range when patrolling
 */
export const PatrolRangeInfinite = 0;

/**
 * no target planet, player, etc
 */
export const None = 0;
export const UnlimitedSpaceDock = -1;
export const NoScanner = -1;
export const NoGate = -1;
export const Infinite = -1;
export const InfiniteGate = 2147483647; /* math.MaxInt32 */
export const ReportAgeUnexplored = -1;
export const Unowned = 0;

/**
 * this mineral packet will decay to nothing before reaching its target
 */
export const MineralPacketDecayToNothing = -1;

export type VictoryCondition = number;
export const VictoryConditionNone = 0;
export const VictoryConditionOwnPlanets: VictoryCondition = 1 << (1 - 1);
export const VictoryConditionAttainTechLevels: VictoryCondition = 1 << (2 - 1);
export const VictoryConditionExceedsScore: VictoryCondition = 1 << (3 - 1);
export const VictoryConditionExceedsSecondPlaceScore: VictoryCondition = 1 << (4 - 1);
export const VictoryConditionProductionCapacity: VictoryCondition = 1 << (5 - 1);
export const VictoryConditionOwnCapitalShips: VictoryCondition = 1 << (6 - 1);
export const VictoryConditionHighestScoreAfterYears: VictoryCondition = 1 << (7 - 1);

export type HullSlotType = number;
export const HullSlotTypeNone = 0;
export const HullSlotTypeEngine: HullSlotType = 1 << 1;
export const HullSlotTypeScanner: HullSlotType = 1 << 2;
export const HullSlotTypeMechanical: HullSlotType = 1 << 3;
export const HullSlotTypeBomb: HullSlotType = 1 << 4;
export const HullSlotTypeMining: HullSlotType = 1 << 5;
export const HullSlotTypeElectrical: HullSlotType = 1 << 6;
export const HullSlotTypeShield: HullSlotType = 1 << 7;
export const HullSlotTypeArmor: HullSlotType = 1 << 8;
export const HullSlotTypeCargo: HullSlotType = 1 << 9;
export const HullSlotTypeSpaceDock: HullSlotType = 1 << 10;
export const HullSlotTypeWeapon: HullSlotType = 1 << 11;
export const HullSlotTypeOrbital: HullSlotType = 1 << 12;
export const HullSlotTypeMineLayer: HullSlotType = 1 << 13;
export const HullSlotTypeElectricalMechanical = HullSlotTypeElectrical | HullSlotTypeMechanical;
export const HullSlotTypeOrbitalElectrical = HullSlotTypeOrbital | HullSlotTypeElectrical;
export const HullSlotTypeShieldElectricalMechanical =
	HullSlotTypeShield | HullSlotTypeElectrical | HullSlotTypeMechanical;
export const HullSlotTypeScannerElectricalMechanical =
	HullSlotTypeScanner | HullSlotTypeElectrical | HullSlotTypeMechanical;
export const HullSlotTypeArmorScannerElectricalMechanical =
	HullSlotTypeArmor | HullSlotTypeScanner | HullSlotTypeElectrical | HullSlotTypeMechanical;
export const HullSlotTypeMineElectricalMechanical =
	HullSlotTypeMineLayer | HullSlotTypeElectrical | HullSlotTypeMechanical;
export const HullSlotTypeShieldArmor = HullSlotTypeShield | HullSlotTypeArmor;
export const HullSlotTypeWeaponShield = HullSlotTypeShield | HullSlotTypeWeapon;
export const HullSlotTypeGeneral =
	HullSlotTypeScanner |
	HullSlotTypeMechanical |
	HullSlotTypeElectrical |
	HullSlotTypeShield |
	HullSlotTypeArmor |
	HullSlotTypeWeapon |
	HullSlotTypeMineLayer;

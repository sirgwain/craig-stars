import type { Cost } from './Cost';
import type { Size } from './Game';
import type { MineFieldStats, MineFieldType } from './MineField';
import type { FreighterGrowth, LRT, PRT, SpendLeftoverPointsOn } from './Race';
import type { TechLevel } from './TechLevel';
import rulesjson from '$lib/ssr/rules.json';

export type RandomEvent = 'Comet' | 'MineralDeposit' | 'PlanetaryChange' | 'AncientArtifact';

export type CometSize = 'Small' | 'Medium' | 'Large' | 'Huge';

export type RepairRate =
	| 'None'
	| 'Moving'
	| 'Stopped'
	| 'Orbiting'
	| 'OrbitingOwnPlanet'
	| 'Starbase';

export type WormholeStability =
	| 'RockSolid'
	| 'Stable'
	| 'MostlyStable'
	| 'Average'
	| 'SlightlyVolatile'
	| 'Volatile'
	| 'ExtremelyVolatile';

export type StartingFleetHull =
	| ''
	| 'Colony Ship'
	| 'Destroyer'
	| 'Medium Freighter'
	| 'Mini-Colony Ship'
	| 'Mini Mine Layer'
	| 'Mini-Miner'
	| 'Midget-Miner'
	| 'Privateer'
	| 'Scout';
export type ShipDesignPurpose =
	| ''
	| 'Scout'
	| 'Colonizer'
	| 'Bomber'
	| 'StructureBomber'
	| 'SmartBomber'
	| 'Fighter'
	| 'FighterScout'
	| 'CapitalShip'
	| 'Freighter'
	| 'ColonistFreighter'
	| 'FuelFreighter'
	| 'MultiPurposeFreighter'
	| 'ArmedFreighter'
	| 'Miner'
	| 'Terraformer'
	| 'DamageMineLayer'
	| 'SpeedMineLayer'
	| 'Starbase'
	| 'FuelDepot'
	| 'StarbaseQuarter'
	| 'StarbaseHalf'
	| 'PacketThrower'
	| 'Stargater'
	| 'Fort'
	| 'StarterColony';

export type Rules = {
	id?: number;
	createdAt?: string; // ISO 8601 string for time
	updatedAt?: string;
	gameId?: number;
	cometStatsBySize?: Record<CometSize, CometStats>;
	fleetSafeSpeedExplosionChance?: number;
	invasionDefenseCoverageFactor?: number;
	lrtSpecs?: Partial<Record<LRT, LRTSpec>>;
	maxPopulation?: number;
	maxTechLevel?: number;
	mineFieldCloak?: number;
	mineFieldStatsByType?: Record<MineFieldType, MineFieldStats>;
	mineralDecayFactor?: number;
	minHabFloor?: number;
	movesToRunAway?: number;
	mysteryTraderRules?: MysteryTraderRules;
	packetDecayRate?: Record<number, number>;
	packetMaxOverwarpSpeed?: number;
	packetMinDecay?: number;
	planetMinDistance?: number;
	populationOvercrowdDieoffRate?: number;
	populationOvercrowdDieoffRateMax?: number;
	populationOvercrowdResourcePenalty?: number;
	populationOvercrowdResourceMax?: number;
	populationScannerError?: number;
	prtSpecs?: Partial<Record<PRT, PRTSpec>>;
	raceStartingPoints?: number;
	radiatingImmune?: number;
	randomArtifactResearchBonusRange?: number[];
	randomCometMinYear?: number;
	randomCometMinYearPlayerWorld?: number;
	randomEventChances?: Record<RandomEvent, number>;
	randomMineralDepositBonusRange?: number[];
	remoteMiningMineOutput?: number;
	repairRates?: Record<RepairRate, number>;
	salvageDecayMin?: number;
	salvageDecayRate?: number;
	salvageFromBattleFactor?: number;
	scrapMineralAmount?: number;
	scrapResourceAmount?: number;
	showPublicScoresAfterYears?: number;
	smartDefenseCoverageFactor?: number;
	stargateMaxHullMassFactor?: number;
	stargateMaxRangeFactor?: number;
	tachyonCloakReduction?: number;
	tachyonMaxCloakReduction?: number;
	techsId?: number;
	techTradeChance?: number;
	torpedoSplashDamage?: number;
	wormholeCloak?: number;
	wormholePairsForSize?: Record<Size, number>;
	wormholeStatsByStability?: Record<WormholeStability, WormholeStats>;
} & CostRules &
	BattleRules &
	UniverseGenerationRules;

export type UniverseGenerationRules = {
	highRadMineralConcentrationBonusThreshold?: number;
	limitMineralConcentration?: number;
	maxExtraWorldDistance?: number;
	maxHab?: number;
	maxMineralConcentration?: number;
	maxStartingMineralConcentration?: number;
	maxStartingMineralSurface?: number;
	minExtraPlanetMineralConcentration?: number;
	minExtraWorldDistance?: number;
	minHab?: number;
	minHomeworldMineralConcentration?: number;
	minMineralConcentration?: number;
	minStartingMineralConcentration?: number;
	minStartingMineralSurface?: number;
	raceLeftoverPointsPerItem?: Record<SpendLeftoverPointsOn, number>;
	startingYear?: number;
	wormholeMinPlanetDistance?: number;
};

export type CostRules = {
	defenseCost?: Cost;
	factoryCostGermanium?: number;
	mineralAlchemyCost?: number;
	planetaryScannerCost?: Cost;
	starbaseComponentCostReduction?: number;
	terraformCost?: Cost;
	techBaseCost?: number[];
};

export type BattleRules = {
	beamRangeDropoff?: number;
	numBattleRounds?: number;
};

export type CometStats = {
	allMinerals?: number;
	allRandomMinerals?: number;
	bonusMinerals?: number;
	bonusRandomMinerals?: number;
	bonusMinConcentration?: number;
	bonusRandomConcentration?: number;
	bonusAffectsMinerals?: number;
	minTerraform?: number;
	randomTerraform?: number;
	affectsHabs?: number;
	popKilledPercent?: number;
};

export type MysteryTraderRules = {
	chanceSpawn?: number[];
	chanceMaxTechGetsPart?: number;
	chanceCourseChange?: number;
	chanceSpeedUpOnly?: number;
	chanceAgain?: number;
	minYear?: number;
	evenYearOnly?: boolean;
	minWarp?: number;
	maxWarp?: number;
	maxMysteryTraders?: number;
	requestedBoon?: number;
	genesisDeviceCost?: Cost;
	techBoon?: MysteryTraderTechBoonRules[];
};

export type MysteryTraderTechBoonRules = {
	techLevels?: number;
	rewards?: MysteryTraderTechBoonMineralsReward[];
};

export type MysteryTraderTechBoonMineralsReward = {
	mineralsGiven?: number;
	reward?: number;
};

export type WormholeStats = {
	yearsToDegrade: number;
	chanceToJump: number;
	jiggleDistance: number;
};

export type PRTSpec = {
	prt?: PRT;
	pointCost?: number;
	startingTechLevels?: TechLevel;
	startingPlanets?: StartingPlanet[];
	techCostOffset?: TechCostOffset;
	mineralsPerSingleMineralPacket?: number;
	mineralsPerMixedMineralPacket?: number;
	packetResourceCost?: number;
	packetMineralCostFactor?: number;
	packetReceiverFactor?: number;
	packetDecayFactor?: number;
	packetOverSafeWarpPenalty?: number;
	packetBuiltInScanner?: boolean;
	detectPacketDestinationStarbases?: boolean;
	detectAllPackets?: boolean;
	packetTerraformChance?: number;
	packetPermaformChance?: number;
	packetPermaTerraformSizeUnit?: number;
	canGateCargo?: boolean;
	canDetectStargatePlanets?: boolean;
	shipsVanishInVoid?: boolean;
	builtInScannerMultiplier?: number;
	techsCostExtraLevel?: number;
	freighterGrowth?: FreighterGrowth;
	growthFactor?: number;
	maxPopulationOffset?: number;
	builtInCloakUnits?: number;
	stealsResearch?: StealsResearch;
	freeCargoCloaking?: boolean;
	mineFieldsAreScanners?: boolean;
	mineFieldRateMoveFactor?: number;
	mineFieldSafeWarpBonus?: number;
	mineFieldMinDecayFactor?: number;
	mineFieldBaseDecayRate?: number;
	mineFieldPlanetDecayRate?: number;
	mineFieldMaxDecayRate?: number;
	canDetonateMineFields?: boolean;
	mineFieldDetonateDecayRate?: number;
	discoverDesignOnScan?: boolean;
	canRemoteMineOwnPlanets?: boolean;
	invasionAttackBonus?: number;
	invasionDefendBonus?: number;
	movementBonus?: number;
	instaforming?: boolean;
	permaformChance?: number;
	permaformPopulation?: number;
	repairFactor?: number;
	starbaseRepairFactor?: number;
	starbaseCostFactor?: number;
	innateMining?: boolean;
	innateMinesFactor?: number;
	innateResources?: boolean;
	innateScanner?: boolean;
	innateScannerFactor?: number;
	canBuildDefenses?: boolean;
	livesOnStarbases?: boolean;
	minHabFloor?: number;
};

export type LRTSpec = {
	lrt?: LRT;
	startingFleets?: StartingFleet[];
	pointCost?: number;
	startingTechLevels?: TechLevel;
	techCostOffset?: TechCostOffset;
	newTechCostFactorOffset?: number;
	miniaturizationMax?: number;
	miniaturizationPerLevel?: number;
	noAdvancedScanners?: boolean;
	scanRangeFactorOffset?: number;
	fuelEfficiencyOffset?: number;
	maxPopulationOffset?: number;
	mineralAlchemyCostOffset?: number;
	scrapMineralOffset?: number;
	scrapMineralOffsetStarbase?: number;
	scrapResourcesOffset?: number;
	scrapResourcesOffsetStarbase?: number;
	startingPopulationFactorDelta?: number;
	starbaseBuiltInCloakUnits?: number;
	starbaseCostFactor?: number;
	researchFactorOffset?: number;
	researchSplashDamage?: number;
	shieldStrengthFactorOffset?: number;
	shieldRegenerationRateOffset?: number;
	armorStrengthFactorOffset?: number;
	engineFailureRateOffset?: number;
	engineReliableSpeed?: number;
};

export type TechCostOffset = {
	engine?: number;
	beamWeapon?: number;
	torpedo?: number;
	bomb?: number;
	planetaryDefense?: number;
	stargate?: number;
	terraforming?: number;
};

export type StartingPlanet = {
	population?: number;
	mines?: number;
	factories?: number;
	defenses?: number;
	habPenaltyFactor?: number;
	hasStargate?: boolean;
	hasMassDriver?: boolean;
	starbaseDesignName?: string;
	starbaseHull?: string;
	startingFleets?: StartingFleet[];
	homeworld?: boolean;
};

export type StartingFleet = {
	name?: string;
	hullName?: StartingFleetHull;
	hullSetNumber?: number;
	purpose?: ShipDesignPurpose;
};

export type StealsResearch = {
	energy?: number;
	weapons?: number;
	propulsion?: number;
	construction?: number;
	electronics?: number;
	biotechnology?: number;
};

export const defaultRules: Rules = rulesjson as Rules;

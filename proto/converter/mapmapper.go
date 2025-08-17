package converter

import (
	cs "github.com/sirgwain/craig-stars/cs"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
)

func MinefieldTypeMapToIntMap(source map[cs.MinefieldType]int) map[int32]int32 {
	m := make(map[int32]int32, len(source))
	for k, v := range source {
		m[int32(CSMinefieldTypeToMinefieldType(k))] = int32(v)
	}

	return m
}

func IntMapToMinefieldTypeMap(source map[int32]int32) map[cs.MinefieldType]int {
	m := make(map[cs.MinefieldType]int, len(source))
	for k, v := range source {
		m[MinefieldTypeToCSMinefieldType(craig_starsv1.MinefieldType(k))] = int(v)
	}

	return m
}

func MinefieldTypeToInt32(source cs.MinefieldType) int32 {
	return int32(CSMinefieldTypeToMinefieldType(source))
}

func Int32ToMinefieldType(source int32) cs.MinefieldType {
	return MinefieldTypeToCSMinefieldType(craig_starsv1.MinefieldType(source))
}

func SpendLeftoverPointsOnMapToIntMap(source map[cs.SpendLeftoverPointsOn]int) map[int32]int32 {
	m := make(map[int32]int32, len(source))
	for k, v := range source {
		m[int32(CSSpendLeftoverPointsOnToSpendLeftoverPointsOn(k))] = int32(v)
	}
	return m
}

func IntMapToSpendLeftoverPointsOnMap(source map[int32]int32) map[cs.SpendLeftoverPointsOn]int {
	m := make(map[cs.SpendLeftoverPointsOn]int, len(source))
	for k, v := range source {
		m[SpendLeftoverPointsOnToCSSpendLeftoverPointsOn(craig_starsv1.SpendLeftoverPointsOn(k))] = int(v)
	}
	return m
}

func CometSizeMapToCometStatsMap(c Converter, source map[cs.CometSize]cs.CometStats) map[int32]*craig_starsv1.CometStats {
	m := make(map[int32]*craig_starsv1.CometStats, len(source))
	for k, v := range source {
		m[int32(CSCometSizeToCometSize(k))] = c.ConvertCSCometStats(v)
	}
	return m
}

func CometStatsMapToCometSizeMap(c Converter, source map[int32]*craig_starsv1.CometStats) map[cs.CometSize]cs.CometStats {
	m := make(map[cs.CometSize]cs.CometStats, len(source))
	for k, v := range source {
		m[CometSizeToCSCometSize(craig_starsv1.CometSize(k))] = c.ConvertCometStats(v)
	}
	return m
}

// prt_specs
func StringMapToPRTSpecMap(c Converter, source map[cs.PRT]cs.PRTSpec) map[int32]*craig_starsv1.PRTSpec {
	m := make(map[int32]*craig_starsv1.PRTSpec, len(source))
	for k, v := range source {
		m[int32(CSPRTToPRT(k))] = c.ConvertCSPRTSpec(v)
	}
	return m
}

func PRTSpecMapToStringMap(c Converter, source map[int32]*craig_starsv1.PRTSpec) map[cs.PRT]cs.PRTSpec {
	m := make(map[cs.PRT]cs.PRTSpec, len(source))
	for k, v := range source {
		m[PRTToCSPRT(craig_starsv1.Prt(k))] = c.ConvertPRTSpec(v)
	}
	return m
}

// lrt_specs
func IntMapToLRTSpecMap(c Converter, source map[cs.LRT]cs.LRTSpec) map[int32]*craig_starsv1.LRTSpec {
	m := make(map[int32]*craig_starsv1.LRTSpec, len(source))
	for k, v := range source {
		m[int32(k)] = c.ConvertCSLRTSpec(v)
	}
	return m
}

func LRTSpecMapToIntMap(c Converter, source map[int32]*craig_starsv1.LRTSpec) map[cs.LRT]cs.LRTSpec {
	m := make(map[cs.LRT]cs.LRTSpec, len(source))
	for k, v := range source {
		m[cs.LRT(k)] = c.ConvertLRTSpec(v)
	}
	return m
}

// minefield_stats_by_type
func MinefieldTypeMapToMinefieldStatsMap(c Converter, source map[cs.MinefieldType]cs.MinefieldStats) map[int32]*craig_starsv1.MinefieldStats {
	m := make(map[int32]*craig_starsv1.MinefieldStats, len(source))
	for k, v := range source {
		m[int32(CSMinefieldTypeToMinefieldType(k))] = c.ConvertCSMinefieldStats(v)
	}
	return m
}

func MinefieldStatsMapToMinefieldTypeMap(c Converter, source map[int32]*craig_starsv1.MinefieldStats) map[cs.MinefieldType]cs.MinefieldStats {
	m := make(map[cs.MinefieldType]cs.MinefieldStats, len(source))
	for k, v := range source {
		m[Int32ToMinefieldType(k)] = c.ConvertMinefieldStats(v)
	}
	return m
}

// packet_decay_rate
func PacketDecayRateMapTocketDecayRateMap(source map[int]float64) map[int32]float64 {
	m := make(map[int32]float64, len(source))
	for k, v := range source {
		m[int32(k)] = v
	}
	return m
}

func ProtoToPacketDecayRateMap(source map[int32]float64) map[int]float64 {
	m := make(map[int]float64, len(source))
	for k, v := range source {
		m[int(k)] = v
	}
	return m
}

// random_event_chances
func RandomEventChancesMapTondomEventChancesMap(source map[cs.RandomEvent]float64) map[int32]float64 {
	m := make(map[int32]float64, len(source))
	for k, v := range source {
		m[int32(CSRandomEventToRandomEvent(k))] = v
	}
	return m
}

func ProtoToRandomEventChancesMap(source map[int32]float64) map[cs.RandomEvent]float64 {
	m := make(map[cs.RandomEvent]float64, len(source))
	for k, v := range source {
		m[RandomEventToCSRandomEvent(craig_starsv1.RandomEvent(k))] = v
	}
	return m
}

// repair_rates
func RepairRatesMapTopairRatesMap(source map[cs.RepairRate]float64) map[int32]float64 {
	m := make(map[int32]float64, len(source))
	for k, v := range source {
		m[int32(CSRepairRateToRepairRate(k))] = v
	}
	return m
}

func ProtoToRepairRatesMap(source map[int32]float64) map[cs.RepairRate]float64 {
	m := make(map[cs.RepairRate]float64, len(source))
	for k, v := range source {
		m[RepairRateToCSRepairRate(craig_starsv1.RepairRate(k))] = v
	}
	return m
}

// wormhole_pairs_for_size
func SizeMapToIntMap(c Converter, source map[cs.Size]int) map[int32]int32 {
	m := make(map[int32]int32, len(source))
	for k, v := range source {
		m[int32(CSSizeToSize(k))] = int32(v)
	}
	return m
}

func IntMapToSizeMap(c Converter, source map[int32]int32) map[cs.Size]int {
	m := make(map[cs.Size]int, len(source))
	for k, v := range source {
		m[SizeToCSSize(craig_starsv1.Size(k))] = int(v)
	}
	return m
}

// wormhole_stats_by_stability
func WormholeStabilityMapToIntMap(c Converter, source map[cs.WormholeStability]cs.WormholeStats) map[int32]*craig_starsv1.WormholeStats {
	m := make(map[int32]*craig_starsv1.WormholeStats, len(source))
	for k, v := range source {
		m[int32(CSWormholeStabilityToWormholeStability(k))] = c.ConvertCSWormholeStats(v)
	}
	return m
}

func IntMapToWormholeStabilityMap(c Converter, source map[int32]*craig_starsv1.WormholeStats) map[cs.WormholeStability]cs.WormholeStats {
	m := make(map[cs.WormholeStability]cs.WormholeStats, len(source))
	for k, v := range source {
		m[WormholeStabilityToCSWormholeStability(craig_starsv1.WormholeStability(k))] = c.ConvertWormholeStats(v)
	}
	return m
}

// tech_cost_offset
func TechCostOffsetToIntMap(c Converter, source cs.TechCostOffset) map[int32]float64 {
	m := make(map[int32]float64, len(source))
	for k, v := range source {
		m[int32(CSTechTagToTechTag(k))] = v
	}
	return m
}

func IntMapToTechCostOffset(c Converter, source map[int32]float64) cs.TechCostOffset {
	m := make(cs.TechCostOffset, len(source))
	for k, v := range source {
		m[TechTagToCSTechTag(craig_starsv1.TechTag(k))] = v
	}
	return m
}

// terraform_hab_type
func TerraformHabTypeMapToIntMap(c Converter, source map[cs.TerraformHabType]*cs.TechTerraform) map[int32]*craig_starsv1.TechTerraform {
	m := make(map[int32]*craig_starsv1.TechTerraform, len(source))
	for k, v := range source {
		m[int32(CSTerraformHabTypeToTerraformHabType(k))] = c.ConvertCSTechTerraform(v)
	}
	return m
}

func IntMapToTerraformHabTypeMap(c Converter, source map[int32]*craig_starsv1.TechTerraform) map[cs.TerraformHabType]*cs.TechTerraform {
	m := make(map[cs.TerraformHabType]*cs.TechTerraform, len(source))
	for k, v := range source {
		m[TerraformHabTypeToCSTerraformHabType(craig_starsv1.TerraformHabType(k))] = c.ConvertTechTerraform(v)
	}
	return m
}

// purposes
func ShipDesignPurposeMapToBoolMap(c Converter, source map[cs.ShipDesignPurpose]bool) map[int32]bool {
	m := make(map[int32]bool, len(source))
	for k, v := range source {
		m[int32(CSShipDesignPurposeToShipDesignPurpose(k))] = v
	}
	return m
}

func BoolMapToShipDesignPurposeMap(c Converter, source map[int32]bool) map[cs.ShipDesignPurpose]bool {
	m := make(map[cs.ShipDesignPurpose]bool, len(source))
	for k, v := range source {
		m[ShipDesignPurposeToCSShipDesignPurpose(craig_starsv1.ShipDesignPurpose(k))] = v
	}
	return m
}

// costs
func QueueItemTypeMapToIntMap(c Converter, source map[cs.QueueItemType]cs.Cost) map[int32]*craig_starsv1.Cost {
	m := make(map[int32]*craig_starsv1.Cost, len(source))
	for k, v := range source {
		m[int32(CSQueueItemTypeToQueueItemType(k))] = c.ConvertCSCost(v)
	}
	return m
}

func IntMapToQueueItemTypeMap(c Converter, source map[int32]*craig_starsv1.Cost) map[cs.QueueItemType]cs.Cost {
	m := make(map[cs.QueueItemType]cs.Cost, len(source))
	for k, v := range source {
		m[QueueItemTypeToCSQueueItemType(craig_starsv1.QueueItemType(k))] = c.ConvertCost(v)
	}
	return m
}

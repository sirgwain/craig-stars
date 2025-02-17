import {
	MysteryTraderRewardArmor,
	MysteryTraderRewardBeamWeapon,
	MysteryTraderRewardBomb,
	MysteryTraderRewardElectrical,
	MysteryTraderRewardEngine,
	MysteryTraderRewardMechanical,
	MysteryTraderRewardMineRobot,
	MysteryTraderRewardShield,
	MysteryTraderRewardTorpedo,
	type MysteryTraderRewardType
} from './cs';

export function isHullComponent(type: MysteryTraderRewardType): boolean {
	switch (type) {
		case MysteryTraderRewardEngine:
		case MysteryTraderRewardBomb:
		case MysteryTraderRewardArmor:
		case MysteryTraderRewardShield:
		case MysteryTraderRewardElectrical:
		case MysteryTraderRewardMechanical:
		case MysteryTraderRewardTorpedo:
		case MysteryTraderRewardMineRobot:
		case MysteryTraderRewardBeamWeapon:
			return true;
	}
	return false;
}

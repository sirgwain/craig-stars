import { MysteryTraderRewardType } from './cs-proto';

export function isHullComponent(type: MysteryTraderRewardType): boolean {
	switch (type) {
		case MysteryTraderRewardType.ENGINE:
		case MysteryTraderRewardType.BOMB:
		case MysteryTraderRewardType.ARMOR:
		case MysteryTraderRewardType.SHIELD:
		case MysteryTraderRewardType.ELECTRICAL:
		case MysteryTraderRewardType.MECHANICAL:
		case MysteryTraderRewardType.TORPEDO:
		case MysteryTraderRewardType.MINE_ROBOT:
		case MysteryTraderRewardType.BEAM_WEAPON:
		case MysteryTraderRewardType.JUMP_GATE:
			return true;
	}
	return false;
}

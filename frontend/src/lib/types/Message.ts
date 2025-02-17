import {
	TargetFleet,
	TargetMineField,
	TargetMineralPacket,
	TargetMysteryTrader,
	TargetPlanet,
	TargetWormhole,
	type PlayerMessage,
	type PlayerMessageTarget
} from './cs';
import { MapObjectType } from './MapObject';
import type { PlayerSettings } from './PlayerSettings';

export function getMapObjectTypeForMessageType(
	targetType: PlayerMessageTarget | MapObjectType | string | undefined
): MapObjectType {
	switch (targetType) {
		case TargetPlanet:
			return MapObjectType.Planet;
		case TargetFleet:
			return MapObjectType.Fleet;
		case TargetWormhole:
			return MapObjectType.Wormhole;
		case TargetMineField:
			return MapObjectType.MineField;
		case TargetMysteryTrader:
			return MapObjectType.MysteryTrader;
		case TargetMineralPacket:
			return MapObjectType.MineralPacket;
	}

	return MapObjectType.None;
}

// get the next visible message taking into account filters
export function getNextVisibleMessageNum(
	num: number,
	showFilteredMessages: boolean,
	messages: PlayerMessage[],
	settings: PlayerSettings
): number {
	for (let i = num + 1; i < messages.length; i++) {
		if (showFilteredMessages || settings.isMessageVisible(messages[i].type)) {
			return i;
		}
	}
	return num;
}

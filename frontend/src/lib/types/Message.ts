import {
	MapObjectTypeFleet,
	MapObjectTypeMinefield,
	MapObjectTypeMineralPacket,
	MapObjectTypeMysteryTrader,
	MapObjectTypeNone,
	MapObjectTypePlanet,
	MapObjectTypeWormhole,
	TargetFleet,
	TargetMinefield,
	TargetMineralPacket,
	TargetMysteryTrader,
	TargetPlanet,
	TargetWormhole,
	type MapObjectType,
	type PlayerMessage,
	type PlayerMessageTarget
} from './cs';
import type { PlayerSettings } from './PlayerSettings';

export function getMapObjectTypeForMessageType(
	targetType: PlayerMessageTarget | MapObjectType | string | undefined
): MapObjectType {
	switch (targetType) {
		case TargetPlanet:
			return MapObjectTypePlanet;
		case TargetFleet:
			return MapObjectTypeFleet;
		case TargetWormhole:
			return MapObjectTypeWormhole;
		case TargetMinefield:
			return MapObjectTypeMinefield;
		case TargetMysteryTrader:
			return MapObjectTypeMysteryTrader;
		case TargetMineralPacket:
			return MapObjectTypeMineralPacket;
	}

	return MapObjectTypeNone;
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

import { MapObjectType, PlayerMessageTargetType, type PlayerMessage } from '$lib/types/cs-proto';
import type { MapObjectTargetLike } from './MapObject';
import type { PlayerSettings } from './PlayerSettings';

export function getMapObjectTarget(message: PlayerMessage): MapObjectTargetLike {
	return {
		...message.target!,
		targetType: getMapObjectTypeForMessageType(message.target?.targetType)
	};
}

export function getMapObjectTypeForMessageType(
	targetType: PlayerMessageTargetType | undefined
): MapObjectType {
	switch (targetType) {
		case PlayerMessageTargetType.PLANET:
			return MapObjectType.PLANET;
		case PlayerMessageTargetType.FLEET:
			return MapObjectType.FLEET;
		case PlayerMessageTargetType.WORMHOLE:
			return MapObjectType.WORMHOLE;
		case PlayerMessageTargetType.MINEFIELD:
			return MapObjectType.MINEFIELD;
		case PlayerMessageTargetType.MYSTERY_TRADER:
			return MapObjectType.MYSTERY_TRADER;
		case PlayerMessageTargetType.MINERAL_PACKET:
			return MapObjectType.MINERAL_PACKET;
	}

	return MapObjectType.UNSPECIFIED;
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

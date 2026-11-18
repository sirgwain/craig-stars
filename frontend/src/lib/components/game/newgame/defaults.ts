import { create } from '@bufbuild/protobuf';
import {
	AiDifficulty,
	Density,
	GameSettingsSchema,
	GameStartMode,
	NewGamePlayerSchema,
	NewGamePlayerType,
	PlayerPositions,
	Size
} from '#lib/types/cs-proto.js';
import {
	VictoryConditionAttainTechLevels,
	VictoryConditionExceedsSecondPlaceScore,
	VictoryConditionOwnPlanets
} from '#lib/types/Consts.js';
import { getColor } from './playerColors';

export const allPlayerTypes = [
	NewGamePlayerType.HOST,
	NewGamePlayerType.OPEN,
	NewGamePlayerType.GUEST,
	NewGamePlayerType.AI
];

export function newGamePlayer(type: NewGamePlayerType, index: number) {
	return create(NewGamePlayerSchema, {
		type,
		color: getColor(index),
		aiDifficulty: type === NewGamePlayerType.HOST ? AiDifficulty.UNSPECIFIED : AiDifficulty.NORMAL
	});
}

export function singlePlayerSettings(name = 'A Barefoot Jaywalk') {
	return create(GameSettingsSchema, {
		name,
		public: false,
		size: Size.SMALL,
		density: Density.NORMAL,
		playerPositions: PlayerPositions.MODERATE,
		randomEvents: true,
		startMode: GameStartMode.UNSPECIFIED,
		players: [
			newGamePlayer(NewGamePlayerType.HOST, 0),
			newGamePlayer(NewGamePlayerType.AI, 1),
			newGamePlayer(NewGamePlayerType.AI, 2)
		],
		victoryConditions: {
			conditions:
				VictoryConditionOwnPlanets |
				VictoryConditionAttainTechLevels |
				VictoryConditionExceedsSecondPlaceScore,
			numCriteriaRequired: 1,
			yearsPassed: 50,
			ownPlanets: 60,
			attainTechLevel: 22,
			attainTechLevelNumFields: 4,
			exceedsScore: 11000,
			exceedsSecondPlaceScore: 100,
			productionCapacity: 100,
			ownCapitalShips: 100,
			highestScoreAfterYears: 100
		}
	});
}

export function hostedGameSettings(name?: string) {
	const settings = singlePlayerSettings(name);
	settings.players[1] = newGamePlayer(NewGamePlayerType.OPEN, 1);
	return settings;
}

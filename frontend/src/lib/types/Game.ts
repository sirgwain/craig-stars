import {
	AIDifficultyCheater,
	AIDifficultyNormal,
	DensityDense,
	DensityNormal,
	DensityPacked,
	DensitySparse,
	GameStartModeAccBBS,
	GameStartModeMax,
	GameStartModeNormal,
	NewGamePlayerTypeAI,
	NewGamePlayerTypeGuest,
	NewGamePlayerTypeHost,
	NewGamePlayerTypeOpen,
	PlayerPositionsClose,
	PlayerPositionsDistant,
	PlayerPositionsFarther,
	PlayerPositionsModerate,
	SizeHuge,
	SizeHugeWide,
	SizeLarge,
	SizeLargeWide,
	SizeMedium,
	SizeMediumWide,
	SizeSmall,
	SizeSmallWide,
	SizeTiny,
	SizeTinyWide,
	type AIDifficulty,
	type Density,
	type GameStartMode,
	type PlayerPositions,
	type Size
} from './cs';

export const Sizes: Size[] = [
	SizeTiny,
	SizeTinyWide,
	SizeSmall,
	SizeSmallWide,
	SizeMedium,
	SizeMediumWide,
	SizeLarge,
	SizeLargeWide,
	SizeHuge,
	SizeHugeWide
];

export const Densities: Density[] = [DensitySparse, DensityNormal, DensityDense, DensityPacked];
export const GameStartModes: GameStartMode[] = [
	GameStartModeNormal,
	GameStartModeAccBBS,
	GameStartModeMax
];
export const PlayerPositionses: PlayerPositions[] = [
	PlayerPositionsClose,
	PlayerPositionsModerate,
	PlayerPositionsFarther,
	PlayerPositionsDistant
];

export const NewGamePlayerTypes = [
	NewGamePlayerTypeHost,
	NewGamePlayerTypeGuest,
	NewGamePlayerTypeOpen,
	NewGamePlayerTypeAI
];

export const GameStartModeFullNames: { [key in GameStartMode]: string } = {
	[GameStartModeNormal]: 'Normal',
	[GameStartModeAccBBS]: 'Accelerated BBS Play',
	[GameStartModeMax]: 'Max Start'
};

export const AIDifficulties: AIDifficulty[] = [
	// AIDifficultyNone,
	// AIDifficultyEasy,
	AIDifficultyNormal,
	// AIDifficultyHard,
	AIDifficultyCheater
];

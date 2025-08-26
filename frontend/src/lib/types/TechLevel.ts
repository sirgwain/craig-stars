import { TechField, TechLevelSchema, type TechLevel } from '$lib/types/cs-proto';
import { create } from '@bufbuild/protobuf';

export const emptyTechLevel = (): TechLevel =>
	create(TechLevelSchema, {
		energy: 0,
		weapons: 0,
		propulsion: 0,
		construction: 0,
		electronics: 0,
		biotechnology: 0
	});

export function hasRequiredLevels(tl: TechLevel, required: TechLevel | undefined): boolean {
	return (
		(tl.energy ?? 0) >= (required?.energy ?? 0) &&
		(tl.weapons ?? 0) >= (required?.weapons ?? 0) &&
		(tl.propulsion ?? 0) >= (required?.propulsion ?? 0) &&
		(tl.construction ?? 0) >= (required?.construction ?? 0) &&
		(tl.electronics ?? 0) >= (required?.electronics ?? 0) &&
		(tl.biotechnology ?? 0) >= (required?.biotechnology ?? 0)
	);
}

export function subtract(tl1: TechLevel | undefined, tl2: TechLevel | undefined): TechLevel {
	return create(TechLevelSchema, {
		energy: (tl1?.energy ?? 0) - (tl2?.energy ?? 0),
		weapons: (tl1?.weapons ?? 0) - (tl2?.weapons ?? 0),
		propulsion: (tl1?.propulsion ?? 0) - (tl2?.propulsion ?? 0),
		construction: (tl1?.construction ?? 0) - (tl2?.construction ?? 0),
		electronics: (tl1?.electronics ?? 0) - (tl2?.electronics ?? 0),
		biotechnology: (tl1?.biotechnology ?? 0) - (tl2?.biotechnology ?? 0)
	});
}

export function sum(tl: TechLevel | undefined): number {
	if (!tl) return 0;
	return (
		(tl.energy ?? 0) +
		(tl.weapons ?? 0) +
		(tl.propulsion ?? 0) +
		(tl.construction ?? 0) +
		(tl.electronics ?? 0) +
		(tl.biotechnology ?? 0)
	);
}

export function levelsAbove(
	req: TechLevel | undefined,
	level: TechLevel | undefined
): number | undefined {
	const diffs: number[] = [];
	if (req?.energy) {
		diffs.push((level?.energy ?? 0) - req?.energy);
	}
	if (req?.weapons) {
		diffs.push((level?.weapons ?? 0) - req?.weapons);
	}
	if (req?.propulsion) {
		diffs.push((level?.propulsion ?? 0) - req?.propulsion);
	}
	if (req?.construction) {
		diffs.push((level?.construction ?? 0) - req?.construction);
	}
	if (req?.electronics) {
		diffs.push((level?.electronics ?? 0) - req?.electronics);
	}
	if (req?.biotechnology) {
		diffs.push((level?.biotechnology ?? 0) - req?.biotechnology);
	}
	return Math.min(...diffs);
}

export function get(tl: TechLevel | undefined, field: TechField): number {
	if (!tl) {
		return 0;
	}
	switch (field) {
		case TechField.ENERGY:
			return tl.energy ?? 0;
		case TechField.WEAPONS:
			return tl.weapons ?? 0;
		case TechField.PROPULSION:
			return tl.propulsion ?? 0;
		case TechField.CONSTRUCTION:
			return tl.construction ?? 0;
		case TechField.ELECTRONICS:
			return tl.electronics ?? 0;
		case TechField.BIOTECHNOLOGY:
			return tl.biotechnology ?? 0;
		default:
			throw new Error('invalid field: ' + field);
	}
}

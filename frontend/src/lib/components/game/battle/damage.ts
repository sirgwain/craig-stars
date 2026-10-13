// How a hit on a token is marked on the battle board: shield damage and armor damage get
// different icons and colors.

// 24x24 icon paths
export const SHIELD_ICON_PATH = 'M12 2 20 5v6c0 5.5-3.4 9.4-8 11-4.6-1.6-8-5.5-8-11V5Z';
export const ARMOR_ICON_PATH =
	'M12 1l2.2 7.2L21 5l-4.4 6 6.4 2.5-7.2 1.1L17.5 22 12 16.8 6.5 22l1.7-7.4L1 13.5l6.4-2.5L3 5l6.8 3.2Z';

export const SHIELD_COLOR = '#8fe3ff';
export const ARMOR_COLOR = '#ffb050';
export const DESTROYED_COLOR = '#ff6060';

export type DamageKind = 'shield' | 'armor' | 'destroyed';

export type DamageMark = {
	kind: DamageKind;
	amount: number;
	text: string;
	color: string;
	label: string;
};

type Hit = {
	damageDoneShields?: number;
	damageDoneArmor?: number;
	tokensDestroyed?: number;
};

export function formatDamage(kind: DamageKind, amount: number): string {
	if (kind === 'destroyed') {
		return amount > 1 ? `${amount} destroyed` : 'destroyed';
	}
	return `-${amount}`;
}

// the marks for one hit, armor first, then shields, then ships destroyed
export function damageMarks(hit: Hit | undefined): DamageMark[] {
	const marks: DamageMark[] = [];
	const armor = hit?.damageDoneArmor ?? 0;
	const shields = hit?.damageDoneShields ?? 0;
	const destroyed = hit?.tokensDestroyed ?? 0;
	if (armor > 0) {
		marks.push({
			kind: 'armor',
			amount: armor,
			text: formatDamage('armor', armor),
			color: ARMOR_COLOR,
			label: `${armor} armor damage`
		});
	}
	if (shields > 0) {
		marks.push({
			kind: 'shield',
			amount: shields,
			text: formatDamage('shield', shields),
			color: SHIELD_COLOR,
			label: `${shields} shield damage`
		});
	}
	if (destroyed > 0) {
		marks.push({
			kind: 'destroyed',
			amount: destroyed,
			text: formatDamage('destroyed', destroyed),
			color: DESTROYED_COLOR,
			label: `${destroyed} destroyed`
		});
	}
	return marks;
}

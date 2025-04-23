import type { DesignFinder } from '$lib/services/Universe';
import { startCase } from 'lodash-es';
import {
	QueueItemTypeAutoDefenses,
	QueueItemTypeAutoFactories,
	QueueItemTypeAutoMaxTerraform,
	QueueItemTypeAutoMineralAlchemy,
	QueueItemTypeAutoMineralPacket,
	QueueItemTypeAutoMines,
	QueueItemTypeAutoMinTerraform,
	QueueItemTypeBoraniumMineralPacket,
	QueueItemTypeDefenses,
	QueueItemTypeFactory,
	QueueItemTypeGenesisDevice,
	QueueItemTypeGermaniumMineralPacket,
	QueueItemTypeIroniumMineralPacket,
	QueueItemTypeMine,
	QueueItemTypeMineralAlchemy,
	QueueItemTypeMixedMineralPacket,
	QueueItemTypePlanetaryScanner,
	QueueItemTypeShipToken,
	QueueItemTypeStarbase,
	QueueItemTypeTerraformEnvironment,
	type ProductionQueueItem,
	type QueueItemType
} from './cs';

export const validQueueItemTypes = new Set([
	QueueItemTypeIroniumMineralPacket,
	QueueItemTypeBoraniumMineralPacket,
	QueueItemTypeGermaniumMineralPacket,
	QueueItemTypeMixedMineralPacket,
	QueueItemTypeFactory,
	QueueItemTypeMine,
	QueueItemTypeDefenses,
	QueueItemTypeMineralAlchemy,
	QueueItemTypeTerraformEnvironment,
	QueueItemTypeAutoMines,
	QueueItemTypeAutoFactories,
	QueueItemTypeAutoDefenses,
	QueueItemTypeAutoMineralAlchemy,
	QueueItemTypeAutoMinTerraform,
	QueueItemTypeAutoMaxTerraform,
	QueueItemTypeAutoMineralPacket,
	QueueItemTypeShipToken,
	QueueItemTypeStarbase,
	QueueItemTypePlanetaryScanner,
	QueueItemTypeGenesisDevice
]);

export const stringToQueueItemType = (value: string): QueueItemType | undefined => {
	return validQueueItemTypes.has(value) ? value : undefined;
};

export const fromQueueItemType = (type: QueueItemType): ProductionQueueItem => ({
	type,
	quantity: 1,
	allocated: {},
	tags: {}
});

/**
 * Determine if a {@linkcode QueueItemType} is an auto item
 * @param type The {@linkcode QueueItemType} to check
 * @returns `true` if type denotes an auto item
 */
export function isAuto(type: QueueItemType): boolean {
	switch (type) {
		case QueueItemTypeAutoMines:
		case QueueItemTypeAutoFactories:
		case QueueItemTypeAutoDefenses:
		case QueueItemTypeAutoMineralAlchemy:
		case QueueItemTypeAutoMinTerraform:
		case QueueItemTypeAutoMaxTerraform:
		case QueueItemTypeAutoMineralPacket:
			return true;
		default:
			return false;
	}
}

/**
 * Check if a {@linkcode QueueItemType} is a concrete or auto planetary item.
 * @param type the {@linkcode QueueItemType} to check
 * @returns `true` if item is a concrete or auto planetary item (i.e. not a ship/starbase)
 */
export function isPlanetary(type: QueueItemType): boolean {
	return type !== QueueItemTypeShipToken && type !== QueueItemTypeStarbase;
}

/**
 * Get the concrete type corresponding to a given {@linkcode QueueItemType}.
 * @param type The {@linkcode QueueItemType} to check
 * @returns The concrete version of the {@linkcode QueueItemType} -
 * Factories for AutoFactories, Mines for AutoMines, etc.
 */
export const concreteType = (type: QueueItemType): QueueItemType => {
	switch (type) {
		case QueueItemTypeAutoMines:
			return QueueItemTypeMine;
		case QueueItemTypeAutoFactories:
			return QueueItemTypeFactory;
		case QueueItemTypeAutoDefenses:
			return QueueItemTypeDefenses;
		case QueueItemTypeAutoMineralAlchemy:
			return QueueItemTypeMineralAlchemy;
		case QueueItemTypeAutoMinTerraform:
		case QueueItemTypeAutoMaxTerraform:
			return QueueItemTypeTerraformEnvironment;
		case QueueItemTypeAutoMineralPacket:
			return QueueItemTypeMixedMineralPacket;
		default:
			return type;
	}
};

export function getFullName(item: ProductionQueueItem, designFinder: DesignFinder): string {
	switch (item.type) {
		case QueueItemTypeStarbase:
		case QueueItemTypeShipToken:
			return designFinder.getMyDesign(item.designNum)?.name ?? '';
		case QueueItemTypeAutoMineralAlchemy:
			return 'Alchemy (Auto Build)';
		case QueueItemTypeMineralAlchemy:
			return 'Mineral Alchemy';
		case QueueItemTypeAutoMines:
			return 'Mine (Auto Build)';
		case QueueItemTypeAutoFactories:
			return 'Factory (Auto Build)';
		case QueueItemTypeAutoDefenses:
			return 'Defense (Auto Build)';
		case QueueItemTypeAutoMinTerraform:
			return 'Minimum Terraform';
		case QueueItemTypeAutoMaxTerraform:
			return 'Maximum Terraform';
		case QueueItemTypeTerraformEnvironment:
			return 'Terraform Environment';
		case QueueItemTypeIroniumMineralPacket:
			return 'Mineral Packet (Ironium)';
		case QueueItemTypeBoraniumMineralPacket:
			return 'Mineral Packet (Boranium)';
		case QueueItemTypeGermaniumMineralPacket:
			return 'Mineral Packet (Germanium)';
		case QueueItemTypeMixedMineralPacket:
			return 'Mixed Mineral Packet';
		case QueueItemTypeAutoMineralPacket:
			return 'Mixed Mineral Packet (Auto)';
		case QueueItemTypePlanetaryScanner:
			return 'Planetary Scanner';
		case QueueItemTypeGenesisDevice:
			return 'Genesis Device';
		default:
			return item.type.toString();
	}
}

export function getShortName(item: ProductionQueueItem, designFinder: DesignFinder): string {
	switch (item.type) {
		case QueueItemTypeStarbase:
		case QueueItemTypeShipToken:
			return designFinder.getMyDesign(item.designNum)?.name ?? '';
		case QueueItemTypeTerraformEnvironment:
			return 'Terraform Environment';
		case QueueItemTypeAutoMines:
			return 'Mine (Auto)';
		case QueueItemTypeAutoFactories:
			return 'Factory (Auto)';
		case QueueItemTypeAutoDefenses:
			return 'Defenses (Auto)';
		case QueueItemTypeAutoMineralAlchemy:
			return 'Alchemy (Auto)';
		case QueueItemTypeAutoMaxTerraform:
			return 'Max Terraform (Auto)';
		case QueueItemTypeAutoMinTerraform:
			return 'Min Terraform (Auto)';
		default:
			return `${startCase(item.type)}`;
	}
}

/**
 * Get the singular name of a {@linkcode QueueItemType} in proper English.
 * @param type the {@linkcode QueueItemType} being checked.
 * @returns The singular form of this {@linkcode QueueItemType}, suitable for use in messages.
 */
export function getSingularName(type: QueueItemType): string {
	switch (type) {
		case QueueItemTypeAutoMineralAlchemy:
			return 'auto mineral alchemy';
		case QueueItemTypeMineralAlchemy:
			return 'mineral alchemy';
		case QueueItemTypeAutoMines:
			return 'auto mine';
		case QueueItemTypeAutoFactories:
			return 'auto factory';
		case QueueItemTypeAutoMinTerraform:
			return 'minimum terraform';
		case QueueItemTypeAutoMaxTerraform:
			return 'maximum terraform';
		case QueueItemTypeAutoDefenses:
			return 'auto defense outpost';
		case QueueItemTypeDefenses:
			return 'defense outpost';
		case QueueItemTypeIroniumMineralPacket:
			return 'ironium mineral packet';
		case QueueItemTypeBoraniumMineralPacket:
			return 'boranium mineral packet';
		case QueueItemTypeGermaniumMineralPacket:
			return 'germanium mineral packet';
		case QueueItemTypeMixedMineralPacket:
			return 'mixed mineral packet';
		case QueueItemTypeTerraformEnvironment:
			return 'terraform environment';
		case QueueItemTypeAutoMineralPacket:
			return 'auto mineral packet';
		case QueueItemTypePlanetaryScanner:
			return 'planetary scanner';
		case QueueItemTypeGenesisDevice:
			return 'genesis device';
		default:
			return `${startCase(type).toLowerCase()}`;
	}
}

/**
 * Get the plural name of a {@linkcode QueueItemType} in proper English.
 * @param type the {@linkcode QueueItemType} being checked.
 * @returns The plural form of this {@linkcode QueueItemType}, suitable for use in messages.
 */
export function getPluralName(type: QueueItemType): string {
	switch (type) {
		case QueueItemTypeAutoMineralAlchemy:
			// yes, the plural of "alchemy" is alchemies. FIGHT ME
			return 'auto mineral alchemies';
		case QueueItemTypeMineralAlchemy:
			return 'mineral alchemies';
		case QueueItemTypeAutoFactories:
			return 'auto factories';
		default:
			return getSingularName(type) + 's';
	}
}

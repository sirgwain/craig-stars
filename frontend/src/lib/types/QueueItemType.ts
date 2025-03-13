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

/**
 * Determine if a ProductionQueueItem is an auto item
 * @param type The type to check
 * @returns
 */
export const isAuto = (type: QueueItemType): boolean => {
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
};

/**
 * Get the concrete type for a queue item type,
 * @param type The QueueItemType
 * @returns Factory for AuotFactories, Mine for AutoMines, etc
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
			return 'Alchemy';
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
		case QueueItemTypeIroniumMineralPacket:
			return 'Mineral Packet (Ironium)';
		case QueueItemTypeBoraniumMineralPacket:
			return 'Mineral Packet (Boranium)';
		case QueueItemTypeGermaniumMineralPacket:
			return 'Mineral Packet (Germanium)';
		case QueueItemTypeTerraformEnvironment:
			return 'Terraform Environment';
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

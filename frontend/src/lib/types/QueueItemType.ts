import type { DesignFinder } from '$lib/services/Universe';
import { startCase } from 'lodash-es';
import type { ProductionQueueItem } from './Production';

export type QueueItemType = (typeof QueueItemTypes)[keyof typeof QueueItemTypes];

export const QueueItemTypes = {
	AutoMines: 'AutoMines',
	AutoFactories: 'AutoFactories',
	AutoDefenses: 'AutoDefenses',
	AutoMineralAlchemy: 'AutoMineralAlchemy',
	AutoMinTerraform: 'AutoMinTerraform',
	AutoMaxTerraform: 'AutoMaxTerraform',
	AutoMineralPacket: 'AutoMineralPacket',
	Factory: 'Factory',
	Mine: 'Mine',
	Defenses: 'Defenses',
	MineralAlchemy: 'MineralAlchemy',
	TerraformEnvironment: 'TerraformEnvironment',
	IroniumMineralPacket: 'IroniumMineralPacket',
	BoraniumMineralPacket: 'BoraniumMineralPacket',
	GermaniumMineralPacket: 'GermaniumMineralPacket',
	MixedMineralPacket: 'MixedMineralPacket',
	ShipToken: 'ShipToken',
	Starbase: 'Starbase',
	PlanetaryScanner: 'PlanetaryScanner',
	GenesisDevice: 'GenesisDevice'
} as const;

export const stringToQueueItemType = (value: string): QueueItemType => {
	return QueueItemTypes[value as keyof typeof QueueItemTypes];
};

/**
 * Determine if a ProductionQueueItem is an auto item
 * @param type The type to check
 * @returns
 */
export const isAuto = (type: QueueItemType): boolean => {
	switch (type) {
		case QueueItemTypes.AutoMines:
		case QueueItemTypes.AutoFactories:
		case QueueItemTypes.AutoDefenses:
		case QueueItemTypes.AutoMineralAlchemy:
		case QueueItemTypes.AutoMinTerraform:
		case QueueItemTypes.AutoMaxTerraform:
		case QueueItemTypes.AutoMineralPacket:
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
		case QueueItemTypes.AutoMines:
			return QueueItemTypes.Mine;
		case QueueItemTypes.AutoFactories:
			return QueueItemTypes.Factory;
		case QueueItemTypes.AutoDefenses:
			return QueueItemTypes.Defenses;
		case QueueItemTypes.AutoMineralAlchemy:
			return QueueItemTypes.MineralAlchemy;
		case QueueItemTypes.AutoMinTerraform:
		case QueueItemTypes.AutoMaxTerraform:
			return QueueItemTypes.TerraformEnvironment;
		case QueueItemTypes.AutoMineralPacket:
			return QueueItemTypes.MixedMineralPacket;
		default:
			return type;
	}
};

export function getFullName(item: ProductionQueueItem, designFinder: DesignFinder): string {
	switch (item.type) {
		case QueueItemTypes.Starbase:
		case QueueItemTypes.ShipToken:
			return designFinder.getMyDesign(item.designNum)?.name ?? '';
		case QueueItemTypes.AutoMineralAlchemy:
			return 'Alchemy (Auto Build)';
		case QueueItemTypes.MineralAlchemy:
			return 'Alchemy';
		case QueueItemTypes.AutoMines:
			return 'Mine (Auto Build)';
		case QueueItemTypes.AutoFactories:
			return 'Factory (Auto Build)';
		case QueueItemTypes.AutoDefenses:
			return 'Defense (Auto Build)';
		case QueueItemTypes.AutoMinTerraform:
			return 'Minimum Terraform';
		case QueueItemTypes.AutoMaxTerraform:
			return 'Maximum Terraform';
		case QueueItemTypes.IroniumMineralPacket:
			return 'Mineral Packet (Ironium)';
		case QueueItemTypes.BoraniumMineralPacket:
			return 'Mineral Packet (Boranium)';
		case QueueItemTypes.GermaniumMineralPacket:
			return 'Mineral Packet (Germanium)';
		case QueueItemTypes.TerraformEnvironment:
			return 'Terraform Environment';
		case QueueItemTypes.MixedMineralPacket:
			return 'Mixed Mineral Packet';
		case QueueItemTypes.AutoMineralPacket:
			return 'Mixed Mineral Packet (Auto)';
		case QueueItemTypes.PlanetaryScanner:
			return 'Planetary Scanner';
		case QueueItemTypes.GenesisDevice:
			return 'Genesis Device';
		default:
			return item.type.toString();
	}
}

export function getShortName(item: ProductionQueueItem, designFinder: DesignFinder): string {
	switch (item.type) {
		case QueueItemTypes.Starbase:
		case QueueItemTypes.ShipToken:
			return designFinder.getMyDesign(item.designNum)?.name ?? '';
		case QueueItemTypes.TerraformEnvironment:
			return 'Terraform Environment';
		case QueueItemTypes.AutoMines:
			return 'Mine (Auto)';
		case QueueItemTypes.AutoFactories:
			return 'Factory (Auto)';
		case QueueItemTypes.AutoDefenses:
			return 'Defenses (Auto)';
		case QueueItemTypes.AutoMineralAlchemy:
			return 'Alchemy (Auto)';
		case QueueItemTypes.AutoMaxTerraform:
			return 'Max Terraform (Auto)';
		case QueueItemTypes.AutoMinTerraform:
			return 'Min Terraform (Auto)';
		default:
			return `${startCase(item.type)}`;
	}
}

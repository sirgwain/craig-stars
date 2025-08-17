import type { DesignFinder } from '$lib/services/Universe';

import { QueueItemType, type ProductionQueueItem } from '$lib/types/cs-proto';
import { enumToString } from './Enums';

export const validQueueItemTypes = new Set([
	QueueItemType.IRONIUM_MINERAL_PACKET,
	QueueItemType.BORANIUM_MINERAL_PACKET,
	QueueItemType.GERMANIUM_MINERAL_PACKET,
	QueueItemType.MIXED_MINERAL_PACKET,
	QueueItemType.FACTORY,
	QueueItemType.MINE,
	QueueItemType.DEFENSES,
	QueueItemType.MINERAL_ALCHEMY,
	QueueItemType.TERRAFORM_ENVIRONMENT,
	QueueItemType.AUTO_MINES,
	QueueItemType.AUTO_FACTORIES,
	QueueItemType.AUTO_DEFENSES,
	QueueItemType.AUTO_MINERAL_ALCHEMY,
	QueueItemType.AUTO_MIN_TERRAFORM,
	QueueItemType.AUTO_MAX_TERRAFORM,
	QueueItemType.AUTO_MINERAL_PACKET,
	QueueItemType.SHIP_TOKEN,
	QueueItemType.STARBASE,
	QueueItemType.PLANETARY_SCANNER,
	QueueItemType.GENESIS_DEVICE
]);

/**
 * Determine if a ProductionQueueItem is an auto item
 * @param type The type to check
 * @returns
 */
export const isAuto = (type: QueueItemType): boolean => {
	switch (type) {
		case QueueItemType.AUTO_MINES:
		case QueueItemType.AUTO_FACTORIES:
		case QueueItemType.AUTO_DEFENSES:
		case QueueItemType.AUTO_MINERAL_ALCHEMY:
		case QueueItemType.AUTO_MIN_TERRAFORM:
		case QueueItemType.AUTO_MAX_TERRAFORM:
		case QueueItemType.AUTO_MINERAL_PACKET:
			return true;
		default:
			return false;
	}
};

/**
 * Get the concrete type for a queue item type,
 * @param type The QueueItemType.
 * @returns Factory for AuotFactories, Mine for AutoMines, etc
 */
export const concreteType = (type: QueueItemType): QueueItemType => {
	switch (type) {
		case QueueItemType.AUTO_MINES:
			return QueueItemType.MINE;
		case QueueItemType.AUTO_FACTORIES:
			return QueueItemType.FACTORY;
		case QueueItemType.AUTO_DEFENSES:
			return QueueItemType.DEFENSES;
		case QueueItemType.AUTO_MINERAL_ALCHEMY:
			return QueueItemType.MINERAL_ALCHEMY;
		case QueueItemType.AUTO_MIN_TERRAFORM:
		case QueueItemType.AUTO_MAX_TERRAFORM:
			return QueueItemType.TERRAFORM_ENVIRONMENT;
		case QueueItemType.AUTO_MINERAL_PACKET:
			return QueueItemType.MIXED_MINERAL_PACKET;
		default:
			return type;
	}
};

export function getFullName(item: ProductionQueueItem, designFinder: DesignFinder): string {
	switch (item.type) {
		case QueueItemType.STARBASE:
		case QueueItemType.SHIP_TOKEN:
			return designFinder.getMyDesign(item.designNum)?.name ?? '';
		case QueueItemType.AUTO_MINERAL_ALCHEMY:
			return 'Alchemy (Auto Build)';
		case QueueItemType.MINERAL_ALCHEMY:
			return 'Alchemy';
		case QueueItemType.AUTO_MINES:
			return 'Mine (Auto Build)';
		case QueueItemType.AUTO_FACTORIES:
			return 'Factory (Auto Build)';
		case QueueItemType.AUTO_DEFENSES:
			return 'Defense (Auto Build)';
		case QueueItemType.AUTO_MIN_TERRAFORM:
			return 'Minimum Terraform';
		case QueueItemType.AUTO_MAX_TERRAFORM:
			return 'Maximum Terraform';
		case QueueItemType.IRONIUM_MINERAL_PACKET:
			return 'Mineral Packet (Ironium)';
		case QueueItemType.BORANIUM_MINERAL_PACKET:
			return 'Mineral Packet (Boranium)';
		case QueueItemType.GERMANIUM_MINERAL_PACKET:
			return 'Mineral Packet (Germanium)';
		case QueueItemType.TERRAFORM_ENVIRONMENT:
			return 'Terraform Environment';
		case QueueItemType.MIXED_MINERAL_PACKET:
			return 'Mixed Mineral Packet';
		case QueueItemType.AUTO_MINERAL_PACKET:
			return 'Mixed Mineral Packet (Auto)';
		case QueueItemType.PLANETARY_SCANNER:
			return 'Planetary Scanner';
		case QueueItemType.GENESIS_DEVICE:
			return 'Genesis Device';
		default:
			return enumToString(QueueItemType, item.type);
	}
}

export function getShortName(item: ProductionQueueItem, designFinder: DesignFinder): string {
	switch (item.type) {
		case QueueItemType.STARBASE:
		case QueueItemType.SHIP_TOKEN:
			return designFinder.getMyDesign(item.designNum)?.name ?? '';
		case QueueItemType.TERRAFORM_ENVIRONMENT:
			return 'Terraform Environment';
		case QueueItemType.AUTO_MINES:
			return 'Mine (Auto)';
		case QueueItemType.AUTO_FACTORIES:
			return 'Factory (Auto)';
		case QueueItemType.AUTO_DEFENSES:
			return 'Defenses (Auto)';
		case QueueItemType.AUTO_MINERAL_ALCHEMY:
			return 'Alchemy (Auto)';
		case QueueItemType.AUTO_MAX_TERRAFORM:
			return 'Max Terraform (Auto)';
		case QueueItemType.AUTO_MIN_TERRAFORM:
			return 'Min Terraform (Auto)';
		default:
			return `${enumToString(QueueItemType, item.type)}`;
	}
}

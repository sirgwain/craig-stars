import type { Cost } from './Cost';
import type { QueueItemType } from './QueueItemType';

export interface ProductionQueueItem {
	type: QueueItemType;
	quantity: number;
	designNum?: number;
	allocated?: Cost;
	costOfOne?: Cost;
	skipped?: boolean;
	yearsToBuildOne?: number;
	yearsToBuildAll?: number;
	yearsToSkipAuto?: number;
	percentComplete?: number;
}


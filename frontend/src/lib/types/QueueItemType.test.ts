import { ProductionQueueItemSchema, QueueItemType } from '#lib/types/cs-proto.js';
import { create } from '@bufbuild/protobuf';
import { describe, expect, it } from 'vitest';
import { Infinite } from './Consts';
import { getCompletionDescription } from './QueueItemType';

describe('completion description', () => {
	const describeItem = (
		type: QueueItemType,
		yearsToBuildOne: number,
		yearsToBuildAll: number,
		yearsToSkipAuto = Infinite
	) =>
		getCompletionDescription(
			create(ProductionQueueItemSchema, {
				type,
				quantity: 10,
				queueItemCompletionEstimate: { yearsToBuildOne, yearsToBuildAll, yearsToSkipAuto }
			})
		);

	it('describes when the whole order is built', () => {
		expect(describeItem(QueueItemType.FACTORY, 1, 1)).toBe('1 year');
		expect(describeItem(QueueItemType.FACTORY, 3, 3)).toBe('3 years');
		expect(describeItem(QueueItemType.FACTORY, 1, 2)).toBe('1 to 2 years');
		expect(describeItem(QueueItemType.FACTORY, 2, Infinite)).toBe('2 to ???');
		expect(describeItem(QueueItemType.FACTORY, Infinite, Infinite)).toBe('never');
	});

	it('describes auto items by when they stop building', () => {
		expect(describeItem(QueueItemType.AUTO_FACTORIES, 1, 1, 4)).toBe('1 to 4 years');
		expect(describeItem(QueueItemType.AUTO_FACTORIES, Infinite, Infinite)).toBe('Skipped');
	});
});

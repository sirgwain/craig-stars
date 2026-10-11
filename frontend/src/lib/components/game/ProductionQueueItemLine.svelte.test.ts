import '../../../css/app.css';
import { Universe } from '#lib/services/Universe.js';
import { ProductionQueueItemSchema, QueueItemType } from '#lib/types/cs-proto.js';
import { Infinite } from '#lib/types/Consts.js';
import { create } from '@bufbuild/protobuf';
import { writable } from 'svelte/store';
import { expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import ProductionQueueItemLine from './ProductionQueueItemLine.svelte';

const state = vi.hoisted(() => ({ context: undefined as unknown }));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));

const classes = ['this-year', 'next-year', 'skipped', 'never'];

// which of the queue item colors a line gets for these estimates
function colors(type: QueueItemType, yearsToBuildOne: number, yearsToBuildAll: number) {
	state.context = { universe: writable(new Universe()) };
	const item = create(ProductionQueueItemSchema, {
		type,
		quantity: 100,
		queueItemCompletionEstimate: { yearsToBuildOne, yearsToBuildAll, yearsToSkipAuto: Infinite }
	});
	const { container, unmount } = render(ProductionQueueItemLine, { index: 0, item });
	const button = container.querySelector('button')!;
	const found = classes.filter((c) => button.classList.contains(`text-queue-item-${c}`));
	unmount();
	return found;
}

// Green means the whole order is built this year, blue means only some of it is.
it.each([QueueItemType.MINE, QueueItemType.AUTO_MINES])(
	'colors item type %s by how much of it is built this year',
	(type) => {
		expect(colors(type, 1, 1)).toEqual(['this-year']);
		expect(colors(type, 1, 3)).toEqual(['this-year', 'next-year']);
		expect(colors(type, 1, Infinite)).toEqual(['this-year', 'next-year']);
		expect(colors(type, 4, 6)).toEqual([]);
	}
);

it('colors items that can never be built', () => {
	expect(colors(QueueItemType.MINE, Infinite, Infinite)).toEqual(['never']);
	// an auto item with nothing to build is skipped
	expect(colors(QueueItemType.AUTO_MINES, Infinite, Infinite)).toEqual(['skipped']);
});

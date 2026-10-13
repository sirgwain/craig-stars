import '../../../css/app.css';
import { ShipDesignSchema } from '#lib/types/cs-proto.js';
import { create } from '@bufbuild/protobuf';
import { tick } from 'svelte';
import { expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import DesignCard from './DesignCard.svelte';

// the card doesn't need a game, but the tooltips it can open import the context
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => undefined }));

// open the delete dialog for a design with these ships built and queued
async function confirmDelete(numInstances: number, numQueued: number) {
	const onDelete = vi.fn(async () => {});
	const design = create(ShipDesignSchema, { num: 3, name: 'Stalwart', spec: { numInstances } });
	const { container, unmount } = render(DesignCard, {
		design,
		href: '/design',
		copyhref: '/copy',
		numQueued,
		onDelete
	});
	const find = (type: string) => document.querySelector<HTMLElement>(`[data-type="${type}"]`);
	container.querySelector<HTMLButtonElement>('.card-actions button.btn')!.click();
	await tick();
	return { design, onDelete, find, unmount };
}

it('confirms deleting an unused design without a warning', async () => {
	const { design, onDelete, find, unmount } = await confirmDelete(0, 0);
	expect(find('confirm-dialog')).not.toBeNull();
	expect(find('delete-design-instances')).toBeNull();
	expect(find('delete-design-queued')).toBeNull();
	expect(find('confirm-dialog-confirm')!.classList.contains('btn-error')).toBe(false);
	expect(onDelete).not.toHaveBeenCalled();

	find('confirm-dialog-confirm')!.click();
	await tick();
	expect(onDelete).toHaveBeenCalledWith(design);
	expect(find('confirm-dialog')).toBeNull();
	unmount();
});

it('warns about the ships that are destroyed and unqueued', async () => {
	const built = await confirmDelete(2, 0);
	expect(built.find('delete-design-instances')!.textContent).toContain('2 ships');
	expect(built.find('delete-design-queued')).toBeNull();
	expect(built.find('confirm-dialog-confirm')!.classList.contains('btn-error')).toBe(true);
	built.unmount();

	const queued = await confirmDelete(0, 1);
	expect(queued.find('delete-design-instances')).toBeNull();
	expect(queued.find('delete-design-queued')!.textContent).toContain('1 ship');
	expect(queued.find('confirm-dialog-confirm')!.classList.contains('btn-error')).toBe(true);
	queued.unmount();
});

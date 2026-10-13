import '../../css/app.css';
import { createRawSnippet, tick } from 'svelte';
import { expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import ConfirmDialog from './ConfirmDialog.svelte';

const find = (type: string) => document.querySelector<HTMLElement>(`[data-type="${type}"]`);
const children = createRawSnippet(() => ({ render: () => '<p data-type="body">Really?</p>' }));

it('shows its body and closes when confirmed', async () => {
	const onConfirm = vi.fn();
	const onCancel = vi.fn();
	const { unmount } = render(ConfirmDialog, {
		show: true,
		title: 'Scrap it?',
		confirmText: 'Scrap',
		children,
		onConfirm,
		onCancel
	});
	expect(find('confirm-dialog')!.textContent).toContain('Scrap it?');
	expect(find('body')!.textContent).toBe('Really?');
	expect(find('confirm-dialog-confirm')!.textContent).toContain('Scrap');
	expect(find('confirm-dialog-confirm')!.classList.contains('btn-error')).toBe(false);

	find('confirm-dialog-confirm')!.click();
	await tick();
	expect(onConfirm).toHaveBeenCalledOnce();
	expect(onCancel).not.toHaveBeenCalled();
	expect(find('confirm-dialog')).toBeNull();
	unmount();
});

it('cancels with the cancel button or escape', async () => {
	for (const cancel of [
		() => find('confirm-dialog-cancel')!.click(),
		() => window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
	]) {
		const onConfirm = vi.fn();
		const onCancel = vi.fn();
		const { unmount } = render(ConfirmDialog, {
			show: true,
			title: 'Scrap it?',
			destructive: true,
			children,
			onConfirm,
			onCancel
		});
		expect(find('confirm-dialog-confirm')!.classList.contains('btn-error')).toBe(true);

		cancel();
		await tick();
		expect(onCancel).toHaveBeenCalledOnce();
		expect(onConfirm).not.toHaveBeenCalled();
		expect(find('confirm-dialog')).toBeNull();
		unmount();
	}
});

it('renders nothing until it is shown', () => {
	const { unmount } = render(ConfirmDialog, { show: false, title: 'Scrap it?', children });
	expect(find('confirm-dialog')).toBeNull();
	unmount();
});

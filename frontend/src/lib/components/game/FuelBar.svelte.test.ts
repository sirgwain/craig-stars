import '../../../css/app.css';
import { page, userEvent } from 'vitest/browser';
import { describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import FuelBar from './FuelBar.svelte';

describe('FuelBar', () => {
	it('is a read only gauge without a press handler', async () => {
		render(FuelBar, { value: 30, capacity: 50 });
		await expect.element(page.getByRole('slider', { name: 'Fuel amount' })).toBeInTheDocument();
		expect(page.getByRole('button', { name: 'Transfer fuel' }).elements()).toHaveLength(0);
	});

	it('opens the transfer dialog when pressed', async () => {
		const onPress = vi.fn<(event: PointerEvent) => void>();
		render(FuelBar, { value: 30, capacity: 50, onPress });
		const bar = page.getByRole('button', { name: 'Transfer fuel' }).element() as HTMLElement;
		bar.dispatchEvent(new PointerEvent('pointerdown', { bubbles: true }));
		expect(onPress).toHaveBeenCalledOnce();
	});

	it('opens the transfer dialog from the keyboard', async () => {
		const onPress = vi.fn<(event: PointerEvent) => void>();
		render(FuelBar, { value: 30, capacity: 50, onPress });
		const bar = page.getByRole('button', { name: 'Transfer fuel' }).element() as HTMLElement;
		bar.focus();
		await userEvent.keyboard('{Enter}');
		expect(onPress).toHaveBeenCalledOnce();
	});

	it('stays a slider while editable, even with a press handler', async () => {
		const onPress = vi.fn<(event: PointerEvent) => void>();
		const valuechanged = vi.fn<(value: number) => void>();
		render(FuelBar, { value: 30, capacity: 50, editable: true, onPress, valuechanged });
		const bar = page.getByRole('slider', { name: 'Fuel amount' }).element() as HTMLElement;
		bar.focus();
		await userEvent.keyboard('{ArrowRight}');
		expect(onPress).not.toHaveBeenCalled();
		expect(valuechanged).toHaveBeenCalled();
	});
});

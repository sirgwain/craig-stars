import { page, userEvent } from 'vitest/browser';
import { describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import FuelBar from './FuelBar.svelte';
import MineralBar from './MineralBar.svelte';
import WarpSpeedGauge from './WarpSpeedGauge.svelte';

describe('Keyboard gauge controls', () => {
	it('commits fuel changes and respects the available capacity', async () => {
		const valuechanged = vi.fn();
		render(FuelBar, { value: 25, capacity: 50, editable: true, valuechanged });
		const slider = page.getByRole('slider', { name: 'Fuel amount' });
		(slider.element() as HTMLElement).focus();
		await userEvent.keyboard('{ArrowRight}');
		await expect.element(slider).toHaveAttribute('aria-valuenow', '26');
		expect(valuechanged).toHaveBeenLastCalledWith(26);
		await userEvent.keyboard('{End}{ArrowRight}');
		await expect.element(slider).toHaveAttribute('aria-valuenow', '50');
	});

	it('uses the accepted mineral transfer amount returned by the callback', async () => {
		const onValueChanged = vi.fn(() => 24);
		render(MineralBar, { value: 25, capacity: 50, onValueChanged });
		const slider = page.getByRole('slider', { name: 'ironium amount' });
		(slider.element() as HTMLElement).focus();
		await userEvent.keyboard('{ArrowLeft}');
		expect(onValueChanged).toHaveBeenCalledWith(24);
		await expect.element(slider).toHaveAttribute('aria-valuenow', '24');
	});

	it('commits warp changes and exposes the stargate choice', async () => {
		const onValueChanged = vi.fn();
		render(WarpSpeedGauge, {
			value: 8,
			min: 0,
			max: 11,
			stargateSpeed: 11,
			useStargate: true,
			onValueChanged
		});
		const slider = page.getByRole('slider', { name: 'Warp speed' });
		(slider.element() as HTMLElement).focus();
		await userEvent.keyboard('{End}');
		expect(onValueChanged).toHaveBeenLastCalledWith(11);
		await expect.element(slider).toHaveAttribute('aria-valuetext', 'Use Stargate');
		await userEvent.keyboard('{Home}');
		await expect.element(slider).toHaveAttribute('aria-valuenow', '0');
	});

	it('keeps readonly fuel out of the tab order and rejects keyboard changes', async () => {
		const valuechanged = vi.fn();
		render(FuelBar, { value: 25, capacity: 50, valuechanged });
		const slider = page.getByRole('slider', { name: 'Fuel amount' });
		await expect.element(slider).toHaveAttribute('aria-readonly', 'true');
		expect(slider.element().hasAttribute('tabindex')).toBe(false);
		slider.element().dispatchEvent(new KeyboardEvent('keydown', { key: 'End', bubbles: true }));
		expect(valuechanged).not.toHaveBeenCalled();
		await expect.element(slider).toHaveAttribute('aria-valuenow', '25');
	});
});

import '../../../../css/app.css';
import { showTooltip, tooltipComponent } from '#lib/services/Stores.js';
import { tick } from 'svelte';
import { get } from 'svelte/store';
import { afterEach, expect, it } from 'vitest';
import { render } from 'vitest-browser-svelte';
import TextTooltip, { type TextTooltipProps } from './TextTooltip.svelte';
import Tooltip from './Tooltip.svelte';

let clicks = 0;
const countClick = () => clicks++;
const click = () => document.body.dispatchEvent(new MouseEvent('click', { bubbles: true }));

afterEach(() => {
	tooltipComponent.set(undefined);
	window.removeEventListener('click', countClick);
	clicks = 0;
});

async function open() {
	showTooltip<TextTooltipProps>(10, 10, TextTooltip, { text: 'Scout' });
	await tick();
}

it('closes on release of the press that opened it without eating its click', async () => {
	const { unmount } = render(Tooltip);
	window.addEventListener('click', countClick);

	await open();
	window.dispatchEvent(new PointerEvent('pointerup'));
	click();

	expect(get(tooltipComponent)).toBeUndefined();
	expect(clicks).toBe(1);
	unmount();
});

it('swallows the click from a tap that dismisses it', async () => {
	const { unmount } = render(Tooltip);
	window.addEventListener('click', countClick);

	await open();
	// tapping the tooltip to dismiss it must not click whatever is underneath
	window.dispatchEvent(new PointerEvent('pointerdown'));
	window.dispatchEvent(new PointerEvent('pointerup'));
	click();

	expect(get(tooltipComponent)).toBeUndefined();
	expect(clicks).toBe(0);

	// later clicks go through
	click();
	expect(clicks).toBe(1);
	unmount();
});

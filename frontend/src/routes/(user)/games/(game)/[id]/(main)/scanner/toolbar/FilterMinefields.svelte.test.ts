import { PlayerSettings } from '#lib/types/PlayerSettings.js';
import { flushSync, mount, unmount } from 'svelte';
import { get, writable } from 'svelte/store';
import { expect, it, vi } from 'vitest';
import FilterMinefields from './FilterMinefields.svelte';

const state = vi.hoisted(() => ({ context: undefined as unknown }));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));

it('chooses whose minefields are shown', () => {
	const settings = writable(new PlayerSettings());
	state.context = { settings };
	const target = document.createElement('div');
	document.body.appendChild(target);
	const menu = mount(FilterMinefields, { target });
	const item = (label: string) =>
		[...target.querySelectorAll('button')].find((b) => b.textContent?.includes(label))!;
	const click = (label: string) => {
		item(label).click();
		flushSync();
	};
	const button = target.querySelector('[data-filter="minefields"]')!;

	expect(item('All Mine Fields').ariaChecked).toBe('true');
	expect(button.classList).toContain('btn-primary');

	click('Mine Fields of Neutrals');
	expect(get(settings).showMinefieldsOf).toEqual(['mine', 'friends', 'enemies']);
	expect(item('All Mine Fields').ariaChecked).toBe('false');
	expect(item('Mine Fields of Neutrals').ariaChecked).toBe('false');

	click('No Mine Fields');
	expect(get(settings).showMinefields).toBe(false);
	expect(item('No Mine Fields').ariaChecked).toBe('true');
	expect(button.classList).not.toContain('btn-primary');

	// picking one turns minefields back on
	click('Mine Fields of Enemies');
	expect(get(settings).showMinefields).toBe(true);
	expect(get(settings).showMinefieldsOf).toEqual(['enemies']);

	click('All Mine Fields');
	expect(get(settings).showMinefieldsOf).toEqual(['mine', 'friends', 'neutrals', 'enemies']);

	unmount(menu);
	target.remove();
});

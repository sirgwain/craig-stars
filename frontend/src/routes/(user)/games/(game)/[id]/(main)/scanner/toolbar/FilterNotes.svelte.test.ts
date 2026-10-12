import { Universe } from '#lib/services/Universe.js';
import { MapObjectType, PlayerNoteSchema } from '#lib/types/cs-proto.js';
import { filterNote } from '#lib/types/Filter.js';
import { PlayerSettings } from '#lib/types/PlayerSettings.js';
import { create } from '@bufbuild/protobuf';
import { flushSync, mount, unmount } from 'svelte';
import { get, writable } from 'svelte/store';
import { expect, it, vi } from 'vitest';
import FilterNotes from './FilterNotes.svelte';

const state = vi.hoisted(() => ({ context: undefined as unknown }));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));

it('chooses which notes are shown', () => {
	const planetNote = (targetNum: number, tags: string[]) =>
		create(PlayerNoteSchema, { targetType: MapObjectType.PLANET, targetNum, tags, text: 'note' });
	const target = planetNote(1, ['Target']);
	const both = planetNote(2, ['target', 'gate hub']);
	const plain = planetNote(3, []);
	const label = create(PlayerNoteSchema, { text: 'the front' });
	const u = new Universe();
	u.setNotes([target, both, plain, label]);

	const settings = writable(new PlayerSettings());
	state.context = { settings, universe: writable(u) };
	const element = document.createElement('div');
	document.body.appendChild(element);
	const menu = mount(FilterNotes, { target: element });
	const item = (text: string) =>
		[...element.querySelectorAll('button')].find((b) => b.textContent?.includes(text))!;
	const click = (text: string) => {
		item(text).click();
		flushSync();
	};
	const button = element.querySelector('[data-filter="notes"]')!;
	const shown = () => [target, both, plain, label].map((n) => filterNote(n, get(settings)));

	expect(item('All Notes').ariaChecked).toBe('true');
	expect(button.classList).toContain('btn-primary');
	expect(shown()).toEqual([true, true, true, true]);

	// a note with another tag that is still shown stays on the map
	click('Target');
	expect(item('Target').ariaChecked).toBe('false');
	expect(item('All Notes').ariaChecked).toBe('false');
	expect(shown()).toEqual([false, true, true, true]);

	click('Map Labels');
	click('Notes Without Tags');
	expect(shown()).toEqual([false, true, false, false]);

	click('No Notes');
	expect(get(settings).showNotes).toBe(false);
	expect(button.classList).not.toContain('btn-primary');
	expect(shown()).toEqual([false, false, false, false]);

	// picking one turns notes back on with only that kind
	click('Map Labels');
	expect(get(settings).showNotes).toBe(true);
	expect(shown()).toEqual([false, false, false, true]);

	click('All Notes');
	expect(shown()).toEqual([true, true, true, true]);

	unmount(menu);
	element.remove();
});

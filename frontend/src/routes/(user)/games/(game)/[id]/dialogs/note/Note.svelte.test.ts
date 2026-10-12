import { Universe } from '#lib/services/Universe.js';
import { FleetSchema, MapObjectType, PlanetSchema, PlayerNoteSchema } from '#lib/types/cs-proto.js';
import { create } from '@bufbuild/protobuf';
import { flushSync, mount, unmount } from 'svelte';
import { writable } from 'svelte/store';
import { afterEach, expect, it, vi } from 'vitest';
import Note from './Note.svelte';

const state = vi.hoisted(() => ({ context: undefined as unknown }));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));

const savePlayerNote = vi.fn(async () => {});
const deletePlayerNote = vi.fn(async () => {});
const onClose = vi.fn();
let cleanup = () => {};

afterEach(() => {
	cleanup();
	vi.clearAllMocks();
});

type NoteData = Parameters<typeof create<typeof PlayerNoteSchema>>[1];

function open(noteData: NoteData, readOnly = false, otherNotes: NoteData[] = []) {
	const universe = new Universe();
	universe.fleets = [
		create(FleetSchema, { mapObject: { type: MapObjectType.FLEET, num: 4, name: 'Scout #4' } })
	];
	universe.setNotes(otherNotes.map((n) => create(PlayerNoteSchema, n)));
	universe.allPlanets = [
		create(PlanetSchema, { mapObject: { type: MapObjectType.PLANET, num: 1, name: 'Alpha' } })
	];
	state.context = {
		universe: writable(universe),
		readOnly: writable(readOnly),
		savePlayerNote,
		deletePlayerNote
	};
	const target = document.createElement('div');
	document.body.appendChild(target);
	const component = mount(Note, {
		target,
		props: { note: create(PlayerNoteSchema, noteData), onClose }
	});
	cleanup = () => {
		unmount(component);
		target.remove();
	};

	const textarea = target.querySelector('textarea')!;
	return {
		target,
		textarea,
		type: (text: string) => {
			textarea.value = text;
			textarea.dispatchEvent(new Event('input', { bubbles: true }));
			flushSync();
		},
		search: (text: string) => {
			const input = target.querySelector<HTMLInputElement>('[data-type="tag-search"]')!;
			input.value = text;
			input.dispatchEvent(new Event('input', { bubbles: true }));
			flushSync();
		},
		enter: () => {
			target
				.querySelector<HTMLInputElement>('[data-type="tag-search"]')!
				.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true }));
			flushSync();
		},
		click: (type: string) => {
			target.querySelector<HTMLButtonElement>(`[data-type="${type}"]`)!.click();
			flushSync();
		}
	};
}

it('saves a new note on a planet', async () => {
	const note = open({ targetType: MapObjectType.PLANET, targetNum: 1 });
	expect(note.target.querySelector('h2')?.textContent).toBe('Notes: Alpha');
	// there's nothing to delete yet
	expect(note.target.querySelector('[data-type="note-delete"]')).toBeNull();

	note.type('colonize next');
	note.target.querySelector<HTMLInputElement>('input[type="checkbox"]')!.click();
	flushSync();
	note.click('note-save');

	expect(savePlayerNote).toHaveBeenCalledWith(
		expect.objectContaining({ targetNum: 1, text: 'colonize next', showOnMap: true })
	);
	await vi.waitFor(() => expect(onClose).toHaveBeenCalledOnce());
});

it('tags a planet without writing a note', () => {
	const note = open({ targetType: MapObjectType.PLANET, targetNum: 1 }, false, [
		{ targetType: MapObjectType.PLANET, targetNum: 2, tags: ['Target', 'gate hub'] },
		{ targetType: MapObjectType.PLANET, targetNum: 3, tags: ['target'] }
	]);
	const toggles = () =>
		[...note.target.querySelectorAll<HTMLButtonElement>('[data-type="tag-toggle"]')].map(
			(b) => `${b.dataset.id}:${b.ariaPressed}`
		);
	// tags used on other planets can be toggled, the most used first
	expect(toggles()).toEqual(['Target:false', 'gate hub:false']);

	note.target
		.querySelector<HTMLButtonElement>('[data-type="tag-toggle"][data-id="Target"]')!
		.click();
	flushSync();
	expect(toggles()).toEqual(['Target:true', 'gate hub:false']);

	// typing filters the tags and offers to create a new one
	note.search('ga');
	expect(toggles()).toEqual(['gate hub:false']);
	expect(note.target.querySelector('[data-type="tag-create"]')?.textContent).toContain('ga');
	note.search(' my  colony ');
	expect(toggles()).toEqual([]);
	note.click('tag-create');
	expect(toggles()).toEqual(['Target:true', 'gate hub:false', 'my colony:true']);

	// a tag that only differs by case is the same tag, enter toggles it off again
	note.search('TARGET');
	expect(note.target.querySelector('[data-type="tag-create"]')).toBeNull();
	note.enter();
	expect(toggles()).toEqual(['Target:false', 'gate hub:false', 'my colony:true']);
	expect(savePlayerNote).not.toHaveBeenCalled();

	note.click('note-save');
	expect(savePlayerNote).toHaveBeenCalledWith(
		expect.objectContaining({ targetNum: 1, text: '', tags: ['my colony'], showOnMap: false })
	);
});

it('deletes a note when its text and tags are cleared', () => {
	const note = open({
		gameDbObject: { id: 7n },
		targetType: MapObjectType.PLANET,
		targetNum: 1,
		tags: ['Target']
	});

	note.target.querySelector<HTMLButtonElement>('[data-type="tag-toggle"]')!.click();
	flushSync();
	note.click('note-save');

	expect(savePlayerNote).not.toHaveBeenCalled();
	expect(deletePlayerNote).toHaveBeenCalledOnce();
});

it("doesn't tag labels or show fleet notes on the map", () => {
	const label = open({ position: { x: 10, y: 20 } });
	expect(label.target.querySelector('[data-type="tag-picker"]')).toBeNull();
	cleanup();

	const fleet = open({ targetType: MapObjectType.FLEET, targetNum: 4, targetPlayerNum: 1 });
	expect(fleet.target.querySelector('h2')?.textContent).toBe('Notes: Scout #4');
	expect(fleet.target.querySelector('[data-type="tag-picker"]')).not.toBeNull();
	expect(fleet.target.querySelector('input[type="checkbox"]')).toBeNull();
});

it('always shows labels on the map', () => {
	const note = open({ position: { x: 10, y: 20 } });
	expect(note.target.querySelector('h2')?.textContent).toBe('New Map Label');
	expect(note.target.querySelector('input[type="checkbox"]')).toBeNull();

	note.type('the front');
	note.textarea.dispatchEvent(
		new KeyboardEvent('keydown', { key: 'Enter', ctrlKey: true, bubbles: true })
	);

	expect(savePlayerNote).toHaveBeenCalledWith(
		expect.objectContaining({ text: 'the front', showOnMap: true })
	);
});

it('deletes a note when its text is cleared', () => {
	const note = open({ gameDbObject: { id: 7n }, position: { x: 10, y: 20 }, text: 'the front' });

	note.type('  ');
	note.click('note-save');

	expect(savePlayerNote).not.toHaveBeenCalled();
	expect(deletePlayerNote).toHaveBeenCalledOnce();
});

it('shows why a note failed to save and stays open', async () => {
	savePlayerNote.mockRejectedValueOnce(new Error('offline'));
	const note = open({ position: { x: 10, y: 20 } });

	note.type('the front');
	note.click('note-save');

	await vi.waitFor(() =>
		expect(note.target.querySelector('[data-type="note-error"]')?.textContent).toBe(
			'Failed to save note'
		)
	);
	expect(onClose).not.toHaveBeenCalled();
});

it("can't be changed while viewing a turn read-only", () => {
	const note = open({ gameDbObject: { id: 7n }, position: { x: 10, y: 20 }, text: 'hi' }, true);

	expect(note.textarea.readOnly).toBe(true);
	expect(note.target.querySelector('[data-type="note-save"]')).toBeNull();
	expect(note.target.querySelector('[data-type="note-delete"]')).toBeNull();
});

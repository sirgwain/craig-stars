import { Universe } from '#lib/services/Universe.js';
import { MapObjectType, PlanetSchema, PlayerNoteSchema } from '#lib/types/cs-proto.js';
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

function open(noteData: Parameters<typeof create<typeof PlayerNoteSchema>>[1], readOnly = false) {
	const universe = new Universe();
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
		click: (type: string) => {
			target.querySelector<HTMLButtonElement>(`[data-type="${type}"]`)!.click();
			flushSync();
		}
	};
}

it('saves a new note on a planet', async () => {
	const note = open({ targetType: MapObjectType.PLANET, targetNum: 1 });
	expect(note.target.querySelector('h2')?.textContent).toBe('Note: Alpha');
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

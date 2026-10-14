import { Universe } from '#lib/services/Universe.js';
import { PlayerIntelSchema, PlayerMailSchema, type PlayerMail } from '#lib/types/cs-proto.js';
import { create } from '@bufbuild/protobuf';
import { flushSync, mount, unmount } from 'svelte';
import { writable } from 'svelte/store';
import { afterEach, expect, it, vi } from 'vitest';
import Mail from './Mail.svelte';

const state = vi.hoisted(() => ({ context: undefined as unknown }));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));

type MailData = Parameters<typeof create<typeof PlayerMailSchema>>[1];

// the server gives new mail an id, and sets who it's from and the year
const savePlayerMail = vi.fn(async (mail: PlayerMail) =>
	create(PlayerMailSchema, {
		gameDbObject: { id: mail.gameDbObject?.id || 99n },
		year: 2401,
		fromPlayerNum: 1,
		toPlayerNum: mail.toPlayerNum,
		text: mail.text
	})
);
const deletePlayerMail = vi.fn(async () => {});
const onClose = vi.fn();
let cleanup = () => {};

afterEach(() => {
	cleanup();
	vi.clearAllMocks();
});

async function open(toPlayerNum: number, existing: MailData[] = [], readOnly = false) {
	const universe = new Universe();
	universe.playerIntels = [
		create(PlayerIntelSchema, { num: 1, seen: true, racePluralName: 'Humanoids' }),
		create(PlayerIntelSchema, { num: 2, seen: true, racePluralName: 'Rabbitoids' }),
		// a player we haven't met
		create(PlayerIntelSchema, { num: 3 })
	];
	state.context = {
		game: writable({ year: 2401 }),
		player: writable({ num: 1 }),
		universe: writable(universe),
		readOnly: writable(readOnly),
		getPlayerMail: vi.fn(async () => existing.map((m) => create(PlayerMailSchema, m))),
		savePlayerMail,
		deletePlayerMail
	};
	const target = document.createElement('div');
	document.body.appendChild(target);
	const component = mount(Mail, { target, props: { toPlayerNum, onClose } });
	cleanup = () => {
		unmount(component);
		target.remove();
	};
	const items = () =>
		[...target.querySelectorAll('[data-type="mail-item"]')].map((li) =>
			li.textContent?.replace(/\s+/g, ' ').trim()
		);
	await vi.waitFor(() => {
		flushSync();
		expect(items()).toHaveLength(existing.length);
	});

	return {
		target,
		items,
		textarea: () => target.querySelector('textarea')!,
		to: () => target.querySelector<HTMLSelectElement>('[data-type="mail-to"]')!,
		type: (text: string) => {
			const textarea = target.querySelector('textarea')!;
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

it('sends mail to the player being replied to', async () => {
	const mail = await open(2);
	expect([...mail.to().options].map((o) => o.textContent)).toEqual([
		'Everybody',
		'Humanoids',
		'Rabbitoids',
		'Player 3'
	]);
	expect(mail.to().value).toBe('2');
	// there's nothing to send yet
	expect(mail.target.querySelector<HTMLButtonElement>('[data-type="mail-save"]')!.disabled).toBe(
		true
	);

	mail.type('truce?');
	mail.click('mail-save');

	expect(savePlayerMail).toHaveBeenCalledWith(
		expect.objectContaining({ toPlayerNum: 2, text: 'truce?' })
	);
	// the dialog stays open and shows the mail waiting to be sent
	await vi.waitFor(() => {
		flushSync();
		expect(mail.items()).toEqual(['You to Rabbitoids, 2401 Sends next turn Edit Delete truce?']);
	});
	expect(mail.textarea().value).toBe('');
	expect(onClose).not.toHaveBeenCalled();
});

it('changes and deletes mail that has not been delivered', async () => {
	const mail = await open(0, [
		{ gameDbObject: { id: 1n }, year: 2400, fromPlayerNum: 2, toPlayerNum: 1, text: 'hello' },
		{ gameDbObject: { id: 2n }, year: 2400, fromPlayerNum: 1, toPlayerNum: 0, text: 'old news' },
		{ gameDbObject: { id: 3n }, year: 2401, fromPlayerNum: 1, toPlayerNum: 2, text: 'trcue?' }
	]);
	// the newest is first, and only undelivered mail can be changed
	expect(mail.items()).toEqual([
		'You to Rabbitoids, 2401 Sends next turn Edit Delete trcue?',
		'You to Everybody, 2400 old news',
		'Rabbitoids to you, 2400 hello'
	]);

	mail.click('mail-edit');
	expect(mail.textarea().value).toBe('trcue?');
	expect(mail.to().value).toBe('2');
	mail.type('truce?');
	mail.click('mail-save');
	expect(savePlayerMail).toHaveBeenCalledWith(
		expect.objectContaining({ gameDbObject: expect.objectContaining({ id: 3n }), text: 'truce?' })
	);
	await vi.waitFor(() => {
		flushSync();
		expect(mail.items()[0]).toContain('truce?');
	});
	expect(mail.items()).toHaveLength(3);

	mail.click('mail-delete');
	expect(deletePlayerMail).toHaveBeenCalledWith(
		expect.objectContaining({ gameDbObject: expect.objectContaining({ id: 3n }) })
	);
	await vi.waitFor(() => {
		flushSync();
		expect(mail.items()).toHaveLength(2);
	});
});

it('only shows mail when the game is read only', async () => {
	const mail = await open(
		0,
		[{ gameDbObject: { id: 3n }, year: 2401, fromPlayerNum: 1, toPlayerNum: 2, text: 'truce?' }],
		true
	);
	expect(mail.target.querySelector('textarea')).toBeNull();
	expect(mail.target.querySelector('[data-type="mail-save"]')).toBeNull();
	expect(mail.items()).toEqual(['You to Rabbitoids, 2401 Sends next turn truce?']);
});

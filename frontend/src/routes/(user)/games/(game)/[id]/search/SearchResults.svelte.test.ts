import '../../../../../../css/app.css';
import { Universe } from '#lib/services/Universe.js';
import { MapObjectType, type Fleet, type MysteryTrader } from '#lib/types/cs-proto.js';
import { CommandedPlayer } from '#lib/types/Player.js';
import { PlayerSettings } from '#lib/types/PlayerSettings.js';
import { writable } from 'svelte/store';
import { describe, expect, it, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import SearchResults from './SearchResults.svelte';

const state = vi.hoisted(() => ({ context: undefined as unknown }));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));

function mapObject(type: MapObjectType, num: number, name: string) {
	return { mapObject: { type, num, name, playerNum: 0, position: { x: 0, y: 0 } } };
}

function setup(searchQuery: string) {
	const player = new CommandedPlayer();
	player.num = 1;
	const settings = new PlayerSettings();
	settings.searchQuery = searchQuery;
	const universe = new Universe();
	universe.fleets = [mapObject(MapObjectType.FLEET, 1, 'Zed Scout')] as unknown as Fleet[];
	universe.mysteryTraders = [
		mapObject(MapObjectType.MYSTERY_TRADER, 1, 'Zed Trader'),
		mapObject(MapObjectType.MYSTERY_TRADER, 2, 'Zed Trader Two')
	] as unknown as MysteryTrader[];
	state.context = {
		player: writable(player),
		settings: writable(settings),
		universe: writable(universe)
	};
}

describe('SearchResults', () => {
	it('has a close button when it can be cancelled', async () => {
		setup('');
		const onCancel = vi.fn();
		render(SearchResults, { onOk: vi.fn(), onCancel });
		await page.getByRole('button', { name: 'Close search' }).click();
		expect(onCancel).toHaveBeenCalledOnce();
	});

	it('has no close button on the search page', async () => {
		setup('');
		render(SearchResults, { onOk: vi.fn() });
		await expect.element(page.getByPlaceholder('search')).toBeInTheDocument();
		expect(page.getByRole('button', { name: 'Close search' }).elements()).toHaveLength(0);
	});

	it('picks the highlighted mystery trader and stops at the last result', async () => {
		setup('zed');
		const onOk = vi.fn();
		render(SearchResults, { onOk });
		// fleet, trader, trader two; extra presses stay on the last one
		await userEvent.keyboard('{ArrowDown}{Enter}');
		expect(onOk.mock.calls[0][0].mapObject.name).toBe('Zed Trader');
		await userEvent.keyboard('{ArrowDown}{ArrowDown}{ArrowDown}{Enter}');
		expect(onOk.mock.calls[1][0].mapObject.name).toBe('Zed Trader Two');
	});
});

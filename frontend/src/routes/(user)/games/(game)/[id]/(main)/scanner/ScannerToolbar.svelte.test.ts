import '../../../../../../../css/app.css';
import { Universe } from '#lib/services/Universe.js';
import { CommandedPlayer } from '#lib/types/Player.js';
import { PlayerSettings } from '#lib/types/PlayerSettings.js';
import { tick } from 'svelte';
import { writable } from 'svelte/store';
import { beforeAll, describe, expect, it, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import ScannerToolbar from './ScannerToolbar.svelte';

const state = vi.hoisted(() => ({ context: undefined as unknown }));
vi.mock('#lib/services/GameContext.js', () => ({ getGameContext: () => state.context }));

type Stage = 'full' | 'views' | 'menu';

// render the toolbar in a container of the given width and report how it laid out
async function layout(width: number, shareMap: boolean) {
	const player = new CommandedPlayer();
	player.num = 1;
	if (shareMap) {
		player.relations = [{ shareMap: true } as never];
	}
	state.context = {
		player: writable(player),
		settings: writable(new PlayerSettings()),
		universe: writable(new Universe())
	};

	const host = document.createElement('div');
	host.style.width = `${width}px`;
	document.body.appendChild(host);
	const { container, unmount } = render(
		ScannerToolbar,
		{
			onShowSearch: vi.fn(),
			onCycleMapObject: vi.fn(),
			onNextMapObject: vi.fn(),
			onPreviousMapObject: vi.fn()
		},
		{ baseElement: host }
	);
	// let the toolbar measure itself and settle on a stage
	await tick();
	await new Promise((resolve) => requestAnimationFrame(() => resolve(undefined)));
	await tick();

	const navbar = container.querySelector<HTMLElement>('.navbar')!;
	const inToolbar = (href: string) => !!navbar.querySelector(`a[href="${href}"]:not(details a)`);
	const stage: Stage = inToolbar('#scanner-toggle')
		? 'full'
		: inToolbar('#normal-view')
			? 'views'
			: 'menu';
	const result = { stage, overflow: navbar.scrollWidth - navbar.clientWidth };
	unmount();
	host.remove();
	return result;
}

describe('ScannerToolbar', () => {
	beforeAll(async () => {
		// wide enough that the cycle button for phones is hidden
		await page.viewport(1280, 800);
	});

	it.each([
		// the filters don't fit on portrait iPads, so they collapse and the views stay
		{ width: 820, shareMap: false, stage: 'views' },
		{ width: 834, shareMap: false, stage: 'views' },
		{ width: 834, shareMap: true, stage: 'views' },
		// landscape iPads fit the whole toolbar
		{ width: 1024, shareMap: false, stage: 'full' },
		{ width: 1024, shareMap: true, stage: 'full' },
		// the toolbar over the scanner, next to a two column command pane
		{ width: 780, shareMap: false, stage: 'views' },
		{ width: 520, shareMap: false, stage: 'views' },
		{ width: 400, shareMap: false, stage: 'menu' }
	])('$width wide, shareMap=$shareMap shows $stage', async ({ width, shareMap, stage }) => {
		const result = await layout(width, shareMap);
		expect(result.stage).toBe(stage);
		expect(result.overflow).toBeLessThanOrEqual(0);
	});

	it.each([false, true])('never overflows, shareMap=%s', async (shareMap) => {
		const stages: Stage[] = [];
		for (let width = 340; width <= 1000; width += 6) {
			const result = await layout(width, shareMap);
			expect(result.overflow, `overflowed at ${width} (${result.stage})`).toBeLessThanOrEqual(0);
			stages.push(result.stage);
		}
		// as it gets wider it only ever expands
		expect(stages).toEqual([...stages].sort((a, b) => order(a) - order(b)));
		expect(new Set(stages)).toEqual(new Set(['menu', 'views', 'full']));
	});
});

function order(stage: Stage) {
	return ['menu', 'views', 'full'].indexOf(stage);
}

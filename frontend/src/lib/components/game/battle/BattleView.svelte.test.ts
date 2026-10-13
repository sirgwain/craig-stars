import '../../../../css/app.css';
import { Universe } from '#lib/services/Universe.js';
import {
	BattleRecordSchema,
	PlayerUniverseSchema,
	ShipDesignSchema,
	type BattleRecordTokenActionJson
} from '#lib/types/cs-proto.js';
import { create, fromJson } from '@bufbuild/protobuf';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';
import BattleView from './BattleView.svelte';

// the viewer works without a game (no game context), like on the battle test page
vi.mock('#lib/services/GameContext.js', () => ({
	gameKey: Symbol('game'),
	designFinderKey: Symbol('designFinder'),
	playerFinderKey: Symbol('playerFinder'),
	getGameContext: () => undefined
}));

// two tokens trading shots: a beam on shields, a beam through to armor, a torpedo salvo that
// destroys a ship
const actions: BattleRecordTokenActionJson[] = [
	{
		type: 'BATTLE_RECORD_TOKEN_ACTION_TYPE_BEAM_FIRE',
		tokenNum: 1,
		round: 1,
		slot: 1,
		from: { x: 1, y: 5 },
		to: { x: 6, y: 5 },
		targetNum: 2,
		damageDoneShields: 50
	},
	{
		type: 'BATTLE_RECORD_TOKEN_ACTION_TYPE_BEAM_FIRE',
		tokenNum: 1,
		round: 1,
		slot: 1,
		from: { x: 1, y: 5 },
		to: { x: 6, y: 5 },
		targetNum: 2,
		damageDoneShields: 10,
		damageDoneArmor: 30
	},
	{
		type: 'BATTLE_RECORD_TOKEN_ACTION_TYPE_TORPEDO_FIRE',
		tokenNum: 2,
		round: 1,
		slot: 1,
		from: { x: 6, y: 5 },
		to: { x: 1, y: 5 },
		targetNum: 1,
		target: { designNum: 1, quantity: 1 },
		damageDoneArmor: 40,
		torpedoHits: 2,
		torpedoMisses: 1,
		tokensDestroyed: 1
	}
];

function renderBattle() {
	const battle = fromJson(BattleRecordSchema, {
		num: 1,
		tokens: [
			{ num: 1, playerNum: 1, designNum: 1, position: { x: 1, y: 5 }, startingQuantity: 2 },
			{ num: 2, playerNum: 2, designNum: 1, position: { x: 6, y: 5 }, startingQuantity: 1 }
		],
		actionsPerRound: [{ actions: [] }, { actions }]
	});
	const universe = new Universe();
	universe.setData(
		1,
		create(PlayerUniverseSchema, {
			designs: [
				create(ShipDesignSchema, { playerNum: 1, num: 1, name: 'Shooter', hull: 'Destroyer' }),
				create(ShipDesignSchema, { playerNum: 2, num: 1, name: 'Target', hull: 'Destroyer' })
			]
		})
	);
	return render(BattleView, {
		battleRecord: battle,
		playerFinder: universe,
		designFinder: universe
	});
}

const phaseInput = () => page.getByRole('spinbutton', { name: 'phase' });
const animations = () => page.getByRole('checkbox', { name: 'Animations' });
const damage = () => page.getByRole('checkbox', { name: 'Damage' });
const play = () => page.getByRole('button', { name: 'Play', exact: true });
const pause = () => page.getByRole('button', { name: 'Pause' });
const next = () => page.getByRole('button', { name: 'next', exact: true });
const savedPrefs = () => JSON.parse(localStorage.getItem('cs.battle.playback') ?? '{}');

describe('BattleView', () => {
	beforeEach(() => localStorage.clear());

	it('waits at the start until it is played', async () => {
		renderBattle();
		await expect.element(play()).toBeVisible();
		await expect.element(phaseInput()).toHaveValue(0);
		await expect.element(animations()).toBeChecked();
		expect(pause().elements()).toHaveLength(0);

		await page.getByRole('button', { name: '4x' }).click();
		await play().click();
		await expect.element(phaseInput(), { timeout: 10_000 }).toHaveValue(3);
		await expect.element(page.getByRole('button', { name: 'Replay' })).toBeVisible();
	});

	it('steps through the battle, the shots of one weapon slot together', async () => {
		renderBattle();
		// the two beam shots of one weapon slot are one step
		await next().click();
		await expect.element(phaseInput()).toHaveValue(2);
		await page.getByRole('button', { name: 'previous' }).click();
		await expect.element(phaseInput()).toHaveValue(1);
		await page.getByRole('button', { name: 'end', exact: true }).click();
		await expect.element(phaseInput()).toHaveValue(3);
		await page.getByRole('button', { name: 'begin' }).click();
		await expect.element(phaseInput()).toHaveValue(0);
	});

	it('marks shield and armor damage on the target', async () => {
		renderBattle();
		// both beam shots hit the same target: their damage is added up
		await next().click();
		await expect.element(page.getByLabelText('30 armor damage')).toBeVisible();
		await expect.element(page.getByLabelText('60 shield damage')).toBeVisible();
		expect(getComputedStyle(page.getByLabelText('60 shield damage').element()).color).not.toEqual(
			getComputedStyle(page.getByLabelText('30 armor damage').element()).color
		);

		// the first shot alone
		await page.getByRole('button', { name: 'previous' }).click();
		await expect.element(page.getByLabelText('50 shield damage')).toBeVisible();
		expect(page.getByLabelText(/armor damage/).elements()).toHaveLength(0);

		// the torpedoes mark their damage when they land
		await page.getByRole('button', { name: 'end', exact: true }).click();
		await expect.element(page.getByLabelText('40 armor damage')).toBeVisible();
	});

	it('hides the damage marks when they are turned off, and remembers it', async () => {
		const { unmount } = renderBattle();
		await expect.element(damage()).toBeChecked();
		await next().click();
		await expect.element(page.getByLabelText('30 armor damage')).toBeVisible();
		await damage().click();
		await expect.element(damage()).not.toBeChecked();
		expect(page.getByTestId('battle-damage').elements()).toHaveLength(0);
		unmount();

		renderBattle();
		await expect.element(damage()).not.toBeChecked();
		await next().click();
		await expect.element(phaseInput()).toHaveValue(2);
		expect(page.getByTestId('battle-damage').elements()).toHaveLength(0);
	});

	it('remembers the speed and whether attacks are animated', async () => {
		const { unmount } = renderBattle();
		await page.getByRole('button', { name: '4x' }).click();
		await animations().click();
		await expect.element(animations()).not.toBeChecked();
		expect(savedPrefs()).toEqual({ animations: false, damage: true, speed: 4 });
		unmount();

		renderBattle();
		await expect.element(animations()).not.toBeChecked();
		await expect
			.element(page.getByRole('button', { name: '4x' }))
			.toHaveAttribute('aria-pressed', 'true');
		// it still waits to be played
		await expect.element(play()).toBeVisible();
	});

	it('leaves space to a focused button', async () => {
		renderBattle();
		const fast = page.getByRole('button', { name: '4x' });
		(fast.element() as HTMLElement).focus();
		await userEvent.keyboard(' ');
		await expect.element(fast).toHaveAttribute('aria-pressed', 'true');
		await expect.element(play()).toBeVisible();

		// elsewhere on the page it plays and pauses
		(document.activeElement as HTMLElement).blur();
		await userEvent.keyboard(' ');
		await expect.element(pause()).toBeVisible();
		await userEvent.keyboard(' ');
		await expect.element(play()).toBeVisible();
	});
});

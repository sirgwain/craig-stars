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
import { page } from 'vitest/browser';
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
const styleSelect = () => page.getByRole('combobox', { name: 'battle style' });
const savedPrefs = () => JSON.parse(localStorage.getItem('cs.battle.playback') ?? '{}');

describe('BattleView', () => {
	beforeEach(() => localStorage.clear());

	it('is the classic step by step viewer by default', async () => {
		renderBattle();
		await expect.element(styleSelect()).toHaveValue('classic');
		const options = (styleSelect().element() as HTMLSelectElement).options;
		expect(Array.from(options).map((o) => o.text)).toEqual(['Classic', 'Original']);
		expect(page.getByRole('button', { name: 'Play battle' }).elements()).toHaveLength(0);
		expect(page.getByRole('button', { name: '2x' }).elements()).toHaveLength(0);

		await page.getByRole('button', { name: 'next', exact: true }).click();
		await expect.element(phaseInput()).toHaveValue(1);
		await page.getByRole('button', { name: 'end' }).click();
		await expect.element(phaseInput()).toHaveValue(3);
		await page.getByRole('button', { name: 'previous' }).click();
		await expect.element(phaseInput()).toHaveValue(2);
	});

	it('marks shield and armor damage on the target', async () => {
		renderBattle();
		await page.getByRole('button', { name: 'next', exact: true }).click();
		await expect.element(page.getByLabelText('50 shield damage')).toBeVisible();
		expect(page.getByLabelText(/armor damage/).elements()).toHaveLength(0);

		await page.getByRole('button', { name: 'next', exact: true }).click();
		await expect.element(page.getByLabelText('30 armor damage')).toBeVisible();
		await expect.element(page.getByLabelText('10 shield damage')).toBeVisible();

		await page.getByRole('button', { name: 'next', exact: true }).click();
		await expect.element(page.getByLabelText('40 armor damage')).toBeVisible();
		await expect.element(page.getByLabelText('1 destroyed')).toBeVisible();

		// the marks look different for shields and armor
		const shield = page.getByLabelText('10 shield damage');
		const armor = page.getByLabelText('30 armor damage');
		await page.getByRole('button', { name: 'previous' }).click();
		expect(getComputedStyle(shield.element()).color).not.toEqual(
			getComputedStyle(armor.element()).color
		);
	});

	it('switches to the original playback, which plays on its own, and remembers it', async () => {
		const { unmount } = renderBattle();
		await styleSelect().selectOptions('Original');
		await expect.element(page.getByRole('button', { name: /pause/i })).toBeVisible();
		// the classic damage marks are not drawn over the original look
		expect(page.getByTestId('battle-damage').elements()).toHaveLength(0);
		await page.getByRole('button', { name: '4x' }).click();
		await expect.element(phaseInput(), { timeout: 10_000 }).toHaveValue(3);
		await expect.element(page.getByRole('button', { name: /replay/i })).toBeVisible();
		expect(savedPrefs()).toEqual({ style: 'original', speed: 4 });
		unmount();

		// it starts playing when a battle opens
		renderBattle();
		await expect.element(styleSelect()).toHaveValue('original');
		await expect.element(page.getByRole('button', { name: /pause/i })).toBeVisible();
		await expect
			.element(page.getByRole('button', { name: '4x' }))
			.toHaveAttribute('aria-pressed', 'true');
	});

	it('steps through the original playback', async () => {
		renderBattle();
		await styleSelect().selectOptions('Original');
		await page.getByRole('button', { name: /pause/i }).click();
		await page.getByRole('button', { name: 'begin' }).click();
		await expect.element(phaseInput()).toHaveValue(0);
		await expect.element(page.getByRole('button', { name: 'Play battle' })).toBeVisible();
		// the two beam shots of one weapon slot are one step
		await page.getByRole('button', { name: 'next', exact: true }).click();
		await expect.element(phaseInput()).toHaveValue(2);
		await page.getByRole('button', { name: 'previous' }).click();
		await expect.element(phaseInput()).toHaveValue(1);
	});

	it('back to classic stops the playback', async () => {
		renderBattle();
		await styleSelect().selectOptions('Original');
		await expect.element(page.getByRole('button', { name: /pause/i })).toBeVisible();
		await styleSelect().selectOptions('Classic');
		await expect.element(page.getByRole('button', { name: 'next attack' })).toBeVisible();
		expect(page.getByRole('button', { name: /pause/i }).elements()).toHaveLength(0);
		expect(savedPrefs().style).toBe('classic');
	});
});

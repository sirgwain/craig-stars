<script lang="ts">
	import { page } from '$app/stores';
	import ErrorPage from '$lib/components/ErrorPage.svelte';
	import ErrorToast from '$lib/components/ErrorToast.svelte';
	import Menu from '$lib/components/Menu.svelte';
	import Tooltip from '$lib/components/game/tooltips/Tooltip.svelte';
	import { bindNavigationHotkeys, unbindNavigationHotkeys } from '$lib/navigationHotkeys';
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
	import { getGameContext, initGameContext, updateGameContext } from '$lib/services/Contexts';
	import type { CSError } from '$lib/services/Errors';
	import { FullGame } from '$lib/services/FullGame';
	import { commandedMapObject, me, nextMapObject, previousMapObject } from '$lib/services/Stores';
	import { Universe } from '$lib/services/Universe';
	import { GameState } from '$lib/types/Game';
	import { Player } from '$lib/types/Player';
	import hotkeys from 'hotkeys-js';
	import { onDestroy, onMount } from 'svelte';
	import GameMenu from './GameMenu.svelte';
	import type { Unsubscriber } from 'svelte/store';

	initGameContext();
	const { game, player, universe } = getGameContext();

	let id = parseInt($page.params.id);
	let error: string | undefined = undefined;

	let fg: FullGame | undefined = undefined;
	let state: GameState;
	let year: number;
	let unsubscribe: Unsubscriber | undefined;

	onMount(async () => {
		try {
			// on first mount, load the game
			await loadGame();
			fg = $game;
		} catch (e) {
			const err = e as CSError;
			if (err?.statusCode === 404) {
				error = 'Game not found';
			} else {
				error = `${e}`;
			}
		}

		// setup the quantityModifier
		bindQuantityModifier();

		// if we are in an active game, bind the navigation hotkeys, i.e. F4 for research, Esc to go back
		if ($game?.state == GameState.WaitingForPlayers) {
			bindNavigationHotkeys(id, page);

			hotkeys('F9', () => {
				onSubmitTurn();
			});
			hotkeys('n', () => {
				nextMapObject();
			});
			hotkeys('p', () => {
				previousMapObject();
			});
		}
	});

	// async onMount means onDestroy needs to live
	// outside
	onDestroy(() => {
		updateGameContext(new FullGame(), new Player(), new Universe());

		unbindQuantityModifier();
		unbindNavigationHotkeys();
		hotkeys.unbind('F9');
		hotkeys.unbind('n');
		hotkeys.unbind('p');
	});

	async function loadGame() {
		// empty the context
		unsubscribe && unsubscribe();

		// load a new game, when this is successful, the $game contenxt will be updated
		const loaded = await $game.load(id);

		// update some local state we use for detecting changes
		state = loaded.state;
		year = loaded.year;

		// subscribe to game update events
		unsubscribe = game.subscribe(onGameChange);

		// if we are in setup mode or we have submitted our turn, start the
		// load player status job
		if (loaded.state == GameState.WaitingForPlayers && !$player.submittedTurn) {
			loaded.commandHomeWorld();
		}
	}

	// method called when game changes
	async function onGameChange() {
		if (state != $game.state) {
			// console.log('game state changed');
			await loadGame();
		}
		if (year != $game.year) {
			// console.log('game year changed');
			await loadGame();
		}
	}

	async function onSubmitTurn() {
		$game = await $game.submitTurn();
		// if a new turn generates, submittedTurn will be reset to false
		if (!$player.submittedTurn) {
			// command our homeworld after a new turn is generated
			$game.commandHomeWorld();
		}
	}
</script>

{#if fg}
	<main class="flex flex-col mb-20 md:mb-0">
		<div class="flex-initial">
			<GameMenu on:submit-turn={onSubmitTurn} />
		</div>
		<ErrorToast />
		<slot>Game</slot>
	</main>
	<Tooltip />
{:else if error}
	<main class="flex flex-col mb-20 md:mb-0">
		<div class="flex-initial">
			<Menu user={$me} />
		</div>
		<ErrorPage {error} />
	</main>
{/if}

<style>
	main {
		height: 100vh; /* Fallback for browsers that do not support Custom Properties */
		height: calc(var(--vh, 1vh) * 100);
	}
</style>

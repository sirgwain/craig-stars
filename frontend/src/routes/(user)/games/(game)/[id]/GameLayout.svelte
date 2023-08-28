<script lang="ts">
	import { page } from '$app/stores';
	import ErrorPage from '$lib/components/ErrorPage.svelte';
	import ErrorToast from '$lib/components/ErrorToast.svelte';
	import LoadingModal from '$lib/components/LoadingModal.svelte';
	import Menu from '$lib/components/Menu.svelte';
	import Popup from '$lib/components/game/tooltips/Popup.svelte';
	import Tooltip from '$lib/components/game/tooltips/Tooltip.svelte';
	import { bindNavigationHotkeys, unbindNavigationHotkeys } from '$lib/navigationHotkeys';
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
	import { clearGameContext, getGameContext, initGameContext } from '$lib/services/Contexts';
	import type { CSError } from '$lib/services/Errors';
	import type { FullGame } from '$lib/services/FullGame';
	import {
		clearLoadingModalText,
		loadingModalText,
		me,
		nextMapObject,
		previousMapObject,
		setLoadingModalText
	} from '$lib/services/Stores';
	import { GameState } from '$lib/types/Game';
	import { wait } from '$lib/wait';
	import hotkeys from 'hotkeys-js';
	import { onDestroy, onMount } from 'svelte';
	import type { Unsubscriber } from 'svelte/store';
	import GameMenu from './GameMenu.svelte';

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
			// empty the context
			clearGameContext();

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
	});

	// async onMount means onDestroy needs to live
	// outside
	onDestroy(() => {
		$game.stopPollingStatus();
		unsubscribe && unsubscribe();
		clearGameContext();

		unbindQuantityModifier();
		unbindNavigationHotkeys();
		hotkeys.unbind('F9', 'root');
		hotkeys.unbind('n', 'root');
		hotkeys.unbind('p', 'root');
		hotkeys.deleteScope('root');
	});

	async function loadGame() {
		setLoadingModalText('Loading game...');
		try {
			// empty the context
			unsubscribe && unsubscribe();

			unbindNavigationHotkeys();
			hotkeys.unbind('F9', 'root');
			hotkeys.unbind('n', 'root');
			hotkeys.unbind('p', 'root');
			hotkeys.deleteScope('root');

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

				bindNavigationHotkeys(id, page);

				hotkeys('F9', 'root', () => {
					onSubmitTurn();
				});
				hotkeys('n', 'root', () => {
					nextMapObject();
				});
				hotkeys('p', 'root', () => {
					previousMapObject();
				});
				hotkeys.setScope('root');
			}
		} finally {
			clearLoadingModalText();
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
		// console.time('onSubmitTurn');

		setLoadingModalText(`Submitting turn for ${$game.year}...`);
		const hideLoadingDelay = wait(500);
		try {
			// turn off game change listening
			unsubscribe && unsubscribe();

			// update the UI
			$game = await $game.submitTurn();

			// update our local state
			state = $game.state;
			year = $game.year;

			// resubscribe to game update events
			unsubscribe = game.subscribe(onGameChange);

			// if a new turn generates, submittedTurn will be reset to false
			if (!$player.submittedTurn) {
				// command our homeworld after a new turn is generated
				$game.commandHomeWorld();
			}
		} catch {
			// reload the game status in the event of an error
			await $game.loadStatus();
		} finally {
			// make sure this finishes
			await hideLoadingDelay;
			clearLoadingModalText();
			// console.timeEnd('onSubmitTurn');
		}
	}
</script>

{#if fg}
	<main class="flex flex-col h-[100dvh]">
		<header class="flex-none z-50">
			<GameMenu on:submit-turn={onSubmitTurn} />
		</header>
		<ErrorToast />
		<LoadingModal text={$loadingModalText} />
		<slot>Game</slot>
	</main>
	<Tooltip />
	<Popup />
{:else if error}
	<main class="flex flex-col">
		<div class="flex-initial">
			<Menu user={$me} />
		</div>
		<ErrorPage {error} />
	</main>
{/if}

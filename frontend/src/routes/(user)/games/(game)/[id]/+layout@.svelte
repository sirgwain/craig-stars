<script lang="ts">
	import { page } from '$app/stores';
	import ErrorPage from '$lib/components/ErrorPage.svelte';
	import Menu from '$lib/components/Menu.svelte';
	import { bindNavigationHotkeys, unbindNavigationHotkeys } from '$lib/navigationHotkeys';
	import { createGameContext, gameKey, type GameContext } from '$lib/services/GameContext';
	import type { CSError } from '$lib/services/Errors';
	import type { FullGame } from '$lib/services/FullGame';
	import { GameService } from '$lib/services/GameService';
	import { clearLoadingModalText, me, setLoadingModalText } from '$lib/services/Stores';
	import { GameState } from '$lib/types/Game';
	import { wait } from '$lib/wait';
	import hotkeys from 'hotkeys-js';
	import { onDestroy, onMount, setContext } from 'svelte';
	import type { Unsubscriber } from 'svelte/store';
	import { get } from 'svelte/store';
	import GameLayout from './GameLayout.svelte';
	import { goto } from '$app/navigation';

	let id = parseInt($page.params.id);

	let context: GameContext | undefined = undefined;
	let error: string | undefined = undefined;
	let contextSetup = false;

	let unsubscribe: Unsubscriber | undefined;
	let state: GameState;
	let year: number;

	onMount(async () => {
		try {
			setLoadingModalText('Loading game...');

			// on mount, load the game and setup the context used by the rest of the children
			const loaded = await GameService.loadFullGame(id);
			context = createGameContext(loaded);

			hotkeys.setScope('root');
		} catch (e) {
			const err = e as CSError;
			if (err?.statusCode === 404) {
				error = 'Game not found';
			} else {
				error = `${e}`;
			}
		} finally {
			clearLoadingModalText();
		}
	});

	onDestroy(() => {
		if (!context) return;

		hotkeys.deleteScope('root');

		unsubscribe && unsubscribe();
	});

	// update the context of the game
	$: {
		if (context && !contextSetup) {
			contextSetup = true;

			// store the latest state/year so we can reload if the game changes
			const game = get(context.game);
			state = game.state;
			year = game.year;
			context.commandHomeWorld();

			// subscribe to game change events so we do a full reload if the year/state changes
			unsubscribe = context.game.subscribe(onGameChange);

			// setup the context for our child components
			setContext(gameKey, context);
		}
	}

	// every time the game updates, check if we have a new year/state change
	// and if so, reset the context
	async function onGameChange(game: FullGame) {
		if (!context) return;

		if (state != game.state || year != game.year) {
			// console.log('game state changed');
			const loaded = await GameService.loadFullGame(id);
			state = loaded.state;
			year = loaded.year;
			context.resetContext(loaded);
			context.commandHomeWorld();
		}

		// if the game is active and we haven't submitted our turn
		// bind the navigation hotkeys
		if (state == GameState.WaitingForPlayers && !get(context.player).submittedTurn) {
			// reset key bindings
			unbindNavigationHotkeys();
			hotkeys.unbind('F9', 'root');

			bindNavigationHotkeys(get(context.game).id, page);
			hotkeys('F9', 'root', () => {
				onSubmitTurn();
			});
		} else {
			unbindNavigationHotkeys();
			hotkeys.unbind('F9', 'root');
		}
	}

	async function onSubmitTurn() {
		if (!context) return;

		// don't allow submitTurn if we're not in an active game
		// waiting to submit our turn
		const g = get(context.game);
		const p = get(context.player);
		if (g.state != GameState.WaitingForPlayers || p.submittedTurn) {
			return;
		}

		// console.time('onSubmitTurn');
		setLoadingModalText(`Submitting turn for ${year}...`);
		const hideLoadingDelay = wait(500);
		try {
			// update the UI
			await context.submitTurn();
			goto(`/games/${g.id}`);
		} catch {
			// reload the game status in the event of an error
			await context.loadStatus();
		} finally {
			// make sure this finishes
			await hideLoadingDelay;
			clearLoadingModalText();
			// console.timeEnd('onSubmitTurn');
		}
	}
</script>

{#if contextSetup}
	<GameLayout on:submit-turn={onSubmitTurn}>
		<slot>Game</slot>
	</GameLayout>
{:else if error}
	<main class="flex flex-col">
		<div class="flex-initial">
			<Menu user={$me} />
		</div>
		<ErrorPage {error} />
	</main>
{/if}

<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import ErrorPage from '$lib/components/ErrorPage.svelte';
	import Menu from '$lib/components/Menu.svelte';
	import { bindNavigationHotkeys, unbindNavigationHotkeys } from '$lib/navigationHotkeys';
	import { gameClient, playerClient } from '$lib/services/connect';
	import { FullGame } from '$lib/services/FullGame';
	import { createGameContext, gameKey, type GameContext } from '$lib/services/GameContext';
	import { clearLoadingModalText, me, setLoadingModalText } from '$lib/services/Stores';
	import { Universe } from '$lib/services/Universe';
	import { GameState } from '$lib/types/cs-proto';
	import { getGameWithPlayersFlat } from '$lib/types/Game';
	import { CommandedPlayer } from '$lib/types/Player';
	import { wait } from '$lib/wait';
	import { loadWasm } from '$lib/wasm';
	import { Code, type ConnectError } from '@connectrpc/connect';
	import hotkeys from 'hotkeys-js';
	import { onDestroy, onMount, setContext, type Snippet } from 'svelte';
	import type { Unsubscriber } from 'svelte/store';
	import { get } from 'svelte/store';
	import GameLayout from './GameLayout.svelte';
	type Props = {
		children?: Snippet;
	};

	let { children }: Props = $props();

	let id = parseInt($page.params.id);

	let context: GameContext | undefined = $state(undefined);
	let error: string | undefined = $state(undefined);
	let contextSetup = $state(false);

	let unsubscribe: Unsubscriber | undefined = $state();
	let gameState: GameState = $state(GameState.SETUP);
	let year: number = $state(2400);

	onMount(async () => {
		try {
			setLoadingModalText('Loading game...');
			const wasmResp = loadWasm();
			// on mount, load the game and setup the context used by the rest of the children
			const { game, universe, player } = await loadFullGame(BigInt(id));

			const cs = await wasmResp;
			context = await createGameContext(cs, game, player, universe);
			if (game.state === GameState.WAITING_FOR_PLAYERS) {
				context.setFullyLoaded(true);
			}

			hotkeys.setScope('root');
		} catch (e) {
			const err = e as ConnectError;
			if (err.code === Code.NotFound) {
				error = 'Game not found';
			} else {
				error = `${e}`;
			}
		} finally {
			clearLoadingModalText();
		}
	});

	onDestroy(() => {
		hotkeys.deleteScope('root');

		if (!context) return;

		unsubscribe?.();
		setContext(gameKey, undefined);
	});

	// if no context is defined, create it
	$effect(() => {
		if (context && !contextSetup) {
			contextSetup = true;

			// store the latest state/year so we can reload if the game changes
			const game = get(context.game);
			gameState = game.state;
			year = game.year;
			context.commandHomeWorld();

			// subscribe to game change events so we do a full reload if the year/state changes
			unsubscribe = context.game.subscribe(onGameChange);

			// setup the context for our child components
			setContext(gameKey, context);
		}
	});

	// load the full game with intel and universe objects
	async function loadFullGame(gameId: bigint) {
		const { game } = await gameClient.getGame({ gameId });
		if (!game) {
			throw Error('failed to load game');
		}
		const fg = Object.assign(new FullGame(), getGameWithPlayersFlat(game));

		// create empty universe/player objects as placeholders
		const u = new Universe();
		const p: CommandedPlayer = new CommandedPlayer();
		if (fg.state != GameState.SETUP) {
			const playerResp = playerClient.getPlayer({ gameId });
			const { universe } = await playerClient.getUniverse({ gameId });
			const { player } = await playerResp;

			if (!universe || !player) {
				throw Error('failed to load player and universe');
			}
			// configure the universe for the player after the player is loaded
			u.setData(player.num, universe);

			// update the player
			Object.assign(p, player);
		}

		return { game: fg, player: p, universe: u };
	}

	// every time the game updates, check if we have a new year/state change
	// and if so, reset the context
	async function onGameChange(game: FullGame) {
		if (!context) return;

		if (gameState != game.state || year != game.year) {
			const { game, universe, player } = await loadFullGame(BigInt(id));

			gameState = game.state;
			year = game.year;
			await context.resetContext(game, player, universe);
			if (game.state === GameState.WAITING_FOR_PLAYERS) {
				context.setFullyLoaded(true);
			}

			context.commandHomeWorld();
		}

		// if the game is active and we haven't submitted our turn
		// bind the navigation hotkeys
		if (gameState === GameState.WAITING_FOR_PLAYERS && !get(context.player).submittedTurn) {
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
		if (g.state !== GameState.WAITING_FOR_PLAYERS || p.submittedTurn) {
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
	<GameLayout {onSubmitTurn}>
		{#if children}{@render children()}{:else}Game{/if}
	</GameLayout>
{:else if error}
	<main class="flex flex-col">
		<div class="flex-initial">
			<Menu user={$me} />
		</div>
		<ErrorPage {error} />
	</main>
{/if}

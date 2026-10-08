<script lang="ts">
	import { afterNavigate, goto, replaceState } from '$app/navigation';
	import { page } from '$app/state';
	import ErrorPage from '#lib/components/ErrorPage.svelte';
	import Menu from '#lib/components/Menu.svelte';
	import { bindNavigationHotkeys, unbindNavigationHotkeys } from '#lib/navigationHotkeys.js';
	import { setAsPlayerNum, setReadOnly } from '#lib/services/asPlayerInterceptor.js';
	import { gameClient, playerClient } from '#lib/services/connect.js';
	import { FullGame } from '#lib/services/FullGame.js';
	import { createGameContext, gameKey, type GameContext } from '#lib/services/GameContext.js';
	import { clearLoadingModalText, me, setLoadingModalText } from '#lib/services/Stores.js';
	import { Universe } from '#lib/services/Universe.js';
	import { GameState } from '#lib/types/cs-proto.js';
	import { getGameWithPlayersFlat } from '#lib/types/Game.js';
	import { getFirstHotSeatPlayer, getUserPlayers, isHotSeat } from '#lib/types/HotSeat.js';
	import { CommandedPlayer, CommandedPlayerRace } from '#lib/types/Player.js';
	import { wait } from '#lib/wait.js';
	import { loadWasm } from '#lib/wasm.js';
	import { Code, type ConnectError } from '@connectrpc/connect';
	import hotkeys from 'hotkeys-js';
	import { onDestroy, onMount, setContext, type Snippet } from 'svelte';
	import type { Unsubscriber } from 'svelte/store';
	import { get } from 'svelte/store';
	import GameLayout from './GameLayout.svelte';

	type Props = { children?: Snippet };

	let { children }: Props = $props();
	let id = parseInt(page.params.id || '0');
	let context: GameContext | undefined = $state(undefined);
	let error: string | undefined = $state(undefined);
	let contextSetup = $state(false);

	let unsubscribe: Unsubscriber | undefined = $state();
	let gameState: GameState = $state(GameState.SETUP);
	let year: number = $state(2400);

	// the player we are playing as, from the ?asPlayer= query param, for hot seat games or admin views
	let asPlayerNum: number | undefined = parseAsPlayerNum(page.url.searchParams.get('asPlayer'));
	// incremented for each load so slow loads for a previous player are discarded
	let loadSeq = 0;
	// true when a hot seat user opens the game without choosing a player, before the context exists
	let choosingPlayer = false;
	setAsPlayerNum(asPlayerNum);

	onMount(async () => {
		try {
			setLoadingModalText('Loading game...');
			const wasmResp = loadWasm();
			// on mount, load the game and setup the context used by the rest of the children
			const { game, universe, player } = await loadFullGame(BigInt(id));

			const cs = await wasmResp;
			context = await createGameContext(cs, game, player, universe);
			context.readOnly.set(defaultReadOnly(player));
			context.choosingPlayer.set(choosingPlayer);
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
		setAsPlayerNum(undefined);
		setReadOnly(false);

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
			const unsubscribeGame = context.game.subscribe(onGameChange);
			// viewing a turn read-only changes which hotkeys are available
			const unsubscribeReadOnly = context.readOnly.subscribe(() => updateHotkeys());
			const unsubscribeChoosing = context.choosingPlayer.subscribe(() => updateHotkeys());
			unsubscribe = () => {
				unsubscribeGame();
				unsubscribeReadOnly();
				unsubscribeChoosing();
			};

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
		if (fg.state != GameState.SETUP && !asPlayerNum) {
			choosePlayer(fg);
		}

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
			p.race = new CommandedPlayerRace(p.race);
		}

		return { game: fg, player: p, universe: u };
	}

	// every time the game updates, check if we have a new year/state change
	// and if so, reset the context
	async function onGameChange(game: FullGame) {
		if (!context) return;

		if (gameState != game.state || year != game.year) {
			gameState = game.state;
			year = game.year;

			// new turns in hot seat games go back to choosing a player
			if (game.state === GameState.WAITING_FOR_PLAYERS && isHotSeat(game.players, $me.id)) {
				context.choosingPlayer.set(true);
				updateUrl(undefined);
			}
			await reload();
		}

		updateHotkeys();
	}

	// if the game is active and we haven't submitted our turn
	// bind the navigation hotkeys
	function updateHotkeys() {
		if (!context) return;
		const submittedTurn = get(context.player).submittedTurn;
		const readOnly = get(context.readOnly);
		const choosing = get(context.choosingPlayer);
		if (gameState === GameState.WAITING_FOR_PLAYERS && !choosing && (!submittedTurn || readOnly)) {
			// reset key bindings
			unbindNavigationHotkeys();
			hotkeys.unbind('F9', 'root');

			bindNavigationHotkeys(get(context.game).id);
			if (!readOnly) {
				hotkeys('F9', 'root', () => {
					onSubmitTurn();
				});
			}
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
		if (
			g.state !== GameState.WAITING_FOR_PLAYERS ||
			p.submittedTurn ||
			get(context.readOnly) ||
			get(context.choosingPlayer)
		) {
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

	// reload the game, player and universe and reset the context
	async function reload() {
		if (!context) return;
		const seq = ++loadSeq;
		const { game, universe, player } = await loadFullGame(BigInt(id));
		if (seq !== loadSeq) {
			// a newer load started (i.e. we switched players), ignore this one
			return;
		}

		gameState = game.state;
		year = game.year;
		await context.resetContext(game, player, universe);
		context.readOnly.set(defaultReadOnly(player));
		if (game.state === GameState.WAITING_FOR_PLAYERS) {
			context.setFullyLoaded(true);
		}

		context.commandHomeWorld();
		updateHotkeys();
	}

	// admins view other players' turns read-only, so submitted turns (like AI players) show the map
	function defaultReadOnly(player: CommandedPlayer): boolean {
		return player.num > 0 && $me.isAdmin() && player.userId !== $me.id;
	}

	// switch to playing (or viewing) as another player
	async function switchPlayer(num: number) {
		asPlayerNum = num;
		setAsPlayerNum(num);
		context?.choosingPlayer.set(false);
		await reload();
	}

	// pick a default player for this game if the url didn't specify one
	function choosePlayer(game: FullGame) {
		if (isHotSeat(game.players, $me.id)) {
			// load the first player, but let the user choose which player to play
			asPlayerNum = getFirstHotSeatPlayer(game.players, $me.id)?.num;
			choosingPlayer = game.state === GameState.WAITING_FOR_PLAYERS;
		} else if (getUserPlayers(game.players, $me.id).length === 0 && $me.isAdmin()) {
			// admins viewing a game they aren't in start as the first player
			asPlayerNum = game.players[0]?.num;
		}
		setAsPlayerNum(asPlayerNum);
		if (!choosingPlayer) {
			updateUrl(asPlayerNum);
		}
	}

	function parseAsPlayerNum(value: string | null): number | undefined {
		const num = parseInt(value ?? '');
		return num > 0 ? num : undefined;
	}

	// keep ?asPlayer= in the url so reloads and links keep the same player
	function updateUrl(num: number | undefined) {
		if (parseAsPlayerNum(page.url.searchParams.get('asPlayer')) !== num) {
			const url = new URL(page.url.href);
			if (num) {
				url.searchParams.set('asPlayer', `${num}`);
			} else {
				url.searchParams.delete('asPlayer');
			}
			try {
				replaceState(url, page.state);
			} catch {
				// the router isn't ready yet, afterNavigate will update the url
			}
		}
	}

	// the url is the source of truth for which player we are, i.e. the Play links in
	// hot seat games or admin View As links. Links that drop the param keep the current player.
	afterNavigate(() => {
		const urlPlayerNum = parseAsPlayerNum(page.url.searchParams.get('asPlayer'));
		const choosing = context ? get(context.choosingPlayer) : choosingPlayer;
		if (urlPlayerNum && urlPlayerNum !== asPlayerNum) {
			switchPlayer(urlPlayerNum);
		} else if (urlPlayerNum && choosing) {
			// the hot seat user chose the player we already loaded
			context?.choosingPlayer.set(false);
		} else if (!choosing) {
			updateUrl(asPlayerNum);
		}
	});
</script>

{#if contextSetup}
	<GameLayout {onSubmitTurn}>
		{#if children}
			{@render children()}
		{:else}
			Game
		{/if}
	</GameLayout>
{:else if error}
	<main class="flex flex-col">
		<div class="flex-initial"><Menu user={$me} /></div>
		<ErrorPage {error} />
	</main>
{/if}

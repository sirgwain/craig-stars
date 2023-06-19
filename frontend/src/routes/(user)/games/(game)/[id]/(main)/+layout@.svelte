<script lang="ts">
	import { page } from '$app/stores';
	import NotFound from '$lib/components/NotFound.svelte';
	import { bindNavigationHotkeys, unbindNavigationHotkeys } from '$lib/navigationHotkeys';
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
	import { getGameContext, initGameContext, updateGameContext } from '$lib/services/Contexts';
	import { FullGame } from '$lib/services/FullGame';
	import { commandedMapObject, techs } from '$lib/services/Stores';
	import { GameState } from '$lib/types/Game';
	import { onMount } from 'svelte';
	import GameMenu from '../GameMenu.svelte';
	import { Universe } from '$lib/services/Universe';
	import { Player } from '$lib/types/Player';

	let id = parseInt($page.params.id);
	let loadAttempted = false;
	let fg: FullGame | undefined = undefined;

	initGameContext();

	const { game, player, universe } = getGameContext();

	onMount(async () => {
		if (!$game || $game.id !== id) {
			// empty the context
			updateGameContext(new FullGame(), new Player(), new Universe());

			try {
				const loading = new FullGame();
				await loading.load(id);

				fg = loading;
				const gameTechs = fg.techs;
				techs.update(() => gameTechs);
			} finally {
				loadAttempted = true;
			}
		}

		if ($game?.state == GameState.WaitingForPlayers) {
			if (!$commandedMapObject || $commandedMapObject.gameId != $game.id) {
				$game.commandHomeWorld();
			}
		}

		// setup the quantityModifier
		bindQuantityModifier();

		// if we are in an active game, bind the navigation hotkeys, i.e. F4 for research, Esc to go back
		if ($game?.state == GameState.WaitingForPlayers) {
			bindNavigationHotkeys(id, page);
		}

		return () => {
			unbindQuantityModifier();
			unbindNavigationHotkeys();
		};
	});

	async function onSubmitTurn() {
		$game = await $game?.submitTurn();
		if ($game?.state == GameState.WaitingForPlayers) {
			// trigger reactivity
			$game = $game;
			$game.commandHomeWorld();
		}
	}
</script>

{#if fg}
	<main class="flex flex-col mb-20 md:mb-0">
		<div class="flex-initial">
			<GameMenu game={$game} on:submit-turn={onSubmitTurn} />
		</div>
		<!-- We want our main game view to only fill the screen (minus the toolbar) -->
		<div class="grow viewport">
			<slot>Game</slot>
		</div>
	</main>
{:else if loadAttempted}
	<NotFound title="Game not found" />
{/if}

<style>
	main {
		height: 100vh; /* Fallback for browsers that do not support Custom Properties */
		height: calc(var(--vh, 1vh) * 100);
	}

	.viewport {
		max-height: calc(100vh-4rem);
		max-height: calc((var(--vh, 1vh) * 100)-4rem);
	}
</style>

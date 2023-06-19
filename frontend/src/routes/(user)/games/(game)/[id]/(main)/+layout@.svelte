<script lang="ts">
	import { page } from '$app/stores';
	import NotFound from '$lib/components/NotFound.svelte';
	import { bindNavigationHotkeys, unbindNavigationHotkeys } from '$lib/navigationHotkeys';
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
	import { commandedMapObject, game, techs } from '$lib/services/Stores';
	import { FullGame } from '$lib/services/FullGame';
	import { GameState } from '$lib/types/Game';
	import { onMount } from 'svelte';
	import GameMenu from '../GameMenu.svelte';

	let id = parseInt($page.params.id);

	let loadAttempted = false;

	onMount(async () => {
		if (!$game || $game.id !== id) {
			game.update(() => undefined);

			try {
				const fg = new FullGame();
				await fg.load(id);

				game.update(() => fg);
				techs.update(() => fg.techs);
			} finally {
				loadAttempted = true;
			}
		}

		if ($game?.state == GameState.WaitingForPlayers) {
			if (!$commandedMapObject || $commandedMapObject.gameId != $game.id) {
				$game.universe.commandHomeWorld();
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
			$game.universe.commandHomeWorld();
		}
	}
</script>

{#if $game}
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

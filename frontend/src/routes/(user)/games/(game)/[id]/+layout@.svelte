<script lang="ts">
	import { page } from '$app/stores';
	import ErrorToast from '$lib/components/ErrorToast.svelte';
	import NotFound from '$lib/components/NotFound.svelte';
	import Tooltip from '$lib/components/game/tooltips/Tooltip.svelte';
	import { bindNavigationHotkeys, unbindNavigationHotkeys } from '$lib/navigationHotkeys';
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
	import { game } from '$lib/services/Context';
	import { FullGame } from '$lib/services/FullGame';
	import { GameState } from '$lib/types/Game';
	import { onMount } from 'svelte';
	import GameMenu from './GameMenu.svelte';

	let id = parseInt($page.params.id);

	let loadAttempted = false;

	onMount(async () => {
		const g = new FullGame();
		await g.load(id);

		game.update(() => g);

		loadAttempted = true;

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
</script>

{#if $game}
	<main class="flex flex-col mb-20 md:mb-0">
		<div class="flex-initial">
			<GameMenu game={$game} />
		</div>
		<ErrorToast />
		<!-- We want our main game view to only fill the screen (minus the toolbar) -->
		<div class="grow viewport">
			<slot>Game</slot>
		</div>
	</main>
	<Tooltip />
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

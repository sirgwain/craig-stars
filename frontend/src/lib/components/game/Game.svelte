<script lang="ts">
	import CommandPane from '$lib/components/game/command/CommandPane.svelte';
	import Viewport from '$lib/components/game/Scanner.svelte';
	import Toolbar from '$lib/components/game/Toolbar.svelte';
	import { setCommandedPlanet, setGameContext } from '$lib/services/Context';
	import { GameService } from '$lib/services/GameService';
	import { PlayerService } from '$lib/services/PlayerService';
	import type { Game } from '$lib/types/Game';
	import type { Player } from '$lib/types/Player';
	import { onMount } from 'svelte';

	export let id: number;
	let player: Player;
	let game: Game;
	let playerService: PlayerService;
	let gameService: GameService = new GameService();

	onMount(async () => {
		// load the game on mount
		({ game, player } = await gameService.loadGame(id));
		playerService = new PlayerService(player);
	});

	// all other components will use this context
	$: if (game && player) {
		setGameContext(game, player);
		const homeworld = player.planets.find((p) => p.homeworld);
		if (homeworld) {
			setCommandedPlanet(homeworld);
		} else {
			setCommandedPlanet(player.planets[0]);
		}
	}

	async function onSubmitTurn() {
		const result = await playerService.submitTurn();
		if (result !== undefined) {
			({ game, player } = result);
		}
	}
</script>

{#if player}
	<div class="flex flex-col h-full">
		<div class="flex-none">
			<Toolbar {game} on:submit-turn={onSubmitTurn} />
		</div>

		<div class="flex-1">
			<div class="flex h-full">
				<div class="flex-none border-gray-700 shadow-sm">
					<CommandPane />
				</div>
				<div class="ml-5 flex-1 h-full border-gray-700 border-2 shadow-sm">
					<Viewport />
				</div>
			</div>
		</div>
	</div>
{/if}

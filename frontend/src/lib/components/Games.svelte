<script lang="ts">
	import GamesTable from '$lib/components/GamesTable.svelte';
	import { GameService } from '$lib/services/GameService';
	import type { Game } from '$lib/types/Game';
	import { onMount } from 'svelte';

	const gameService = new GameService();

	let myGames: Game[];
	let hostedGames: Game[];
	let openGames: Game[];

	onMount(async () => {
		const sorter = (a: Game, b: Game) =>
			b.createdAt && a.createdAt ? b.createdAt.localeCompare(a.createdAt) : 0;

		myGames = await gameService.loadPlayerGames();
		myGames.sort(sorter);

		hostedGames = await gameService.loadHostedGames();
		hostedGames.sort(sorter);

		openGames = await gameService.loadOpenGames();
		openGames.sort(sorter);
	});
</script>

<h2 class="font-semibold text-xl">
	Hosted Games <a href="/host-game" class="btn justify-end">Host Game</a>
</h2>
<div class="mt-2">
	<GamesTable games={hostedGames} />
</div>

<h2 class="font-semibold text-xl">My Games</h2>
<div class="mt-2">
	<GamesTable games={myGames} />
</div>

<h2 class="font-semibold text-xl">Open Games</h2>
<div class="mt-2">
	<GamesTable games={openGames} />
</div>

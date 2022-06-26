<script lang="ts">
	import GamesTable from '$lib/components/GamesTable.svelte';

	import Games from '$lib/components/GamesTable.svelte';
	import { GameService } from '$lib/services/GameService';
	import type { Game } from '$lib/types/Game';
	import { onMount } from 'svelte';

	const gameService = new GameService();

	let myGames: Game[];
	let hostedGames: Game[];

	onMount(async () => {
		myGames = await gameService.loadPlayerGames();

		myGames.sort((a, b) =>
			b.createdAt && a.createdAt ? b.createdAt.localeCompare(a.createdAt) : 0
		);

		hostedGames = await gameService.loadHostedGames();

		hostedGames.sort((a, b) =>
			b.createdAt && a.createdAt ? b.createdAt.localeCompare(a.createdAt) : 0
		);
	});
</script>

<div class="prose">
	<h2>Hosted Games <a href="/host-game" class="btn justify-end">Host Game</a></h2>
</div>
<div class="mt-2">
	<GamesTable games={hostedGames} />
</div>

<div class="prose">
	<h2>My Games</h2>
</div>
<div class="mt-2">
	<GamesTable games={myGames} />
</div>

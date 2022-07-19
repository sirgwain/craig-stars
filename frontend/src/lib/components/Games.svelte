<script lang="ts">
	import { GameService } from '$lib/services/GameService';
	import { GameState, type Game } from '$lib/types/Game';
	import { onMount } from 'svelte';
	import { me } from '$lib/services/Context';

	const gameService = new GameService();

	let myGames: Game[];
	let hostedGames: Game[];
	let openGames: Game[];

	onMount(async () => {
		await loadGames();
	});

	const loadGames = async () => {
		const sorter = (a: Game, b: Game) =>
			b.createdAt && a.createdAt ? b.createdAt.localeCompare(a.createdAt) : 0;

		myGames = (await gameService.loadPlayerGames()).filter((g) => g.hostId != $me?.id);
		myGames.sort(sorter);

		hostedGames = await gameService.loadHostedGames();
		hostedGames.sort(sorter);

		openGames = await gameService.loadOpenGames();
		openGames.sort(sorter);
	};

	const deleteGame = async (game: Game) => {
		await gameService.deleteGame(game.id);
		await loadGames();
	};
</script>

<div class="container mx-auto grid grid-cols-12 gap-1">
	<h2 class="font-semibold text-xl col-span-full">
		Hosted Games <a href="/host-game" class="btn justify-end">Host Game</a>
	</h2>
	<div class="col-span-6 text-secondary">Name</div>
	<div class="col-span-3 text-secondary">Year</div>
	<div class="col-span-3 text-secondary" />

	{#if hostedGames}
		{#each hostedGames as game}
			<div class="col-span-6">
				<a class="text-primary text-2xl hover:text-accent w-full" href="/games/{game.id}"
					>{game.name}</a
				>
			</div>
			<div class="col-span-3 text-2xl">
				{#if game.state == GameState.WaitingForPlayers}
					{game.year}
				{:else if game.state == GameState.Setup}
					Setting Up
				{:else}
					{game.state}
				{/if}
			</div>
			<div class="col-span-3">
				<button
					on:click={() => deleteGame(game)}
					class="float-right btn btn-error btn-danger btn-sm">Delete</button
				>
			</div>
		{/each}
	{/if}

	{#if myGames?.length > 0}
		<h2 class="font-semibold text-xl col-span-full">My Games</h2>
		<div class="col-span-6 text-secondary">Name</div>
		<div class="col-span-3 text-secondary">Years</div>
		<div class="col-span-3 text-secondary">Players</div>

		{#each myGames as game}
			<div class="col-span-6">
				<a class="text-primary text-2xl hover:text-accent w-full" href="/games/{game.id}"
					>{game.name}</a
				>
			</div>
			<div class="col-span-3 text-2xl">
				{game.year}
			</div>
			<div class="col-span-3 text-2xl">
				{game.numPlayers}
			</div>
		{/each}
	{/if}

	{#if openGames}
		<h2 class="font-semibold text-xl col-span-full">Open Games</h2>
		<div class="col-span-6 text-secondary">Name</div>
		<div class="col-span-6 text-secondary">Players</div>

		{#each openGames as game}
			<div class="col-span-6">
				<a class="text-primary text-2xl hover:text-accent w-full" href="/join-game/{game.id}"
					>{game.name}</a
				>
			</div>
			<div class="col-span-3 text-2xl">
				{game.numPlayers} / {game.openPlayerSlots + game.numPlayers}
			</div>
		{/each}
	{/if}
</div>

<!-- <h2 class="font-semibold text-xl">My Games</h2>
<div class="mt-2">
	<GamesTable games={myGames} />
</div>

<h2 class="font-semibold text-xl">Open Games</h2>
<div class="mt-2">
	<GamesTable games={openGames} />
</div> -->

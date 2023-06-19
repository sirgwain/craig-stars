<script lang="ts">
	import { GameService } from '$lib/services/GameService';
	import type { Game } from '$lib/types/Game';
	import { onMount } from 'svelte';
	import { me } from '$lib/services/Context';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { XMark } from '@steeze-ui/heroicons';
	import Galaxy from '$lib/components/icons/Galaxy.svelte';
	import Processor from '$lib/components/icons/Processor.svelte';

	const gameService = new GameService();

	let myGames: Game[];
	let openGames: Game[];

	onMount(() => {
		loadGames();
	});

	const loadGames = () => {
		const sorter = (a: Game, b: Game) =>
			b.createdAt && a.createdAt ? b.createdAt.localeCompare(a.createdAt) : 0;

		gameService.loadPlayerGames().then((games) => {
			myGames = games;
			myGames.sort(sorter);
		});

		gameService.loadOpenGames().then((games) => {
			openGames = games;
			openGames.sort(sorter);
		});
	};

	const deleteGame = async (game: Game) => {
		if (confirm(`Are you sure you want to delete ${game.name}?`)) {
			await gameService.deleteGame(game.id);
			await loadGames();
		}
	};
</script>

<h2 class="font-semibold text-xl my-2">Games</h2>
<div class="flex justify-evenly">
	<a class="btn gap-2" href="/host-game">
		<Galaxy class="fill-base-content w-12 h-12" />
		Host
	</a>
	<button class="btn gap-2" href="/single-player-game">
		<Processor class="fill-base-content w-12 h-12" />
		Single Player
	</button>
</div>
<div class="mt-2 grid grid-cols-12 gap-1">
	{#if myGames?.length > 0}
		<div class="col-span-6 text-secondary">Name</div>
		<div class="col-span-2 text-secondary">Year</div>
		<div class="col-span-2 text-secondary">Players</div>
		<div class="col-span-2" />

		{#each myGames as game}
			<div class="col-span-6">
				<a class="text-primary text-2xl hover:text-accent w-full" href="/games/{game.id}"
					>{game.name}</a
				>
			</div>
			<div class="col-span-2 text-2xl">
				{game.year}
			</div>
			<div class="col-span-2 text-2xl">
				{game.numPlayers}
			</div>
			{#if game.hostId == $me?.id}
				<div class="col-span-2">
					<button
						on:click={() => deleteGame(game)}
						class="float-right btn btn-error btn-danger btn-sm"
					>
						<Icon src={XMark} size="16" class="hover:stroke-accent md:hidden" />
						<span class="hidden md:inline-block">Delete</span></button
					>
				</div>
			{:else}
				<div class="col-span-2" />
			{/if}
		{/each}
	{/if}

	{#if openGames?.length > 0}
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
				{game.numPlayers - game.openPlayerSlots} / {game.numPlayers}
			</div>
		{/each}
	{/if}
</div>

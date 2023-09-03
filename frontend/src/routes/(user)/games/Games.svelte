<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import Galaxy from '$lib/components/icons/Galaxy.svelte';
	import Processor from '$lib/components/icons/Processor.svelte';
	import { GameService } from '$lib/services/GameService';
	import { me } from '$lib/services/Stores';
	import { GameState, type Game } from '$lib/types/Game';
	import { onMount } from 'svelte';
	import ActiveGameRow from './ActiveGameRow.svelte';

	let myGames: Game[];
	let myGamesWaitingToGenerate: Game[];
	let myGamesWaitingToStart: Game[];
	let gamesWaitingToStart: Game[];
	let openGames: Game[];
	let pendingGames: Game[];
	let singlePlayerGames: Game[];
	let submittedGames: Game[];

	onMount(() => {
		loadGames();
	});

	function loadGames() {
		const sorter = (a: Game, b: Game) =>
			b.createdAt && a.createdAt ? b.createdAt.localeCompare(a.createdAt) : 0;

		GameService.loadPlayerGames().then((games) => {
			myGames = games;
			myGames.sort(sorter);
			// find all multiplayer games where we haven't submitted a turn yet
			pendingGames = games
				.filter((g) => g.state != GameState.Setup)
				.filter((g) => !g.players.find((p) => p.userId == $me.id)?.submittedTurn)
				.filter((g) => g.players.find((p) => p.userId != $me.id && !p.aiControlled));

			// find all multiplayer games where we haven't submitted a turn yet
			gamesWaitingToStart = games
				.filter((g) => g.state == GameState.Setup)
				.filter((g) => g.players.find((p) => p.userId != $me.id));

			// find all games where we've submitted our turn
			submittedGames = games.filter(
				(g) => g.players.find((p) => p.userId == $me.id)?.submittedTurn
			);

			// find all single player games
			singlePlayerGames = games.filter(
				(g) => !g.players.find((p) => p.userId != $me.id && !p.aiControlled)
			);

			myGamesWaitingToGenerate = games.filter(
				(g) => g.hostId == $me.id && g.state == GameState.Setup
			);
			pendingGames.push(...myGamesWaitingToGenerate);
		});

		GameService.loadOpenGames().then((games) => {
			openGames = games;
			openGames.sort(sorter);
		});
	}

	function ready(game: Game): boolean {
		return game.players.find((p) => p.userId == $me.id)?.ready ?? false;
	}

	const deleteGame = async (game: Game) => {
		if (confirm(`Are you sure you want to delete ${game.name}?`)) {
			await GameService.deleteGame(game.id);
			await loadGames();
		}
	};
</script>

<div class="flex justify-evenly">
	<a class="btn gap-2" href="/host-game">
		<Galaxy class="fill-current w-12 h-12" />
		Host
	</a>
	<a class="btn gap-2" href="/single-player-game">
		<Processor class="fill-current w-12 h-12" />
		Single Player
	</a>
</div>

{#if pendingGames?.length > 0}
	<ItemTitle>Games Waiting on You</ItemTitle>

	<div class="mt-2 grid grid-cols-12 gap-1">
		{#if pendingGames?.length > 0}
			<div class="col-span-6 text-secondary">Name</div>
			<div class="col-span-2 text-secondary">Year</div>
			<div class="col-span-2 text-secondary">Players</div>
			<div class="col-span-2" />

			{#each pendingGames as game}
				<ActiveGameRow {game} on:delete={() => deleteGame(game)} />
			{/each}
		{/if}
	</div>
{/if}

{#if submittedGames?.length > 0}
	<ItemTitle>Games Waiting on Others</ItemTitle>

	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-6 text-secondary">Name</div>
		<div class="col-span-2 text-secondary">Year</div>
		<div class="col-span-2 text-secondary">Players</div>
		<div class="col-span-2" />

		{#each submittedGames as game}
			<ActiveGameRow {game} on:delete={() => deleteGame(game)} />
		{/each}
	</div>
{/if}

{#if singlePlayerGames?.length > 0}
	<ItemTitle>Single Player Games</ItemTitle>

	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-6 text-secondary">Name</div>
		<div class="col-span-2 text-secondary">Year</div>
		<div class="col-span-2 text-secondary">Players</div>
		<div class="col-span-2" />
		{#each singlePlayerGames as game}
			<ActiveGameRow {game} showNumSubmitted={false} on:delete={() => deleteGame(game)} />
		{/each}
	</div>
{/if}

{#if gamesWaitingToStart?.length > 0}
	<ItemTitle>Waiting to Start</ItemTitle>
	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-6 text-secondary">Name</div>
		<div class="col-span-6 text-secondary">Players</div>

		{#each gamesWaitingToStart as game}
			<div class="col-span-6">
				<a
					class="text-primary text-2xl hover:text-accent w-full"
					href={ready(game) ? `/games/${game.id}` : `/join-game/${game.id}`}>{game.name}</a
				>
			</div>
			<div class="col-span-3 text-md">
				{(game.numPlayers ?? 0) - (game.openPlayerSlots ?? 0)} / {game.numPlayers}
			</div>
		{/each}
	</div>
{/if}

{#if openGames?.length > 0}
	<ItemTitle>New Open Games</ItemTitle>
	<div class="mt-2 grid grid-cols-12 gap-1">
		<h2 class="font-semibold text-xl col-span-full">Open Games</h2>
		<div class="col-span-6 text-secondary">Name</div>
		<div class="col-span-6 text-secondary">Players</div>

		{#each openGames as game}
			<div class="col-span-6">
				<a
					class="text-primary text-2xl hover:text-accent w-full"
					href={game.hostId == $me.id ? `/games/${game.id}` : `/join-game/${game.id}`}
					>{game.name}</a
				>
			</div>
			<div class="col-span-3 text-md">
				{game.numPlayers - game.openPlayerSlots} / {game.numPlayers}
			</div>
		{/each}
	</div>
{/if}

<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import Galaxy from '$lib/components/icons/Galaxy.svelte';
	import Processor from '$lib/components/icons/Processor.svelte';
	import { GameService } from '$lib/services/GameService';
	import { me } from '$lib/services/Stores';
	import { GameState, type Game } from '$lib/types/Game';
	import { onMount } from 'svelte';
	import ActiveGameRow from './ActiveGameRow.svelte';
	import SetupGameRow from './SetupGameRow.svelte';
	import { PlayerService } from '$lib/services/PlayerService';

	let myGames: Game[];
	let gamesWaitingToStart: Game[];
	let openGames: Game[];
	let newTurnGames: Game[];
	let singlePlayerGames: Game[];
	let submittedTurnGames: Game[];
	let archivedGames: Game[];

	onMount(() => {
		loadGames();
	});

	function loadGames() {
		const sorter = (a: Game, b: Game) =>
			b.createdAt && a.createdAt ? b.createdAt.localeCompare(a.createdAt) : 0;

		GameService.loadPlayerGames().then((games) => {
			archivedGames = games.filter(
				(g) => g.archived || g.players.find((p) => p.userId == $me.id)?.archived
			);
			myGames = games.filter(
				(g) => !(g.archived || g.players.find((p) => p.userId == $me.id)?.archived)
			);
			myGames.sort(sorter);
			// find all multiplayer games where we haven't submitted a turn yet
			newTurnGames = myGames
				.filter((g) => g.state != GameState.Setup)
				.filter((g) => !g.players.find((p) => p.userId == $me.id)?.submittedTurn)
				.filter((g) => g.players.find((p) => p.userId != $me.id && !p.aiControlled));

			// find all games where we've submitted our turn
			submittedTurnGames = myGames.filter(
				(g) => g.players.find((p) => p.userId == $me.id)?.submittedTurn
			);

			// find all single player games
			singlePlayerGames = myGames.filter(
				(g) => !g.players.find((p) => p.userId != $me.id && !p.aiControlled)
			);

			// find all multiplayer games we are part of in setup
			gamesWaitingToStart = myGames
				.filter((g) => g.state == GameState.Setup)
				.filter((g) => g.players.find((p) => p.userId != $me.id));
		});

		GameService.loadOpenGames().then((games) => {
			openGames = games;
			openGames.sort(sorter);
			openGames = openGames.filter((g) => g.hostId != $me.id);
		});
	}

	const deleteGame = async (game: Game) => {
		if (confirm(`Are you sure you want to delete ${game.name}?`)) {
			await GameService.deleteGame(game.id);
			await loadGames();
		}
	};
	const archiveGame = async (game: Game) => {
		if (confirm(`Are you sure you want to archive ${game.name}?`)) {
			await PlayerService.archiveGame(game.id);
			await loadGames();
		}
	};
	const unArchiveGame = async (game: Game) => {
		if (confirm(`Are you sure you want to archive ${game.name}?`)) {
			await PlayerService.archiveGame(game.id);
			await loadGames();
		}
	};
</script>

<div class="flex justify-evenly">
	{#if !$me.isGuest()}
		<a class="btn gap-2" href="/host-game">
			<Galaxy class="fill-current w-12 h-12" />
			Host
		</a>
	{/if}
	<a class="btn gap-2" href="/single-player-game">
		<Processor class="fill-current w-12 h-12" />
		Single Player
	</a>
</div>

{#if newTurnGames?.length > 0}
	<ItemTitle>New Turns</ItemTitle>

	<div class="mt-2 grid grid-cols-12 gap-1">
		{#if newTurnGames?.length > 0}
			<div class="col-span-5 text-secondary">Name</div>
			<div class="col-span-2 text-secondary">Year</div>
			<div class="col-span-3 text-secondary">Players</div>
			<div class="col-span-2" />

			{#each newTurnGames as game}
				<ActiveGameRow
					{game}
					on:delete={() => deleteGame(game)}
					on:archive={() => archiveGame(game)}
				/>
			{/each}
		{/if}
	</div>
{/if}

{#if singlePlayerGames?.length > 0}
	<ItemTitle>Single Player Games</ItemTitle>

	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-5 text-secondary">Name</div>
		<div class="col-span-2 text-secondary">Year</div>
		<div class="col-span-3 text-secondary">Players</div>
		<div class="col-span-2" />
		{#each singlePlayerGames as game}
			<ActiveGameRow
				{game}
				showNumSubmitted={false}
				on:delete={() => deleteGame(game)}
				on:archive={() => archiveGame(game)}
			/>
		{/each}
	</div>
{/if}

{#if submittedTurnGames?.length > 0}
	<ItemTitle>Submitted</ItemTitle>

	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-5 text-secondary">Name</div>
		<div class="col-span-2 text-secondary">Year</div>
		<div class="col-span-3 text-secondary">Players</div>
		<div class="col-span-2" />

		{#each submittedTurnGames as game}
			<ActiveGameRow
				{game}
				on:delete={() => deleteGame(game)}
				on:archive={() => archiveGame(game)}
			/>
		{/each}
	</div>
{/if}
{#if gamesWaitingToStart?.length > 0}
	<ItemTitle>Waiting to Start</ItemTitle>
	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-5 text-secondary">Name</div>
		<div class="col-span-5 text-secondary">Players</div>
		<div class="col-span-2" />
		{#each gamesWaitingToStart as game}
			<SetupGameRow {game} on:delete={() => deleteGame(game)} />
		{/each}
	</div>
{/if}

{#if openGames?.length > 0 && !$me.isGuest()}
	<ItemTitle>New Open Games</ItemTitle>
	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-5 text-secondary">Name</div>
		<div class="col-span-5 text-secondary">Players</div>
		<div class="col-span-2" />

		{#each openGames as game}
			<SetupGameRow {game} on:delete={() => deleteGame(game)} />
		{/each}
	</div>
{/if}

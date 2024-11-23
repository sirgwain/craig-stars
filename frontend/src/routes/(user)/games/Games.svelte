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

	const sorter = (a: Game, b: Game) =>
		b.createdAt && a.createdAt ? b.createdAt.localeCompare(a.createdAt) : 0;

	let games: Game[] = $state([]);
	let openGames: Game[] = $state([]);

	// all games where I am a player
	let myGames = $derived(
		games
			.filter((g) => !(g.archived || g.players.find((p) => p.userId == $me.id)?.archived))
			.sort(sorter)
	);
	// find all multiplayer games we are part of in setup
	let gamesWaitingToStart = $derived(
		myGames
			.filter((g) => g.state == GameState.Setup)
			.filter((g) => g.players.find((p) => p.userId != $me.id))
	);

	// find all multiplayer games where we haven't submitted a turn yet
	let newTurnGames = $derived(
		myGames
			.filter((g) => g.state != GameState.Setup)
			.filter((g) => !g.players.find((p) => p.userId == $me.id)?.submittedTurn)
			.filter((g) => g.players.find((p) => p.userId != $me.id && !p.aiControlled))
	);

	// find all single player games
	let singlePlayerGames = $derived(
		myGames.filter((g) => !g.players.find((p) => p.userId != $me.id && !p.aiControlled))
	);

	// find all games where we've submitted our turn
	let submittedTurnGames = $derived(
		myGames.filter((g) => g.players.find((p) => p.userId == $me.id)?.submittedTurn)
	);

	let archivedGames = $derived(
		games.filter((g) => g.archived || g.players.find((p) => p.userId == $me.id)?.archived)
	);

	onMount(async () => {
		games = await GameService.loadPlayerGames();
		openGames = (await GameService.loadOpenGames()).filter((g) => g.hostId != $me.id).sort(sorter);
	});

	function removeGame(game: Game) {
		games = games.filter((g) => g.id !== game.id);
		openGames = openGames.filter((g) => g.id !== game.id);
	}

	async function deleteGame(game: Game) {
		if (confirm(`Are you sure you want to delete ${game.name}?`)) {
			await GameService.deleteGame(game.id);
			removeGame(game);
		}
	}
	async function archiveGame(game: Game) {
		if (confirm(`Are you sure you want to archive ${game.name}?`)) {
			await PlayerService.archiveGame(game.id);
			removeGame(game);
		}
	}
	async function unArchiveGame(game: Game) {
		if (confirm(`Are you sure you want to archive ${game.name}?`)) {
			await PlayerService.archiveGame(game.id);
			removeGame(game);
		}
	}
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
			<div class="col-span-2"></div>

			{#each newTurnGames as game}
				<ActiveGameRow
					{game}
					onDelete={() => deleteGame(game)}
					onArchive={() => archiveGame(game)}
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
		<div class="col-span-2"></div>
		{#each singlePlayerGames as game}
			<ActiveGameRow
				{game}
				showNumSubmitted={false}
				onDelete={() => deleteGame(game)}
				onArchive={() => archiveGame(game)}
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
		<div class="col-span-2"></div>

		{#each submittedTurnGames as game}
			<ActiveGameRow {game} onDelete={() => deleteGame(game)} onArchive={() => archiveGame(game)} />
		{/each}
	</div>
{/if}
{#if gamesWaitingToStart?.length > 0}
	<ItemTitle>Waiting to Start</ItemTitle>
	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-5 text-secondary">Name</div>
		<div class="col-span-5 text-secondary">Players</div>
		<div class="col-span-2"></div>
		{#each gamesWaitingToStart as game}
			<SetupGameRow {game} onDelete={() => deleteGame(game)} />
		{/each}
	</div>
{/if}

{#if openGames?.length > 0 && !$me.isGuest()}
	<ItemTitle>New Open Games</ItemTitle>
	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-5 text-secondary">Name</div>
		<div class="col-span-5 text-secondary">Players</div>
		<div class="col-span-2"></div>

		{#each openGames as game}
			<SetupGameRow {game} onDelete={() => deleteGame(game)} />
		{/each}
	</div>
{/if}

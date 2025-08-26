<script lang="ts">
	import Galaxy from '$lib/components/icons/Galaxy.svelte';
	import Processor from '$lib/components/icons/Processor.svelte';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import { gameClient } from '$lib/services/connect';
	import { me } from '$lib/services/Stores';
	import { GameState, type GameWithPlayers } from '$lib/types/cs-proto';
	import { onMount } from 'svelte';
	import ActiveGameRow from './ActiveGameRow.svelte';
	import SetupGameRow from './SetupGameRow.svelte';

	const sorter = (a: GameWithPlayers, b: GameWithPlayers) =>
		b.game?.createdAt && a.game?.createdAt
			? new Date(
					Number(b.game.createdAt.seconds) * 1000 + Math.floor(b.game.createdAt.nanos / 1000000)
				)
					.toISOString()
					.localeCompare(
						new Date(
							Number(a.game.createdAt.seconds) * 1000 + Math.floor(a.game.createdAt.nanos / 1000000)
						).toISOString()
					)
			: 0;

	let games: GameWithPlayers[] = $state([]);
	let openGames: GameWithPlayers[] = $state([]);

	// all games where I am a player
	let myGames = $derived(
		games
			.filter(
				(g) => !(g.game?.archived || g.players.find((p) => p.userId == BigInt($me.id))?.archived)
			)
			.sort(sorter)
	);
	// find all multiplayer games we are part of in setup
	let gamesWaitingToStart = $derived(
		myGames
			.filter((g) => g.game?.state == GameState.SETUP)
			.filter((g) => g.players.find((p) => p.userId != BigInt($me.id)))
	);

	// find all multiplayer games where we haven't submitted a turn yet
	let newTurnGames = $derived(
		myGames
			.filter((g) => g.game?.state != GameState.SETUP)
			.filter((g) => !g.players.find((p) => p.userId == BigInt($me.id))?.submittedTurn)
			.filter((g) => g.players.find((p) => p.userId != BigInt($me.id) && !p.aiControlled))
	);

	// find all single player games
	let singlePlayerGames = $derived(
		myGames.filter((g) => !g.players.find((p) => p.userId != BigInt($me.id) && !p.aiControlled))
	);

	// find all games where we've submitted our turn
	let submittedTurnGames = $derived(
		myGames.filter((g) => g.players.find((p) => p.userId == BigInt($me.id))?.submittedTurn)
	);

	onMount(async () => {
		const playerGamesResponse = await gameClient.getGames({ open: false });
		games = playerGamesResponse.games;

		const openGamesResponse = await gameClient.getGames({ open: true });
		openGames = openGamesResponse.games
			.filter((g) => g.game?.hostId !== BigInt($me.id))
			.sort(sorter);
	});

	function removeGame(game: GameWithPlayers) {
		games = games.filter((g) => g.game?.id !== game.game?.id);
		openGames = openGames.filter((g) => g.game?.id !== game.game?.id);
	}

	async function deleteGame(game: GameWithPlayers) {
		if (game.game?.id && confirm(`Are you sure you want to delete ${game.game.name}?`)) {
			await gameClient.deleteGame({ gameId: game.game.id });
			removeGame(game);
		}
	}
	async function archiveGame(game: GameWithPlayers) {
		if (game.game?.id && confirm(`Are you sure you want to archive ${game.game.name}?`)) {
			await gameClient.archiveGame({ gameId: game.game.id });
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

{#if newTurnGames.length > 0}
	<ItemTitle>New Turns</ItemTitle>

	<div class="mt-2 grid grid-cols-12 gap-1">
		{#if newTurnGames.length > 0}
			<div class="col-span-5 text-secondary">Name</div>
			<div class="col-span-2 text-secondary">Year</div>
			<div class="col-span-3 text-secondary">Players</div>
			<div class="col-span-2"></div>

			{#each newTurnGames as game (game.game?.id)}
				<ActiveGameRow
					{game}
					onDelete={() => deleteGame(game)}
					onArchive={() => archiveGame(game)}
				/>
			{/each}
		{/if}
	</div>
{/if}

{#if singlePlayerGames.length > 0}
	<ItemTitle>Single Player Games</ItemTitle>

	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-5 text-secondary">Name</div>
		<div class="col-span-2 text-secondary">Year</div>
		<div class="col-span-3 text-secondary">Players</div>
		<div class="col-span-2"></div>
		{#each singlePlayerGames as game (game.game?.id)}
			<ActiveGameRow
				{game}
				showNumSubmitted={false}
				onDelete={() => deleteGame(game)}
				onArchive={() => archiveGame(game)}
			/>
		{/each}
	</div>
{/if}

{#if submittedTurnGames.length > 0}
	<ItemTitle>Submitted</ItemTitle>

	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-5 text-secondary">Name</div>
		<div class="col-span-2 text-secondary">Year</div>
		<div class="col-span-3 text-secondary">Players</div>
		<div class="col-span-2"></div>

		{#each submittedTurnGames as game (game.game?.id)}
			<ActiveGameRow {game} onDelete={() => deleteGame(game)} onArchive={() => archiveGame(game)} />
		{/each}
	</div>
{/if}
{#if gamesWaitingToStart.length > 0}
	<ItemTitle>Waiting to Start</ItemTitle>
	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-5 text-secondary">Name</div>
		<div class="col-span-5 text-secondary">Players</div>
		<div class="col-span-2"></div>
		{#each gamesWaitingToStart as game (game.game?.id)}
			<SetupGameRow {game} onDelete={() => deleteGame(game)} />
		{/each}
	</div>
{/if}

{#if openGames.length > 0 && !$me.isGuest()}
	<ItemTitle>New Open Games</ItemTitle>
	<div class="mt-2 grid grid-cols-12 gap-1">
		<div class="col-span-5 text-secondary">Name</div>
		<div class="col-span-5 text-secondary">Players</div>
		<div class="col-span-2"></div>

		{#each openGames as game (game.game?.id)}
			<SetupGameRow {game} onDelete={() => deleteGame(game)} />
		{/each}
	</div>
{/if}

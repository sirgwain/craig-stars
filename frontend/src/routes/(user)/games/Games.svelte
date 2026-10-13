<script lang="ts">
	import Galaxy from '#lib/components/icons/Galaxy.svelte';
	import Processor from '#lib/components/icons/Processor.svelte';
	import ItemTitle from '#lib/components/ItemTitle.svelte';
	import { gameClient } from '#lib/services/connect.js';
	import { me } from '#lib/services/Stores.js';
	import { GameState, type GameWithPlayers } from '#lib/types/cs-proto.js';
	import { onMount } from 'svelte';
	import GameListCard from './GameListCard.svelte';
	import { addError } from '#lib/services/Errors.js';
	import type { ConnectError } from '@connectrpc/connect';

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

	let loading = $state(true);
	let failed = $state(false);
	let showAllSolo = $state(false);
	let soloSearch = $state('');
	const myPlayer = (g: GameWithPlayers) => g.players.find((p) => p.userId === BigInt($me.id));
	const multiplayer = (g: GameWithPlayers) =>
		g.players.some((p) => p.userId !== BigInt($me.id) && !p.aiControlled);
	let myGames = $derived(
		games.filter((g) => !g.game?.archived && !myPlayer(g)?.archived).sort(sorter)
	);
	let gamesWaitingToStart = $derived(myGames.filter((g) => g.game?.state === GameState.SETUP));
	let activeGames = $derived(myGames.filter((g) => g.game?.state !== GameState.SETUP));
	let newTurnGames = $derived(
		activeGames.filter((g) => multiplayer(g) && !myPlayer(g)?.submittedTurn)
	);
	let submittedTurnGames = $derived(
		activeGames.filter((g) => multiplayer(g) && myPlayer(g)?.submittedTurn)
	);
	let singlePlayerGames = $derived(activeGames.filter((g) => !multiplayer(g)));
	let filteredSolo = $derived(
		singlePlayerGames.filter((g) => g.game?.name.toLowerCase().includes(soloSearch.toLowerCase()))
	);

	async function loadGames() {
		loading = true;
		failed = false;
		try {
			const [mine, open] = await Promise.all([
				gameClient.getGames({ open: false }),
				gameClient.getGames({ open: true })
			]);
			games = mine.games;
			openGames = open.games
				.filter(
					(g) =>
						g.game?.hostId !== BigInt($me.id) && !g.players.some((p) => p.userId === BigInt($me.id))
				)
				.sort(sorter);
		} catch (e) {
			failed = true;
			addError(e as ConnectError);
		} finally {
			loading = false;
		}
	}
	onMount(loadGames);

	function removeGame(game: GameWithPlayers) {
		games = games.filter((g) => g.game?.id !== game.game?.id);
		openGames = openGames.filter((g) => g.game?.id !== game.game?.id);
	}

	async function deleteGame(game: GameWithPlayers) {
		if (game.game?.id && confirm(`Are you sure you want to delete ${game.game.name}?`)) {
			try {
				await gameClient.deleteGame({ gameId: game.game.id });
				removeGame(game);
			} catch (e) {
				addError(e as ConnectError);
			}
		}
	}
	async function archiveGame(game: GameWithPlayers) {
		if (game.game?.id && confirm(`Are you sure you want to archive ${game.game.name}?`)) {
			try {
				await gameClient.archiveGame({ gameId: game.game.id });
				removeGame(game);
			} catch (e) {
				addError(e as ConnectError);
			}
		}
	}
</script>

<div class="py-4 space-y-6">
	<header class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<div>
			<h1 class="text-3xl font-bold">Games</h1>
			<p class="text-base-content/70 mt-1">
				{loading
					? 'Loading your galaxies…'
					: `${newTurnGames.length} multiplayer ${newTurnGames.length === 1 ? 'turn' : 'turns'} ready to play`}
			</p>
		</div>
		<div class="flex flex-wrap gap-2">
			{#if !$me.isGuest()}<a class="btn btn-primary" href="/host-game"
					><Galaxy class="fill-current size-6" />Host game</a
				>{/if}
			<a class="btn btn-primary" href="/single-player-game"
				><Processor class="fill-current size-6" />Single Player</a
			>
		</div>
	</header>
	{#if loading}<div class="flex justify-center py-10">
			<span class="loading loading-spinner loading-lg" aria-label="Loading games"></span>
		</div>
	{:else if failed}<div class="alert alert-error">
			<span>Unable to load games.</span><button class="btn btn-sm" onclick={loadGames}>Retry</button
			>
		</div>
	{:else}
		{#if newTurnGames.length}<section aria-label="New turns">
				<ItemTitle>New Turns ({newTurnGames.length})</ItemTitle>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-3 mt-3">
					{#each newTurnGames as game (game.game?.id)}<GameListCard
							{game}
							attention
							onDelete={() => deleteGame(game)}
							onArchive={() => archiveGame(game)}
						/>{/each}
				</div>
			</section>{/if}
		{#if submittedTurnGames.length}<section aria-label="Submitted turns">
				<ItemTitle>Submitted ({submittedTurnGames.length})</ItemTitle>
				<div class="space-y-2 mt-3">
					{#each submittedTurnGames as game (game.game?.id)}<GameListCard
							{game}
							compact
							onDelete={() => deleteGame(game)}
							onArchive={() => archiveGame(game)}
						/>{/each}
				</div>
			</section>{/if}
		{#if singlePlayerGames.length}<section aria-label="Single player games">
				<ItemTitle>Single Player Games ({singlePlayerGames.length})</ItemTitle>
				{#if singlePlayerGames.length > 5}<input
						class="input input-sm w-full mt-3"
						aria-label="Search single player games"
						placeholder="Search single player games…"
						bind:value={soloSearch}
					/>{/if}
				<div class="space-y-2 mt-3">
					{#each showAllSolo || soloSearch ? filteredSolo : filteredSolo.slice(0, 5) as game (game.game?.id)}<GameListCard
							{game}
							compact
							solo
							onDelete={() => deleteGame(game)}
							onArchive={() => archiveGame(game)}
						/>{:else}<p class="text-base-content/70">No matching games.</p>{/each}
				</div>
				{#if singlePlayerGames.length > 5 && !soloSearch}<button
						class="btn btn-ghost mt-2"
						onclick={() => (showAllSolo = !showAllSolo)}
						>{showAllSolo ? 'Show fewer' : `Show all ${singlePlayerGames.length} games`}</button
					>{/if}
			</section>{/if}
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			{#if gamesWaitingToStart.length}<section>
					<ItemTitle>Waiting to Start ({gamesWaitingToStart.length})</ItemTitle>
					<div class="space-y-2 mt-3">
						{#each gamesWaitingToStart as game (game.game?.id)}<GameListCard
								{game}
								compact
								setup
								onDelete={() => deleteGame(game)}
							/>{/each}
					</div>
				</section>{/if}
			{#if openGames.length && !$me.isGuest()}<section>
					<ItemTitle>New Open Games ({openGames.length})</ItemTitle>
					<div class="space-y-2 mt-3">
						{#each openGames as game (game.game?.id)}<GameListCard
								{game}
								compact
								setup
								open
							/>{/each}
					</div>
				</section>{/if}
		</div>
		{#if !myGames.length}<div class="card bg-base-200 border border-base-300">
				<div class="card-body items-center text-center py-10">
					<h2 class="text-2xl font-semibold">Your first galaxy awaits</h2>
					<p>Start a single player game or join an open lobby.</p>
					<a class="btn btn-accent mt-3" href="/single-player-game">Start single player</a>
				</div>
			</div>{/if}
	{/if}
	<a class="btn btn-ghost btn-sm" href="/games/archived">Archived games →</a>
</div>

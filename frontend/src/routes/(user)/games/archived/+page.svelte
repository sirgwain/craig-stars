<script lang="ts">
	import Unarchive from '$lib/components/icons/Unarchive.svelte';
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { defaultSortBy, type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { gameClient } from '$lib/services/connect';
	import { addError } from '$lib/services/Errors';
	import { me } from '$lib/services/Stores';
	import { Size, type Game } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import { getGameWithPlayersFlat, type GameWithPlayersFlat } from '$lib/types/Game';
	import { timestampToString } from '$lib/types/Timestamp';
	import { XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';

	type TableGame = GameWithPlayersFlat & { action?: never };
	const columns: TableColumn<TableGame>[] = [
		{
			key: 'id',
			title: 'Num'
		},
		{
			key: 'name',
			title: 'Name'
		},
		{
			key: 'hostId',
			title: 'Host'
		},
		{
			key: 'createdAt',
			title: 'Created'
		},
		{
			key: 'updatedAt',
			title: 'Updated'
		},
		{
			key: 'year',
			title: 'Year'
		},
		{
			key: 'size',
			title: 'Size'
		},
		{
			key: 'players',
			title: 'Players',
			sortBy: (a, b) => a.players.length - b.players.length
		},
		{
			key: 'action',
			title: '',
			sortable: false
		}
	];

	// filterable games
	let games: GameWithPlayersFlat[] = $state([]);
	let search = $state('');
	let sortKey: keyof TableGame = $state('updatedAt');
	let filteredGames = $derived(
		games
			.filter((i) => i.name.toLowerCase().indexOf(search.toLowerCase()) != -1)
			.map<TableGame>((r) => r as TableGame)
			.sort((a, b) => defaultSortBy(a, b, sortKey, descending))
	);
	let descending = $state(true);

	async function archiveGame(game: Game) {
		if (game.id && confirm(`Are you sure you want to archive ${game.name}?`)) {
			await gameClient.unarchiveGame({ gameId: game.id });
			games = games.filter((g) => g.id !== game.id);
		}
	}

	async function deleteGame(game: Game) {
		if (
			game.id &&
			confirm(`Are you sure you want to delete ${game.name}? This operation cannot be undone.`)
		) {
			await gameClient.deleteGame({ gameId: game.id });
			games = games.filter((g) => g.id !== game.id);
		}
	}

	onMount(async () => {
		try {
			const { games: playerGames } = await gameClient.getGames({});
			games = playerGames
				.map((gwp) => getGameWithPlayersFlat(gwp))
				.filter((g) => g.archived || g.players.find((p) => p.userId == $me.id)?.archived);
		} catch (err) {
			addError(`${err}`);
		}
	});
</script>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<TableSearchInput bind:value={search} />
	</div>
	<Table
		{columns}
		rows={filteredGames}
		classes={{
			table: 'table table-compact table-auto w-full',
			td: 'first:table-cell [&:nth-child(2)]:table-cell [&:nth-child(3)]:table-cell hidden sm:table-cell',
			th: 'first:table-cell [&:nth-child(2)]:table-cell [&:nth-child(3)]:table-cell hidden sm:table-cell'
		}}
	>
		{#snippet head({ isSorted, sortDescending, column })}
			<span>
				<SortableTableHeader
					{column}
					isSorted={isSorted || sortKey === column.key}
					sortDescending={sortDescending || (sortKey === column.key && descending)}
					onSorted={(col, sortDescending) => {
						sortKey = col.key;
						descending = sortDescending;
					}}
				/>
			</span>
		{/snippet}

		{#snippet cell({ column, row, cell })}
			<span>
				{#if column.key == 'name'}
					<a class="cs-link text-xl" href="/games/{row.id}">{cell}</a>
				{:else if column.key == 'createdAt'}
					{timestampToString(row.createdAt)}
				{:else if column.key == 'updatedAt'}
					{timestampToString(row.updatedAt)}
				{:else if column.key == 'hostId'}
					{row.players.find((p) => p.userId === row.hostId)?.name}
				{:else if column.key == 'size'}
					{enumToString(Size, row.size)}
				{:else if column.key == 'players'}
					{row.players.length}
				{:else if column.key == 'action'}
					{#if row.hostId == $me.id}
						<div class="col-span-2 flex justify-center join">
							<button
								onclick={() => archiveGame(row)}
								class="btn btn-info btn-sm rounded-l-md"
								title="Unarchive Game"
							>
								<Unarchive class="hover:stroke-accent w-4 h-4 stroke-base-content fill-none" />
							</button>
							<button
								onclick={() => deleteGame(row)}
								class="btn btn-error btn-sm border-l-secondary rounded-r-md"
								title="Delete Game"
							>
								<Icon src={XMark} size="16" class="hover:stroke-accent" />
							</button>
						</div>
					{:else if row.archived}
						<span class="text-warning">Archived by host</span>
					{:else}
						<div class="col-span-2 flex justify-center">
							<button
								onclick={() => archiveGame(row)}
								class="btn btn-info btn-sm rounded-md"
								title="Unarchive Game"
							>
								<Unarchive class="hover:stroke-accent w-4 h-4 stroke-base-content fill-none" />
							</button>
						</div>
					{/if}
				{:else}
					{cell}
				{/if}
			</span>
		{/snippet}
	</Table>
</div>

<script lang="ts">
	import Unarchive from '$lib/components/icons/Unarchive.svelte';
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { addError } from '$lib/services/Errors';
	import { GameService } from '$lib/services/GameService';
	import { PlayerService } from '$lib/services/PlayerService';
	import { me } from '$lib/services/Stores';
	import type { Game } from '$lib/types/Game';
	import { XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { format, parseJSON } from 'date-fns';
	import { reverse, sortBy } from 'lodash-es';
	import { onMount } from 'svelte';

	const columns: TableColumn<Game>[] = [
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
	let games: Game[];
	let filteredGames: Game[] = [];
	let search = '';
	let sortKey = 'updatedAt';
	let descending = true;

	$: filteredGames = games;

	$: filteredGames = sortBy(
		games?.filter((i) => i.name.toLowerCase().indexOf(search.toLowerCase()) != -1),
		sortKey
	);

	$: descending && (filteredGames = reverse(filteredGames));

	async function archiveGame(game: Game) {
		if (confirm(`Are you sure you want to unarchive ${game.name}?`)) {
			await PlayerService.unArchiveGame(game.id);
			games = games.filter((g) => g.id !== game.id);
		}
	}

	async function deleteGame(game: Game) {
		if (confirm(`Are you sure you want to delete ${game.name}? This operation cannot be undone.`)) {
			await GameService.deleteGame(game.id);
			games = games.filter((g) => g.id !== game.id);
		}
	}

	onMount(async () => {
		try {
			games = (await GameService.loadPlayerGames()).filter(
				(g) => g.archived || g.players.find((p) => p.userId == $me.id)?.archived
			);
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
		<span slot="head" let:isSorted let:sortDescending let:column>
			<SortableTableHeader
				{column}
				isSorted={isSorted || sortKey === column.key}
				sortDescending={sortDescending || (sortKey === column.key && descending)}
				on:sorted={(e) => {
					sortKey = column.key;
					descending = e.detail.sortDescending;
				}}
			/>
		</span>

		<span slot="cell" let:column let:row let:cell>
			{#if column.key == 'name'}
				<a class="cs-link text-xl" href="/games/{row.id}">{cell}</a>
			{:else if column.key == 'createdAt'}
				{format(parseJSON(cell), 'E, MMM do yyyy hh:mm aaa')}
			{:else if column.key == 'updatedAt'}
				{format(parseJSON(cell), 'E, MMM do yyyy hh:mm aaa')}
			{:else if column.key == 'hostId'}
				{row.players.find((p) => p.userId === row.hostId)?.name}
			{:else if column.key == 'players'}
				{row.numPlayers}
			{:else if column.key == 'action'}
				{#if row.hostId == $me.id}
					<div class="col-span-2 flex justify-center join">
						<button
							on:click={() => archiveGame(row)}
							class="btn btn-info btn-sm rounded-l-md"
							title="Unarchive Game"
						>
							<Unarchive class="hover:stroke-accent w-4 h-4 stroke-base-content fill-none" />
						</button>
						<button
							on:click={() => deleteGame(row)}
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
							on:click={() => archiveGame(row)}
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
	</Table>
</div>

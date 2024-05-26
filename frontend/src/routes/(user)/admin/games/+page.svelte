<script lang="ts">
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { AdminService } from '$lib/services/AdminService';
	import type { Game } from '$lib/types/Game';
	import type { User } from '$lib/types/User';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';
	import { format, parseJSON } from 'date-fns';
	import { reverse, sortBy } from 'lodash-es';
	import { onMount } from 'svelte';

	const columns: TableColumn<Game>[] = [
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
			title: 'Players'
		}
	];

	// filterable games
	let games: Game[];
	let usersById: Map<number, User> = new Map();
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

	onMount(async () => {
		try {
			const users = await AdminService.loadUsers();
			users.forEach((u) => {
				usersById.set(u.id, u);
			});
			games = await AdminService.loadGames();
		} catch (err) {
			// TODO: show error
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
			td: 'first:table-cell nth-child(2):table-cell hidden sm:table-cell',
			th: 'first:table-cell nth-child(2):table-cell hidden sm:table-cell'
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
				{usersById.get(cell)?.username ?? 'unknown'}
			{:else if column.key == 'players'}
				{row.numPlayers}
			{:else}
				{cell}
			{/if}
		</span>
	</Table>
</div>

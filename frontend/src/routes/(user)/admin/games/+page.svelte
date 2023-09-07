<script lang="ts">
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import { AdminService } from '$lib/services/AdminService';
	import type { Game } from '$lib/types/Game';
	import type { User } from '$lib/types/User';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';
	import { format, parseJSON } from 'date-fns';
	import { onMount } from 'svelte';

	const columns: SvelteTableColumn[] = [
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
			key: 'players',
			title: 'Players'
		}
	];

	// filterable games
	let games: Game[];
	let usersById: Map<number, User> = new Map();
	let filteredGames: Game[] = [];
	let search = '';

	$: filteredGames = games;

	$: filteredGames = games?.filter((i) => i.name.toLowerCase().indexOf(search.toLowerCase()) != -1);

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
	<SvelteTable
		{columns}
		rows={filteredGames}
		classes={{
			table: 'table table-compact table-auto w-full',
			td: 'first:table-cell nth-child(2):table-cell hidden sm:table-cell',
			th: 'first:table-cell nth-child(2):table-cell hidden sm:table-cell'
		}}
		let:column
		let:cell
		let:row
	>
		<span slot="head" let:isSorted let:sortDescending>
			<SortableTableHeader {column} {isSorted} {sortDescending} />
		</span>

		<span slot="cell">
			{#if column.key == 'name'}
				<a class="cs-link text-2xl" href="/games/{row.id}">{cell}</a>
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
	</SvelteTable>
</div>

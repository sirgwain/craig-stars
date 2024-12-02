<script lang="ts">
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { defaultSortBy, type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { AdminService } from '$lib/services/AdminService';
	import { addError, CSError } from '$lib/services/Errors';
	import type { Game } from '$lib/types/Game';
	import type { User } from '$lib/types/User';
	import { format, parseJSON } from 'date-fns';
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
		}
	];

	// filterable games
	let games: Game[] = $state([]);
	let usersById: Map<number, User> = $state(new Map<number, User>());
	let sortKey = $state(localStorage.getItem('allGamesSortKey') ?? 'updatedAt') as keyof Game;
	let sortDescending: boolean = $state(
		(localStorage.getItem('allGamesSortDescending') ?? 'true') === 'true'
	);
	let filteredGames: Game[] = $derived(
		games
			?.filter((i) => i.name.toLowerCase().indexOf(search.toLowerCase()) != -1)
			.sort((a, b) => defaultSortBy(a, b, sortKey, sortDescending))
	);
	let search = $state('');

	function onSorted(column: TableColumn<Game>, descending: boolean) {
		sortDescending = descending;
		sortKey = column.key;

		localStorage.setItem('allGamesSortKey', sortKey);
		localStorage.setItem('allGamesSortDescending', `${sortDescending}`);
	}

	onMount(async () => {
		try {
			const users = await AdminService.loadUsers();
			usersById = new Map(users.map((u) => [u.id, u]));
			games = await AdminService.loadGames();
		} catch (e) {
			addError(e as CSError);
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
		externalSortAndFilter={true}
		classes={{
			table: 'table table-compact table-auto w-full',
			td: 'first:table-cell [&:nth-child(2)]:table-cell [&:nth-child(3)]:table-cell hidden sm:table-cell',
			th: 'first:table-cell [&:nth-child(2)]:table-cell [&:nth-child(3)]:table-cell hidden sm:table-cell'
		}}
	>
		{#snippet head({ column })}
			<span>
				<SortableTableHeader
					{column}
					isSorted={sortKey === column.key}
					sortDescending={sortDescending || (sortKey === column.key && sortDescending)}
					{onSorted}
				/>
			</span>
		{/snippet}

		{#snippet cell({ column, row, cell })}
			<span>
				{#if column.key == 'name'}
					<a class="cs-link text-xl" href="/games/{row.id}">{cell}</a>
				{:else if column.key == 'createdAt'}
					{format(parseJSON(cell), 'E, MMM do yyyy hh:mm aaa')}
				{:else if column.key == 'updatedAt'}
					{format(parseJSON(cell), 'E, MMM do yyyy hh:mm aaa')}
				{:else if column.key == 'hostId'}
					{usersById.get(cell)?.username ?? 'unknown'}
				{:else if column.key == 'players'}
					{row.players.length}
				{:else}
					{cell}
				{/if}
			</span>
		{/snippet}
	</Table>
</div>

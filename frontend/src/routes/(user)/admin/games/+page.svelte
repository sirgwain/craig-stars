<script lang="ts">
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { defaultSortBy, type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { adminClient } from '$lib/services/connect';
	import { addError } from '$lib/services/Errors';
	import { Size, type User } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import { getGameWithPlayersFlat, type GameWithPlayersFlat } from '$lib/types/Game';
	import { compare, timestampToString } from '$lib/types/Timestamp';
	import type { ConnectError } from '@connectrpc/connect';
	import { onMount } from 'svelte';

	const columns: TableColumn<GameWithPlayersFlat>[] = [
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
			title: 'Created',
			sortBy: (a, b) => compare(a.createdAt, b.createdAt)
		},
		{
			key: 'updatedAt',
			title: 'Updated',
			sortBy: (a, b) => compare(a.updatedAt, b.updatedAt)
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
	let games: GameWithPlayersFlat[] = $state([]);
	let usersById: Map<bigint, User> = $state(new Map<bigint, User>());
	let sortKey = $state(
		localStorage.getItem('allGamesSortKey') ?? 'updatedAt'
	) as keyof GameWithPlayersFlat;
	let sortDescending: boolean = $state(
		(localStorage.getItem('allGamesSortDescending') ?? 'true') === 'true'
	);
	let filteredGames: GameWithPlayersFlat[] = $derived(
		games
			.filter((i) => i.name.toLowerCase().indexOf(search.toLowerCase()) != -1)
			.sort((a, b) => defaultSortBy(a, b, sortKey, sortDescending, columns))
	);
	let search = $state('');

	function onSorted(column: TableColumn<GameWithPlayersFlat>, descending: boolean) {
		sortDescending = descending;
		sortKey = column.key;

		localStorage.setItem('allGamesSortKey', sortKey);
		localStorage.setItem('allGamesSortDescending', `${sortDescending}`);
	}

	onMount(async () => {
		try {
			const { users } = await adminClient.getUsers({});
			usersById = new Map(users.map((u) => [u.id, u]));
			const resp = await adminClient.getAllGames({});
			games = resp.games.map((g) => getGameWithPlayersFlat(g));
		} catch (e) {
			addError(e as ConnectError);
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
					{timestampToString(row.createdAt)}
				{:else if column.key == 'updatedAt'}
					{timestampToString(row.updatedAt)}
				{:else if column.key == 'hostId'}
					{usersById.get(row.hostId)?.username ?? 'unknown'}
				{:else if column.key == 'size'}
					{enumToString(Size, row.size)}
				{:else if column.key == 'players'}
					{row.players.length}
				{:else}
					{cell}
				{/if}
			</span>
		{/snippet}
	</Table>
</div>

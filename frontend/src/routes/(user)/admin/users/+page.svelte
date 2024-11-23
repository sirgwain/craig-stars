<script lang="ts">
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { defaultSortBy, type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { AdminService } from '$lib/services/AdminService';
	import { addError, CSError } from '$lib/services/Errors';
	import type { User } from '$lib/types/User';
	import { format, parseJSON } from 'date-fns';
	import { onMount } from 'svelte';

	type UserWithNum = User & { num: number };

	const columns: TableColumn<UserWithNum>[] = [
		{
			key: 'num',
			title: 'Num'
		},
		{
			key: 'username',
			title: 'Username'
		},
		{
			key: 'lastLogin',
			title: 'Last Login'
		},
		{
			key: 'role',
			title: 'Role'
		},
		{
			key: 'createdAt',
			title: 'Created'
		}
	];

	// filterable users
	let users: User[] = $state([]);
	let search = $state('');
	let sortKey: keyof UserWithNum = $state(localStorage.getItem('usersSortKey') ?? 'num') as keyof User;
	let sortDescending: boolean = $state(localStorage.getItem('usersSortDescending') === 'true');

	let filteredUsers: UserWithNum[] = $derived(
		users
			.map((u, i) => Object.assign(u, { num: i + 1 }))
			.sort((a, b) => defaultSortBy(a, b, sortKey, sortDescending))
			.filter((i) => i.username.toLowerCase().indexOf(search.toLowerCase()) != -1)
	);

	function onSorted(column: TableColumn<UserWithNum>, descending: boolean) {
		sortDescending = descending;
		sortKey = column.key;

		localStorage.setItem('usersSortKey', sortKey);
		localStorage.setItem('usersSortDescending', `${sortDescending}`);
	}

	onMount(async () => {
		try {
			users = await AdminService.loadUsers();
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
		rows={filteredUsers}
		externalSortAndFilter={true}
		classes={{
			table: 'table table-compact table-auto w-full',
			td: 'first:table-cell [&:nth-child(2)]:table-cell hidden sm:table-cell',
			th: 'first:table-cell [&:nth-child(2)]:table-cell hidden sm:table-cell'
		}}
	>
		{#snippet head({ column })}
			<span>
				<SortableTableHeader
					{column}
					isSorted={sortKey === column.key}
					{sortDescending}
					{onSorted}
				/>
			</span>
		{/snippet}

		{#snippet cell({ column, row, cell })}
			<span>
				{#if column.key == 'username'}
					{cell}
					{#if row.isGuest()}<a
							href={`/admin/users/convert-guest/${row.id}`}
							class="btn btn-ghost btn-outline btn-sm">Convert Guest</a
						>{/if}
				{:else if column.key == 'createdAt'}
					{format(parseJSON(cell), 'E, MMM do yyyy hh:mm aaa')}
				{:else if column.key == 'lastLogin' && cell}
					{format(parseJSON(cell), 'E, MMM do yyyy hh:mm aaa')}
				{:else}
					{cell}
				{/if}
			</span>
		{/snippet}
	</Table>
</div>

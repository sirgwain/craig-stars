<script lang="ts">
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { AdminService } from '$lib/services/AdminService';
	import type { User } from '$lib/types/User';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';
	import { format, parseJSON } from 'date-fns';
	import { onMount } from 'svelte';

	type UserWithNum = User & { num: number };

	const columns: TableColumn<User>[] = [
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
	let users: User[];
	let filteredUsers: UserWithNum[] = [];
	let search = '';

	$: filteredUsers = users?.map((u, i) => Object.assign(u, { num: i + 1 }));

	$: filteredUsers = users
		?.map((u, i) => Object.assign(u, { num: i + 1 }))
		.filter((i) => i.username.toLowerCase().indexOf(search.toLowerCase()) != -1);

	onMount(async () => {
		try {
			users = await AdminService.loadUsers();
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
		rows={filteredUsers}
		classes={{
			table: 'table table-compact table-auto w-full',
			td: 'first:table-cell [&:nth-child(2)]:table-cell hidden sm:table-cell',
			th: 'first:table-cell [&:nth-child(2)]:table-cell hidden sm:table-cell'
		}}
	>
		<span slot="head" let:isSorted let:sortDescending let:column>
			<SortableTableHeader {column} {isSorted} {sortDescending} />
		</span>

		<span slot="cell" let:column let:row let:cell>
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
	</Table>
</div>

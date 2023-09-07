<script lang="ts">
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import { AdminService } from '$lib/services/AdminService';
	import type { User } from '$lib/types/User';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';
	import { format, parseJSON } from 'date-fns';
	import { onMount } from 'svelte';

	const columns: SvelteTableColumn[] = [
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
	let filteredUsers: User[] = [];
	let search = '';

	$: filteredUsers = users;

	$: filteredUsers = users?.filter(
		(i) => i.username.toLowerCase().indexOf(search.toLowerCase()) != -1
	);

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
	<SvelteTable
		{columns}
		rows={filteredUsers}
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
			{#if column.key == 'createdAt'}
				{format(parseJSON(cell), 'E, MMM do yyyy hh:mm aaa')}
			{:else if column.key == 'lastLogin' && cell}
				{format(parseJSON(cell), 'E, MMM do yyyy hh:mm aaa')}
			{:else}
				{cell}
			{/if}
		</span>
	</SvelteTable>
</div>

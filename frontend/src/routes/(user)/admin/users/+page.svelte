<script lang="ts">
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { defaultSortBy, type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { adminClient } from '$lib/services/connect';
	import { addError } from '$lib/services/Errors';
	import { UserRole, UserSchema, type User } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import { timestampToString } from '$lib/types/Timestamp';
	import { clone } from '@bufbuild/protobuf';
	import type { ConnectError } from '@connectrpc/connect';
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
	let sortKey: keyof UserWithNum = $state(
		localStorage.getItem('usersSortKey') ?? 'num'
	) as keyof UserWithNum;
	let sortDescending: boolean = $state(localStorage.getItem('usersSortDescending') === 'true');

	let filteredUsers: UserWithNum[] = $derived(
		users
			.map((u, i) => Object.assign(clone(UserSchema, u), { num: i + 1 }))
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
			const resp = await adminClient.getUsers({});
			users = resp.users;
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
					{#if row.role === UserRole.GUEST}<a
							href={`/admin/users/convert-guest/${row.id}`}
							class="btn btn-ghost btn-outline btn-sm">Convert Guest</a
						>{/if}
				{:else if column.key == 'role'}
					{enumToString(UserRole, row.role)}
				{:else if column.key == 'createdAt'}
					{timestampToString(row.createdAt)}
				{:else if column.key == 'lastLogin' && cell}
					{timestampToString(row.updatedAt)}
				{:else}
					{cell}
				{/if}
			</span>
		{/snippet}
	</Table>
</div>

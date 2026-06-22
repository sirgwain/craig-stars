<script lang="ts">
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table from '$lib/components/table/Table.svelte';
	import type { TableColumn } from '$lib/components/table/Table';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { raceClient } from '$lib/services/connect';
	import { addError } from '$lib/services/Errors';
	import { type Race } from '$lib/types/cs-proto';
	import { getLabelForPRT } from '$lib/types/Race';
	import { timestampToString } from '$lib/types/Timestamp';
	import type { ConnectError } from '@connectrpc/connect';
	import { XCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';

	type TableRace = Race & { action?: never };
	const columns: TableColumn<TableRace>[] = [
		{
			key: 'pluralName',
			title: 'Race'
		},
		{
			key: 'prt',
			title: 'PRT'
		},
		{
			key: 'createdAt',
			title: 'Created'
		},
		{
			key: 'action',
			title: '',
			sortable: false
		}
	];

	// filterable races
	let races: Race[] = $state([]);
	let search = $state('');
	let filteredRaces: Race[] = $derived(
		races.filter((i) => i.name.toLowerCase().indexOf(search.toLowerCase()) != -1)
	);

	async function removeItem(item: Race) {
		if (item.id && confirm(`Are you sure you want to delete ${item.name}`)) {
			await raceClient.deleteRace({ raceId: item.id });
			races = races.filter((b) => b.id != item.id);
		}
	}

	onMount(async () => {
		try {
			const resp = await raceClient.getRaces({});
			races = resp.races;
		} catch (e) {
			addError(e as ConnectError);
		}
	});
</script>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<TableSearchInput bind:value={search} />
		<a href="/races/new" class="btn btn-secondary">Create</a>
	</div>
	<Table
		{columns}
		rows={filteredRaces}
		classes={{
			table: 'table table-compact table-auto w-full',
			td: 'first:table-cell [&:nth-child(4)]:table-cell hidden sm:table-cell',
			th: 'first:table-cell [&:nth-child(4)]:table-cell hidden sm:table-cell'
		}}
	>
		{#snippet head({ isSorted, sortDescending, column })}
			<span>
				<SortableTableHeader {column} {isSorted} {sortDescending} />
			</span>
		{/snippet}

		{#snippet cell({ column, row, cell })}
			<span>
				{#if column.key == 'pluralName'}
					<a class="cs-link text-2xl" href="/races/{row.id}">{cell}</a>
				{:else if column.key == 'prt'}
					{getLabelForPRT(row.prt)}
				{:else if column.key == 'createdAt'}
					{#if row.createdAt}
						{timestampToString(row.createdAt)}
					{/if}
				{:else if column.key == 'action'}
					<button
						onclick={() => removeItem(row)}
						type="button"
						data-type="delete-button"
						data-id={row.id}
						><Icon
							class="h-10 align-middle hover:stroke-primary-focus stroke-error"
							src={XCircle}
							size="24"
						/></button
					>
				{:else}
					{cell}
				{/if}
			</span>
		{/snippet}
	</Table>
</div>

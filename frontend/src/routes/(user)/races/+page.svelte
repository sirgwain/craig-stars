<script lang="ts">
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import { RaceService } from '$lib/services/RaceService';
	import type { Race } from '$lib/types/Race';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';
	import { XCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { format, parseJSON } from 'date-fns';
	import { onMount } from 'svelte';

	const columns: SvelteTableColumn[] = [
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
		}
	];

	// filterable races
	let races: Race[];
	let filteredRaces: Race[] = [];
	let search = '';

	$: filteredRaces = races;

	$: filteredRaces = races?.filter((i) => i.name.toLowerCase().indexOf(search.toLowerCase()) != -1);

	export const removeItem = async (item: Race) => {
		if (item.id && confirm(`Are you sure you want to delete ${item.name}`)) {
			await RaceService.delete(item);
			races = races.filter((b) => b.id != item.id);
		}
	};

	onMount(async () => {
		try {
			races = await RaceService.load();
		} catch (err) {
			// TODO: show error
		}
	});
</script>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<TableSearchInput bind:value={search} />
		<a href="/races/new" class="btn btn-secondary">Create</a>
	</div>
	<SvelteTable
		{columns}
		rows={filteredRaces}
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
			{#if column.key == 'pluralName'}
				<a class="cs-link text-2xl" href="/races/{row.id}">{cell}</a>
			{:else if column.key == 'createdAt'}
				{format(parseJSON(cell), 'E, MMM do yyyy hh:mm aaa')}
			{:else if column.key == 'action'}
				<button on:click|preventDefault={() => removeItem(row)} type="button"
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
	</SvelteTable>
</div>

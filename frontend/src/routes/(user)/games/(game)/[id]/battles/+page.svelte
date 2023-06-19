<script lang="ts">
	import { page } from '$app/stores';
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { game } from '$lib/services/Context';
	import type { BattleRecord } from '$lib/types/Battle';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';

	let id = parseInt($page.params.id);

	// filterable battles
	let filteredBattles: BattleRecord[] = [];
	let search = '';

	$: filteredBattles = $game?.player.battles ?? [];
	$: filteredBattles =
		$game?.player.battles.filter(
			(i) => $game?.player.getBattleLocation(i).toLowerCase().indexOf(search.toLowerCase()) != -1
		) ?? [];

	const columns: SvelteTableColumn<BattleRecord>[] = [
		{
			key: 'location',
			title: 'Location'
		}
	];
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li>Battles</li>
	</svelte:fragment>
</Breadcrumb>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<TableSearchInput bind:value={search} />
	</div>
	<SvelteTable
		{columns}
		rows={filteredBattles}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full',
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
			{#if column.key == 'location'}
				<a class="cs-link text-2xl" href={`/games/${id}/battles/${row.num}`}
					>{$game?.player.getBattleLocation(row)}</a
				>
			{:else}
				{cell}
			{/if}
		</span>
	</SvelteTable>
</div>

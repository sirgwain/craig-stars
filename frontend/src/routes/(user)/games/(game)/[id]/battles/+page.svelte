<script lang="ts">
	import { page } from '$app/stores';
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { player } from '$lib/services/Context';
	import { GameService } from '$lib/services/GameService';
	import type { BattleRecord } from '$lib/types/Battle';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';
	import { onMount } from 'svelte';

	let id = parseInt($page.params.id);

	onMount(async () => {
		const p = await GameService.loadFullPlayer(id);
		player.update(() => p);
	});

	// filterable battles
	let filteredBattles: BattleRecord[] = [];
	let search = '';

	$: filteredBattles = $player?.battles ?? [];
	$: filteredBattles =
		$player?.battles.filter(
			(i) => $player?.getBattleLocation(i).toLowerCase().indexOf(search.toLowerCase()) != -1
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
					>{$player?.getBattleLocation(row)}</a
				>
			{:else}
				{cell}
			{/if}
		</span>
	</SvelteTable>
</div>

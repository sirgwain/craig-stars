<script lang="ts">
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import type { BattleRecord } from '$lib/types/Battle';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';

	const { game, player, universe } = getGameContext();

	// filterable battles
	let filteredBattles: BattleRecord[] = [];
	let search = '';

	$: filteredBattles = $universe.battles ?? [];
	$: filteredBattles =
		$universe.battles.filter(
			(i) =>
				$universe.getBattleLocation(i).toLowerCase().indexOf(search.toLowerCase()) != -1
		) ?? [];

	$: allies = new Set($player.getAllies());

	function getNumShips(record: BattleRecord): number {
		return Object.values(record.stats.numShipsByPlayer ?? {}).reduce(
			(count, num) => count + num,
			0
		);
	}

	function getOurShips(record: BattleRecord): number {
		let count = 0;
		allies.forEach(
			(ally) => (count += record.stats?.numShipsByPlayer ? record.stats?.numShipsByPlayer[ally] : 0)
		);
		return count;
	}

	function getTheirShips(record: BattleRecord): number {
		let count = 0;
		allies.forEach(
			(ally) =>
				(count +=
					record.stats?.numShipsByPlayer && !record.stats?.numShipsByPlayer[ally]
						? record.stats?.numShipsByPlayer[ally]
						: 0)
		);
		return count;
	}

	function getOurDead(record: BattleRecord): number {
		let count = 0;
		allies.forEach(
			(ally) =>
				(count += record.stats?.shipsDestroyedByPlayer
					? record.stats?.shipsDestroyedByPlayer[ally]
					: 0)
		);
		return count;
	}

	function getTheirDead(record: BattleRecord): number {
		let count = 0;
		allies.forEach(
			(ally) =>
				(count +=
					record.stats?.shipsDestroyedByPlayer && !record.stats?.shipsDestroyedByPlayer[ally]
						? record.stats?.shipsDestroyedByPlayer[ally]
						: 0)
		);
		return count;
	}

	const columns: SvelteTableColumn<BattleRecord>[] = [
		{
			key: 'location',
			title: 'Location'
		},
		{
			key: 'players',
			title: 'Players'
		},
		{
			key: 'ships',
			title: 'Ships'
		},
		{
			key: 'ours',
			title: 'Ours'
		},
		{
			key: 'theirs',
			title: 'Theirs'
		},
		{
			key: 'ourDead',
			title: 'Our Dead'
		},
		{
			key: 'theirDead',
			title: 'Their Dead'
		},
		{
			key: 'oursLeft',
			title: 'Ours Left'
		},
		{
			key: 'theirsLeft',
			title: 'Theirs Left'
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
				<a class="cs-link text-2xl" href={`/games/${$game.id}/battles/${row.num}`}
					>{$universe.getBattleLocation(row)}</a
				>
			{:else if column.key == 'players'}
				{row.stats?.numPlayers ?? 0}
			{:else if column.key == 'ships'}
				{getNumShips(row)}
			{:else if column.key == 'ours'}
				{getOurShips(row)}
			{:else if column.key == 'theirs'}
				{getTheirShips(row)}
			{:else if column.key == 'ourDead'}
				{getOurDead(row)}
			{:else if column.key == 'theirDead'}
				{getTheirDead(row)}
			{:else if column.key == 'oursLeft'}
				{getOurShips(row) - getOurDead(row)}
			{:else if column.key == 'theirsLeft'}
				{getTheirShips(row) - getTheirDead(row)}
			{:else}
				{cell ?? ''}
			{/if}
		</span>
	</SvelteTable>
</div>

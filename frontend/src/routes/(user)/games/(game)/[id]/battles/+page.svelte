<script lang="ts">
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import type { BattleRecord } from '$lib/types/Battle';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';

	const { game, player, universe } = getGameContext();

	type BattleRow = {
		num: number;
		location: string;
		numPlayers: number;
		numShips: number;
		ours: number;
		theirs: number;
		ourDead: number;
		theirDead: number;
		oursLeft: number;
		theirsLeft: number;
	};

	// filterable battles
	let filteredBattles: BattleRow[] = [];
	let search = '';
	$: allies = new Set($player.getAllies());

	$: battleRows = $universe.battles.map((b) => ({
		num: b.num,
		location: $universe.getBattleLocation(b),
		numPlayers: Object.keys(b.stats?.numShipsByPlayer ?? {}).length,
		numShips: getNumShips(b),
		ours: getOurShips(b),
		theirs: getTheirDead(b),
		ourDead: getOurDead(b),
		theirDead: getTheirDead(b),
		oursLeft: getOurShips(b) - getOurDead(b),
		theirsLeft: getTheirShips(b) - getTheirDead(b)
	}));

	$: filteredBattles =
		battleRows.filter((i) => i.location.toLowerCase().indexOf(search.toLowerCase()) != -1) ?? [];

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
		Object.entries(record.stats?.numShipsByPlayer ?? {}).forEach((entry) => {
			const playerNum: number = parseInt(entry[0]);
			const numShips = entry[1];
			if (!allies.has(playerNum)) {
				count += numShips;
			}
		});
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
		Object.entries(record.stats?.shipsDestroyedByPlayer ?? {}).forEach((entry) => {
			const playerNum: number = parseInt(entry[0]);
			const numShips = entry[1];
			if (!allies.has(playerNum)) {
				count += numShips;
			}
		});
		return count;
	}

	const columns: SvelteTableColumn<BattleRecord>[] = [
		{
			key: 'location',
			title: 'Location'
		},
		{
			key: 'numPlayers',
			title: 'Players'
		},
		{
			key: 'numShips',
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
				<a class="cs-link text-2xl" href={`/games/${$game.id}/battles/${row.num}`}>{row.location}</a
				>
			{:else}
				{cell ?? ''}
			{/if}
		</span>
	</SvelteTable>
</div>

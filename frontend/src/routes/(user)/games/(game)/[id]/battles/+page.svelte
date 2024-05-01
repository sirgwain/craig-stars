<script lang="ts">
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import {
		getNumShips,
		getOurDead,
		getOurShips,
		getTheirDead,
		getTheirShips
	} from '$lib/types/Battle';

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
		ours: getOurShips(b, allies),
		theirs: getTheirDead(b, allies),
		ourDead: getOurDead(b, allies),
		theirDead: getTheirDead(b, allies),
		oursLeft: getOurShips(b, allies) - getOurDead(b, allies),
		theirsLeft: getTheirShips(b, allies) - getTheirDead(b, allies)
	}));

	$: filteredBattles =
		battleRows.filter((i) => i.location.toLowerCase().indexOf(search.toLowerCase()) != -1) ?? [];

	const columns: TableColumn<BattleRow>[] = [
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
	<Table
		{columns}
		rows={filteredBattles}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full',
			td: 'first:table-cell nth-child(2):table-cell hidden sm:table-cell',
			th: 'first:table-cell nth-child(2):table-cell hidden sm:table-cell'
		}}
	>
		<span slot="head" let:isSorted let:sortDescending let:column>
			<SortableTableHeader {column} {isSorted} {sortDescending} />
		</span>

		<span slot="cell" let:column let:row let:cell>
			{#if column.key == 'location'}
				<a class="cs-link text-xl text-left" href={`/games/${$game.id}/battles/${row.num}`}>{row.location}</a
				>
			{:else}
				{cell ?? ''}
			{/if}
		</span>
	</Table>
</div>

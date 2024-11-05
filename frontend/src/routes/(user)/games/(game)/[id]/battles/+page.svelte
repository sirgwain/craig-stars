<script lang="ts">
	import { goto } from '$app/navigation';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { type BattleRecordDetails } from '$lib/types/Battle';
	import { Check } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { game, player, universe, settings, gotoBattle } = getGameContext();

	// filterable battles
	let filteredBattles: BattleRecordDetails[] = [];
	let search = '';

	$: battleRows = $universe.getBattles(
		$settings.sortBattlesKey,
		$settings.sortBattlesDescending,
		$player
	);

	$: filteredBattles =
		battleRows.filter((i) => i.location.toLowerCase().indexOf(search.toLowerCase()) != -1) ?? [];


	const columns: TableColumn<BattleRecordDetails>[] = [
		{
			key: 'location',
			title: 'Location'
		},
		{
			key: 'present',
			title: 'Present',
			hidden: $player.getAllies().length === 1 // hide if we're only friends with ourself
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

	function gotoTarget(row: BattleRecordDetails) {
		gotoBattle(row.num);
		goto(`/games/${$game.id}`);
	}

	function onSorted(column: TableColumn<BattleRecordDetails>, sortDescending: boolean) {
		$settings.sortBattlesDescending = sortDescending;
		$settings.sortBattlesKey = column.key;
	}
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
		externalSortAndFilter={true}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full',
			td: 'first:table-cell hidden sm:table-cell',
			th: 'first:table-cell hidden sm:table-cell'
		}}
	>
		<span slot="head" let:column>
			<SortableTableHeader
				{column}
				isSorted={$settings.sortBattlesKey === column.key}
				sortDescending={$settings.sortBattlesDescending}
				on:sorted={(e) => {
					onSorted(column, e.detail.sortDescending);
				}}
			/>
		</span>

		<span slot="cell" let:column let:row let:cell>
			{#if column.key == 'location'}
				<div class="flex flex-row justify-between">
					<a class="cs-link text-xl text-left" href={`/games/${$game.id}/battles/${row.num}`}
						>{row.location}</a
					>
					<button
						on:click={(e) => gotoTarget(row)}
						class="btn btn-outline btn-sm normal-case btn-secondary p-2 mx-1"
						title="goto">Goto</button
					>
				</div>
			{:else if column.key == 'present'}
				{#if row.present}
					<Icon src={Check} size="24" class="stroke-success" />
				{/if}
			{:else}
				{cell ?? ''}
			{/if}
		</span>
	</Table>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import CargoMini from '$lib/components/game/CargoMini.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { fleetsSortBy, getLocation, type Fleet } from '$lib/types/Fleet';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';

	const { game, player, universe, settings, commandMapObject, zoomToMapObject } = getGameContext();

	const selectFleet = (fleet: Fleet) => {
		commandMapObject(fleet);
		zoomToMapObject(fleet);
		goto(`/games/${$game.id}`);
	};

	// filterable fleets
	let filteredFleets: Fleet[] = [];
	let search = '';

	$: filteredFleets =
		$universe
			.getMyFleets($settings.sortFleetsKey, $settings.sortFleetsDescending)
			.filter((i) => i.name.toLowerCase().indexOf(search.toLowerCase()) != -1) ?? [];

	const columns: TableColumn<Fleet>[] = [
		{
			key: 'name',
			title: 'Name',
			sortBy: fleetsSortBy('name', $universe)
		},
		{
			key: 'num',
			title: 'ID',
			sortBy: fleetsSortBy('num', $universe)
		},
		{
			key: 'location',
			title: 'Location',
			sortBy: fleetsSortBy('location', $universe)
		},
		{
			key: 'destination',
			title: 'Destination',
			sortBy: fleetsSortBy('destination', $universe)
		},
		{
			key: 'eta',
			title: 'ETA',
			sortBy: fleetsSortBy('eta', $universe)
		},
		{
			key: 'fuel',
			title: 'Fuel',
			sortBy: fleetsSortBy('fuel', $universe)
		},
		{
			key: 'cargo',
			title: 'Cargo',
			sortBy: fleetsSortBy('cargo', $universe)
		},
		{
			key: 'composition',
			title: 'Composition',
			sortable: false,
			sortBy: fleetsSortBy('composition', $universe)
		},
		{
			key: 'cloak',
			title: 'Cloak',
			sortBy: fleetsSortBy('cloak', $universe)
		},
		{
			key: 'battlePlanNum',
			title: 'Battle Plan',
			sortBy: fleetsSortBy('battlePlanNum', $universe)
		},
		{
			key: 'mass',
			title: 'Mass',
			sortBy: fleetsSortBy('mass', $universe)
		}
	];

	function onSorted(column: TableColumn<Fleet>, sortDescending: boolean) {
		$settings.sortFleetsDescending = sortDescending;
		$settings.sortFleetsKey = column.key;
	}
</script>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<TableSearchInput bind:value={search} />
	</div>
	<Table
		{columns}
		rows={filteredFleets}
		externalSortAndFilter={true}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full'
		}}
	>
		<span slot="head" let:column>
			<SortableTableHeader
				{column}
				isSorted={$settings.sortFleetsKey === column.key}
				sortDescending={$settings.sortFleetsDescending}
				on:sorted={(e) => {
					onSorted(column, e.detail.sortDescending);
				}}
			/>
		</span>

		<span slot="cell" let:column let:row let:cell>
			{#if column.key == 'name'}
				<button class="cs-link text-xl text-left" on:click={() => selectFleet(row)}>{cell}</button>
			{:else if column.key == 'location'}
				{getLocation(row, $universe)}
			{:else if column.key == 'destination'}
				{row.waypoints && row.waypoints.length > 1
					? $universe.getTargetName(row.waypoints[1])
					: '--'}
			{:else if column.key == 'eta'}
				-- <!-- TODO: fleet class?  -->
			{:else if column.key == 'fuel'}
				{row.fuel}mg
			{:else if column.key == 'cargo'}
				<CargoMini cargo={row.cargo} />
			{:else if column.key == 'composition'}
				{@const design = $game
					? $universe.getDesign(
							$player.num,
							row.tokens && row.tokens.length ? row.tokens[0].designNum : 0
						)
					: undefined}
				<div class="flex flex-row justify-between">
					<div>
						{design ? design.name : ''}
					</div>
					<div>
						{row.tokens && row.tokens.length ? row.tokens[0].quantity : 0}
					</div>
				</div>
			{:else if column.key == 'cloak'}
				{row.spec && row.spec.cloakPercent ? row.spec.cloakPercent + '%' : ''}
			{:else if column.key == 'battlePlanNum'}
				{@const battlePlan = $game ? $player.getBattlePlan(row.battlePlanNum ?? 0) : undefined}
				{battlePlan?.name ?? ''}
			{:else if column.key == 'mass'}
				{row.spec?.mass ?? 0}
			{:else}
				{cell}
			{/if}
		</span>
	</Table>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import CargoMini from '$lib/components/game/CargoMini.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { commandMapObject, zoomToMapObject } from '$lib/services/Stores';
	import { totalCargo } from '$lib/types/Cargo';
	import type { Fleet } from '$lib/types/Fleet';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';

	const { game, player, universe } = getGameContext();

	const selectFleet = (fleet: Fleet) => {
		commandMapObject(fleet);
		zoomToMapObject(fleet);
		goto(`/games/${$game.id}`);
	};

	const getLocation = (fleet: Fleet) =>
		fleet.orbitingPlanetNum
			? $universe.getPlanet(fleet.orbitingPlanetNum)?.name ?? 'unknown'
			: `Space: (${fleet.position.x}, ${fleet.position.y})`;

	const getDestination = (fleet: Fleet) => {
		if (fleet.waypoints?.length && fleet.waypoints?.length > 1) {
			return $universe.getTargetName(fleet.waypoints[1]);
		}
		return '--';
	};

	// filterable fleets
	let filteredFleets: Fleet[] = [];
	let search = '';

	$: filteredFleets =
		$universe.fleets
			.sort((a, b) => a.num - b.num)
			.filter(
				(i) =>
					i.playerNum === $player.num && i.name.toLowerCase().indexOf(search.toLowerCase()) != -1
			) ?? [];

	const columns: SvelteTableColumn<Fleet>[] = [
		{
			key: 'name',
			title: 'Name'
		},
		{
			key: 'num',
			title: 'ID'
		},
		{
			key: 'location',
			title: 'Location',
			sortBy: (a, b) => getLocation(a).localeCompare(getLocation(b))
		},
		{
			key: 'destination',
			title: 'Destination',
			sortBy: (a, b) => getDestination(a).localeCompare(getDestination(b))
		},
		{
			key: 'eta',
			title: 'ETA'
		},
		{
			key: 'fuel',
			title: 'Fuel'
		},
		{
			key: 'cargo',
			title: 'Cargo',
			sortBy: (a, b) => totalCargo(a.cargo) - totalCargo(b.cargo)
		},
		{
			key: 'composition',
			title: 'Composition'
		},
		{
			key: 'cloak',
			title: 'Cloak'
		},
		{
			key: 'battlePlanNum',
			title: 'Battle Plan'
		},
		{
			key: 'mass',
			title: 'Mass',
			sortBy: (a, b) => (a.spec?.mass ?? 0) - (b.spec?.mass ?? 0)
		}
	];
</script>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<TableSearchInput bind:value={search} />
	</div>
	<SvelteTable
		{columns}
		rows={filteredFleets}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full',
		}}
		let:column
		let:cell
		let:row
	>
		<span slot="head" let:isSorted let:sortDescending>
			<SortableTableHeader {column} {isSorted} {sortDescending} />
		</span>

		<span slot="cell">
			{#if column.key == 'name'}
				<button class="cs-link text-2xl" on:click={() => selectFleet(row)}>{cell}</button>
			{:else if column.key == 'location'}
				{getLocation(row)}
			{:else if column.key == 'destination'}
				{row.waypoints.length > 1 ? $universe.getTargetName(row.waypoints[1]) : '--'}
			{:else if column.key == 'eta'}
				-- <!-- TODO: fleet class?  -->
			{:else if column.key == 'fuel'}
				{row.fuel}mg
			{:else if column.key == 'cargo'}
				<CargoMini cargo={row.cargo} />
			{:else if column.key == 'composition'}
				{@const design = $game
					? $universe.getDesign($player.num, row.tokens[0].designNum)
					: undefined}
				<div class="flex flex-row justify-between">
					<div>
						{design ? design.name : ''}
					</div>
					<div>
						{row.tokens[0].quantity}
					</div>
				</div>
			{:else if column.key == 'cloak'}
				{row.spec.cloak ? row.spec.cloak + '%' : ''}
			{:else if column.key == 'battlePlanNum'}
				{@const battlePlan = $game ? $player.getBattlePlan(row.battlePlanNum) : undefined}
				{battlePlan?.name ?? ''}
			{:else if column.key == 'mass'}
				{row.spec.mass}
			{:else}
				{cell}
			{/if}
		</span>
	</SvelteTable>
</div>

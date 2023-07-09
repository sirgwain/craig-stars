<script lang="ts">
	import { goto } from '$app/navigation';
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import MineralMini from '$lib/components/game/MineralMini.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { commandMapObject, zoomToMapObject } from '$lib/services/Stores';
	import { getQueueItemShortName, planetsSortBy, type Planet } from '$lib/types/Planet';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';

	const { game, player, universe, settings } = getGameContext();

	const selectPlanet = (planet: Planet) => {
		commandMapObject(planet);
		zoomToMapObject(planet);
		goto(`/games/${$game.id}`);
	};

	// filterable planets
	let filteredPlanets: Planet[] = [];
	let search = '';

	$: filteredPlanets =
		$universe
			.getMyPlanets()
			.filter(
				(i) =>
					i.playerNum == $player.num && i.name.toLowerCase().indexOf(search.toLowerCase()) != -1
			) ?? [];

	const columns: SvelteTableColumn<Planet>[] = [
		{
			key: 'name',
			title: 'Name',
			sortBy: planetsSortBy('name')
		},
		{
			key: 'starbase',
			title: 'Starbase',
			sortBy: planetsSortBy('starbase')
		},
		{
			key: 'population',
			title: 'Population',
			sortBy: planetsSortBy('population')
		},
		{
			key: 'populationDensity',
			title: 'Cap',
			sortBy: planetsSortBy('populationDensity')
		},
		{
			key: 'habitability',
			title: 'Value',
			sortBy: planetsSortBy('habitability')
		},
		{
			key: 'production',
			title: 'Production',
			sortBy: planetsSortBy('production')
		},
		{
			key: 'mines',
			title: 'Mine',
			sortBy: planetsSortBy('mines')
		},
		{
			key: 'factories',
			title: 'Factories',
			sortBy: planetsSortBy('factories')
		},
		{
			key: 'defense',
			title: 'Defense',
			sortBy: planetsSortBy('defense')
		},
		{
			key: 'minerals',
			title: 'Minerals',
			sortBy: planetsSortBy('minerals')
		},
		{
			key: 'miningRate',
			title: 'Mining Rate',
			sortBy: planetsSortBy('miningRate')
		},
		{
			key: 'mineralConcentration',
			title: 'Mineral Concentration',
			sortBy: planetsSortBy('mineralConcentration')
		},
		{
			key: 'resources',
			title: 'Resources',
			sortBy: planetsSortBy('resources')
		},
		{
			key: 'driverDest',
			title: 'Driver Destination',
			sortable: false
		},
		{
			key: 'routingDestination',
			title: 'routing Destination',
			sortable: false
		}
	];

	function onSorted(column: SvelteTableColumn, sortDescending: boolean) {
		$settings.sortPlanetsDescending = sortDescending;
		$settings.sortPlanetsKey = column.key;
	}
</script>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<TableSearchInput bind:value={search} />
	</div>
	<SvelteTable
		{columns}
		rows={filteredPlanets}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full'
		}}
		let:column
		let:cell
		let:row
	>
		<span slot="head" let:isSorted let:sortDescending>
			<SortableTableHeader
				{column}
				isSorted={isSorted || $settings.sortPlanetsKey === column.key}
				sortDescending={sortDescending ||
					($settings.sortPlanetsKey === column.key && $settings.sortPlanetsDescending)}
				on:sorted={(e) => {
					onSorted(column, e.detail.sortDescending);
				}}
			/>
		</span>

		<span slot="cell">
			{#if column.key == 'name'}
				<button class="cs-link text-2xl" on:click={() => selectPlanet(row)}>{cell}</button>
			{:else if column.key == 'starbase'}
				{row.spec.starbaseDesignName ?? ''}
			{:else if column.key == 'population'}
				{((row.cargo.colonists ?? 0) * 100).toLocaleString()}
			{:else if column.key == 'populationDensity'}
				{((row.spec.populationDensity ?? 0) * 100).toFixed(1)}%
			{:else if column.key == 'habitability'}
				{row.spec.habitability ?? 0}%
			{:else if column.key == 'production'}
				{#if row.productionQueue?.length}
					<div class="flex justify-between">
						<div>{getQueueItemShortName(row.productionQueue[0], $universe)}</div>
						<div>{row.productionQueue[0].quantity}</div>
					</div>
				{:else}
					--- Queue is Empty ---
				{/if}
			{:else if column.key == 'mines'}
				{row.mines ?? 0}
			{:else if column.key == 'factories'}
				{row.factories ?? 0}
			{:else if column.key == 'defense'}
				{((row.spec.defenseCoverage ?? 0) * 100).toFixed(1)}%
			{:else if column.key == 'minerals'}
				<MineralMini mineral={row.cargo} />
			{:else if column.key == 'miningRate'}
				<MineralMini mineral={row.spec.miningOutput} />
			{:else if column.key == 'mineralConcentration'}
				<MineralMini mineral={row.mineralConcentration} />
			{:else if column.key == 'resources'}
				{row.spec.resourcesPerYearAvailable ?? 0} / {row.spec.resourcesPerYear ?? 0}
			{:else if column.key == 'driverDest'}
				--
			{:else if column.key == 'routingDestination'}
				--
			{:else}
				{cell}
			{/if}
		</span>
	</SvelteTable>
</div>

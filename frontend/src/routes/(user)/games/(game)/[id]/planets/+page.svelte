<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import SortableTableHeader from '$lib/components/SortableTableHeader.svelte';
	import TableSearchInput from '$lib/components/TableSearchInput.svelte';
	import MineralMini from '$lib/components/game/MineralMini.svelte';
	import { commandMapObject, game, zoomToMapObject } from '$lib/services/Context';
	import { totalMinerals } from '$lib/types/Mineral';
	import { getQueueItemShortName, type Planet } from '$lib/types/Planet';
	import { SvelteTable, type SvelteTableColumn } from '@hurtigruten/svelte-table';

	let id = parseInt($page.params.id);

	const selectPlanet = (planet: Planet) => {
		commandMapObject(planet);
		zoomToMapObject(planet);
		goto(`/games/${id}`);
	};

	// filterable planets
	let filteredPlanets: Planet[] = [];
	let search = '';

	$: filteredPlanets = $game?.universe.planets ?? [];
	$: filteredPlanets =
		$game?.universe.planets.filter(
			(i) => i.name.toLowerCase().indexOf(search.toLowerCase()) != -1
		) ?? [];

	const columns: SvelteTableColumn<Planet>[] = [
		{
			key: 'name',
			title: 'Name'
		},
		{
			key: 'starbase',
			title: 'Starbase',
			sortBy: (a, b) =>
				(a.spec.starbaseDesignName ?? '').localeCompare(b.spec.starbaseDesignName ?? '')
		},
		{
			key: 'population',
			title: 'Population',
			sortBy: (a, b) => (a.cargo?.colonists ?? 0) - (b.cargo?.colonists ?? 0)
		},
		{
			key: 'populationDensity',
			title: 'Cap',
			sortBy: (a, b) => (a.spec.populationDensity ?? 0) - (b.spec.populationDensity ?? 0)
		},
		{
			key: 'habitability',
			title: 'Value',
			sortBy: (a, b) => (a.spec.habitability ?? 0) - (b.spec.habitability ?? 0)
		},
		{
			key: 'production',
			title: 'Production',
			sortable: false
		},
		{
			key: 'mines',
			title: 'Mine',
			sortBy: (a, b) => (a.mines ?? 0) - (b.mines ?? 0)
		},
		{
			key: 'factories',
			title: 'Factories',
			sortBy: (a, b) => (a.factories ?? 0) - (b.factories ?? 0)
		},
		{
			key: 'defense',
			title: 'Defense',
			sortBy: (a, b) => (a.spec.defenseCoverage ?? 0) - (b.spec.defenseCoverage ?? 0)
		},
		{
			key: 'minerals',
			title: 'Minerals',
			sortBy: (a, b) => totalMinerals(a.cargo) - totalMinerals(b.cargo)
		},
		{
			key: 'miningRate',
			title: 'Mining Rate',
			sortBy: (a, b) => totalMinerals(a.spec.mineralOutput) - totalMinerals(b.spec.mineralOutput)
		},
		{
			key: 'mineralConcentration',
			title: 'Mineral Concentration',
			sortBy: (a, b) =>
				totalMinerals(a.mineralConcentration) - totalMinerals(b.mineralConcentration)
		},
		{
			key: 'resources',
			title: 'Resources',
			sortBy: (a, b) =>
				(a.spec.resourcesPerYearAvailable ?? 0) - (b.spec.resourcesPerYearAvailable ?? 0)
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
</script>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<TableSearchInput bind:value={search} />
	</div>
	<SvelteTable
		{columns}
		rows={filteredPlanets}
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
						<div>{getQueueItemShortName(row.productionQueue[0])}</div>
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
				<MineralMini mineral={row.spec.mineralOutput} />
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

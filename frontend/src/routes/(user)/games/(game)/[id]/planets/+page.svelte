<script lang="ts">
	import { goto } from '$app/navigation';
	import MineralMini from '$lib/components/game/MineralMini.svelte';
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { getQueueItemShortName, planetsSortBy, type Planet } from '$lib/types/Planet';
	import { Check } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import ProductionQueueDialog from '../dialogs/production/ProductionQueueDialog.svelte';
	import ProductionQueueItemLine from '$lib/components/game/ProductionQueueItemLine.svelte';

	const { game, universe, settings, commandMapObject, selectMapObject, zoomToMapObject } =
		getGameContext();

	const selectPlanet = (planet: Planet) => {
		commandMapObject(planet);
		selectMapObject(planet);
		zoomToMapObject(planet);
		goto(`/games/${$game.id}`);
	};

	// filterable planets
	let filteredPlanets: Planet[] = [];
	let search = '';

	// production queue dialog
	let showProductionQueueDialog = false;

	$: filteredPlanets =
		$universe
			.getMyPlanets($settings.sortPlanetsKey, $settings.sortPlanetsDescending)
			.filter((i) => i.name.toLowerCase().indexOf(search.toLowerCase()) != -1) ?? [];

	const columns: TableColumn<Planet>[] = [
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
			key: 'contributesOnlyLeftoverToResearch',
			title: 'Contributes Only Leftover To Research',
			sortBy: planetsSortBy('contributesOnlyLeftoverToResearch')
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

	function onSorted(column: TableColumn<Planet>, sortDescending: boolean) {
		$settings.sortPlanetsDescending = sortDescending;
		$settings.sortPlanetsKey = column.key;
	}

	function onProductionQueueDialog(planet: Planet) {
		commandMapObject(planet);
		showProductionQueueDialog = true;
	}
</script>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<TableSearchInput bind:value={search} />
	</div>
	<Table
		{columns}
		rows={filteredPlanets}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full'
		}}
	>
		<div slot="head" let:isSorted let:sortDescending let:column>
			<SortableTableHeader
				{column}
				isSorted={isSorted || $settings.sortPlanetsKey === column.key}
				sortDescending={sortDescending ||
					($settings.sortPlanetsKey === column.key && $settings.sortPlanetsDescending)}
				on:sorted={(e) => {
					onSorted(column, e.detail.sortDescending);
				}}
			/>
		</div>

		<span slot="cell" let:row let:column let:cell>
			{#if column.key == 'name'}
				<button class="cs-link text-2xl" on:click={() => selectPlanet(row)}>{cell}</button>
			{:else if column.key == 'starbase'}
				{row.spec.starbaseDesignName ?? ''}
			{:else if column.key == 'population'}
				{((row.cargo?.colonists ?? 0) * 100).toLocaleString()}
			{:else if column.key == 'populationDensity'}
				{((row.spec.populationDensity ?? 0) * 100).toFixed(1)}%
			{:else if column.key == 'habitability'}
				{#if row.spec.canTerraform}
					<span
						class:text-habitable={(row.spec.habitability ?? 0) > 0}
						class:text-uninhabitable={(row.spec.habitability ?? 0) < 0}
						>{row.spec.habitability ?? 0}%</span
					>
					/ <span class="text-terraformable">{row.spec.terraformedHabitability ?? 0}%</span>
				{:else}
					<span
						class:text-habitable={(row.spec.habitability ?? 0) > 0}
						class:text-uninhabitable={(row.spec.habitability ?? 0) < 0}
					>
						{row.spec.habitability ?? 0}%</span
					>
				{/if}
			{:else if column.key == 'production'}
				{#if row.productionQueue?.length}
					<button
						on:click={() => onProductionQueueDialog(row)}
						class="text-base w-32 flex justify-between text-left cursor-pointer"
					>
						<ProductionQueueItemLine item={row.productionQueue[0]} index={0} shortName={true} />
					</button>
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
			{:else if column.key == 'contributesOnlyLeftoverToResearch'}
				{#if row.contributesOnlyLeftoverToResearch}
					<Icon src={Check} size="24" class="stroke-success" />
				{/if}
			{:else if column.key == 'driverDest'}
				--
			{:else if column.key == 'routingDestination'}
				--
			{:else}
				{cell}
			{/if}
		</span>
	</Table>
</div>

<ProductionQueueDialog bind:show={showProductionQueueDialog} />

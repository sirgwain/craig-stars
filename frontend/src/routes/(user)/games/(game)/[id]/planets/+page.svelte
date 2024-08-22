<script lang="ts">
	import { goto } from '$app/navigation';
	import MineralMini from '$lib/components/game/MineralMini.svelte';
	import ProductionQueueItemLine from '$lib/components/game/ProductionQueueItemLine.svelte';
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { owned, ownedBy } from '$lib/types/MapObject';
	import { totalMinerals } from '$lib/types/Mineral';
	import { Unexplored, planetsSortBy, type Planet } from '$lib/types/Planet';
	import { Check } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import ProductionQueueDialog from '../dialogs/production/ProductionQueueDialog.svelte';
	import PopulationTooltip from '$lib/components/game/tooltips/PopulationTooltip.svelte';
	import type { PopulationTooltipProps } from '$lib/components/game/tooltips/PopulationTooltip.svelte';
	import { showTooltip } from '$lib/services/Stores';

	const { game, player, universe, settings, commandMapObject, selectMapObject, zoomToMapObject } =
		getGameContext();

	const selectPlanet = (planet: Planet) => {
		if (ownedBy(planet, $player.num)) {
			commandMapObject(planet);
		}
		selectMapObject(planet);
		zoomToMapObject(planet);
		goto(`/games/${$game.id}`);
	};

	// filterable planets
	let filteredPlanets: Planet[] = [];
	let search = '';

	// production queue dialog
	let showProductionQueueDialog = false;

	$: filteredPlanets = $settings.showAllPlanets
		? $universe
				.getPlanets($settings.sortPlanetsKey, $settings.sortPlanetsDescending)
				.filter(
					(i) =>
						i.name.toLowerCase().indexOf(search.toLowerCase()) != -1 ||
						$universe.getPlayerPluralName(i.playerNum)?.toLowerCase().indexOf(search.toLowerCase()) != -1
				) ?? []
		: $universe
				.getMyPlanets($settings.sortPlanetsKey, $settings.sortPlanetsDescending)
				.filter((i) => i.name.toLowerCase().indexOf(search.toLowerCase()) != -1) ?? [];

	$: columns = [
		{
			key: 'name',
			title: 'Name',
			sortBy: planetsSortBy('name')
		},
		{
			key: 'owner',
			title: 'Owner',
			hidden: !$settings.showAllPlanets,
			sortBy: (a, b) =>
				$universe.getPlayerPluralName(a.playerNum)?.localeCompare($universe.getPlayerPluralName(b.playerNum))
		},
		{
			key: 'reportAge',
			title: 'Report Age',
			hidden: !$settings.showAllPlanets,
			sortBy: (a, b) => (a.reportAge ?? 0) - (b.reportAge ?? 0)
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
			hidden: $settings.showAllPlanets,
			sortBy: planetsSortBy('populationDensity')
		},
		{
			key: 'populationGrowth',
			title: 'Growth',
			hidden: $settings.showAllPlanets,
			sortBy: planetsSortBy('populationGrowth')
		},
		{
			key: 'habitability',
			title: 'Value',
			sortBy: planetsSortBy('habitability')
		},
		{
			key: 'production',
			title: 'Production',
			hidden: $settings.showAllPlanets,
			sortBy: planetsSortBy('production')
		},
		{
			key: 'mines',
			title: 'Mine',
			hidden: $settings.showAllPlanets,
			sortBy: planetsSortBy('mines')
		},
		{
			key: 'factories',
			title: 'Factories',
			hidden: $settings.showAllPlanets,
			sortBy: planetsSortBy('factories')
		},
		{
			key: 'defense',
			title: 'Defense',
			hidden: $settings.showAllPlanets,
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
			hidden: $settings.showAllPlanets,
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
			hidden: $settings.showAllPlanets,
			sortBy: planetsSortBy('resources')
		},
		{
			key: 'contributesOnlyLeftoverToResearch',
			title: 'Contributes Only Leftover To Research',
			hidden: $settings.showAllPlanets,
			sortBy: planetsSortBy('contributesOnlyLeftoverToResearch')
		},
		{
			key: 'driverDest',
			title: 'Driver Destination',
			hidden: $settings.showAllPlanets,
			sortable: false
		},
		{
			key: 'routingDestination',
			title: 'routing Destination',
			hidden: $settings.showAllPlanets,
			sortable: false
		}
	] as TableColumn<Planet>[];

	function onSorted(column: TableColumn<Planet>, sortDescending: boolean) {
		$settings.sortPlanetsDescending = sortDescending;
		$settings.sortPlanetsKey = column.key;
	}

	function onProductionQueueDialog(planet: Planet) {
		commandMapObject(planet);
		showProductionQueueDialog = true;
	}

	function onPopulationTooltip(e: PointerEvent, planet: Planet) {
		showTooltip<PopulationTooltipProps>(e.x, e.y, PopulationTooltip, {
			playerFinder: $universe,
			player: $player,
			planet
		});
	}
</script>

<div class="w-full">
	<div class="flex flex-row justify-between m-2">
		<div><TableSearchInput bind:value={search} /></div>
		<div class="form-control">
			<label class="label cursor-pointer">
				<span class="label-text mr-1">Show All</span>
				<input
					type="checkbox"
					class="toggle"
					class:toggle-accent={$settings.showAllPlanets}
					bind:checked={$settings.showAllPlanets}
				/>
			</label>
		</div>
	</div>
	<Table
		{columns}
		rows={filteredPlanets}
		externalSortAndFilter={true}
		classes={{
			table: 'table table-zebra table-compact table-auto w-full'
		}}
	>
		<div slot="head" let:column>
			<SortableTableHeader
				{column}
				isSorted={$settings.sortPlanetsKey === column.key}
				sortDescending={$settings.sortPlanetsDescending}
				on:sorted={(e) => {
					onSorted(column, e.detail.sortDescending);
				}}
			/>
		</div>

		<span slot="cell" let:row let:column let:cell>
			{#if column.key == 'name'}
				<button class="cs-link text-xl text-left" on:click={() => selectPlanet(row)}>{cell}</button>
			{:else if column.key == 'owner'}
				<span style={`color: ${$universe.getPlayerColor(row.playerNum)};`}>
					{owned(row) ? $universe.getPlayerPluralName(row.playerNum) ?? '' : ''}
				</span>
			{:else if column.key == 'reportAge'}
				{#if row.reportAge == 0 || row.reportAge === undefined}
					current
				{:else if row.reportAge == Unexplored}
					unexplored
				{:else}
					{row.reportAge} years old
				{/if}
			{:else if column.key == 'starbase'}
				{row.spec.starbaseDesignName ?? ''}
			{:else if column.key == 'population'}
				<div class="cursor-help" on:pointerdown|preventDefault={(e) => onPopulationTooltip(e, row)}>
					{row.spec.population ? row.spec.population.toLocaleString() : ''}
				</div>
			{:else if column.key == 'populationDensity'}
				<div class="cursor-help" on:pointerdown|preventDefault={(e) => onPopulationTooltip(e, row)}>
					{((row.spec.populationDensity ?? 0) * 100).toFixed(1)}%
				</div>
			{:else if column.key == 'populationGrowth'}
				<div class="cursor-help" on:pointerdown|preventDefault={(e) => onPopulationTooltip(e, row)}>
					{(row.spec.growthAmount ?? 0).toLocaleString()}
				</div>
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
				{:else if ownedBy(row, $player.num)}
					--- Queue is Empty ---
				{/if}
			{:else if column.key == 'mines'}
				{row.mines ?? 0}
			{:else if column.key == 'factories'}
				{row.factories ?? 0}
			{:else if column.key == 'defense'}
				{((row.spec.defenseCoverage ?? 0) * 100).toFixed(1)}%
			{:else if column.key == 'minerals'}
				{#if totalMinerals(row.cargo) != 0}
					<MineralMini mineral={row.cargo} />
				{/if}
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

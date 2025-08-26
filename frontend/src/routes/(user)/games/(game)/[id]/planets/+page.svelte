<script lang="ts">
	import { goto } from '$app/navigation';
	import MineralMini from '$lib/components/game/MineralMini.svelte';
	import ProductionQueueItemLine from '$lib/components/game/ProductionQueueItemLine.svelte';
	import FactoriesTooltip, {
		type FactoriesTooltipProps
	} from '$lib/components/game/tooltips/FactoriesTooltip.svelte';
	import MinesTooltip, {
		type MinesTooltipProps
	} from '$lib/components/game/tooltips/MinesTooltip.svelte';
	import type { PopulationTooltipProps } from '$lib/components/game/tooltips/PopulationTooltip.svelte';
	import PopulationTooltip from '$lib/components/game/tooltips/PopulationTooltip.svelte';
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import { onTechTooltip } from '$lib/components/game/tooltips/TechTooltip.svelte';
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { showTooltip, techs } from '$lib/services/Stores';
	import { population } from '$lib/types/Cargo';
	import { ReportAgeUnexplored } from '$lib/types/Consts';
	import type { Planet, ShipDesign } from '$lib/types/cs-proto';
	import { MapObjectTargetSchema, MapObjectType, MineralSchema } from '$lib/types/cs-proto';
	import { owned, ownedBy, type MapObjectLike } from '$lib/types/MapObject';
	import { getGrowth, planetsSortBy } from '$lib/types/Planet';
	import { emptyVector } from '$lib/types/Vector';
	import { create } from '@bufbuild/protobuf';
	import { Check } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import ProductionQueueDialog from '../dialogs/production/ProductionQueueDialog.svelte';

	const {
		game,
		player,
		universe,
		settings,
		commandedPlanet,
		commandMapObject,
		selectMapObject,
		zoomToMapObject,
		nextMapObject,
		previousMapObject,
		updatePlanetOrders
	} = getGameContext();

	// filterable planets
	let search = $state('');

	// production queue dialog
	let showProductionQueueDialog = $state(false);

	let filteredPlanets: Planet[] = $derived(
		$settings.showAllPlanets
			? $universe
					.getPlanets($settings.sortPlanetsKey, $settings.sortPlanetsDescending)
					.map<TablePlanet>( // convert Planet to a TablePlanet, so populate all the planet fields as empty
						(r) =>
							({
								...r,
								mines: 0,
								factories: 0,
								mineYears: create(MineralSchema),
								defenses: 0
							}) as TablePlanet
					)
					.filter(
						(i) =>
							i.mapObject?.name.toLowerCase().indexOf(search.toLowerCase()) != -1 ||
							$universe
								.getPlayerPluralName(i.mapObject.playerNum)
								.toLowerCase()
								.indexOf(search.toLowerCase()) != -1
					)
			: $universe
					.getMyPlanets($settings.sortPlanetsKey, $settings.sortPlanetsDescending)
					.map<TablePlanet>((r) => r as TablePlanet)
					.filter((i) => i.mapObject?.name.toLowerCase().indexOf(search.toLowerCase()) != -1)
	);

	// columns change based on whether we are showing all planets or just the player planets
	type TablePlanet = Planet & {
		name?: never;
		owner?: never;
		population?: never;
		populationDensity?: never;
		populationGrowth?: never;
		habitability?: never;
		production?: never;
		defense?: never;
		minerals?: never;
		miningRate?: never;
		resources?: never;
		driverDest?: never;
		routingDestination?: never;
		reportAge?: number;
		starbase?: never;
		mines?: number;
		factories?: number;
		contributesOnlyLeftoverToResearch?: boolean;
	};
	let columns: TableColumn<TablePlanet>[] = $derived([
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
				$universe
					.getPlayerPluralName(a.mapObject?.playerNum)
					.localeCompare($universe.getPlayerPluralName(b.mapObject?.playerNum))
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
			title: 'Defense Coverage',
			sortBy: planetsSortBy('defense')
		},
		{
			key: 'minerals',
			title: 'Surface Minerals',
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
			title: 'Routing Destination',
			hidden: $settings.showAllPlanets,
			sortable: false
		}
	]);

	function onSorted(column: TableColumn<TablePlanet>, sortDescending: boolean) {
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

	function onMinesTooltip(e: PointerEvent, planet: Planet) {
		e.preventDefault();
		showTooltip<MinesTooltipProps>(e.x, e.y, MinesTooltip, {
			planetName: planet.mapObject?.name ?? '',
			mines: planet.mines,
			maxMines: planet.spec?.maxMines ?? 0,
			maxPossibleMines: planet.spec?.maxPossibleMines ?? 0,
			canBuildMines: $player.race.spec.innateMining
		});
	}

	function onFactoriesTooltip(e: PointerEvent, planet: Planet) {
		e.preventDefault();
		showTooltip<FactoriesTooltipProps>(e.x, e.y, FactoriesTooltip, {
			planetName: planet.mapObject?.name ?? '',
			factories: planet.factories,
			maxFactories: planet.spec?.maxFactories ?? 0,
			maxPossibleFactories: planet.spec?.maxPossibleFactories ?? 0,
			canBuildFactories: $player.race.spec.innateResources
		});
	}

	function showDesign(e: PointerEvent, planet: Planet) {
		e.preventDefault();
		onShipDesignTooltip(
			e,
			$universe.getDesign(
				planet.mapObject?.playerNum ?? 0,
				planet.spec?.planetStarbaseSpec?.starbaseDesignNum ?? 0
			) as ShipDesign | undefined
		);
	}

	function onDefenseTooltip(e: PointerEvent, planet: Planet) {
		e.preventDefault();
		onTechTooltip(e, $techs.getTech(planet.spec?.defense ?? ''));
	}

	function gotoMapObject(mo: MapObjectLike) {
		if (
			(ownedBy(mo, $player.num) && mo.mapObject?.type === MapObjectType.PLANET) ||
			mo.mapObject?.type === MapObjectType.FLEET
		) {
			commandMapObject(mo);
		}
		selectMapObject(mo);
		zoomToMapObject(mo);
		goto(`/games/${$game.id}`);
	}

	async function onNextPlanet(updateOrders: boolean) {
		if (!$commandedPlanet) {
			return;
		}
		if (updateOrders) {
			await updatePlanetOrders($commandedPlanet);
		}

		nextMapObject();
	}

	async function onPrevPlanet(updateOrders: boolean) {
		if (!$commandedPlanet) {
			return;
		}
		if (updateOrders) {
			await updatePlanetOrders($commandedPlanet);
		}

		previousMapObject();
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
			table: 'table table-zebra table-compact table-auto w-full',
			th: 'sticky top-0 bg-base-200 z-10'
		}}
	>
		{#snippet head({ column })}
			<div>
				<SortableTableHeader
					{column}
					isSorted={$settings.sortPlanetsKey === column.key}
					sortDescending={$settings.sortPlanetsDescending}
					{onSorted}
				/>
			</div>
		{/snippet}

		{#snippet cell({ row, column, cell })}
			{@const planet = row as Planet}
			<span>
				{#if column.key == 'name'}
					<button class="cs-link text-xl text-left" onclick={() => gotoMapObject(row)}
						>{row.mapObject?.name}</button
					>
				{:else if column.key == 'owner'}
					<span style={`color: ${$universe.getPlayerColor(row.mapObject?.playerNum)};`}>
						{owned(row) ? $universe.getPlayerPluralName(row.mapObject?.playerNum) : ''}
					</span>
				{:else if column.key == 'reportAge'}
					{#if row.mapObject?.reportAge === 0 || row.mapObject?.reportAge === undefined}
						current
					{:else if row.mapObject?.reportAge === ReportAgeUnexplored}
						unexplored
					{:else}
						{row.mapObject?.reportAge} years old
					{/if}
				{:else if column.key == 'starbase'}
					{#if row.spec?.planetStarbaseSpec?.starbaseDesignName}
						<span class="cursor-help" onpointerdown={(e) => showDesign(e, row)}>
							{row.spec.planetStarbaseSpec.starbaseDesignName}
						</span>
					{/if}
				{:else if column.key == 'population'}
					<div class="cursor-help" onpointerdown={(e) => onPopulationTooltip(e, row)}>
						{population(row.cargo) ? population(row.cargo).toLocaleString() : ''}
					</div>
				{:else if column.key == 'populationDensity'}
					<div class="cursor-help" onpointerdown={(e) => onPopulationTooltip(e, row)}>
						{((row.spec?.populationDensity ?? 0) * 100).toFixed(1)}%
					</div>
				{:else if column.key == 'populationGrowth'}
					<div class="cursor-help" onpointerdown={(e) => onPopulationTooltip(e, row)}>
						{getGrowth(row).toLocaleString()}
					</div>
				{:else if column.key == 'habitability'}
					{#if row.spec?.canTerraform}
						<div class="cursor-help" onpointerdown={(e) => onPopulationTooltip(e, row)}>
							<span
								class:text-habitable={(row.spec?.habitability ?? 0) > 0}
								class:text-uninhabitable={(row.spec?.habitability ?? 0) < 0}
								>{row.spec?.habitability ?? 0}%</span
							>
							/ <span class="text-terraformable">{row.spec?.terraformedHabitability ?? 0}%</span>
						</div>
					{:else}
						<span
							class="cursor-help"
							onpointerdown={(e) => onPopulationTooltip(e, row)}
							class:text-habitable={(row.spec?.habitability ?? 0) > 0}
							class:text-uninhabitable={(row.spec?.habitability ?? 0) < 0}
						>
							{row.spec?.habitability ?? 0}%</span
						>
					{/if}
				{:else if column.key == 'production'}
					<button
						onclick={() => onProductionQueueDialog(planet)}
						class="text-base w-32 flex justify-between text-left cursor-pointer"
					>
						{#if planet.planetOrders?.productionQueue.length}
							<ProductionQueueItemLine
								item={planet.planetOrders.productionQueue[0]}
								index={0}
								shortName={true}
							/>
						{:else if ownedBy(row, $player.num)}
							-- Queue is Empty --
						{/if}
					</button>
				{:else if column.key == 'mines'}
					<span class="cursor-help" onpointerdown={(e) => onMinesTooltip(e, planet)}>
						{planet.mines} / {planet.spec?.maxMines ?? 0}</span
					>
				{:else if column.key == 'factories'}
					<span class="cursor-help" onpointerdown={(e) => onFactoriesTooltip(e, planet)}>
						{planet.factories}/ {planet.spec?.maxFactories ?? 0}
					</span>
				{:else if column.key == 'defense'}
					{#if row.spec?.defenseCoverage}
						<span
							class:cursor-help={planet.mapObject?.playerNum ?? 0 === $player.num}
							onpointerdown={(e) =>
								planet.mapObject?.playerNum === $player.num && onDefenseTooltip(e, planet)}
							>{((row.spec?.defenseCoverage ?? 0) * 100).toFixed(1)}%
						</span>
					{:else}
						none
					{/if}
				{:else if column.key == 'minerals'}
					<MineralMini mineral={row.cargo} {planet} />
				{:else if column.key == 'miningRate'}
					<MineralMini mineral={row.spec?.miningOutput} {planet} />
				{:else if column.key == 'mineralConcentration'}
					<MineralMini mineral={row.mineralConcentration} {planet} />
				{:else if column.key == 'resources'}
					{row.spec?.resourcesPerYearAvailable ?? 0} / {row.spec?.resourcesPerYear ?? 0}
				{:else if column.key == 'contributesOnlyLeftoverToResearch'}
					{#if planet.planetOrders?.contributesOnlyLeftoverToResearch}
						<Icon src={Check} size="24" class="stroke-success" />
					{/if}
				{:else if column.key == 'driverDest'}
					{@const targetPlanet = planet.planetOrders?.packetTargetNum
						? $universe.getPlanet(planet.planetOrders.packetTargetNum)
						: undefined}
					{#if targetPlanet}
						<button class="cs-link text-xl text-left" onclick={() => gotoMapObject(targetPlanet)}
							>{targetPlanet.mapObject?.name ?? ''}</button
						>
					{:else}
						--
					{/if}
				{:else if column.key == 'routingDestination'}
					{@const routeTarget = planet.planetOrders?.routeTargetNum
						? $universe.getMapObject(
								create(MapObjectTargetSchema, {
									targetPosition: emptyVector(),
									targetType: planet.planetOrders.routeTargetType,
									targetNum: planet.planetOrders.routeTargetNum,
									targetPlayerNum: planet.planetOrders.routeTargetPlayerNum
								})
							)
						: undefined}
					{#if routeTarget}
						<button class="cs-link text-xl text-left" onclick={() => gotoMapObject(routeTarget)}
							>{routeTarget.mapObject?.name}</button
						>
					{:else}
						--
					{/if}
				{:else}
					{cell}
				{/if}
			</span>
		{/snippet}
	</Table>
</div>

<ProductionQueueDialog
	show={showProductionQueueDialog}
	onNext={() => onNextPlanet(true)}
	onPrev={() => onPrevPlanet(true)}
	onOk={(planet) => {
		showProductionQueueDialog = false;
		updatePlanetOrders(planet);
	}}
	onCancel={() => (showProductionQueueDialog = false)}
/>

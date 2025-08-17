<script lang="ts">
	import { goto } from '$app/navigation';
	import CargoBar from '$lib/components/game/CargoBar.svelte';
	import FuelBar from '$lib/components/game/FuelBar.svelte';
	import SortableTableHeader from '$lib/components/table/SortableTableHeader.svelte';
	import Table, { type TableColumn } from '$lib/components/table/Table.svelte';
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import type { AnyFleet } from '$lib/services/Universe';
	import { WaypointTask, type Fleet } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import { fleetsSortBy, getEta, getLocation } from '$lib/types/Fleet';
	import { getMapObjectName } from '$lib/types/MapObject';

	const {
		game,
		player,
		universe,
		settings,
		commandMapObject,
		selectMapObject,
		zoomToMapObject,
		selectWaypoint
	} = getGameContext();

	const selectFleet = (fleet: Fleet) => {
		commandMapObject(fleet);
		zoomToMapObject(fleet);
		goto(`/games/${$game.id}`);
	};

	const selectTarget = (fleet: Fleet, waypointIndex: number) => {
		if (
			!fleet.fleetOrders ||
			waypointIndex < 0 ||
			waypointIndex > (fleet.fleetOrders?.waypoints?.length ?? 0) - 1
		) {
			return;
		}
		const target = $universe.getMapObject(
			fleet.fleetOrders?.waypoints[waypointIndex].mapObjectTarget
		);
		if (target) {
			selectMapObject(target);
			zoomToMapObject(target);
		}

		selectWaypoint(fleet.fleetOrders.waypoints[waypointIndex]);

		commandMapObject(fleet);
		goto(`/games/${$game.id}`);
	};

	// filterable fleets
	let search = $state('');
	let filteredFleets = $derived(
		$universe
			.getMyFleets($settings.sortFleetsKey, $settings.sortFleetsDescending)
			.filter((i) => i.mapObject?.name.toLowerCase().indexOf(search.toLowerCase()) != -1) ?? []
	);

	type TableFleet = AnyFleet & {
		name?: never;
		num?: never;
		battlePlanNum?: never;
		location?: never;
		destination?: never;
		task?: never;
		eta?: never;
		composition?: never;
		cloak?: never;
		mass?: never;
	};
	const columns: TableColumn<TableFleet>[] = [
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
			key: 'task',
			title: 'Task',
			sortBy: fleetsSortBy('task', $universe)
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

	function onSorted(column: TableColumn<TableFleet>, sortDescending: boolean) {
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
			table: 'table table-zebra table-compact table-auto w-full',
			th: 'sticky top-0 bg-base-200 z-10'
		}}
	>
		{#snippet head({ column })}
			<span>
				<SortableTableHeader
					{column}
					isSorted={$settings.sortFleetsKey === column.key}
					sortDescending={$settings.sortFleetsDescending}
					{onSorted}
				/>
			</span>
		{/snippet}

		{#snippet cell({ column, row, cell })}
			<span>
				{#if column.key == 'name'}
					<button class="cs-link text-xl text-left" onclick={() => selectFleet(row)}
						>{getMapObjectName(row)}</button
					>
				{:else if column.key == 'num'}
					{row.mapObject?.num ?? 0}
				{:else if column.key == 'location'}
					{@const location = getLocation(row, $universe)}
					<button class="cs-link text-xl text-left" onclick={() => selectTarget(row, 0)}
						>{location}</button
					>
				{:else if column.key == 'destination'}
					{@const targetName =
						row.fleetOrders?.waypoints && row.fleetOrders?.waypoints.length > 1
							? $universe.getTargetName(row.fleetOrders?.waypoints[1])
							: '--'}

					{#if targetName !== '--'}
						<button class="cs-link text-xl text-left" onclick={() => selectTarget(row, 1)}
							>{targetName}</button
						>
					{:else}
						--
					{/if}
				{:else if column.key == 'task'}
					{row.fleetOrders?.waypoints &&
					row.fleetOrders?.waypoints.length > 1 &&
					row.fleetOrders?.waypoints[1].task
						? enumToString(WaypointTask, row.fleetOrders?.waypoints[1].task)
						: '(no task here)'}
				{:else if column.key == 'eta'}
					{#if getEta(row) == -1}
						<span class="text-error"> Never </span>
					{:else if getEta(row) == 0}
						--
					{:else}
						{getEta(row)}y
					{/if}
				{:else if column.key == 'fuel'}
					<div class="w-32 leading-[1rem]">
						<FuelBar value={row.fuel} capacity={row.spec?.shipDesignSpec?.fuelCapacity} />
					</div>
				{:else if column.key == 'cargo'}
					<div class="w-32 leading-[1rem]">
						<CargoBar value={row.cargo} capacity={row.spec?.shipDesignSpec?.cargoCapacity} />
					</div>
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
							{#if row.tokens.length > 1}
								+
							{/if}
						</div>
					</div>
				{:else if column.key == 'cloak'}
					{row.spec && row.spec.shipDesignSpec?.cloakPercent
						? row.spec.shipDesignSpec?.cloakPercent + '%'
						: '--'}
				{:else if column.key == 'battlePlanNum'}
					{@const battlePlan = $game
						? $player.getBattlePlan(row.fleetOrders?.battlePlanNum ?? 0)
						: undefined}
					{battlePlan?.name ?? ''}
				{:else if column.key == 'mass'}
					{row.spec?.shipDesignSpec?.mass ?? 0}
				{:else}
					{cell}
				{/if}
			</span>
		{/snippet}
	</Table>
</div>

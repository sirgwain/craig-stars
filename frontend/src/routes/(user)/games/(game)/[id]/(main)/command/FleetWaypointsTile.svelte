<script lang="ts">
	import WarpSpeedGauge from '$lib/components/game/WarpSpeedGauge.svelte';
	import type {
		ChangeWaypointProps,
		DeleteWaypointProps,
		SelectWaypointProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { MapObjectType, WaypointSchema, type Waypoint } from '$lib/types/cs-proto';
	import type { CommandedFleet } from '$lib/types/Fleet';
	import { StargateWarpSpeed } from '$lib/types/Consts';
	import { distance } from '$lib/types/Vector';
	import CommandTile from './CommandTile.svelte';
	import { create } from '@bufbuild/protobuf';

	const { player, universe } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
		selectedWaypointIndex: number;
	} & ChangeWaypointProps &
		SelectWaypointProps &
		DeleteWaypointProps;

	let {
		fleet: propFleet,
		selectedWaypointIndex,
		onSelectWaypoint,
		onChangeWaypoint,
		onDeleteWaypoint
	}: Props = $props();

	// local state for the ui components
	let fleet = $state(propFleet);
	let waypoint: Waypoint = $derived(
		fleet.fleetOrders.waypoints[selectedWaypointIndex] ?? create(WaypointSchema)
	);

	$effect(() => {
		// update state when the props change
		fleet = propFleet;
		waypoint = propFleet.fleetOrders?.waypoints[selectedWaypointIndex];
	});

	let previousWaypoint: Waypoint | undefined = $derived.by(() => {
		if (selectedWaypointIndex > 0) {
			return fleet.fleetOrders?.waypoints[selectedWaypointIndex - 1];
		}
	});
	let nextWaypoint: Waypoint | undefined = $derived.by(() => {
		if (selectedWaypointIndex < fleet.fleetOrders?.waypoints.length) {
			return fleet.fleetOrders?.waypoints[selectedWaypointIndex + 1];
		}
	});

	let waypointPlanet = $derived(
		waypoint.mapObjectTarget?.targetType === MapObjectType.PLANET &&
			waypoint.mapObjectTarget?.targetNum
			? $universe.getPlanet(waypoint.mapObjectTarget?.targetNum)
			: undefined
	);
	let waypointPlanetFriendly = $derived(
		waypointPlanet && $player.isFriend(waypointPlanet.mapObject?.playerNum)
	);
	let dist = $derived(
		nextWaypoint || previousWaypoint
			? distance(
					waypoint?.position,
					previousWaypoint ? previousWaypoint?.position : nextWaypoint?.position
				)
			: 0
	);

	// calculate the fuel used per leg of each waypoint, starting at wp1
	let fuelUsagePerLeg = $derived(
		fleet.fleetOrders?.waypoints.slice(1).map((wp1) => wp1.estFuelUsage ?? 0)
	);

	// get the total fuel usage, but accounting for fueling stations
	let fuelUsageTotal = $derived(
		fuelUsagePerLeg.reduce(
			(total, wpUsage, i) =>
				fleet.fleetOrders?.waypoints.length < i + 1 &&
				fleet.fleetOrders?.waypoints[i + 1].mapObjectTarget?.targetType === MapObjectType.PLANET &&
				fleet.canFuel(
					$player,
					$universe.getPlanet(fleet.fleetOrders?.waypoints[i + 1].mapObjectTarget?.targetNum)
				)
					? 0
					: total + wpUsage,
			0
		)
	);

	// will we run out of fuel at any leg of our journey or the last leg that we are currently updating?
	let runOutOfFuel = $derived(fleet.willRunOutOfFuel($player, $universe));

	function onRepeatOrdersChanged(repeat: boolean) {
		fleet.fleetOrders.repeatOrders = repeat;
		onChangeWaypoint?.({ fleet, waypoint, waypointIndex: selectedWaypointIndex });
	}

	function onWarpSpeedChanged(speed: number) {
		waypoint.warpSpeed = speed;
		onChangeWaypoint?.({ fleet, waypoint, waypointIndex: selectedWaypointIndex });
	}

	function onWarpSpeedDragged(speed: number) {
		waypoint.warpSpeed = speed;
	}
</script>

{#if fleet.fleetOrders?.waypoints}
	<CommandTile title="Fleet Waypoints">
		<div class="bg-base-100 h-20 overflow-y-auto">
			<ul class="w-full h-full">
				{#each fleet.fleetOrders?.waypoints as wp, index (index)}
					<li class="pl-1 {selectedWaypointIndex == index ? 'bg-primary-focus' : ''}">
						<button
							type="button"
							class="text-left w-full h=full"
							onclick={() => onSelectWaypoint?.({ fleet, waypoint: wp })}
						>
							{$universe.getTargetName(wp)}
						</button>
					</li>
				{/each}
			</ul>
		</div>
		{#if previousWaypoint}
			<div class="flex justify-between my-1">
				<button
					name="deleteWaypoint"
					class="btn btn-outline btn-sm normal-case btn-secondary"
					onclick={() => {
						onDeleteWaypoint?.({ fleet, waypoint });
					}}
					>Delete
				</button>
			</div>

			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Coming From</span>
				<span>{$universe.getTargetName(previousWaypoint)}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Distance</span>
				<span>{`${dist.toFixed(1)}`} l.y.</span>
			</div>
			<div class="flex mt-1">
				<span class="text-tile-item-title">Warp Factor</span>
				<span class="flex-1 ml-1">
					{#if waypointPlanet && waypointPlanetFriendly && waypointPlanet.spec?.planetStarbaseSpec?.hasStargate}
						<WarpSpeedGauge
							onValueChanged={(value) => onWarpSpeedChanged(value)}
							onValueDragged={(value) => onWarpSpeedDragged(value)}
							value={waypoint.warpSpeed}
							warnSpeed={fleet.spec.shipDesignSpec?.engine?.maxSafeSpeed
								? fleet.spec.shipDesignSpec?.engine.maxSafeSpeed + 1
								: undefined}
							max={StargateWarpSpeed}
							useStargate={true}
						/>
					{:else}
						<WarpSpeedGauge
							onValueChanged={(value) => onWarpSpeedChanged(value)}
							onValueDragged={(value) => onWarpSpeedDragged(value)}
							warnSpeed={fleet.spec.shipDesignSpec?.engine?.maxSafeSpeed
								? fleet.spec.shipDesignSpec?.engine.maxSafeSpeed + 1
								: undefined}
							value={waypoint.warpSpeed}
						/>
					{/if}
				</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Travel Time</span>
				<span>
					{#if waypoint.warpSpeed === StargateWarpSpeed}
						1 year
					{:else if waypoint.warpSpeed === 0}
						Never
					{:else}
						{Math.ceil(Math.floor(dist) / (waypoint.warpSpeed * waypoint.warpSpeed))} years
					{/if}
				</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Leg Fuel Usage</span>
				<span>{fuelUsagePerLeg[selectedWaypointIndex - 1]}mg</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Total Fuel Usage</span>
				<span class:text-error={runOutOfFuel}>{fuelUsageTotal}mg</span>
			</div>

			<label>
				<input
					onchange={(e) => onRepeatOrdersChanged(e.currentTarget.checked ? true : false)}
					checked={fleet.fleetOrders.repeatOrders}
					class="checkbox-xs"
					type="checkbox"
				/> Repeat Orders
			</label>
		{:else if nextWaypoint}
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Going to</span>
				<span>{$universe.getTargetName(nextWaypoint)}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Distance</span>
				<span>{`${dist.toFixed(1)}`} l.y.</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Warp Factor</span>
				<span>{nextWaypoint.warpSpeed}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Travel Time</span>
				<span
					>{Math.ceil(Math.floor(dist) / (nextWaypoint.warpSpeed * nextWaypoint.warpSpeed))} years</span
				>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Total Fuel Usage</span>
				<span class:text-error={runOutOfFuel}>{fuelUsageTotal}mg</span>
			</div>
			<label>
				<input
					onchange={(e) => onRepeatOrdersChanged(e.currentTarget.checked ? true : false)}
					checked={fleet.fleetOrders.repeatOrders}
					class="checkbox-xs"
					type="checkbox"
				/> Repeat Orders
			</label>
		{/if}
	</CommandTile>
{/if}

<script lang="ts">
	import WarpSpeedGauge from '$lib/components/game/WarpSpeedGauge.svelte';
	import type {
		ChangeWaypointProps,
		DeleteWaypointProps,
		SelectWaypointProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import type { Waypoint } from '$lib/types/cs';
	import { MapObjectTypePlanet, StargateWarpSpeed } from '$lib/types/cs';
	import type { CommandedFleet } from '$lib/types/Fleet';
	import { distance } from '$lib/types/Vector';
	import CommandTile from './CommandTile.svelte';

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
	let waypoint = $state(propFleet.waypoints[selectedWaypointIndex]);

	$effect(() => {
		// update state when the props change
		fleet = propFleet;
		waypoint = propFleet.waypoints[selectedWaypointIndex];
	});

	let previousWaypoint: Waypoint | undefined = $derived.by(() => {
		if (selectedWaypointIndex > 0) {
			return fleet.waypoints[selectedWaypointIndex - 1];
		}
	});
	let nextWaypoint: Waypoint | undefined = $derived.by(() => {
		if (selectedWaypointIndex < fleet.waypoints.length) {
			return fleet.waypoints[selectedWaypointIndex + 1];
		}
	});

	let waypointPlanet = $derived(
		waypoint.targetType == MapObjectTypePlanet && waypoint.targetNum
			? $universe.getPlanet(waypoint.targetNum)
			: undefined
	);
	let waypointPlanetFriendly = $derived(
		waypointPlanet && $player.isFriend(waypointPlanet.playerNum)
	);
	let dist = $derived(
		nextWaypoint || previousWaypoint
			? distance(
					waypoint.position,
					previousWaypoint ? previousWaypoint.position : nextWaypoint?.position
				)
			: 0
	);

	// calculate the fuel used per leg of each waypoint, starting at wp1
	let fuelUsagePerLeg = $derived(
		fleet.waypoints.slice(1).map((wp1, index) =>
			fleet.getFuelCost(
				$universe,
				$player.race.spec?.fuelEfficiencyOffset ?? 0,
				// use the warp speed of the currently selected waypoint if we're dragging it around
				// otherwise use the waypoint from the fleet waypoints
				selectedWaypointIndex === index + 1 ? waypoint.warpSpeed : (wp1.warpSpeed ?? 0),
				distance(fleet.waypoints[index].position, wp1.position),
				fleet.spec.cargoCapacity ?? 0
			)
		)
	);

	// get the total fuel usage, but accounting for fueling stations
	let fuelUsageTotal = $derived(
		fuelUsagePerLeg.reduce(
			(total, wpUsage, i) =>
				fleet.waypoints[i + 1].targetType === MapObjectTypePlanet &&
				fleet.canFuel($player, $universe.getPlanet(fleet.waypoints[i + 1].targetNum ?? 0))
					? 0
					: total + wpUsage,
			0
		)
	);

	// will we run out of fuel at any leg of our journey or the last leg that we are currently updating?
	let runOutOfFuel = $derived(fleet.willRunOutOfFuel($player, $universe));

	function onRepeatOrdersChanged(repeat: boolean) {
		fleet.repeatOrders = repeat;
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

{#if fleet.waypoints}
	<CommandTile title="Fleet Waypoints">
		<div class="bg-base-100 h-20 overflow-y-auto">
			<ul class="w-full h-full">
				{#each fleet.waypoints as wp, index}
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
					{#if waypointPlanet && waypointPlanetFriendly && waypointPlanet.spec.hasStargate}
						<WarpSpeedGauge
							onValueChanged={(value) => onWarpSpeedChanged(value)}
							onValueDragged={(value) => onWarpSpeedDragged(value)}
							value={waypoint.warpSpeed}
							warnSpeed={fleet.spec.engine.maxSafeSpeed
								? fleet.spec.engine.maxSafeSpeed + 1
								: undefined}
							max={StargateWarpSpeed}
							useStargate={true}
						/>
					{:else}
						<WarpSpeedGauge
							onValueChanged={(value) => onWarpSpeedChanged(value)}
							onValueDragged={(value) => onWarpSpeedDragged(value)}
							warnSpeed={fleet.spec.engine.maxSafeSpeed
								? fleet.spec.engine.maxSafeSpeed + 1
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
					checked={fleet.repeatOrders}
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
					checked={fleet.repeatOrders}
					class="checkbox-xs"
					type="checkbox"
				/> Repeat Orders
			</label>
		{/if}
	</CommandTile>
{/if}

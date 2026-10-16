<script lang="ts">
	import WarpSpeedGauge from '#lib/components/game/WarpSpeedGauge.svelte';
	import type {
		ChangeWaypointProps,
		DeleteWaypointProps,
		SelectWaypointProps
	} from '#lib/services/Events.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { StargateWarpSpeed } from '#lib/types/Consts.js';
	import { MapObjectType, WaypointSchema, type Waypoint } from '#lib/types/cs-proto.js';
	import { CommandedFleet } from '#lib/types/Fleet.js';
	import { create } from '@bufbuild/protobuf';
	import CommandTile from './CommandTile.svelte';

	const { player, universe, cs } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
		selectedWaypointIndex: number;
	} & ChangeWaypointProps &
		SelectWaypointProps &
		DeleteWaypointProps;

	let {
		selectedWaypointIndex,
		onSelectWaypoint,
		onChangeWaypoint,
		onDeleteWaypoint,
		...rest
	}: Props = $props();

	// local state for the ui components
	let fleet = $derived(rest.fleet);
	let waypoint: Waypoint = $derived(
		fleet.fleetOrders.waypoints[selectedWaypointIndex] ?? create(WaypointSchema)
	);

	let previousWaypoint: Waypoint | undefined = $derived.by(() => {
		if (selectedWaypointIndex > 0) {
			return fleet.fleetOrders.waypoints[selectedWaypointIndex - 1];
		}
	});
	let nextWaypoint: Waypoint | undefined = $derived.by(() => {
		if (selectedWaypointIndex < fleet.fleetOrders.waypoints.length) {
			return fleet.fleetOrders.waypoints[selectedWaypointIndex + 1];
		}
	});

	// the planet at this waypoint, targeted or under a fleet or position we're going to
	let waypointPlanet = $derived(
		waypoint.mapObjectTarget?.targetType === MapObjectType.PLANET &&
			waypoint.mapObjectTarget.targetNum
			? $universe.getPlanet(waypoint.mapObjectTarget.targetNum)
			: waypoint.position && $universe.getPlanetAtPosition(waypoint.position)
	);
	let waypointPlanetFriendly = $derived(
		waypointPlanet && $player.isFriend(waypointPlanet.mapObject?.playerNum)
	);

	// calculate the fuel used per leg of each waypoint, starting at wp1
	let fuelUsagePerLeg = $derived(
		fleet.fleetOrders.waypoints.slice(1).map((wp1) => wp1.estFuelUsage)
	);

	// travel time and fuel are shown for the leg into the selected waypoint and for the whole trip up
	// to it, or to the next waypoint when the fleet's current position is selected
	let tripIndex = $derived(previousWaypoint ? selectedWaypointIndex : 1);
	let dist = $derived(fleet.getDistanceForLeg($universe, tripIndex));
	const formatYears = (years: number | undefined) =>
		years === undefined ? 'Never' : years === 1 ? '1 year' : `${years} years`;
	let legTravelTime = $derived(formatYears(fleet.getYearsForLeg($universe, tripIndex)));
	let totalTravelTime = $derived(formatYears(fleet.getYearsToWaypoint($universe, tripIndex)));

	// the most fuel the trip needs at once, which starts over wherever the fleet refuels
	let fuelUsage = $derived(fleet.getFuelUsageToWaypoint($player, $universe, tripIndex));
	let runOutOfFuel = $derived(fuelUsage > fleet.fuel);
	// the fleet runs dry before it finishes the leg to the selected waypoint
	let runOutOfFuelOnLeg = $derived(fleet.getFuelOnArrival($player, $universe, tripIndex) < 0);

	function onRepeatOrdersChanged(repeat: boolean) {
		fleet.fleetOrders.repeatOrders = repeat;
		onChangeWaypoint?.({ fleet, waypoint, waypointIndex: selectedWaypointIndex });
	}

	function onWarpSpeedChanged(speed: number) {
		waypoint.warpSpeed = speed;
		onChangeWaypoint?.({ fleet, waypoint, waypointIndex: selectedWaypointIndex });
	}

	async function onWarpSpeedDragged(speed: number) {
		waypoint.warpSpeed = speed;
		const result = await cs.wasmService.updateWaypointSpeed({
			fleet,
			waypointIndex: selectedWaypointIndex,
			warpSpeed: speed
		});
		if (!result.fleet) {
			return;
		}
		fleet = new CommandedFleet(result.fleet);
	}
</script>

{#if fleet.fleetOrders.waypoints}
	<CommandTile title="Fleet Waypoints">
		<div class="bg-base-100 h-20 overflow-y-auto">
			<ul class="size-full">
				{#each fleet.fleetOrders.waypoints as wp, index (index)}
					<li class="pl-1 {selectedWaypointIndex == index ? 'bg-primary-focus' : ''}">
						<button
							type="button"
							class="text-left size-full truncate"
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
				<span class="shrink-0 text-tile-item-title">Coming From</span>
				<span class="min-w-0 truncate pl-1">{$universe.getTargetName(previousWaypoint)}</span>
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
								? fleet.spec.shipDesignSpec.engine.maxSafeSpeed + 1
								: undefined}
							max={StargateWarpSpeed}
							useStargate={true}
						/>
					{:else}
						<WarpSpeedGauge
							onValueChanged={(value) => onWarpSpeedChanged(value)}
							onValueDragged={(value) => onWarpSpeedDragged(value)}
							warnSpeed={fleet.spec.shipDesignSpec?.engine?.maxSafeSpeed
								? fleet.spec.shipDesignSpec.engine.maxSafeSpeed + 1
								: undefined}
							value={waypoint.warpSpeed}
						/>
					{/if}
				</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Travel Time</span>
				<span>{legTravelTime}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Total Travel Time</span>
				<span>{totalTravelTime}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Leg Fuel Usage</span>
				<span class:text-error={runOutOfFuelOnLeg}
					>{fuelUsagePerLeg[selectedWaypointIndex - 1]}mg</span
				>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Est Fuel Usage</span>
				<span class:text-error={runOutOfFuel}>{fuelUsage}mg</span>
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
				<span class="shrink-0 text-tile-item-title">Going to</span>
				<span class="min-w-0 truncate pl-1">{$universe.getTargetName(nextWaypoint)}</span>
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
				<span>{legTravelTime}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Est Fuel Usage</span>
				<span class:text-error={runOutOfFuel}>{fuelUsage}mg</span>
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

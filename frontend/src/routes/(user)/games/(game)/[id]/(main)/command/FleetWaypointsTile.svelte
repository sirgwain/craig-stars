<script lang="ts">
	import WarpSpeedGauge from '$lib/components/game/WarpSpeedGauge.svelte';
	import type {
		ChangeWaypointProps,
		DeleteWaypointProps,
		SelectWaypointProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { StargateWarpSpeed } from '$lib/types/Constants';
	import type { CommandedFleet, Waypoint } from '$lib/types/Fleet';
	import { MapObjectType, type MapObject } from '$lib/types/MapObject';
	import { distance } from '$lib/types/Vector';
	import CommandTile from './CommandTile.svelte';

	const { player, universe } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
		selectedWaypoint: Waypoint | undefined;
	} & ChangeWaypointProps &
		SelectWaypointProps &
		DeleteWaypointProps;

	let { fleet, selectedWaypoint, onSelectWaypoint, onChangeWaypoint, onDeleteWaypoint }: Props =
		$props();

	// local state for the ui components
	let warpSpeed = $state(
		selectedWaypoint ? (selectedWaypoint.warpSpeed ?? 0) : (fleet.waypoints[0].warpSpeed ?? 0)
	);
	let repeatOrders = $state(fleet.repeatOrders);
	let waypointRefs: (HTMLLIElement | null)[] = $state([]);

	// if our selectedWaypoint or fleet changes, update the state
	$effect(() => {
		warpSpeed = selectedWaypoint
			? (selectedWaypoint.warpSpeed ?? warpSpeed)
			: (fleet.waypoints[0].warpSpeed ?? warpSpeed);
		repeatOrders = fleet.repeatOrders;
	});

	let selectedWaypointIndex = $derived.by(() => {
		const index = fleet.waypoints.findIndex((wp) => wp == selectedWaypoint);
		return index === -1 ? 0 : index;
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

	let selectedWaypointPlanet = $derived(
		selectedWaypoint?.targetType == MapObjectType.Planet && selectedWaypoint?.targetNum
			? $universe.getPlanet(selectedWaypoint?.targetNum)
			: undefined
	);
	let selectedWaypointPlanetFriendly = $derived(
		selectedWaypointPlanet && $player.isFriend(selectedWaypointPlanet.playerNum)
	);
	let dist = $derived(
		selectedWaypoint && (nextWaypoint || previousWaypoint)
			? distance(
					selectedWaypoint.position,
					previousWaypoint ? previousWaypoint.position : nextWaypoint?.position
				)
			: 0
	);

	// calculate the fuel used per leg of each waypoint, starting at wp1
	let fuelUsagePerLeg = $derived(
		fleet.waypoints
			.slice(1)
			.map((wp1, index) =>
				fleet.getFuelCost(
					$universe,
					$player.race.spec?.fuelEfficiencyOffset ?? 0,
					selectedWaypoint === wp1 ? warpSpeed : (wp1.warpSpeed ?? 0),
					distance(fleet.waypoints[index].position, wp1.position),
					fleet.spec.cargoCapacity ?? 0
				)
			)
	);

	// get the total fuel usage, but accounting for fueling stations
	let fuelUsageTotal = $derived(
		fuelUsagePerLeg.reduce(
			(total, wpUsage, i) =>
				fleet.waypoints[i + 1].targetType === MapObjectType.Planet &&
				fleet.canFuel($player, $universe.getPlanet(fleet.waypoints[i + 1].targetNum ?? 0))
					? 0
					: total + wpUsage,
			0
		)
	);

	// will we run out of fuel at any leg of our journey or the last leg that we are currently updating?
	let runOutOfFuel = $derived(fleet.willRunOutOfFuel($player, $universe));

	function getWaypointTarget(wp: Waypoint): MapObject | undefined {
		if (wp && wp.targetType && wp.targetNum) {
			return $universe.getMapObject(wp);
		}
	}

	function onRepeatOrdersChanged(repeat: boolean) {
		if (selectedWaypoint) {
			repeatOrders = repeat;
			fleet.repeatOrders = repeat;
			onChangeWaypoint?.({ fleet, waypoint: selectedWaypoint });
		}
	}

	function onWarpSpeedChanged(speed: number) {
		if (selectedWaypoint) {
			warpSpeed = speed;
			selectedWaypoint.warpSpeed = speed;
			onChangeWaypoint?.({ fleet, waypoint: selectedWaypoint });
		}
	}

	function onWarpSpeedDragged(speed: number) {
		if (selectedWaypoint) {
			warpSpeed = speed;
			selectedWaypoint.warpSpeed = speed;
		}
	}
</script>

{#if fleet.waypoints && selectedWaypoint}
	<CommandTile title="Fleet Waypoints">
		<div class="bg-base-100 h-20 overflow-y-auto">
			<ul class="w-full h-full">
				{#each fleet.waypoints as wp, index}
					<li
						bind:this={waypointRefs[index]}
						class="pl-1 {selectedWaypointIndex == index ? 'bg-primary-focus' : ''}"
					>
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
					onclick={(e) => {
						e.preventDefault();
						onDeleteWaypoint?.({ fleet, waypoint: selectedWaypoint });
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
					{#if selectedWaypointPlanet && selectedWaypointPlanetFriendly && selectedWaypointPlanet.spec.hasStargate}
						<WarpSpeedGauge
							onvaluechanged={(value) => onWarpSpeedChanged(value)}
							onvaluedragged={(value) => onWarpSpeedDragged(value)}
							bind:value={warpSpeed}
							warnSpeed={fleet.spec.engine.maxSafeSpeed
								? fleet.spec.engine.maxSafeSpeed + 1
								: undefined}
							max={StargateWarpSpeed}
							useStargate={true}
						/>
					{:else}
						<WarpSpeedGauge
							onvaluechanged={(value) => onWarpSpeedChanged(value)}
							onvaluedragged={(value) => onWarpSpeedDragged(value)}
							warnSpeed={fleet.spec.engine.maxSafeSpeed
								? fleet.spec.engine.maxSafeSpeed + 1
								: undefined}
							bind:value={warpSpeed}
						/>
					{/if}
				</span>
			</div>
			<div class="flex justify-between mt-1">
				<span class="text-tile-item-title">Travel Time</span>
				<span>
					{#if warpSpeed === StargateWarpSpeed}
						1 year
					{:else if warpSpeed === 0}
						Never
					{:else}
						{Math.ceil(Math.floor(dist) / (warpSpeed * warpSpeed))} years
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
					bind:checked={repeatOrders}
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
					checked={repeatOrders}
					class="checkbox-xs"
					type="checkbox"
				/> Repeat Orders
			</label>
		{/if}
	</CommandTile>
{/if}

<script lang="ts">
	import {
		commandedMapObjectName,
		getMapObject,
		selectedWaypoint,
		selectMapObject,
		selectWaypoint
	} from '$lib/services/Context';
	import { FleetService } from '$lib/services/FleetService';
	import type { Fleet, Waypoint } from '$lib/types/Fleet';
	import type { MapObject } from '$lib/types/MapObject';
	import type { Player } from '$lib/types/Player';
	import { distance } from '$lib/types/Vector';
	import hotkeys from 'hotkeys-js';
	import { merge } from 'lodash-es';
	import WarpFactorBar from '$lib/components/game/WarpFactorBar.svelte';
	import CommandTile from './CommandTile.svelte';

	export let fleet: Fleet;

	const fleetService = new FleetService();
	let selectedWaypointIndex = 0;
	let previousWaypoint: Waypoint | undefined;
	let previousWaypointMO: MapObject | undefined;
	let nextWaypoint: Waypoint | undefined;
	let nextWaypointMO: MapObject | undefined;

	const getTargetName = (wp: Waypoint) =>
		wp.targetName ?? `Space: (${wp.position.x}, ${wp.position.y})`;

	const getWaypointTarget = (wp: Waypoint): MapObject | undefined => {
		if (wp && wp.targetType && wp.targetNum) {
			return getMapObject(wp.targetType, wp.targetNum, wp.targetPlayerNum);
		}
	};

	const updateNextPrevWaypoints = () => {
		// find the next/previous waypoint
		previousWaypoint = previousWaypointMO = nextWaypoint = nextWaypointMO = undefined;
		if (fleet.waypoints) {
			if (selectedWaypointIndex > 0) {
				previousWaypoint = fleet.waypoints[selectedWaypointIndex - 1];
				previousWaypointMO = getWaypointTarget(previousWaypoint);
			}
			if (selectedWaypointIndex < fleet.waypoints.length) {
				nextWaypoint = fleet.waypoints[selectedWaypointIndex + 1];
				nextWaypointMO = getWaypointTarget(nextWaypoint);
			}
		}
	};

	const onSelectWaypoint = (wp: Waypoint, index: number) => {
		selectedWaypointIndex = index;
		selectWaypoint(wp);
		const mo = getWaypointTarget(wp);
		if (mo) {
			selectMapObject(mo);
		}

		updateNextPrevWaypoints();
	};

	$: dist =
		$selectedWaypoint && (nextWaypoint || previousWaypoint)
			? distance(
					$selectedWaypoint.position,
					previousWaypoint ? previousWaypoint.position : nextWaypoint?.position
			  )
			: 0;

	const onRepeatOrdersChanged = async (repeatOrders: boolean) => {
		if ($selectedWaypoint) {
			fleet.repeatOrders = repeatOrders;
			const f = await FleetService.updateFleetOrders(fleet);

			// update the player fleet
			merge(fleet, f);

			// update the commanded object
			updateNextPrevWaypoints();
		}
	};

	const onWarpFactorChanged = async (warpFactor: number) => {
		if (fleet && $selectedWaypoint) {
			$selectedWaypoint.warpFactor = warpFactor;
			const f = await FleetService.updateFleetOrders(fleet);

			// update the player fleet
			merge(fleet, f);

			// update the commanded object
			updateNextPrevWaypoints();
		}
	};

	const deleteWaypoint = () => {
		if (selectedWaypointIndex != 0 && fleet.waypoints) {
			fleet.waypoints = fleet.waypoints?.filter((wp) => wp != $selectedWaypoint);
			selectedWaypointIndex--;

			onSelectWaypoint(fleet.waypoints[selectedWaypointIndex], selectedWaypointIndex);

			FleetService.updateFleetOrders(fleet).then((fleet) => {
				// update the player fleet
				merge(fleet, fleet);

				updateNextPrevWaypoints();
			});
		}
	};

	const onNextWaypoint = () => {
		if (fleet.waypoints && selectedWaypointIndex + 1 < fleet.waypoints.length) {
			onSelectWaypoint(fleet.waypoints[selectedWaypointIndex + 1], selectedWaypointIndex + 1);
		}
	};

	const onPrevWaypoint = () => {
		if (fleet.waypoints && selectedWaypointIndex > 0) {
			onSelectWaypoint(fleet.waypoints[selectedWaypointIndex - 1], selectedWaypointIndex - 1);
		}
	};

	commandedMapObjectName.subscribe(() => {
		selectedWaypointIndex = 0;
		updateNextPrevWaypoints();
	});

	selectedWaypoint?.subscribe(() => {
		if (fleet.waypoints) {
			selectedWaypointIndex = fleet.waypoints.findIndex((wp) => wp == $selectedWaypoint);
			if (selectedWaypointIndex == -1) {
				selectedWaypointIndex = 0;
			}
			updateNextPrevWaypoints();

			// if (waypointRefs.length > selectedWaypointIndex) {
			// 	// TODO: this is making small screens jump by scrolling
			// 	// to the waypoint
			// 	// waypointRefs[selectedWaypointIndex]?.scrollIntoView();
			// }
		}
	});

	let waypointRefs: (HTMLLIElement | null)[] = [];

	hotkeys('Delete', () => deleteWaypoint());
	hotkeys('Backspace', () => deleteWaypoint());
	hotkeys('down', () => onNextWaypoint());
	hotkeys('up', () => onPrevWaypoint());
</script>

{#if fleet.waypoints && $selectedWaypoint}
	<CommandTile title="Fleet Waypoints">
		<div class="bg-base-100 h-20 overflow-y-auto">
			<ul class="w-full h-full">
				{#each fleet.waypoints as wp, index}
					<li
						on:click={() => onSelectWaypoint(wp, index)}
						bind:this={waypointRefs[index]}
						class="pl-1 {selectedWaypointIndex == index ? 'bg-primary-focus' : ''}"
					>
						{getTargetName(wp)}
					</li>
				{/each}
			</ul>
		</div>
		{#if previousWaypoint}
			<div class="flex justify-between my-1">
				<button
					name="deleteWaypoint"
					class="btn btn-outline btn-sm normal-case btn-secondary"
					on:click={deleteWaypoint}
					>Delete
				</button>
			</div>

			<div class="flex justify-between mt-1">
				<span>Coming From</span>
				<span>{getTargetName(previousWaypoint)}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Distance</span>
				<span>{`${dist.toFixed(1)}`} l.y.</span>
			</div>
			<div class="flex mt-1">
				<span>Warp Factor</span>
				<span class="flex-1 ml-1"
					><WarpFactorBar
						on:valuechanged={(e) => onWarpFactorChanged(e.detail)}
						bind:value={$selectedWaypoint.warpFactor}
					/></span
				>
			</div>
			<div class="flex justify-between mt-1">
				<span>Travel Time</span>
				<span
					>{Math.ceil(dist / ($selectedWaypoint.warpFactor * $selectedWaypoint.warpFactor))} years</span
				>
			</div>
			<div class="flex justify-between mt-1">
				<span>Est Fuel Usage</span>
				<span>{$selectedWaypoint.estFuelUsage ?? 0}mg</span>
			</div>
			<label>
				<input
					on:change={(e) => onRepeatOrdersChanged(e.currentTarget.checked ? true : false)}
					bind:checked={fleet.repeatOrders}
					class="checkbox-xs"
					type="checkbox"
				/> Repeate Orders
			</label>
		{:else if nextWaypoint}
			<div class="flex justify-between mt-1">
				<span>Going to</span>
				<span>{getTargetName(nextWaypoint)}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Distance</span>
				<span>{`${dist.toFixed(1)}`} l.y.</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Warp Factor</span>
				<span>{nextWaypoint.warpFactor}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Travel Time</span>
				<span>{Math.ceil(dist / (nextWaypoint.warpFactor * nextWaypoint.warpFactor))} years</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Est Fuel Usage</span>
				<span>{nextWaypoint.estFuelUsage ?? 0}mg</span>
			</div>
			<label>
				<input
					on:change={(e) => onRepeatOrdersChanged(e.currentTarget.checked ? true : false)}
					checked={fleet.repeatOrders}
					class="checkbox-xs"
					type="checkbox"
				/> Repeate Orders
			</label>
		{/if}
	</CommandTile>
{/if}

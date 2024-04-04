<script lang="ts">
	import WarpSpeedGauge from '$lib/components/game/WarpSpeedGauge.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import type { CommandedFleet, Waypoint } from '$lib/types/Fleet';
	import { MapObjectType, type MapObject, StargateWarpSpeed } from '$lib/types/MapObject';
	import { distance } from '$lib/types/Vector';
	import hotkeys from 'hotkeys-js';
	import { onDestroy, onMount } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	const {
		player,
		universe,
		commandedMapObjectKey,
		selectedWaypoint,
		selectMapObject,
		selectWaypoint,
		updateFleetOrders
	} = getGameContext();

	export let fleet: CommandedFleet;

	let selectedWaypointIndex = 0;
	let previousWaypoint: Waypoint | undefined;
	let previousWaypointMO: MapObject | undefined;
	let nextWaypoint: Waypoint | undefined;
	let nextWaypointMO: MapObject | undefined;
	let waypointRefs: (HTMLLIElement | null)[] = [];

	$: selectedWaypointPlanet =
		$selectedWaypoint?.targetType == MapObjectType.Planet && $selectedWaypoint?.targetNum
			? $universe.getPlanet($selectedWaypoint?.targetNum)
			: undefined;
	$: selectedWaypointPlanetFriendly =
		selectedWaypointPlanet && $player.isFriend(selectedWaypointPlanet.playerNum);

	function getWaypointTarget(wp: Waypoint): MapObject | undefined {
		if (wp && wp.targetType && wp.targetNum) {
			return $universe.getMapObject(wp);
		}
	}

	function updateNextPrevWaypoints() {
		// find the next/previous waypoint
		previousWaypoint = previousWaypointMO = nextWaypoint = nextWaypointMO = undefined;
		if (selectedWaypointIndex > 0) {
			previousWaypoint = fleet.waypoints[selectedWaypointIndex - 1];
			previousWaypointMO = getWaypointTarget(previousWaypoint);
		}
		if (selectedWaypointIndex < fleet.waypoints.length) {
			nextWaypoint = fleet.waypoints[selectedWaypointIndex + 1];
			nextWaypointMO = getWaypointTarget(nextWaypoint);
		}
	}

	function onSelectWaypoint(wp: Waypoint, index: number) {
		selectedWaypointIndex = index;
		selectWaypoint(wp);
		const mo = getWaypointTarget(wp);
		if (mo) {
			selectMapObject(mo);
		}

		updateNextPrevWaypoints();
	}

	$: dist =
		$selectedWaypoint && (nextWaypoint || previousWaypoint)
			? distance(
					$selectedWaypoint.position,
					previousWaypoint ? previousWaypoint.position : nextWaypoint?.position
				)
			: 0;

	// calculate the fuel used per leg of each waypoint, starting at wp1
	$: fuelUsagePerLeg = fleet.waypoints
		.slice(1)
		.map((wp1, index) =>
			fleet.getFuelCost(
				$universe,
				$player.race.spec?.fuelEfficiencyOffset ?? 0,
				$selectedWaypoint === wp1 ? $selectedWaypoint.warpSpeed : wp1.warpSpeed ?? 0,
				distance(fleet.waypoints[index].position, wp1.position),
				fleet.spec.cargoCapacity ?? 0
			)
		);

	$: fuelUsageToSelectedWaypoint = fuelUsagePerLeg.reduce((total, wpUsage) => total + wpUsage, 0);
	$: fuelUsageTotal = fuelUsagePerLeg.reduce((total, wpUsage) => total + wpUsage, 0);

	async function onRepeatOrdersChanged(repeatOrders: boolean) {
		if ($selectedWaypoint) {
			fleet.repeatOrders = repeatOrders;
			await updateFleetOrders(fleet);

			// update the commanded object
			updateNextPrevWaypoints();
		}
	}

	async function onWarpSpeedChanged(warpSpeed: number) {
		if ($selectedWaypoint) {
			$selectedWaypoint.warpSpeed = warpSpeed;
			await updateFleetOrders(fleet);

			// update the commanded object
			updateNextPrevWaypoints();
		}
	}

	async function onWarpSpeedDragged(warpSpeed: number) {
		if ($selectedWaypoint) {
			$selectedWaypoint.warpSpeed = warpSpeed;
		}
	}

	async function deleteWaypoint() {
		if (selectedWaypointIndex != 0 && fleet.waypoints) {
			fleet.waypoints = fleet.waypoints.filter((wp) => wp != $selectedWaypoint);
			
			// select the previous waypoint
			selectedWaypointIndex--;
			onSelectWaypoint(fleet.waypoints[selectedWaypointIndex], selectedWaypointIndex);

			await updateFleetOrders(fleet).then(() => updateNextPrevWaypoints());
		}
	}

	function onNextWaypoint() {
		if (selectedWaypointIndex + 1 < fleet.waypoints.length) {
			onSelectWaypoint(fleet.waypoints[selectedWaypointIndex + 1], selectedWaypointIndex + 1);
		}
	}

	function onPrevWaypoint() {
		if (selectedWaypointIndex > 0) {
			onSelectWaypoint(fleet.waypoints[selectedWaypointIndex - 1], selectedWaypointIndex - 1);
		}
	}

	onMount(() => {
		// TODO: these hotkeys can't be on the component... they are wired up twice because we render the command pane twice
		hotkeys('Delete', 'root', () => {
			deleteWaypoint();
		});
		hotkeys('Backspace', 'root', () => {
			deleteWaypoint();
		});
		// hotkeys('down', () => onNextWaypoint());
		// hotkeys('up', () => onPrevWaypoint());

		const unsubscribeSelectedWaypoint = selectedWaypoint?.subscribe(() => {
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
		});

		// reset the waypoint index every time the commanded mapobject changes
		const unsubscribeCommandedMapObject = commandedMapObjectKey.subscribe(() => {
			selectedWaypointIndex = 0;
			updateNextPrevWaypoints();
		});

		return () => {
			unsubscribeCommandedMapObject();
			unsubscribeSelectedWaypoint();

			hotkeys.unbind('Delete', 'root');
			hotkeys.unbind('Backspace', 'root');
			// hotkeys.unbind('down');
			// hotkeys.unbind('up');
		};
	});
</script>

{#if fleet.waypoints && $selectedWaypoint}
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
							on:click={() => onSelectWaypoint(wp, index)}
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
					on:click={deleteWaypoint}
					>Delete
				</button>
			</div>

			<div class="flex justify-between mt-1">
				<span>Coming From</span>
				<span>{$universe.getTargetName(previousWaypoint)}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Distance</span>
				<span>{`${dist.toFixed(1)}`} l.y.</span>
			</div>
			<div class="flex mt-1">
				<span>Warp Factor</span>
				<span class="flex-1 ml-1">
					{#if selectedWaypointPlanet && selectedWaypointPlanetFriendly && selectedWaypointPlanet.spec.hasStargate}
						<WarpSpeedGauge
							on:valuechanged={(e) => onWarpSpeedChanged(e.detail)}
							on:valuedragged={(e) => onWarpSpeedDragged(e.detail)}
							bind:value={$selectedWaypoint.warpSpeed}
							max={StargateWarpSpeed}
							useStargate={true}
						/>
					{:else}
						<WarpSpeedGauge
							on:valuechanged={(e) => onWarpSpeedChanged(e.detail)}
							on:valuedragged={(e) => onWarpSpeedDragged(e.detail)}
							bind:value={$selectedWaypoint.warpSpeed}
						/>
					{/if}
				</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Travel Time</span>
				<span
					>{Math.ceil(dist / ($selectedWaypoint.warpSpeed * $selectedWaypoint.warpSpeed))} years</span
				>
			</div>
			<div class="flex justify-between mt-1">
				<span>Leg Fuel Usage</span>
				<span>{fuelUsagePerLeg[selectedWaypointIndex - 1]}mg</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Total Fuel Usage</span>
				<span class:text-error={fuelUsageTotal > fleet.fuel}>{fuelUsageTotal}mg</span>
			</div>

			<label>
				<input
					on:change={(e) => onRepeatOrdersChanged(e.currentTarget.checked ? true : false)}
					bind:checked={fleet.repeatOrders}
					class="checkbox-xs"
					type="checkbox"
				/> Repeat Orders
			</label>
		{:else if nextWaypoint}
			<div class="flex justify-between mt-1">
				<span>Going to</span>
				<span>{$universe.getTargetName(nextWaypoint)}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Distance</span>
				<span>{`${dist.toFixed(1)}`} l.y.</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Warp Factor</span>
				<span>{nextWaypoint.warpSpeed}</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Travel Time</span>
				<span>{Math.ceil(dist / (nextWaypoint.warpSpeed * nextWaypoint.warpSpeed))} years</span>
			</div>
			<div class="flex justify-between mt-1">
				<span>Total Fuel Usage</span>
				<span class:text-error={fuelUsageTotal > fleet.fuel}>{fuelUsageTotal}mg</span>
			</div>
			<label>
				<input
					on:change={(e) => onRepeatOrdersChanged(e.currentTarget.checked ? true : false)}
					checked={fleet.repeatOrders}
					class="checkbox-xs"
					type="checkbox"
				/> Repeat Orders
			</label>
		{/if}
	</CommandTile>
{/if}

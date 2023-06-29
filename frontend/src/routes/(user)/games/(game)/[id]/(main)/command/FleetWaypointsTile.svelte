<script lang="ts">
	import WarpSpeedGauge from '$lib/components/game/WarpSpeedGauge.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import {
		commandedMapObjectName,
		selectedWaypoint,
		selectMapObject,
		selectWaypoint
	} from '$lib/services/Stores';
	import type { CommandedFleet, Waypoint } from '$lib/types/Fleet';
	import { MapObjectType, type MapObject, StargateWarpSpeed } from '$lib/types/MapObject';
	import { distance } from '$lib/types/Vector';
	import hotkeys from 'hotkeys-js';
	import { onMount } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	const { game, player, universe } = getGameContext();

	export let fleet: CommandedFleet;

	let selectedWaypointIndex = 0;
	let previousWaypoint: Waypoint | undefined;
	let previousWaypointMO: MapObject | undefined;
	let nextWaypoint: Waypoint | undefined;
	let nextWaypointMO: MapObject | undefined;

	$: selectedWaypointPlanet =
		$selectedWaypoint?.targetType == MapObjectType.Planet && $selectedWaypoint?.targetNum
			? $universe.getPlanet($selectedWaypoint?.targetNum)
			: undefined;
	$: selectedWaypointPlanetFriendly =
		selectedWaypointPlanet && $player.isFriend(selectedWaypointPlanet.playerNum);

	const getWaypointTarget = (wp: Waypoint): MapObject | undefined => {
		if (wp && wp.targetType && wp.targetNum) {
			return $universe.getMapObject(wp);
		}
	};

	const updateNextPrevWaypoints = () => {
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
			await $game.updateFleetOrders(fleet);

			// update the commanded object
			updateNextPrevWaypoints();
		}
	};

	const onWarpSpeedChanged = async (warpSpeed: number) => {
		if ($selectedWaypoint) {
			$selectedWaypoint.warpSpeed = warpSpeed;
			await $game.updateFleetOrders(fleet);

			// update the commanded object
			updateNextPrevWaypoints();
		}
	};

	const deleteWaypoint = async () => {
		if (selectedWaypointIndex != 0 && fleet.waypoints) {
			fleet.waypoints = fleet.waypoints.filter((wp) => wp != $selectedWaypoint);
			selectedWaypointIndex--;

			await $game.updateFleetOrders(fleet).then(() => updateNextPrevWaypoints());
			onSelectWaypoint(fleet.waypoints[selectedWaypointIndex], selectedWaypointIndex);
		}
	};

	const onNextWaypoint = () => {
		if (selectedWaypointIndex + 1 < fleet.waypoints.length) {
			onSelectWaypoint(fleet.waypoints[selectedWaypointIndex + 1], selectedWaypointIndex + 1);
		}
	};

	const onPrevWaypoint = () => {
		if (selectedWaypointIndex > 0) {
			onSelectWaypoint(fleet.waypoints[selectedWaypointIndex - 1], selectedWaypointIndex - 1);
		}
	};

	commandedMapObjectName.subscribe(() => {
		selectedWaypointIndex = 0;
		updateNextPrevWaypoints();
	});

	selectedWaypoint?.subscribe(() => {
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

	let waypointRefs: (HTMLLIElement | null)[] = [];

	onMount(() => {
		// TODO: these hotkeys can't be on the component... they are wired up twice because we render the command pane twice
		hotkeys('Delete', () => {
			deleteWaypoint();
		});
		hotkeys('Backspace', () => {
			deleteWaypoint();
		});
		hotkeys('down', () => onNextWaypoint());
		hotkeys('up', () => onPrevWaypoint());

		return () => {
			hotkeys.unbind('Delete');
			hotkeys.unbind('Backspace');
			hotkeys.unbind('down');
			hotkeys.unbind('up');
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
							bind:value={$selectedWaypoint.warpSpeed}
							max={StargateWarpSpeed}
							useStargate={true}
						/>
					{:else}
						<WarpSpeedGauge
							on:valuechanged={(e) => onWarpSpeedChanged(e.detail)}
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
				<span>Est Fuel Usage</span>
				<span class:text-error={($selectedWaypoint.estFuelUsage ?? 0) > fleet.spec.fuelCapacity}
					>{$selectedWaypoint.estFuelUsage ?? 0}mg</span
				>
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
				<span>Est Fuel Usage</span>
				<span>{nextWaypoint.estFuelUsage ?? 0}mg</span>
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

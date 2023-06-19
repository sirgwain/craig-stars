<script lang="ts">
	import { commandedFleet, player } from '$lib/services/Context';

	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	type WaypointLine = {
		path: string;
		props: any;
	};

	let lines: WaypointLine[] = [];

	$: {
		if ($data && $player) {
			const fleets = $player.fleets.filter((fleet) => fleet.waypoints.length > 1);
			lines = fleets.map((fleet) => {
				return {
					path: 'M' + fleet.waypoints.map((wp) => `${$xGet(wp)}, ${$yGet(wp)}`).join('L'),
					props: {
						class: fleet == $commandedFleet ? 'waypoint-line-commanded' : 'waypoint-line'
					}
				};
			});
		}
	}
</script>

{#each lines as line}
	<path d={line.path} {...line.props} />
{/each}

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
				const coords = fleet.waypoints.map((wp) => {
					return { position: { x: wp.position.x, y: wp.position.y } };
				});
				// move the first coord along the heading a bit so the line starts after our icon
				const heading = fleet.heading ?? { x: 0, y: 0 };
				coords[0].position.x += heading.x * 5;
				coords[0].position.y += heading.y * 5;

				return {
					path: 'M' + coords.map((coord) => `${$xGet(coord)}, ${$yGet(coord)}`).join('L'),
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

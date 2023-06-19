<script lang="ts">
	import { commandedFleet } from '$lib/services/Stores';
	import { emptyVector } from '$lib/types/Vector';

	import { getGameContext } from '$lib/services/Contexts';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const { game, player, universe } = getGameContext();
	const scale = getContext<Writable<number>>('scale');
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	type WaypointLine = {
		path: string;
		props: any;
	};

	let lines: WaypointLine[] = [];

	$: strokeWidth = 1;

	$: {
		if ($data) {
			const fleets = $universe.fleets.filter((fleet) => fleet.waypoints?.length);
			lines = fleets
				.filter((fleet) => (fleet.waypoints?.length ?? 0) > 1)
				.map((fleet) => {
					const coords = fleet.waypoints?.map((wp) => {
						return { position: { x: wp.position.x, y: wp.position.y } };
					}) ?? [{ position: emptyVector }];
					// move the first coord along the heading a bit so the line starts after our icon
					const heading = fleet.heading ?? { x: 0, y: 0 };
					coords[0].position.x += (heading.x * 5) / $scale;
					coords[0].position.y += (heading.y * 5) / $scale;

					strokeWidth = (fleet.num === $commandedFleet?.num ? 5 : 3) / $scale;

					return {
						path: 'M' + coords.map((coord) => `${$xGet(coord)}, ${$yGet(coord)}`).join('L'),
						props: {
							class:
								fleet.num === $commandedFleet?.num ? 'waypoint-line-commanded' : 'waypoint-line'
						}
					};
				});
		}
	}
</script>

{#each lines as line}
	<path d={line.path} {...line.props} stroke-width={strokeWidth} />
{/each}

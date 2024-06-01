<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { minus, normalized } from '$lib/types/Vector';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	type Line = {
		path: string;
		props: any;
	};

	const { universe } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	const strokeWidth = 1;

	$: wormholes = $universe.wormholes.filter((w) => w.destinationNum);

	let lines: Line[] = [];

	$: {
		const numsUsed = new Set<number>();
		lines = wormholes
			.filter((wormhole) => {
				const used = numsUsed.has(wormhole.num) || numsUsed.has(wormhole.destinationNum ?? 0);
				numsUsed.add(wormhole.num);
				if (wormhole.destinationNum) {
					numsUsed.add(wormhole.destinationNum);
				}
				return !used;
			})
			.map((wormhole) => {
				// get the target, if it's empty, just point to our planet position (which will render an empty line)
				// it should not be empty...
				const target = $universe.getWormhole(wormhole.destinationNum ?? 0);
				const coords = [
					{ position: wormhole.position },
					{ position: target?.position ?? wormhole.position }
				];

				const heading = normalized(minus(coords[0].position, coords[1].position));
				coords[0].position = {
					x: (coords[0].position.x ?? 0) - heading.x * 3,
					y: coords[0].position.y - heading.y * 3
				};
				coords[1].position = {
					x: coords[1].position.x + heading.x * 3,
					y: coords[1].position.y + heading.y * 3
				};

				return {
					path: 'M' + coords.map((coord) => `${$xGet(coord)}, ${$yGet(coord)}`).join('L'),
					props: {
						'stroke-width': strokeWidth
					}
				};
			});
	}
</script>

{#each lines as line}
	<path d={line.path} {...line.props} class="wormhole-link" />
{/each}

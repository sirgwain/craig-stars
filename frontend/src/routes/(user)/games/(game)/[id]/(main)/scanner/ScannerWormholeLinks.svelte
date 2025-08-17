<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { VectorSchema } from '$lib/types/cs-proto';
	import { emptyVector, normalized, subtract } from '$lib/types/Vector';
	import { create } from '@bufbuild/protobuf';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { SVGAttributes } from 'svelte/elements';
	import { SvelteSet } from 'svelte/reactivity';

	type Line = {
		path: string;
		props: SVGAttributes<SVGPathElement>;
	};

	const { universe } = getGameContext();
	const { xGet, yGet } = getContext<LayerCake>('LayerCake');

	const strokeWidth = 1;

	let lines: Line[] = $derived.by(() => {
		let wormholes = $universe.wormholeIntels.filter((w) => w.destinationNum);
		const numsUsed = new SvelteSet<number>();
		return wormholes
			.filter((wormhole) => {
				const used =
					numsUsed.has(wormhole.mapObject?.num ?? 0) || numsUsed.has(wormhole.destinationNum ?? 0);
				numsUsed.add(wormhole.mapObject?.num ?? 0);
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
					{ position: wormhole.mapObject?.position ?? emptyVector() },
					{ position: target?.mapObject?.position ?? wormhole.mapObject?.position ?? emptyVector() }
				];

				const heading = normalized(subtract(coords[0].position, coords[1].position));
				coords[0].position = create(VectorSchema, {
					x: (coords[0].position.x ?? 0) - heading.x * 3,
					y: coords[0].position.y - heading.y * 3
				});
				coords[1].position = create(VectorSchema, {
					x: coords[1].position.x + heading.x * 3,
					y: coords[1].position.y + heading.y * 3
				});

				return {
					path: 'M' + coords.map((coord) => `${$xGet(coord)}, ${$yGet(coord)}`).join('L'),
					props: {
						'stroke-width': strokeWidth
					}
				};
			});
	});
</script>

{#each lines as line (line)}
	<path d={line.path} {...line.props} class="wormhole-link" />
{/each}

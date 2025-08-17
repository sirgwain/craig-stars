<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { None } from '$lib/types/Consts';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { universe, commandedPlanet } = getGameContext();
	const { xGet, yGet, xScale } = getContext<LayerCake>('LayerCake');

	let planets = $derived(
		$universe.planets.filter(
			(planet) =>
				planet.planetOrders?.packetTargetNum && planet.planetOrders?.packetTargetNum != None
		)
	);

	let lines = $derived(
		planets.map((planet) => {
			// get the target, if it's empty, just point to our planet position (which will render an empty line)
			// it should not be empty...
			const target = $universe.getPlanet(planet.planetOrders?.packetTargetNum ?? None);
			const coords = [
				{ position: planet.mapObject?.position },
				{ position: target?.mapObject?.position ?? planet.mapObject?.position }
			];

			const strokeWidth = planet.mapObject?.num === $commandedPlanet?.mapObject?.num ? 1.5 : 1;
			const dist =
				(planet.planetOrders?.packetSpeed ?? 0) * (planet.planetOrders?.packetSpeed ?? 0);

			return {
				path: 'M' + coords.map((coord) => `${$xGet(coord)}, ${$yGet(coord)}`).join('L'),
				props: {
					'stroke-width': strokeWidth,
					'stroke-dasharray': `${$xScale(dist) - $xScale(5)} ${$xScale(5)}`,
					'stroke-dashoffset': `${$xScale(dist / 2) - $xScale(5)}`
				}
			};
		})
	);
</script>

<svg>
	<defs>
		<marker
			id="packet-arrow"
			class="packet-arrow"
			viewBox="0 0 10 10"
			refX="13"
			refY="5"
			markerUnits="strokeWidth"
			markerWidth="3"
			markerHeight="3"
			orient="auto"
		>
			<path d="M 0 0 L 10 5 L 0 10 z" stroke="context-stroke" fill="context-fill" />
		</marker>
	</defs>
</svg>
{#each lines as line (line)}
	<path d={line.path} {...line.props} class="packet-dest-line" marker-end="url(#packet-arrow)" />
{/each}

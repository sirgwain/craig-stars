<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { MapObjectTargetSchema, MapObjectType } from '$lib/types/cs-proto';
	import { None } from '$lib/types/Consts';
	import { emptyVector } from '$lib/types/Vector';
	import { create } from '@bufbuild/protobuf';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const { universe, commandedPlanet } = getGameContext();
	const { xGet, yGet, xScale } = getContext<LayerCake>('LayerCake');

	let planets = $derived(
		$universe.planets.filter((planet) => (planet.planetOrders?.routeTargetNum ?? None) != None)
	);

	let lines = $derived(
		planets.map((planet) => {
			// get the target, if it's empty, just point to our planet position (which will render an empty line)
			// it should not be empty...
			const target = $universe.getMapObject(
				create(MapObjectTargetSchema, {
					targetPosition: emptyVector(),
					targetType: planet.planetOrders?.routeTargetType ?? MapObjectType.UNSPECIFIED,
					targetNum: planet.planetOrders?.routeTargetNum ?? 0,
					targetPlayerNum: planet.planetOrders?.routeTargetPlayerNum ?? 0
				})
			);
			const coords = [
				{ position: planet.mapObject?.position },
				{ position: target?.mapObject?.position ?? planet.mapObject?.position }
			];

			const strokeWidth = planet.mapObject?.num === $commandedPlanet?.mapObject.num ? 1.5 : 1;
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
			id="route-arrow"
			class="route-arrow"
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
	<path d={line.path} {...line.props} class="route-line" marker-end="url(#route-arrow)" />
{/each}

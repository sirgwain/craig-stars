<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import type { Fleet } from '$lib/types/Fleet';
	import { MapObjectType, None } from '$lib/types/MapObject';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const { game, player, universe, settings } = getGameContext();
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

	export let planet: Planet;

	let props = {};
	let flagColor = '#555';

	$: planetX = $xGet(planet);
	$: planetY = $yGet(planet);

	const fullyHabitableRadius = 7;
	let orbitingTokens = 0;
	let radius = 3;
	let tokenCountOffset = 8;

	$: {
		// green for us, gray for unexplored, white for explored
		let color = '#555';
		let strokeWidth = 0;
		let strokeColor = '#888';
		let minRadius = 2;
		radius = 3;
		tokenCountOffset = 8;

		if (planet.reportAge !== Unexplored) {
			strokeWidth = 1;
			let habitability = planet.spec?.habitability ?? 0;
			let habitabilityTerraformed = planet.spec?.terraformedHabitability ?? 0;
			if (habitability > 0) {
				color = '#00FF00';
				radius = Math.max((habitability / 100.0) * fullyHabitableRadius, minRadius);
				strokeWidth = (habitability / 100.0) * strokeWidth;
			} else {
				if (habitabilityTerraformed > 0) {
					color = '#FFFF00';
					radius = Math.max((habitabilityTerraformed / 100.0) * fullyHabitableRadius, minRadius);
					strokeWidth = (habitabilityTerraformed / 100.0) * strokeWidth;
				} else {
					color = '#FF0000';
					radius = Math.max((-habitability / 45.0) * fullyHabitableRadius, minRadius);
					strokeWidth = (-habitability / 45.0) * strokeWidth;
				}
			}

			if (planet.playerNum) {
				flagColor = $universe.getPlayerColor(planet.playerNum) ?? '#FF0000';
			}
		}

		tokenCountOffset = 3 + radius / $scale;

		const orbitingFleets = $universe
			.getMapObjectsByPosition(planet)
			.filter((mo) => mo.type === MapObjectType.Fleet);

		orbitingTokens = orbitingFleets
			.map((of) => of as Fleet)
			.reduce(
				(count, f) =>
					count + (f.tokens ? f.tokens.reduce((tokenCount, t) => tokenCount + t.quantity, 0) : 0),
				0
			);

		// setup the properties of our planet circle
		props = {
			r: radius,
			fill: color,
			stroke: strokeColor,
			'stroke-width': strokeWidth
		};
	}
</script>

<circle cx={planetX} cy={planetY} {...props} />

{#if planet.reportAge !== Unexplored && planet.playerNum != None}
	<!-- draw the flag  -->
	<rect width="8" height="6" x={planetX} y={planetY - fullyHabitableRadius * 2} fill={flagColor} />
	<path
		d={`M${planetX}, ${planetY}L${planetX}, ${planetY - fullyHabitableRadius * 2}`}
		stroke={flagColor}
	/>
{/if}
{#if $settings.showFleetTokenCounts && orbitingTokens}
	<!-- translate the group to the location of the fleet so when we scale the text it is around the center-->
	<g
		transform={`translate(${planetX - (8/$scale)} ${
			planetY + tokenCountOffset * 2.5
		})`}
	>
		<text transform={`scale(${1 / $scale})`} class="fill-base-content">{orbitingTokens}</text>
	</g>
{/if}

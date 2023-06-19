<script lang="ts">
	import type { FullGame } from '$lib/services/FullGame';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';

	const game = getContext<FullGame>('game');
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	$: scale = getContext<number>('scale');

	export let planet: Planet;
	export let commanded = false;
	export let orbitingFleets = false;

	let props = {};
	let ringProps: any | undefined = undefined;

	$: hasStarbase = planet.spec?.hasStarbase;
	$: hasMassDriver = planet.spec?.hasMassDriver;
	$: hasStargate = planet.spec?.hasStargate;

	$: planetX = $xGet(planet);
	$: planetY = $yGet(planet);

	$: starbaseWidth = commanded ? 5 : 3;
	$: starbaseXOffset = commanded ? 5 : 2;
	$: starbaseYOffset = commanded ? 11 : 6;

	$: {
		// green for us, gray for unexplored, white for explored
		let color = '#555';
		let strokeWidth = 0;
		let strokeColor = '#555';

		if (planet.playerNum === game?.player.num) {
			color = '#00FF00';
			strokeWidth = 1;
		} else if (planet.playerNum) {
			color = game?.getPlayerColor(planet.playerNum) ?? '#FF0000';
		} else if (planet.reportAge !== Unexplored && !planet.playerNum) {
			color = '#FFF';
		}

		// if anything is orbiting our planet, put a ring on it
		if (orbitingFleets) {
			ringProps = {
				cx: $xGet(planet),
				cy: $yGet(planet),
				stroke: '#fff',
				'stroke-width': 1,
				r: 1 * (commanded ? 10 : 5),
				'fill-opacity': 0
			};
		} else {
			ringProps = undefined;
		}

		// setup the properties of our planet circle
		props = {
			r: (1 * (commanded ? 7 : 3)),
			fill: color,
			stroke: strokeColor,
			'stroke-width': strokeWidth
		};
	}
</script>

{#if ringProps}
	<circle {...ringProps} />
{/if}
<circle cx={planetX} cy={planetY} {...props} />
{#if hasStarbase}
	<rect
		class="starbase"
		width={starbaseWidth}
		height={starbaseWidth}
		rx={0.5}
		x={planetX + starbaseXOffset}
		y={planetY - starbaseYOffset}
	/>
{/if}
{#if hasStargate}
	<rect
		class="stargate"
		width={starbaseWidth}
		height={starbaseWidth}
		rx={0.5}
		x={planetX - starbaseXOffset - starbaseWidth}
		y={planetY - starbaseYOffset}
	/>
{/if}
{#if hasMassDriver}
	<rect
		class="massdriver"
		width={starbaseWidth}
		height={starbaseWidth}
		rx={0.5}
		x={planetX - starbaseWidth / 2}
		y={planetY - starbaseYOffset - starbaseWidth / 2}
	/>
{/if}

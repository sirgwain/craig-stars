<script lang="ts">
	import type { FullGame } from '$lib/services/FullGame';
	import { clamp } from '$lib/services/Math';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const game = getContext<FullGame>('game');
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');
	const scale = getContext<Writable<number>>('scale');

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

	$: starbaseWidth = (commanded ? 6 : 4) / $scale;
	$: starbaseXOffset = (commanded ? 5 : 2) / $scale;
	$: starbaseYOffset = (commanded ? 11 : 6) / $scale;

	$: {
		// green for us, gray for unexplored, white for explored
		let color = '#555';
		let strokeWidth = 0;
		let strokeColor = '#555';

		if (planet.playerNum === game?.player.num) {
			color = '#00FF00';
			strokeWidth = 1 / $scale;
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
				'stroke-width': 1 / $scale,
				r: 1 * (commanded ? 10 : 5) / $scale,
				'fill-opacity': 0
			};
		} else {
			ringProps = undefined;
		}

		// setup the properties of our planet circle
		props = {
			r: 1 * (commanded ? 7 : 3) / $scale,
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

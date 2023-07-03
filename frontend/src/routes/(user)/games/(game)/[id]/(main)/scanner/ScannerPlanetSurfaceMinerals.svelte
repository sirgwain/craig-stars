<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { clamp } from '$lib/services/Math';
	import type { Planet } from '$lib/types/Planet';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';

	const { game, player, universe, settings } = getGameContext();
	const scale = getContext<Writable<number>>('scale');
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	let max = $settings.mineralScale; // 100% concentration

	export let planet: Planet;

	const size = 8; // the size of the mineral bars
	const abovePlanetY = 2;

	let barPercent = {
		ironium: 0,
		boranium: 0,
		germanium: 0
	};

	$: {
		if (planet.cargo) {
			barPercent = {
				ironium: clamp(planet.cargo.ironium ? planet.cargo.ironium / max : 0, 0, 1),
				boranium: clamp(planet.cargo.boranium ? planet.cargo.boranium / max : 0, 0, 1),
				germanium: clamp(planet.cargo.germanium ? planet.cargo.germanium / max : 0, 0, 1)
			};
		}
	}

	$: planetX = $xGet(planet);
	$: planetY = $yGet(planet);
</script>

<ScannerPlanetNormal {planet} />
{#if planet.cargo}
	<rect
		class="ironium-bar"
		width={2}
		height={barPercent.ironium * size}
		x={planetX - 3}
		y={planetY - abovePlanetY - size + (size - barPercent.ironium * size)}
	/>
	<rect
		class="boranium-bar"
		width={2}
		height={barPercent.boranium * size}
		x={planetX}
		y={planetY - abovePlanetY - size + (size - barPercent.boranium * size)}
	/>

	<rect
		class="germanium-bar"
		width={2}
		height={barPercent.germanium * size}
		x={planetX + 3}
		y={planetY - abovePlanetY - size + (size - barPercent.germanium * size)}
	/>
	<path
		class="stroke-white"
		stroke-width={1}
		fill="none"
		d={`M${planetX + 5},${planetY - abovePlanetY}L${planetX - 4},${planetY - abovePlanetY} L${
			planetX - 4
		},${planetY - abovePlanetY - size}`}
	/>
{/if}

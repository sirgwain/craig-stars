<script lang="ts">
	import type { Planet } from '$lib/types/cs-proto';
	import { clamp } from '$lib/services/Math';
	import { ReportAgeUnexplored } from '$lib/types/Consts';
	import MapObjectScaler from './MapObjectScaler.svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';

	let max = 100; // 100% concentration

	type Props = {
		planet: Planet;
	};

	let { planet }: Props = $props();

	const size = 25; // the size of the mineral bars
	const abovePlanetY = 5;

	let barPercent = $derived.by(() => {
		if (planet.mineralConcentration) {
			return {
				ironium: clamp(
					planet.mineralConcentration.ironium ? planet.mineralConcentration.ironium / max : 0,
					0,
					1
				),
				boranium: clamp(
					planet.mineralConcentration.boranium ? planet.mineralConcentration.boranium / max : 0,
					0,
					1
				),
				germanium: clamp(
					planet.mineralConcentration.germanium ? planet.mineralConcentration.germanium / max : 0,
					0,
					1
				)
			};
		} else {
			return { ironium: 0, boranium: 0, germanium: 0 };
		}
	});
</script>

<ScannerPlanetNormal {planet} />
{#if planet.mapObject?.reportAge !== ReportAgeUnexplored}
	<MapObjectScaler mapObject={planet}>
		<rect
			class="ironium-bar"
			width={size / 4}
			height={barPercent.ironium * size}
			x={-size / 3}
			y={-abovePlanetY - size + (size - barPercent.ironium * size)}
		/>
		<rect
			class="boranium-bar"
			width={size / 4}
			height={barPercent.boranium * size}
			x={0}
			y={-abovePlanetY - size + (size - barPercent.boranium * size)}
		/>

		<rect
			class="germanium-bar"
			width={size / 4}
			height={barPercent.germanium * size}
			x={size / 3}
			y={-abovePlanetY - size + (size - barPercent.germanium * size)}
		/>
		<path
			class="stroke-white"
			stroke-width={1}
			fill="none"
			d={`M${size / 1.5},${-abovePlanetY}L${-size / 2},${-abovePlanetY} L${-size / 2},${-abovePlanetY - size}`}
		/>
	</MapObjectScaler>
{/if}

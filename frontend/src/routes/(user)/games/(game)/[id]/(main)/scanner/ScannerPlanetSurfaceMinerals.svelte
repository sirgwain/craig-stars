<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { clamp } from '$lib/services/Math';
	import { totalCargo } from '$lib/types/Cargo';
	import { ReportAgeUnexplored } from '$lib/types/cs';
	import { type Planet } from '$lib/types/cs';
	import MapObjectScaler from './MapObjectScaler.svelte';
	import ScannerPlanetNormal from './ScannerPlanetNormal.svelte';

	const { settings } = getGameContext();

	type Props = {
		planet: Planet;
	};

	let { planet }: Props = $props();

	const size = 25; // the size of the mineral bars
	const abovePlanetY = 5;

	let barPercent = $derived.by(() => {
		let max = $settings.mineralScale; // 100% concentration
		if (!planet.cargo) {
			return {
				ironium: 0,
				boranium: 0,
				germanium: 0
			};
		}
		return {
			ironium: clamp(planet.cargo.ironium ? planet.cargo.ironium / max : 0, 0, 1),
			boranium: clamp(planet.cargo.boranium ? planet.cargo.boranium / max : 0, 0, 1),
			germanium: clamp(planet.cargo.germanium ? planet.cargo.germanium / max : 0, 0, 1)
		};
	});
</script>

<ScannerPlanetNormal {planet} />
{#if planet.reportAge !== ReportAgeUnexplored && totalCargo(planet.cargo) != 0}
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

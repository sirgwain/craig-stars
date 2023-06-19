<script lang="ts">
	import { clamp } from '$lib/services/Math';

	import type { Mineral } from '$lib/types/Mineral';

	import type { Planet } from '$lib/types/Planet';

	export let planet: Planet;
	export let scale = 1.0;
	export let max = 5000.0;
	let numDivisions = 6; // gridlines show 20% on line class
	let divisions: string[] = ['0'];

	for (let i = 1; i < numDivisions; i++) {
		divisions[i] = (i * (max / (numDivisions - 1))).toFixed();
	}

	let barPercent: Mineral = {
		ironium: 0,
		boranium: 0,
		germanium: 0
	};

	$: {
		if (planet && planet.cargo) {
			barPercent = {
				ironium: clamp(planet.cargo.ironium ? (planet.cargo.ironium / max) * 100 : 0, 0, 100),
				boranium: clamp(planet.cargo.boranium ? (planet.cargo.boranium / max) * 100 : 0, 0, 100),
				germanium: clamp(planet.cargo.germanium ? (planet.cargo.germanium / max) * 100 : 0, 0, 100)
			};
		}
	}
</script>

<div class="flex flex-row">
	<div class="text-right flex flex-col justify-evenly w-[5.5rem]">
		<div class="text-ironium">Ironium</div>
		<div class="text-boranium">Boranium</div>
		<div class="text-germanium">Germanium</div>
	</div>
	<div
		class="grow flex flex-col justify-evenly mx-1 px-0.5 py-1 bg-black line gap-2"
	>
		<div style={`width: ${barPercent.ironium?.toFixed()}%`} class="ironium-bar h-full" />
		<div style={`width: ${barPercent.boranium?.toFixed()}%`} class="boranium-bar h-full" />
		<div style={`width: ${barPercent.germanium?.toFixed()}%`} class="germanium-bar h-full" />
	</div>
	<div class="w-[3rem]" />
</div>
<div class="flex flex-row">
	<div class="text-right flex flex-col justify-evenly w-[5.5rem] pr-1">kT</div>
	<div class="grow flex flex-row justify-between">
		{#each divisions as division}
			<div>{division}</div>
		{/each}
		<!-- spacer -->
		<div class="w-[3rem]" />
	</div>
</div>

<style>
	.line {
		background: repeating-linear-gradient(to right, #222, #222 1px, #000 1px, #000 20%);
	}
</style>

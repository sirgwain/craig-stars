<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import type { Mineral } from '$lib/types/Mineral';
	import type { Planet } from '$lib/types/Planet';
	import MineralConcentrationPoint from '$lib/components/game/MineralConcentrationPoint.svelte';
	import MineralTooltip, {
		type MineralTooltipProps
	} from '$lib/components/game/tooltips/MineralTooltip.svelte';
	import { showTooltip } from '$lib/services/Stores';

	export let planet: Planet;
	// export let scale = 1.0;
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

	let concentrationPercent: Mineral = {
		ironium: 0,
		boranium: 0,
		germanium: 0
	};

	$: {
		if (planet.cargo) {
			barPercent = {
				ironium: clamp(planet.cargo.ironium ? (planet.cargo.ironium / max) * 100 : 0, 0, 100),
				boranium: clamp(planet.cargo.boranium ? (planet.cargo.boranium / max) * 100 : 0, 0, 100),
				germanium: clamp(planet.cargo.germanium ? (planet.cargo.germanium / max) * 100 : 0, 0, 100)
			};
		} else {
			barPercent = { ironium: 0, boranium: 0, germanium: 0 };
		}
	}

	$: {
		if (planet.mineralConcentration) {
			concentrationPercent = {
				ironium: clamp(
					planet.mineralConcentration.ironium
						? (planet.mineralConcentration.ironium / 100) * 100
						: 0,
					0,
					100
				),
				boranium: clamp(
					planet.mineralConcentration.boranium
						? (planet.mineralConcentration.boranium / 100) * 100
						: 0,
					0,
					100
				),
				germanium: clamp(
					planet.mineralConcentration.germanium
						? (planet.mineralConcentration.germanium / 100) * 100
						: 0,
					0,
					100
				)
			};
		}
	}

	function onIroniumTooltip(e: PointerEvent) {
		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Ironium',
			surfaceAmount: planet.cargo?.ironium ?? 0,
			concentration: planet.mineralConcentration?.ironium ?? 0,
			miningRate: planet.spec.miningOutput.ironium ?? 0
		});
	}
	function onBoraniumTooltip(e: PointerEvent) {
		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Boranium',
			surfaceAmount: planet.cargo?.boranium ?? 0,
			concentration: planet.mineralConcentration?.boranium ?? 0,
			miningRate: planet.spec.miningOutput.boranium ?? 0
		});
	}
	function onGermaniumTooltip(e: PointerEvent) {
		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Germanium',
			surfaceAmount: planet.cargo?.germanium ?? 0,
			concentration: planet.mineralConcentration?.germanium ?? 0,
			miningRate: planet.spec.miningOutput.germanium ?? 0
		});
	}
</script>

<div class="flex flex-row">
	<div class="text-right flex flex-col justify-evenly w-[5.5rem]">
		<div class="text-ironium">Ironium</div>
		<div class="text-boranium">Boranium</div>
		<div class="text-germanium">Germanium</div>
	</div>
	<div class="grow flex flex-col justify-evenly mx-1 px-0.5 py-1 bg-black line gap-2">
		<div class="h-full relative cursor-help" on:pointerdown|preventDefault={onIroniumTooltip}>
			<MineralConcentrationPoint
				style={`left: ${concentrationPercent.ironium?.toFixed()};`}
				class="absolute ironium-concentration w-auto h-full ironium"
			/>
			<div style={`width: ${barPercent.ironium?.toFixed()}%`} class="ironium-bar h-full" />
		</div>
		<div class="h-full relative cursor-help" on:pointerdown|preventDefault={onBoraniumTooltip}>
			<MineralConcentrationPoint
				style={`left: ${concentrationPercent.boranium?.toFixed()};`}
				class="absolute boranium-concentration w-auto h-full boranium"
			/>
			<div style={`width: ${barPercent.boranium?.toFixed()}%`} class="boranium-bar h-full" />
		</div>
		<div class="h-full relative cursor-help" on:pointerdown|preventDefault={onGermaniumTooltip}>
			<MineralConcentrationPoint
				style={`left: ${concentrationPercent.germanium?.toFixed()};`}
				class="absolute germanium-concentration  h-full germanium"
			/>
			<div style={`width: ${barPercent.germanium?.toFixed()}%`} class="germanium-bar h-full" />
		</div>
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

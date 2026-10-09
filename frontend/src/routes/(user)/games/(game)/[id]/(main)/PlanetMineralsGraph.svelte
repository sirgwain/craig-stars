<script lang="ts">
	import { onPointerKeyDown, onPointerKeyUp } from '#lib/services/Events.js';
	import MineralConcentrationPoint from '#lib/components/game/MineralConcentrationPoint.svelte';
	import MineralTooltip, {
		type MineralTooltipProps
	} from '#lib/components/game/tooltips/MineralTooltip.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { clamp } from '#lib/services/Math.js';
	import { showTooltip } from '#lib/services/Stores.js';
	import type { MineralJson, Planet } from '#lib/types/cs-proto.js';
	import { ChevronDown } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { settings } = getGameContext();

	type Props = {
		planet: Planet;
	};

	let { planet }: Props = $props();

	// the scales players can choose for the mineral graph, like Stars!
	const mineralScales = [100, 500, 1000, 2500, 5000, 7500, 10000, 20000, 30000];

	let max = $derived($settings.mineralScale); // i.e. 0 to 5000 minerals
	let numDivisions = 6; // gridlines show 20% online class
	let divisions: string[] = $derived(
		Array.from({ length: numDivisions }, (_, i) => (i * (max / (numDivisions - 1))).toFixed())
	);

	function percentOfMax(amount: number | undefined): number {
		return clamp(amount ? (amount / max) * 100 : 0, 0, 100);
	}

	let barPercent: MineralJson = $derived({
		ironium: percentOfMax(planet.cargo?.ironium),
		boranium: percentOfMax(planet.cargo?.boranium),
		germanium: percentOfMax(planet.cargo?.germanium)
	});

	// what the surface minerals will be next year, after mining, shown as a shaded bar
	let nextYearMinerals = $derived({
		ironium: (planet.cargo?.ironium ?? 0) + (planet.spec?.miningOutput?.ironium ?? 0),
		boranium: (planet.cargo?.boranium ?? 0) + (planet.spec?.miningOutput?.boranium ?? 0),
		germanium: (planet.cargo?.germanium ?? 0) + (planet.spec?.miningOutput?.germanium ?? 0)
	});
	let nextYearPercent: MineralJson = $derived({
		ironium: percentOfMax(nextYearMinerals.ironium),
		boranium: percentOfMax(nextYearMinerals.boranium),
		germanium: percentOfMax(nextYearMinerals.germanium)
	});

	function setMineralScale(scale: number) {
		$settings.mineralScale = scale;
		// close the dropdown
		(document.activeElement as HTMLElement | null)?.blur();
	}

	let concentrationPercent: MineralJson = $derived(
		planet.mineralConcentration
			? {
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
				}
			: { ironium: 0, boranium: 0, germanium: 0 }
	);

	function onIroniumTooltip(e: PointerEvent) {
		e.preventDefault();

		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Ironium',
			surfaceAmount: planet.cargo?.ironium ?? 0,
			concentration: planet.mineralConcentration?.ironium ?? 0,
			miningRate: planet.spec?.miningOutput?.ironium ?? 0,
			homeworld: !!planet.homeworld
		});
	}
	function onBoraniumTooltip(e: PointerEvent) {
		e.preventDefault();

		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Boranium',
			surfaceAmount: planet.cargo?.boranium ?? 0,
			concentration: planet.mineralConcentration?.boranium ?? 0,
			miningRate: planet.spec?.miningOutput?.boranium ?? 0,
			homeworld: !!planet.homeworld
		});
	}
	function onGermaniumTooltip(e: PointerEvent) {
		e.preventDefault();

		showTooltip<MineralTooltipProps>(e.x, e.y, MineralTooltip, {
			mineralType: 'Germanium',
			surfaceAmount: planet.cargo?.germanium ?? 0,
			concentration: planet.mineralConcentration?.germanium ?? 0,
			miningRate: planet.spec?.miningOutput?.germanium ?? 0,
			homeworld: !!planet.homeworld
		});
	}
</script>

<div class="flex flex-row select-none">
	<div class="text-right flex flex-col justify-evenly w-[5.5rem]">
		<div class="text-ironium">Ironium</div>
		<div class="text-boranium">Boranium</div>
		<div class="text-germanium">Germanium</div>
	</div>
	<div class="grow flex flex-col justify-evenly mx-1 py-1 bg-black line gap-2">
		<div
			role="button"
			tabindex={0}
			aria-label="Show Ironium details"
			onkeydown={onPointerKeyDown}
			onkeyup={onPointerKeyUp}
			class="h-full relative cursor-help"
			onpointerdown={onIroniumTooltip}
		>
			<MineralConcentrationPoint
				style={`left: ${concentrationPercent.ironium?.toFixed()}%;`}
				class="absolute z-10 ironium-concentration w-auto h-full ironium"
			/>
			<div
				style={`width: ${nextYearPercent.ironium}%`}
				class="ironium-bar h-full absolute opacity-50"
			></div>
			<div style={`width: ${barPercent.ironium}%`} class="ironium-bar h-full relative"></div>
			{#if nextYearMinerals.ironium > max}
				<span
					class="absolute left-full top-1/2 -translate-y-1/2 ml-1 text-base-content"
					title="Next year's Ironium exceeds the mineral scale">+</span
				>
			{/if}
		</div>
		<div
			role="button"
			tabindex={0}
			aria-label="Show Boranium details"
			onkeydown={onPointerKeyDown}
			onkeyup={onPointerKeyUp}
			class="h-full relative cursor-help"
			onpointerdown={onBoraniumTooltip}
		>
			<MineralConcentrationPoint
				style={`left: ${concentrationPercent.boranium?.toFixed()}%;`}
				class="absolute z-10 boranium-concentration w-auto h-full boranium"
			/>
			<div
				style={`width: ${nextYearPercent.boranium}%`}
				class="boranium-bar h-full absolute opacity-50"
			></div>
			<div style={`width: ${barPercent.boranium}%`} class="boranium-bar h-full relative"></div>
			{#if nextYearMinerals.boranium > max}
				<span
					class="absolute left-full top-1/2 -translate-y-1/2 ml-1 text-base-content"
					title="Next year's Boranium exceeds the mineral scale">+</span
				>
			{/if}
		</div>
		<div
			role="button"
			tabindex={0}
			aria-label="Show Germanium details"
			onkeydown={onPointerKeyDown}
			onkeyup={onPointerKeyUp}
			class="h-full relative cursor-help"
			onpointerdown={onGermaniumTooltip}
		>
			<MineralConcentrationPoint
				style={`left: ${concentrationPercent.germanium?.toFixed()}%;`}
				class="absolute z-10 germanium-concentration h-full germanium"
			/>
			<div
				style={`width: ${nextYearPercent.germanium}%`}
				class="germanium-bar h-full absolute opacity-50"
			></div>
			<div style={`width: ${barPercent.germanium}%`} class="germanium-bar h-full relative"></div>
			{#if nextYearMinerals.germanium > max}
				<span
					class="absolute left-full top-1/2 -translate-y-1/2 ml-1 text-base-content"
					title="Next year's Germanium exceeds the mineral scale">+</span
				>
			{/if}
		</div>
	</div>
	<div class="w-[3rem] shrink-0"></div>
</div>
<div class="flex flex-row">
	<div class="text-right flex flex-col justify-evenly w-[5.5rem] pr-1">kT</div>
	<div class="grow flex flex-row justify-between">
		{#each divisions as division, index (index)}
			<div>{division}</div>
		{/each}
		<div class="dropdown dropdown-end dropdown-top w-[3rem] text-right">
			<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
			<label
				tabindex="0"
				class="cursor-pointer inline-block"
				aria-label="Mineral scale"
				title="Mineral scale"
				><Icon src={ChevronDown} size="16" class="hover:stroke-accent inline-block" /></label
			>
			<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
			<ul
				tabindex="0"
				class="menu menu-sm dropdown-content p-1 shadow-sm bg-base-300 z-20"
				data-testid="mineral-scale-menu"
			>
				{#each mineralScales as scale (scale)}
					<li>
						<button
							type="button"
							class:menu-active={scale === max}
							onclick={() => setMineralScale(scale)}>{scale.toLocaleString()}kT</button
						>
					</li>
				{/each}
			</ul>
		</div>
	</div>
</div>

<style>
	.line {
		background: repeating-linear-gradient(to right, #222, #222 1px, #000 1px, #000 20%);
	}
</style>

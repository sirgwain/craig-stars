<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { getGravString, getRadString, getTempString, HabType } from '$lib/types/Hab';
	import type { Race } from '$lib/types/Race';
	import { Icon } from '@steeze-ui/svelte-icon';
	import {
		ArrowLeft,
		ChevronDoubleLeft,
		ChevronDoubleRight,
		ChevronLeft,
		ChevronRight
	} from '@steeze-ui/heroicons';
	import HabBar from './HabBar.svelte';

	export let race: Race;

	$: habWidth = {
		grav: (race.habHigh.grav ?? 0) - (race.habLow.grav ?? 0),
		temp: (race.habHigh.temp ?? 0) - (race.habLow.temp ?? 0),
		rad: (race.habHigh.rad ?? 0) - (race.habLow.rad ?? 0)
	};

	$: habChance = ((): number => {
		// do a straight calc of hab width, so if we have a hab with widths of 50, 50% of planets will be habitable
		// so we get (.5 * .5 * .5) = .125, or 1 in 8 planets
		const gravChance = race.immuneGrav ? 1.0 : habWidth.grav / 100.0;
		const tempChance = race.immuneTemp ? 1.0 : habWidth.temp / 100.0;
		const radChance = race.immuneRad ? 1.0 : habWidth.rad / 100.0;
		return gravChance * tempChance * radChance;
	})();
	$: approximateHabitablePlanetRatio = Math.floor(1 / habChance);
</script>

<div class="flex flex-col gap-2">
	<HabBar
		habType={HabType.Gravity}
		bind:habLow={race.habLow.grav}
		bind:habHigh={race.habHigh.grav}
		bind:immune={race.immuneGrav}
	/>
	<HabBar
		habType={HabType.Temperature}
		bind:habLow={race.habLow.temp}
		bind:habHigh={race.habHigh.temp}
		bind:immune={race.immuneTemp}
	/>
	<HabBar
		habType={HabType.Radiation}
		bind:habLow={race.habLow.rad}
		bind:habHigh={race.habHigh.rad}
		bind:immune={race.immuneRad}
	/>
	<label>
		Maxium Colonist Growth Rate Per Year
		<input
			class="input input-sm input-bordered"
			type="number"
			name="growthRate"
			min={1}
			max={100}
			bind:value={race.growthRate}
		/></label
	>

	{#if habChance == 1}
		All planets will be habitable to your race.
	{:else if approximateHabitablePlanetRatio == 1}
		Virtually all planets will be habitable to your race.
	{:else}
		{`You can expect that 1 in ${approximateHabitablePlanetRatio} planets will be habitable to your race.`}
	{/if}
</div>

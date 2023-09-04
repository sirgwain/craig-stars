<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import HabChance from '$lib/components/game/race/HabChance.svelte';
	import LRTsDescriptions from '$lib/components/game/race/LRTsDescriptions.svelte';
	import PRTDescription from '$lib/components/game/race/PRTDescription.svelte';
	import Population from '$lib/components/icons/Population.svelte';
	import { HabType } from '$lib/types/Hab';
	import type { Race } from '$lib/types/Race';
	import { getLabelForPRT } from '$lib/types/Race';
	import HabBar from './HabBar.svelte';
	import PlanetaryProduction from './PlanetaryProduction.svelte';
	import Research from './Research.svelte';

	export let race: Race;
</script>

<div
	class="stats stats-vertical sm:stats-horizontal sm:flex shadow border border-base-200 w-full sm:w-48 sm:mx-auto"
>
	<div class="stat place-items-center sm:grow">
		<div class="stat-title">Growth Rate</div>
		<div class="stat-figure"><Population class="w-8 h-8 fill-base-content" /></div>
		<div class="stat-value">
			{race.growthRate * (race.spec?.growthFactor ?? 0)}%
		</div>
	</div>
</div>
<ItemTitle>Primary Racial Trait</ItemTitle>
<div class="card bg-base-200 shadow w-full">
	<div class="card-body">
		<div class="card-title text-lg">
			{getLabelForPRT(race.prt)}
		</div>
		<div>
			<PRTDescription prt={race.prt} />
		</div>
	</div>
</div>

<ItemTitle>Lesser Racial Traits</ItemTitle>
{#if race.lrts}
	<LRTsDescriptions {race} />
{:else}
	None
{/if}

<ItemTitle>Habitability</ItemTitle>

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
	<HabChance {race} />
</div>

<ItemTitle>Planetary Production</ItemTitle>
<PlanetaryProduction {race} />
<ItemTitle>Research</ItemTitle>
<Research {race} />

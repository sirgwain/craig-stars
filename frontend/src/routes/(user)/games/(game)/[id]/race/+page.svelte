<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import HabChance from '$lib/components/game/race/HabChance.svelte';
	import LRTsDescriptions from '$lib/components/game/race/LRTsDescriptions.svelte';
	import PRTDescription from '$lib/components/game/race/PRTDescription.svelte';
	import Population from '$lib/components/icons/Population.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { getLabelForPRT } from '$lib/types/Race';
	import { Grav, Rad, Temp } from '$lib/types/cs';
	import HabBar from './HabBar.svelte';
	import PlanetaryProduction from './PlanetaryProduction.svelte';
	import Research from './Research.svelte';
	const { player } = getGameContext();

	let race = $state($player.race);
</script>

<div class="w-full mx-auto md:max-w-2xl pb-8">
	<ItemTitle>{race.name}/{race.pluralName}</ItemTitle>
	<div
		class="stats stats-vertical sm:stats-horizontal sm:flex shadow border border-base-200 w-full sm:w-48 sm:mx-auto"
	>
		<div class="stat place-items-center sm:grow">
			<div class="stat-title">Growth Rate</div>
			<div class="stat-figure"><Population class="w-8 h-8 fill-base-content" /></div>
			<div class="stat-value">
				{race.growthRate * (race.spec?.growthFactor ?? 1)}%
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
			habType={Grav}
			habLow={race.habLow.grav}
			habHigh={race.habHigh.grav}
			immune={race.immuneGrav}
		/>
		<HabBar
			habType={Temp}
			habLow={race.habLow.temp}
			habHigh={race.habHigh.temp}
			immune={race.immuneTemp}
		/>
		<HabBar
			habType={Rad}
			habLow={race.habLow.rad}
			habHigh={race.habHigh.rad}
			immune={race.immuneRad}
		/>
		<HabChance {race} />
	</div>

	<ItemTitle>Planetary Production</ItemTitle>
	<PlanetaryProduction {race} />
	<ItemTitle>Research</ItemTitle>
	<Research {race} />
</div>

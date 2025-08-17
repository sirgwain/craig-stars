<script lang="ts">
	import HabChance from '$lib/components/game/race/HabChance.svelte';
	import LRTsDescriptions from '$lib/components/game/race/LRTsDescriptions.svelte';
	import PRTDescription from '$lib/components/game/race/PRTDescription.svelte';
	import Population from '$lib/components/icons/Population.svelte';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import type { Race, RaceSpec } from '$lib/types/cs-proto';
	import { Grav, Rad, Temp } from '$lib/types/Hab';
	import { getLabelForPRT } from '$lib/types/Race';
	import type { WasmClient } from '$lib/wasm';
	import HabBar from './HabBar.svelte';
	import PlanetaryProduction from './PlanetaryProduction.svelte';
	import Research from './Research.svelte';

	type Props = {
		wasmClient: WasmClient;
		race: Race;
	};

	let { wasmClient, race }: Props = $props();

	let spec: RaceSpec | undefined = $state();

	$effect(() => {
		wasmClient.computeRaceSpec({ race }).then((resp) => (spec = resp.spec));
	});
</script>

<div
	class="stats stats-vertical sm:stats-horizontal sm:flex shadow border border-base-200 w-full sm:w-48 sm:mx-auto"
>
	<div class="stat place-items-center sm:grow">
		<div class="stat-title">Growth Rate</div>
		<div class="stat-figure"><Population class="w-8 h-8 fill-base-content" /></div>
		<div class="stat-value">
			{race.growthRate * (spec?.growthFactor ?? 0)}%
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
		habLow={race.habLow?.grav ?? 0}
		habHigh={race.habHigh?.grav ?? 0}
		immune={race.immuneGrav}
	/>
	<HabBar
		habType={Temp}
		habLow={race.habLow?.temp ?? 0}
		habHigh={race.habHigh?.temp ?? 0}
		immune={race.immuneTemp}
	/>
	<HabBar
		habType={Rad}
		habLow={race.habLow?.rad ?? 0}
		habHigh={race.habHigh?.rad ?? 0}
		immune={race.immuneRad}
	/>
	<HabChance {race} />
</div>

<ItemTitle>Planetary Production</ItemTitle>
<PlanetaryProduction {race} />
<ItemTitle>Research</ItemTitle>
<Research {race} />

<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import LRTsDescriptions from '$lib/components/game/race/LRTsDescriptions.svelte';
	import PRTDescription from '$lib/components/game/race/PRTDescription.svelte';
	import { getLabelForPRT } from '$lib/types/Race';
	import {
		AR,
		CA,
		HE,
		IS,
		IT,
		JoaT,
		PP,
		SD,
		SpendLeftoverPointsOnDefenses,
		SpendLeftoverPointsOnFactories,
		SpendLeftoverPointsOnMineralConcentrations,
		SpendLeftoverPointsOnMines,
		SpendLeftoverPointsOnSurfaceMinerals,
		SS,
		WM,
		type Race
	} from '$lib/types/cs';
	import Habitability from './Habitability.svelte';
	import LRTs from './LRTs.svelte';
	import PlanetaryProduction from './PlanetaryProduction.svelte';
	import Research from './Research.svelte';

	type Props = {
		race: Race;
	};

	let { race = $bindable() }: Props = $props();
</script>

<TextInput name="name" bind:value={race.name} />
<TextInput name="pluralName" bind:value={race.pluralName} />
<EnumSelect
	name="spendLeftoverPointsOn"
	options={[
		SpendLeftoverPointsOnSurfaceMinerals,
		SpendLeftoverPointsOnMineralConcentrations,
		SpendLeftoverPointsOnMines,
		SpendLeftoverPointsOnFactories,
		SpendLeftoverPointsOnDefenses
	]}
	bind:value={race.spendLeftoverPointsOn}
/>

<SectionHeader>Primary Racial Trait</SectionHeader>
<EnumSelect
	name="prt"
	options={[HE, SS, WM, CA, IS, SD, PP, IT, AR, JoaT]}
	title="Primary Racial Trait"
	typeTitle={(prt) => getLabelForPRT(prt)}
	bind:value={race.prt}
/>
<div class="card bg-base-200 shadow">
	<div class="card-body">
		<PRTDescription prt={race.prt} />
	</div>
</div>

<SectionHeader>Lesser Racial Traits</SectionHeader>
<LRTs bind:race />
<LRTsDescriptions {race} />

<SectionHeader>Habitability</SectionHeader>
<Habitability bind:race />

<SectionHeader>Planetary Production</SectionHeader>
<PlanetaryProduction bind:race />

<SectionHeader>Research</SectionHeader>
<Research bind:race />

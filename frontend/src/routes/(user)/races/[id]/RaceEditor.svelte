<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import LRTsDescriptions from '$lib/components/game/race/LRTsDescriptions.svelte';
	import PRTDescription from '$lib/components/game/race/PRTDescription.svelte';
	import { PRT, SpendLeftoverPointsOn, getLabelForPRT, type Race } from '$lib/types/Race';
	import Habitability from './Habitability.svelte';
	import LRTs from './LRTs.svelte';
	import PlanetaryProduction from './PlanetaryProduction.svelte';
	import Research from './Research.svelte';

	export let race: Race;
</script>

<TextInput name="name" bind:value={race.name} />
<TextInput name="pluralName" bind:value={race.pluralName} />
<EnumSelect
	name="spendLeftoverPointsOn"
	enumType={SpendLeftoverPointsOn}
	bind:value={race.spendLeftoverPointsOn}
/>

<SectionHeader>Primary Racial Trait</SectionHeader>
<EnumSelect
	name="prt"
	title="Primary Racial Trait"
	enumType={PRT}
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
<LRTsDescriptions bind:race />

<SectionHeader>Habitability</SectionHeader>
<Habitability bind:race />

<SectionHeader>Planetary Production</SectionHeader>
<PlanetaryProduction bind:race />

<SectionHeader>Research</SectionHeader>
<Research bind:race />

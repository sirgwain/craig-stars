<script lang="ts">
	import { MaxNameLength } from '#lib/types/Consts.js';
	import EnumSelect from '#lib/components/EnumSelect.svelte';
	import LRTsDescriptions from '#lib/components/game/race/LRTsDescriptions.svelte';
	import PRTDescription from '#lib/components/game/race/PRTDescription.svelte';
	import SectionHeader from '#lib/components/SectionHeader.svelte';
	import TextInput from '#lib/components/TextInput.svelte';
	import { Prt, type Race, SpendLeftoverPointsOn } from '#lib/types/cs-proto.js';
	import { getLabelForPRT } from '#lib/types/Race.js';
	import { defaultRules } from '#lib/types/Rules.js';

	import Habitability from './Habitability.svelte';
	import LRTs from './LRTs.svelte';
	import PlanetaryProduction from './PlanetaryProduction.svelte';
	import Research from './Research.svelte';

	type Props = {
		race: Race;
		// the race's advantage points
		points?: number;
	};

	let { race = $bindable(), points = 0 }: Props = $props();

	// races saved without a choice spend leftover points on surface minerals
	if (race.spendLeftoverPointsOn === SpendLeftoverPointsOn.UNSPECIFIED) {
		race.spendLeftoverPointsOn = SpendLeftoverPointsOn.SURFACE_MINERALS;
	}

	// Leftover points are spent on a bonus for the homeworld. A bonus the race can't use, or doesn't
	// have the points for, is spent on surface minerals instead, so warn about it.
	const leftoverPointsItems: { [on: number]: string } = {
		[SpendLeftoverPointsOn.MINES]: 'mines',
		[SpendLeftoverPointsOn.FACTORIES]: 'factories',
		[SpendLeftoverPointsOn.DEFENSES]: 'defenses'
	};
	let leftoverPointsWarning = $derived.by(() => {
		const item = leftoverPointsItems[race.spendLeftoverPointsOn];
		const perItem =
			defaultRules.universeGenerationRules?.raceLeftoverPointsPerItem[race.spendLeftoverPointsOn];
		if (!item || !perItem || points <= 0) {
			return undefined;
		}
		const instead = 'Leftover points will be spent on surface minerals instead.';
		// Alternate Reality lives on starbases, with no mines, factories or defenses to add to
		if (race.prt === Prt.AR) {
			return `${getLabelForPRT(race.prt)} races can't use starting ${item}. ${instead}`;
		}
		// the first defense is half price
		const needed =
			race.spendLeftoverPointsOn === SpendLeftoverPointsOn.DEFENSES
				? Math.floor(perItem / 2)
				: perItem;
		return points < needed
			? `Starting ${item} take at least ${needed} points. ${instead}`
			: undefined;
	});
</script>

<TextInput name="name" bind:value={race.name} maxlength={MaxNameLength} />
<TextInput name="pluralName" bind:value={race.pluralName} maxlength={MaxNameLength} />
<EnumSelect
	name="spendLeftoverPointsOn"
	enumType={SpendLeftoverPointsOn}
	bind:value={race.spendLeftoverPointsOn}
/>
{#if leftoverPointsWarning}
	<div class="text-warning ml-35 mb-1" data-type="leftover-points-warning">
		{leftoverPointsWarning}
	</div>
{/if}

<SectionHeader>Primary Racial Trait</SectionHeader>
<EnumSelect
	name="prt"
	enumType={Prt}
	title="Primary Racial Trait"
	typeTitle={(prt) => getLabelForPRT(prt)}
	bind:value={race.prt}
/>
<div class="card bg-base-200 shadow-sm">
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

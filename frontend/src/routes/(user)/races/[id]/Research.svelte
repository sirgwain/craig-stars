<script lang="ts">
	import {
		Biotechnology,
		Construction,
		Electronics,
		Energy,
		JoaT,
		Propulsion,
		Weapons,
		type Race
	} from '$lib/types/cs';
	import ResearchCostField from './ResearchCostField.svelte';

	type Props = {
		race: Race;
	};

	let { race = $bindable() }: Props = $props();
</script>

<div class="flex flex-row flex-wrap justify-center gap-2">
	<ResearchCostField bind:value={race.researchCost.energy} field={Energy} />
	<ResearchCostField bind:value={race.researchCost.weapons} field={Weapons} />
	<ResearchCostField bind:value={race.researchCost.propulsion} field={Propulsion} />
	<ResearchCostField bind:value={race.researchCost.construction} field={Construction} />
	<ResearchCostField bind:value={race.researchCost.electronics} field={Electronics} />
	<ResearchCostField bind:value={race.researchCost.biotechnology} field={Biotechnology} />
</div>

<label class="label justify-start mt-2">
	<input
		class="checkbox"
		type="checkbox"
		name="techsStartHigh"
		bind:checked={race.techsStartHigh}
	/>
	<span class="ml-2"
		>All 'Costs 75% extra' research fields start at Tech {race.prt == JoaT ? '4' : '3'}</span
	>
</label>

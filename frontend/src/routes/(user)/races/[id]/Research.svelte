<script lang="ts">
	import type { Race } from '$lib/types/Race';
	import { TechField } from '$lib/types/TechLevel';
	import ResearchCostField from './ResearchCostField.svelte';

	export let race: Race;
</script>

<div class="flex flex-row flex-wrap justify-center gap-2">
	<ResearchCostField bind:value={race.researchCost.energy} field={TechField.Energy} />
	<ResearchCostField bind:value={race.researchCost.weapons} field={TechField.Weapons} />
	<ResearchCostField bind:value={race.researchCost.propulsion} field={TechField.Propulsion} />
	<ResearchCostField bind:value={race.researchCost.construction} field={TechField.Construction} />
	<ResearchCostField bind:value={race.researchCost.electronics} field={TechField.Electronics} />
	<ResearchCostField bind:value={race.researchCost.biotechnology} field={TechField.Biotechnology} />
</div>

<label class="label justify-start mt-2">
	<input
		class="checkbox"
		type="checkbox"
		name="techsStartHigh"
		bind:checked={race.techsStartHigh}
	/>
	<span class="ml-2">All 'Costs 75% extra' research fields start at Tech 4</span>
</label>

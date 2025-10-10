<script lang="ts">
	import { Prt, type Race, ResearchCostSchema, TechField } from '$lib/types/cs-proto';
	import { create } from '@bufbuild/protobuf';
	import ResearchCostField from './ResearchCostField.svelte';

	type Props = {
		race: Race;
	};

	let { race = $bindable() }: Props = $props();

	let researchCost = $state(create(ResearchCostSchema, race.researchCost));

	$effect(() => {
		race.researchCost = researchCost;
	});
</script>

<div class="flex flex-row flex-wrap justify-center gap-2">
	<ResearchCostField bind:value={researchCost.energy} field={TechField.ENERGY} />
	<ResearchCostField bind:value={researchCost.weapons} field={TechField.WEAPONS} />
	<ResearchCostField bind:value={researchCost.propulsion} field={TechField.PROPULSION} />
	<ResearchCostField bind:value={researchCost.construction} field={TechField.CONSTRUCTION} />
	<ResearchCostField bind:value={researchCost.electronics} field={TechField.ELECTRONICS} />
	<ResearchCostField bind:value={researchCost.biotechnology} field={TechField.BIOTECHNOLOGY} />
</div>

<label class="label justify-start mt-2">
	<input
		class="checkbox"
		type="checkbox"
		name="techsStartHigh"
		bind:checked={race.techsStartHigh}
	/>
	<span class="ml-2"
		>All 'Costs 75% extra' research fields start at Tech {race.prt == Prt.JOAT ? '4' : '3'}</span
	>
</label>

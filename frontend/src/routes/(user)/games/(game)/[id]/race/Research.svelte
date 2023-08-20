<script lang="ts">
	import { PRT, type Race } from '$lib/types/Race';
	import { TechField } from '$lib/types/TechLevel';
	import ResearchCostCard from './ResearchCostCard.svelte';

	export let race: Race;
</script>

<div class="flex flex-row flex-wrap justify-center gap-2">
	<ResearchCostCard value={race.researchCost.energy} field={TechField.Energy} />
	<ResearchCostCard value={race.researchCost.weapons} field={TechField.Weapons} />
	<ResearchCostCard value={race.researchCost.propulsion} field={TechField.Propulsion} />
	<ResearchCostCard value={race.researchCost.construction} field={TechField.Construction} />
	<ResearchCostCard value={race.researchCost.electronics} field={TechField.Electronics} />
	<ResearchCostCard value={race.researchCost.biotechnology} field={TechField.Biotechnology} />
</div>

<div class="mt-2"
	>All 'Costs 75% extra' research fields start at <span class="text-accent"
		>Tech {race.prt == PRT.JoaT ? '4' : '3'}</span
	></div
>

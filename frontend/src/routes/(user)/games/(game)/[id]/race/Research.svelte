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
	import ResearchCostCard from './ResearchCostCard.svelte';

	type Props = {
		race: Race;
	};

	let { race }: Props = $props();
</script>

<div class="flex flex-row flex-wrap justify-center gap-2">
	<ResearchCostCard value={race.researchCost.energy} field={Energy} />
	<ResearchCostCard value={race.researchCost.weapons} field={Weapons} />
	<ResearchCostCard value={race.researchCost.propulsion} field={Propulsion} />
	<ResearchCostCard value={race.researchCost.construction} field={Construction} />
	<ResearchCostCard value={race.researchCost.electronics} field={Electronics} />
	<ResearchCostCard value={race.researchCost.biotechnology} field={Biotechnology} />
</div>

{#if race.techsStartHigh}
	<div class="mt-2">
		All 'Costs 75% extra' research fields start at <span class="text-primary font-semibold"
			>Tech {race.prt == JoaT ? '4' : '3'}</span
		>
	</div>
{/if}

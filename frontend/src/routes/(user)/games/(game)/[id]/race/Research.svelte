<script lang="ts">
	import { Prt, ResearchCostLevel, TechField, type Race } from '$lib/types/cs-proto';
	import ResearchCostCard from './ResearchCostCard.svelte';

	type Props = {
		race: Race;
	};

	let { race }: Props = $props();
</script>

<div class="flex flex-row flex-wrap justify-center gap-2">
	<ResearchCostCard
		value={race.researchCost?.energy ?? ResearchCostLevel.STANDARD}
		field={TechField.ENERGY}
	/>
	<ResearchCostCard
		value={race.researchCost?.weapons ?? ResearchCostLevel.STANDARD}
		field={TechField.WEAPONS}
	/>
	<ResearchCostCard
		value={race.researchCost?.propulsion ?? ResearchCostLevel.STANDARD}
		field={TechField.PROPULSION}
	/>
	<ResearchCostCard
		value={race.researchCost?.construction ?? ResearchCostLevel.STANDARD}
		field={TechField.CONSTRUCTION}
	/>
	<ResearchCostCard
		value={race.researchCost?.electronics ?? ResearchCostLevel.STANDARD}
		field={TechField.ELECTRONICS}
	/>
	<ResearchCostCard
		value={race.researchCost?.biotechnology ?? ResearchCostLevel.STANDARD}
		field={TechField.BIOTECHNOLOGY}
	/>
</div>

{#if race.techsStartHigh}
	<div class="mt-2">
		All 'Costs 75% extra' research fields start at <span class="text-primary font-semibold"
			>Tech {race.prt == Prt.JOAT ? '4' : '3'}</span
		>
	</div>
{/if}

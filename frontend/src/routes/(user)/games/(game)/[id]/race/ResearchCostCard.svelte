<script lang="ts">
	import { startCase } from 'lodash-es';
	import { ResearchCostLevel, type Race } from '$lib/types/Race';
	import type { TechField } from '$lib/types/TechLevel';

	export let field: TechField;
	export let value: ResearchCostLevel;
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300 w-full md:w-48">
	<div class="card-body p-3 gap-0">
		<h2 class="text-lg font-semibold text-center mb-1 text-secondary">
			{startCase(field)} Research
		</h2>
		<div>
			{#if value === ResearchCostLevel.Extra}
				<span class="">Costs 75% extra</span>
			{:else if value === ResearchCostLevel.Standard}
				<span class="">Costs standard amount</span>
			{:else}
				<span class="">Costs 50% less</span>
			{/if}
		</div>
	</div>
</div>

<script lang="ts">
	import {
		ResearchCostExtra,
		ResearchCostStandard,
		type ResearchCostLevel,
		type TechField
	} from '$lib/types/cs';
	import { startCase } from 'lodash-es';

	type Props = {
		field: TechField;
		value: ResearchCostLevel;
	};

	let { field, value }: Props = $props();
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300 w-full md:w-48">
	<div class="card-body p-3 gap-0">
		<h2 class="text-lg font-semibold text-center mb-1 text-secondary">
			{startCase(field)} Research
		</h2>
		<div>
			{#if value === ResearchCostExtra}
				<span class="">Costs 75% extra</span>
			{:else if value === ResearchCostStandard}
				<span class="">Costs standard amount</span>
			{:else}
				<span class="">Costs 50% less</span>
			{/if}
		</div>
	</div>
</div>

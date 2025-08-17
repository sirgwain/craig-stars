<script lang="ts">
	import { ResearchCostLevel, TechField } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';

	type Props = {
		field: TechField;
		value: ResearchCostLevel;
	};

	let { field, value }: Props = $props();
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300 w-full md:w-48">
	<div class="card-body p-3 gap-0">
		<h2 class="text-lg font-semibold text-center mb-1 text-secondary">
			{enumToString(TechField, field)} Research
		</h2>
		<div>
			{#if value === ResearchCostLevel.EXTRA}
				<span class="">Costs 75% extra</span>
			{:else if value === ResearchCostLevel.STANDARD}
				<span class="">Costs standard amount</span>
			{:else}
				<span class="">Costs 50% less</span>
			{/if}
		</div>
	</div>
</div>

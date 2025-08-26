<script lang="ts">
	import { ResearchCostLevel, TechField } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';

	type Props = {
		field: TechField;
		value: ResearchCostLevel;
	};

	let { field, value = $bindable() }: Props = $props();
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300 w-full md:w-auto">
	<div class="card-body p-3 gap-0">
		<h2 class="text-lg font-semibold text-center mb-1 text-secondary">
			{enumToString(TechField, field)} Research
		</h2>
		<div>
			<label class="label cursor-pointer justify-start">
				<input
					type="radio"
					name={`${enumToString(TechField, field)}ResearchCost`}
					value={ResearchCostLevel.EXTRA}
					class="radio radio-sm checked:bg-primary"
					bind:group={value}
				/>
				<span class="label-text ml-2">Costs 75% extra</span>
			</label>
			<label class="label cursor-pointer justify-start">
				<input
					type="radio"
					name={`${enumToString(TechField, field)}ResearchCost`}
					value={ResearchCostLevel.STANDARD}
					class="radio radio-sm checked:bg-primary"
					bind:group={value}
				/>
				<span class="label-text ml-2">Costs standard amount</span>
			</label>
			<label class="label cursor-pointer justify-start">
				<input
					type="radio"
					name={`${enumToString(TechField, field)}ResearchCost`}
					value={ResearchCostLevel.LESS}
					class="radio radio-sm checked:bg-primary"
					bind:group={value}
				/>
				<span class="label-text ml-2">Costs 50% less</span>
			</label>
		</div>
	</div>
</div>

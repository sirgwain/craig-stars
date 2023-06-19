<script lang="ts">
	import { startCase } from 'lodash-es';
	import type { TechField } from '$lib/types/Player';
	import { ResearchCostLevel, type Race } from '$lib/types/Race';

	export let field: TechField;
	export let value: ResearchCostLevel;
</script>

<div class="card bg-base-200 shadow rounded-sm border-2 border-base-300 w-full md:w-auto">
	<div class="card-body p-3 gap-0">
		<h2 class="text-lg font-semibold text-center mb-1 text-secondary">
			{startCase(field)} Research
		</h2>
		<div>
			<label class="label cursor-pointer justify-start">
				<input
					type="radio"
					name={`${field}ResearchCost`}
					value={ResearchCostLevel.Extra}
					class="radio radio-sm checked:bg-primary"
					bind:group={value}
				/>
				<span class="label-text ml-2">Costs 75% extra</span>
			</label>
			<label class="label cursor-pointer justify-start">
				<input
					type="radio"
					name={`${field}ResearchCost`}
					value={ResearchCostLevel.Standard}
					class="radio radio-sm checked:bg-primary"
					bind:group={value}
				/>
				<span class="label-text ml-2">Costs stamdard amount</span>
			</label>
			<label class="label cursor-pointer justify-start">
				<input
					type="radio"
					name={`${field}ResearchCost`}
					value={ResearchCostLevel.Less}
					class="radio radio-sm checked:bg-primary"
					bind:group={value}
				/>
				<span class="label-text ml-2">Costs 50% less</span>
			</label>
		</div>
	</div>
</div>

<script lang="ts">
	import TextInput from '$lib/components/TextInput.svelte';
	import type { DesignFinder } from '$lib/services/Universe';
	import { fromQueueItemType } from '$lib/types/Planet';
	import type { ProductionPlan } from '$lib/types/Player';
	import type { ProductionQueueItem } from '$lib/types/Production';
	import { QueueItemTypes } from '$lib/types/QueueItemType';
	import Production from './Production.svelte';

	export let designFinder: DesignFinder;
	export let plan: ProductionPlan;

	let availableItems: ProductionQueueItem[] = [
		fromQueueItemType(QueueItemTypes.AutoFactories),
		fromQueueItemType(QueueItemTypes.AutoMines),
		fromQueueItemType(QueueItemTypes.AutoDefenses),
		fromQueueItemType(QueueItemTypes.AutoMineralAlchemy),
		fromQueueItemType(QueueItemTypes.AutoMaxTerraform),
		fromQueueItemType(QueueItemTypes.AutoMinTerraform)
	];
</script>

<TextInput name="name" bind:value={plan.name} required />

<!-- edit production -->
<Production {designFinder} {availableItems} bind:queueItems={plan.items} />
<div class="w-1/2 mr-14">
	<label>
		<input
			bind:checked={plan.contributesOnlyLeftoverToResearch}
			class="checkbox checkbox-xs"
			type="checkbox"
		/> Contributes Only Leftover to Research
	</label>
</div>

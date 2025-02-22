<script lang="ts">
	import TextInput from '$lib/components/TextInput.svelte';
	import type { DesignFinder } from '$lib/services/Universe';
	import type { ProductionQueueItem } from '$lib/types/cs';
	import {
		QueueItemTypeAutoDefenses,
		QueueItemTypeAutoFactories,
		QueueItemTypeAutoMaxTerraform,
		QueueItemTypeAutoMineralAlchemy,
		QueueItemTypeAutoMines,
		QueueItemTypeAutoMinTerraform,
		type ProductionPlan
	} from '$lib/types/cs';
	import { fromQueueItemType } from '$lib/types/Planet';
	import Production from './Production.svelte';

	type Props = {
		designFinder: DesignFinder;
		plan: ProductionPlan;
	};

	let { designFinder, plan = $bindable() }: Props = $props();

	let availableItems: ProductionQueueItem[] = [
		fromQueueItemType(QueueItemTypeAutoFactories),
		fromQueueItemType(QueueItemTypeAutoMines),
		fromQueueItemType(QueueItemTypeAutoDefenses),
		fromQueueItemType(QueueItemTypeAutoMineralAlchemy),
		fromQueueItemType(QueueItemTypeAutoMaxTerraform),
		fromQueueItemType(QueueItemTypeAutoMinTerraform)
	];
</script>

<TextInput name="name" bind:value={plan.name} required />

<!-- edit production -->
<Production {designFinder} {availableItems} bind:queueItems={plan.items as ProductionQueueItem[]} />
<div class="w-1/2 mr-14">
	<label>
		<input
			bind:checked={plan.contributesOnlyLeftoverToResearch}
			class="checkbox checkbox-xs"
			type="checkbox"
		/> Contributes Only Leftover to Research
	</label>
</div>

<script lang="ts">
	import TextInput from '$lib/components/TextInput.svelte';
	import type { DesignFinder } from '$lib/services/Universe';
	import { QueueItemType, fromQueueItemType, type ProductionQueueItem } from '$lib/types/Planet';
	import type { ProductionPlan } from '$lib/types/Player';
	import Production from './Production.svelte';

	export let designFinder: DesignFinder;
	export let plan: ProductionPlan;

	let availableItems: ProductionQueueItem[] = [
		fromQueueItemType(QueueItemType.AutoFactories),
		fromQueueItemType(QueueItemType.AutoMines),
		fromQueueItemType(QueueItemType.AutoDefenses),
		fromQueueItemType(QueueItemType.AutoMineralAlchemy),
		fromQueueItemType(QueueItemType.AutoMaxTerraform),
		fromQueueItemType(QueueItemType.AutoMinTerraform)
	];
</script>

<TextInput name="name" bind:value={plan.name} required />

<!-- edit production -->
<Production {designFinder} {availableItems} bind:queueItems={plan.items} />

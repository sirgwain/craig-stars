<script lang="ts">
	import TextInput from '$lib/components/TextInput.svelte';
	import type { DesignFinder } from '$lib/services/Universe';
	import { QueueItemType, type ProductionPlan, type ProductionPlanItem } from '$lib/types/cs-proto';
	import { planItemFromQueueItemType } from '$lib/types/Player';
	import Production from './Production.svelte';

	type Props = {
		designFinder: DesignFinder;
		plan: ProductionPlan;
	};

	// plan is bindable from parent
	let { designFinder, plan = $bindable() }: Props = $props();

	// Local Svelte 5 runes state mirroring plan fields to allow binding to reactive values
	let name: string = $state(plan.name);
	let items: ProductionPlanItem[] = $state(plan.items);
	let contributesOnlyLeftoverToResearch: boolean = $state(plan.contributesOnlyLeftoverToResearch);

	// Keep parent prop in sync with local state (runes-compliant)
	$effect(() => {
		plan.name = name;
		plan.items = items;
		plan.contributesOnlyLeftoverToResearch = contributesOnlyLeftoverToResearch;
	});

	let availableItems: ProductionPlanItem[] = [
		planItemFromQueueItemType(QueueItemType.AUTO_FACTORIES),
		planItemFromQueueItemType(QueueItemType.AUTO_MINES),
		planItemFromQueueItemType(QueueItemType.AUTO_DEFENSES),
		planItemFromQueueItemType(QueueItemType.AUTO_MINERAL_ALCHEMY),
		planItemFromQueueItemType(QueueItemType.AUTO_MAX_TERRAFORM),
		planItemFromQueueItemType(QueueItemType.AUTO_MIN_TERRAFORM)
	];
</script>

<TextInput name="name" bind:value={name} required />

<!-- edit production -->
<Production {designFinder} {availableItems} bind:queueItems={items} />
<div class="w-1/2 mr-14">
	<label>
		<input
			bind:checked={contributesOnlyLeftoverToResearch}
			class="checkbox checkbox-xs"
			type="checkbox"
		/> Contributes Only Leftover to Research
	</label>
</div>

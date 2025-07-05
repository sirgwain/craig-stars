<script lang="ts">
	import type { DesignFinder } from '$lib/services/Universe';
	import { getQueueItemShortName } from '$lib/types/Planet';
	import type { ProductionPlan, ProductionQueueItem } from '$lib/types/cs';
	import { isAuto } from '$lib/types/QueueItemType';
	import { Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	type Props = {
		designFinder: DesignFinder;
		plan: ProductionPlan;
		href: string;
		showDelete?: boolean;
		onDelete?: (plan: ProductionPlan) => void;
	};

	let { designFinder, plan, href, showDelete = true, onDelete }: Props = $props();

	const deletePlan = async (plan: ProductionPlan) => {
		if (plan.name != undefined && confirm(`Are you sure you want to delete ${plan.name}?`)) {
			onDelete?.(plan);
		}
	};
</script>

<div
	class="card bg-base-200 shadow rounded-sm border-2 border-base-300 pt-2 m-1 w-full sm:w-[350px]"
>
	<div class="card-body">
		<h2 class="card-title">
			<a class="cs-link" {href}>{plan.name}</a>
		</h2>
		<div class="flex flex-col gap-2">
			<div class="flex flex-row">
				<div class="font-semibold mr-2">Name</div>
				<div>{plan.name}</div>
			</div>
			<ul class="w-full h-full">
				{#each plan.items as queueItem, index (index)}
					<li class="pl-1">
						<div class="flex flex-row justify-between" class:italic={isAuto(queueItem.type)}>
							<div>
								{getQueueItemShortName(queueItem as ProductionQueueItem, designFinder)}
							</div>
							<div>
								{queueItem.quantity}
							</div>
						</div>
					</li>
				{/each}
			</ul>
			<div>
				{#if plan.contributesOnlyLeftoverToResearch}
					Planet contributes only leftover resources to research
				{/if}
			</div>
		</div>
		{#if showDelete}
			<div class="card-actions justify-start">
				<div>
					<button
						type="button"
						class="btn"
						onclick={() => deletePlan(plan)}
						data-type="delete-button"
						data-id={`${plan.name}`}
					>
						<Icon src={Trash} size="24" class="hover:stroke-accent" />
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>

<script lang="ts">
	import type { DesignFinder } from '$lib/services/Universe';
	import { getQueueItemShortName } from '$lib/types/Planet';
	import type { ProductionPlan } from '$lib/types/Player';
	import { isAuto } from '$lib/types/QueueItemType';
	import { Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let designFinder: DesignFinder;
	export let plan: ProductionPlan;
	export let href: string;
	export let showDelete = true;

	const deletePlan = async (plan: ProductionPlan) => {
		if (plan.name != undefined && confirm(`Are you sure you want to delete ${plan.name}?`)) {
			dispatch('delete', { plan });
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
				{#each plan.items as queueItem}
					<li class="pl-1">
						<div class="flex flex-row justify-between" class:italic={isAuto(queueItem.type)}>
							<div>
								{getQueueItemShortName(queueItem, designFinder)}
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
					<button class="btn" on:click={(e) => deletePlan(plan)}>
						<Icon src={Trash} size="24" class="hover:stroke-accent" />
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>

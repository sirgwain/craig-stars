<script lang="ts">
	import { WaypointTaskTransportAction } from '$lib/types/Fleet';
	import type { TransportPlan } from '$lib/types/Player';
	import { Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';
	import TransportActionDescription from './TransportActionDescription.svelte';

	const dispatch = createEventDispatcher();

	export let plan: TransportPlan;
	export let href: string;
	export let showDelete = true;

	const deletePlan = async (plan: TransportPlan) => {
		if (plan.name != undefined && confirm(`Are you sure you want to delete ${plan.name}?`)) {
			dispatch('delete', { plan });
		}
	};

	const isEmpty = (plan: TransportPlan) =>
		(plan.tasks.fuel.action ?? WaypointTaskTransportAction.None) ==
			WaypointTaskTransportAction.None &&
		(plan.tasks.ironium.action ?? WaypointTaskTransportAction.None) ==
			WaypointTaskTransportAction.None &&
		(plan.tasks.boranium.action ?? WaypointTaskTransportAction.None) ==
			WaypointTaskTransportAction.None &&
		(plan.tasks.germanium.action ?? WaypointTaskTransportAction.None) ==
			WaypointTaskTransportAction.None &&
		(plan.tasks.colonists.action ?? WaypointTaskTransportAction.None) ==
			WaypointTaskTransportAction.None;
</script>

<div
	class="card bg-base-200 shadow rounded-sm border-2 border-base-300 pt-2 m-1 w-full sm:w-[350px]"
>
	<div class="card-body">
		<h2 class="card-title">
			<a class="cs-link" {href}>{plan.name}</a>
		</h2>
		<div class="flex flex-col">
			<div class="flex flex-row">
				<div class="text-right font-semibold mr-2 w-28">Name</div>
				<div>{plan.name}</div>
			</div>
			{#if isEmpty(plan)}
				<div class="flex flex-row">
					<div class="text-right font-semibold mr-2 w-28">Actions</div>
					<div>None</div>
				</div>
			{:else}
				<TransportActionDescription
					action={plan.tasks.fuel.action}
					amount={plan.tasks.fuel.amount}
					units="mg"
					title="Fuel"
					titleTextClass="text-fuel"
				/>
				<TransportActionDescription
					action={plan.tasks.ironium.action}
					amount={plan.tasks.ironium.amount}
					units="kT"
					title="Ironium"
					titleTextClass="text-ironium"
				/>
				<TransportActionDescription
					action={plan.tasks.boranium.action}
					amount={plan.tasks.boranium.amount}
					units="kT"
					title="Boranium"
					titleTextClass="text-boranium"
				/>
				<TransportActionDescription
					action={plan.tasks.germanium.action}
					amount={plan.tasks.germanium.amount}
					units="kT"
					title="Germanium"
					titleTextClass="text-germanium"
				/>
				<TransportActionDescription
					action={plan.tasks.colonists.action}
					amount={plan.tasks.colonists.amount}
					units=""
					title="Colonists"
					titleTextClass="text-colonists"
				/>
			{/if}
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

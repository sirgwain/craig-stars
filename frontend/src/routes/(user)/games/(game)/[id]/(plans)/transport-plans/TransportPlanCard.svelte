<script lang="ts">
	import type { TransportPlan } from '$lib/types/cs';
	import {
		TransportActionFillPercent,
		TransportActionNone,
		TransportActionWaitForPercent
	} from '$lib/types/cs';
	import { Trash } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import TransportActionDescription from './TransportActionDescription.svelte';

	type Props = {
		plan: TransportPlan;
		href: string;
		showDelete?: boolean;
		onDelete?: (plan: TransportPlan) => void;
	};

	let { plan, href, showDelete = true, onDelete }: Props = $props();

	const deletePlan = async (plan: TransportPlan) => {
		if (plan.name != undefined && confirm(`Are you sure you want to delete ${plan.name}?`)) {
			onDelete?.(plan);
		}
	};

	const isEmpty = (plan: TransportPlan) =>
		(plan.tasks?.fuel.action ?? TransportActionNone) == TransportActionNone &&
		(plan.tasks?.ironium.action ?? TransportActionNone) == TransportActionNone &&
		(plan.tasks?.boranium.action ?? TransportActionNone) == TransportActionNone &&
		(plan.tasks?.germanium.action ?? TransportActionNone) == TransportActionNone &&
		(plan.tasks?.colonists.action ?? TransportActionNone) == TransportActionNone;
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
					action={plan.tasks?.fuel.action}
					amount={plan.tasks?.fuel.amount}
					units="mg"
					title="Fuel"
					titleTextClass="text-fuel"
				/>
				<TransportActionDescription
					action={plan.tasks?.ironium.action}
					amount={plan.tasks?.ironium.amount}
					units={[TransportActionWaitForPercent, TransportActionFillPercent].indexOf(
						plan.tasks?.ironium.action ?? TransportActionNone
					) != -1
						? '%'
						: 'kT'}
					title="Ironium"
					titleTextClass="text-ironium"
				/>
				<TransportActionDescription
					action={plan.tasks?.boranium.action}
					amount={plan.tasks?.boranium.amount}
					units={[TransportActionWaitForPercent, TransportActionFillPercent].indexOf(
						plan.tasks?.boranium.action ?? TransportActionNone
					) != -1
						? '%'
						: 'kT'}
					title="Boranium"
					titleTextClass="text-boranium"
				/>
				<TransportActionDescription
					action={plan.tasks?.germanium.action}
					amount={plan.tasks?.germanium.amount}
					units={[TransportActionWaitForPercent, TransportActionFillPercent].indexOf(
						plan.tasks?.germanium.action ?? TransportActionNone
					) != -1
						? '%'
						: 'kT'}
					title="Germanium"
					titleTextClass="text-germanium"
				/>
				<TransportActionDescription
					action={plan.tasks?.colonists.action}
					amount={plan.tasks?.colonists.amount}
					units={[TransportActionWaitForPercent, TransportActionFillPercent].indexOf(
						plan.tasks?.colonists.action ?? TransportActionNone
					) != -1
						? '%'
						: '00'}
					title="Colonists"
					titleTextClass="text-colonists"
				/>
			{/if}
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

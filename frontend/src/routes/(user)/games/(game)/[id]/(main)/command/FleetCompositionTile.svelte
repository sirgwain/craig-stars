<script lang="ts">
	import { EventManager } from '$lib/EventManager';
	import { selectedWaypoint, showDesignPopup } from '$lib/services/Context';
	import type { CommandedFleet } from '$lib/types/Fleet';
	import type { Player } from '$lib/types/Player';
	import { createEventDispatcher } from 'svelte';
	import CommandTile from './CommandTile.svelte';

	const dispatch = createEventDispatcher();

	export let fleet: CommandedFleet;
	export let player: Player;

	const split = () => {
		EventManager.publishSplitFleetDialogRequestedEvent(fleet);
	};

	const splitAll = async () => {
		dispatch('splitAll');
	};
	const merge = () => {
		EventManager.publishMergeFleetDialogRequestedEvent(fleet);
	};
</script>

{#if fleet.waypoints && $selectedWaypoint}
	<CommandTile title="Fleet Composition">
		<div class="bg-base-100 h-20 overflow-y-auto">
			<ul class="w-full h-full">
				{#each fleet.tokens as token, index}
					<li class="pl-1">
						<button
							type="button"
							class="w-full cursor-help"
							on:pointerdown|preventDefault={(e) =>
								showDesignPopup(player.getDesign(player.num, token.designNum), e.x, e.y)}
						>
							<div class="flex flex-row justify-between">
								<div>
									{player.getDesign(player.num, token.designNum)?.name}
								</div>
								<div>
									{token.quantity}
								</div>
							</div>
						</button>
					</li>
				{/each}
			</ul>
		</div>
		<div class="flex justify-between my-1">
			<div class="my-auto">Battle Plan:</div>
			<div>
				<select
					class="select select-outline select-secondary select-sm"
					name="battlePlan"
					bind:value={fleet.battlePlanName}
				>
					{#each player.battlePlans as battlePlan}
						<option value={battlePlan.name}>{battlePlan.name}</option>
					{/each}
				</select>
			</div>
		</div>
		<div class="flex justify-between my-1">
			<div>Est Range:</div>
			<div>--</div>
		</div>
		<div class="flex justify-between my-1">
			<div>Percent Cloacked</div>
			<div>{fleet.spec.cloakPercent ? fleet.spec.cloakPercent * 100 + '%' : 'none'}</div>
		</div>
		<div class="flex justify-between">
			<button on:click={split} class="btn btn-outline btn-sm normal-case btn-secondary"
				>Split</button
			>
			<button on:click={splitAll} class="btn btn-outline btn-sm normal-case btn-secondary"
				>Split All</button
			>
			<button on:click={merge} class="btn btn-outline btn-sm normal-case btn-secondary"
				>Merge</button
			>
		</div>
	</CommandTile>
{/if}

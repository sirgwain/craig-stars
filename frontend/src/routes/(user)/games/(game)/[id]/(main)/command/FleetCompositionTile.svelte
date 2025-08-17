<script lang="ts">
	import { onShipDesignTooltip } from '$lib/components/game/tooltips/ShipDesignTooltip.svelte';
	import type { Waypoint } from '$lib/types/cs-proto';
	import type {
		BattlePlanChangedProps,
		ShowMergeFleetsDialogProps,
		ShowSplitFleetDialogProps,
		SplitAllProps
	} from '$lib/services/Events';
	import { getGameContext } from '$lib/services/GameContext';
	import { Infinite } from '$lib/types/Consts';
	import { getDamagePercentForToken, type CommandedFleet } from '$lib/types/Fleet';
	import CommandTile from './CommandTile.svelte';

	const { commandedFleet, player, universe } = getGameContext();

	type Props = {
		fleet: CommandedFleet;
		selectedWaypoint: Waypoint | undefined;
	} & ShowSplitFleetDialogProps &
		ShowMergeFleetsDialogProps &
		SplitAllProps &
		BattlePlanChangedProps;

	let {
		fleet,
		selectedWaypoint,
		onShowMergeFleetDialog,
		onShowSplitFleetDialog,
		onSplitAll,
		onBattlePlanChanged
	}: Props = $props();

	function split() {
		if (!onShowSplitFleetDialog || !$commandedFleet) {
			return;
		}
		onShowSplitFleetDialog({ src: $commandedFleet });
	}

	function splitAll() {
		if (!onSplitAll) {
			return;
		}
		onSplitAll({ fleet });
	}

	function merge() {
		if (!onShowMergeFleetDialog) {
			return;
		}
		onShowMergeFleetDialog({
			fleet,
			otherFleetsHere: $universe
				.getMyFleetsByPosition(fleet.mapObject.position)
				.filter((f) => f.mapObject?.num !== fleet.mapObject.num)
		});
	}

	function updateBattlePlan(battlePlanNum: number) {
		fleet.fleetOrders.battlePlanNum = battlePlanNum;
		onBattlePlanChanged?.({ fleet, battlePlanNum });
	}
</script>

{#if fleet.fleetOrders?.waypoints && selectedWaypoint}
	<CommandTile title="Fleet Composition">
		<div class="bg-base-100 h-20 overflow-y-auto">
			<ul class="w-full h-full">
				{#each fleet.tokens as token, index (index)}
					<li class="pl-1">
						<button
							type="button"
							class="w-full cursor-help"
							onpointerdown={(e) =>
								onShipDesignTooltip(e, $universe.getDesign($player.num, token.designNum))}
						>
							<div class="flex flex-row justify-between relative">
								{#if (token.damage ?? 0) > 0 && (token.quantityDamaged ?? 0) > 0}
									<div
										style={`width: ${getDamagePercentForToken(
											token,
											$universe.getMyDesign(token.designNum)
										).toFixed()}%`}
										class="damage-bar h-full absolute opacity-50"
									></div>
								{/if}
								<div>
									{$universe.getDesign($player.num, token.designNum)?.name}
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
			<div class="my-auto text-tile-item-title">Battle Plan:</div>
			<div>
				<select
					class="select select-outline select-secondary select-sm text-sm"
					name="battlePlan"
					value={fleet.fleetOrders?.battlePlanNum ?? 0}
					onchange={(e) => updateBattlePlan(parseInt(e.currentTarget.value))}
				>
					{#each $player.playerPlans.battlePlans as battlePlan (battlePlan.num)}
						<option value={battlePlan.num}>{battlePlan.name}</option>
					{/each}
				</select>
			</div>
		</div>
		<div class="flex justify-between my-1">
			<div class="text-tile-item-title">Est Range:</div>
			<div>
				{fleet.spec.shipDesignSpec?.estimatedRange
					? fleet.spec.shipDesignSpec?.estimatedRange === Infinite
						? 'Infinite'
						: `${fleet.spec.shipDesignSpec?.estimatedRange} l.y.`
					: '--'}
			</div>
		</div>
		<div class="flex justify-between my-1">
			<div class="text-tile-item-title">Percent Cloaked</div>
			<div>
				{fleet.spec.shipDesignSpec?.cloakPercent
					? fleet.spec.shipDesignSpec?.cloakPercent + '%'
					: 'none'}
			</div>
		</div>
		<div class="flex justify-between">
			<button onclick={split} class="btn btn-outline btn-sm normal-case btn-secondary">Split</button
			>
			<button onclick={splitAll} class="btn btn-outline btn-sm normal-case btn-secondary"
				>Split All</button
			>
			<button onclick={merge} class="btn btn-outline btn-sm normal-case btn-secondary">Merge</button
			>
		</div>
	</CommandTile>
{/if}

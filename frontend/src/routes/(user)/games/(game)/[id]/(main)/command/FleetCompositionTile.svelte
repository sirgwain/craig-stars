<script lang="ts">
	import TagBadge from '#lib/components/game/TagBadge.svelte';
	import { onShipDesignTooltip } from '#lib/components/game/tooltips/ShipDesignTooltip.js';
	import type { Waypoint } from '#lib/types/cs-proto.js';
	import type {
		BattlePlanChangedProps,
		ShowMergeFleetsDialogProps,
		ShowSplitFleetDialogProps,
		SplitAllProps
	} from '#lib/services/Events.js';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { Infinite } from '#lib/types/Consts.js';
	import {
		getDamagePercentForToken,
		getTokenDamage,
		type CommandedFleet
	} from '#lib/types/Fleet.js';
	import CommandTile from './CommandTile.svelte';
	import { emptyVector } from '#lib/types/Vector.js';
	import { PencilSquare } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { commandedFleet, player, universe, readOnly, editNote } = getGameContext();

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

	// A fleet in transit can't be selected on the map, only its waypoints, so its notes are edited here
	let note = $derived($universe.getNote(fleet));

	function split() {
		if (!onShowSplitFleetDialog || !$commandedFleet) {
			return;
		}
		onShowSplitFleetDialog({ src: $commandedFleet });
	}

	function splitAll() {
		onSplitAll({ fleet });
	}

	function merge() {
		if (!onShowMergeFleetDialog) {
			return;
		}
		onShowMergeFleetDialog({
			fleet,
			otherFleetsHere: $universe
				.getMyFleetsByPosition(fleet.mapObject.position ?? emptyVector())
				.filter((f) => f.mapObject?.num !== fleet.mapObject.num)
		});
	}

	function updateBattlePlan(battlePlanNum: number) {
		fleet.fleetOrders.battlePlanNum = battlePlanNum;
		onBattlePlanChanged?.({ fleet, battlePlanNum });
	}
</script>

{#if selectedWaypoint}
	<CommandTile title="Fleet Composition">
		<div class="bg-base-100 h-20 overflow-y-auto">
			<ul class="size-full">
				{#each fleet.tokens as token, index (index)}
					<li class="pl-1">
						<button
							type="button"
							class="w-full cursor-help"
							onpointerdown={(e) => {
								const design = $universe.getDesign($player.num, token.designNum);
								onShipDesignTooltip(e, design, getTokenDamage(token, design));
							}}
						>
							<div class="flex flex-row justify-between relative">
								{#if token.damage > 0 && token.quantityDamaged > 0}
									<div
										style={`width: ${getDamagePercentForToken(
											token,
											$universe.getMyDesign(token.designNum)
										).toFixed()}%`}
										class="damage-bar h-full absolute opacity-50"
									></div>
								{/if}
								<div class="min-w-0 truncate">
									{$universe.getDesign($player.num, token.designNum)?.name}
								</div>
								<div class="shrink-0 pl-1">
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
					value={fleet.fleetOrders.battlePlanNum}
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
					? fleet.spec.shipDesignSpec.estimatedRange === Infinite
						? 'Infinite'
						: `${fleet.spec.shipDesignSpec.estimatedRange} l.y.`
					: '--'}
			</div>
		</div>
		<div class="flex justify-between my-1">
			<div class="text-tile-item-title">Percent Cloaked</div>
			<div>
				{fleet.spec.shipDesignSpec?.cloakPercent
					? fleet.spec.shipDesignSpec.cloakPercent + '%'
					: 'none'}
			</div>
		</div>
		{#if note || !$readOnly}
			<button
				type="button"
				class="flex flex-row items-center gap-1 w-full my-1 overflow-hidden text-left cursor-pointer hover:text-accent"
				title={note?.text || undefined}
				aria-label={note ? 'Edit notes and tags' : 'Add notes and tags'}
				data-type="fleet-note-button"
				onclick={() => editNote(fleet)}
			>
				<Icon src={PencilSquare} class={`size-4 flex-none ${note ? 'stroke-accent' : ''}`} />
				{#each note?.tags ?? [] as tag (tag)}
					<TagBadge {tag} class="flex-none" />
				{/each}
				<span class="truncate italic" data-type="fleet-note"
					>{note ? note.text : 'Add notes or tags'}</span
				>
			</button>
		{/if}
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

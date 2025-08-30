<script lang="ts">
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { designFinderKey, getGameContext, playerFinderKey } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';
	import type { DesignFinder, PlayerFinder } from '$lib/services/Universe';
	import type { Battle, PhaseToken } from '$lib/types/Battle';
	import { enumToString } from '$lib/types/Enums';
	import {
		BattleTactic,
		BattleTacticSchema,
		BattleTarget,
		BattleTargetSchema
	} from '$lib/types/cs-proto';
	import { enumFromJson } from '@bufbuild/protobuf';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { getContext } from 'svelte';
	import { onShipDesignTooltip } from '../tooltips/ShipDesignTooltip.svelte';
	import { getDisplayColor } from '$lib/utils/colorUtils';

	const { player, universe, settings } = getGameContext();

	const designFinder = getContext<DesignFinder>(designFinderKey);
	const playerFinder = getContext<PlayerFinder>(playerFinderKey);

	type Props = {
		battle: Battle;
		phase: number;
		token: PhaseToken | undefined;
	};

	let { battle, phase, token }: Props = $props();

	let design = $derived(token && designFinder.getDesign(token.playerNum, token.designNum));
	let raceName = $derived(token && playerFinder.getPlayerIntel(token.playerNum)?.racePluralName);
	let tokenState = $derived(token && battle.getTokenForPhase(token.num, phase));
	let armor = $derived(design?.spec?.armor ?? 0);
	let totalArmor = $derived(armor * (tokenState?.quantity ?? 0));
	let currentArmor = $derived(
		token
			? armor * (tokenState?.quantity ?? 0) -
					(tokenState?.damage ?? 0) * (tokenState?.quantityDamaged ?? 0)
			: 0
	);
</script>

<div class="w-full">
	{#if token && design && tokenState}
		<div>
			The {raceName ?? ''}
		</div>
		<div>
			Location: ({token.x}, {token.y})
		</div>
		<div class="text-primary mt-1">
			<button
				type="button"
				class="w-full h-full cursor-help"
				onpointerdown={(e) => onShipDesignTooltip(e, design)}
			>
				<div class="flex flex-col">
					<div>
						{design.name}
						{#if tokenState.quantity > 1}
							x{tokenState.quantity}
						{/if}
						<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
					</div>
					<div class="flex flex-row justify-center">
						<div
							class="border"
							style={`border-color: ${getDisplayColor(design.playerNum, $player, $universe, $settings)};`}
						>
							<TechAvatar tech={$techs.getHull(design.hull)} hullSetNumber={design.hullSetNumber} />
						</div>
					</div>
				</div></button
			>
		</div>
		<div class="flex justify-between">
			<div>
				Initiative: {token.initiative ?? 0}
			</div>
			<div>
				Movement: {token.movement ?? 0}
			</div>
		</div>
		<div class="flex justify-between">
			<div>
				Armor: {currentArmor.toFixed(0)}/{totalArmor}dp
			</div>
			{#if tokenState.destroyedPhase && phase >= tokenState.destroyedPhase}
				<div class="text-error">Destroyed</div>
			{:else if tokenState.damage && armor}
				<div class="text-error">
					Damage: {tokenState.quantityDamaged}@{((tokenState.damage / armor) * 100).toFixed(1)}%
				</div>
			{/if}
		</div>
		<div>
			Shields: {tokenState.stackShields || 'none'}
		</div>
		<div>
			Tactic: {enumToString(
				BattleTactic,
				token.tactic ? enumFromJson(BattleTacticSchema, token.tactic) : BattleTactic.UNSPECIFIED
			)}
		</div>
		<div>
			Primary Target: {enumToString(
				BattleTarget,
				token.primaryTarget
					? enumFromJson(BattleTargetSchema, token.primaryTarget)
					: BattleTarget.UNSPECIFIED
			)}
		</div>
		<div>
			Secondary Target: {enumToString(
				BattleTarget,
				token.secondaryTarget
					? enumFromJson(BattleTargetSchema, token.secondaryTarget)
					: BattleTarget.UNSPECIFIED
			)}
		</div>
	{/if}
</div>

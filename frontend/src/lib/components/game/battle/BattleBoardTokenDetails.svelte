<script lang="ts">
	import { designFinderKey, playerFinderKey } from '$lib/services/GameContext';
	import type { DesignFinder, PlayerFinder } from '$lib/services/Universe';
	import type { Battle, PhaseToken } from '$lib/types/Battle';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { startCase } from 'lodash-es';
	import { getContext } from 'svelte';
	import { onShipDesignTooltip } from '../tooltips/ShipDesignTooltip.svelte';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { techs } from '$lib/services/Stores';

	const designFinder = getContext<DesignFinder>(designFinderKey);
	const playerFinder = getContext<PlayerFinder>(playerFinderKey);

	export let battle: Battle;
	export let phase: number;
	export let token: PhaseToken | undefined;

	$: design = token && designFinder.getDesign(token.playerNum, token.designNum);
	$: raceName = token && playerFinder.getPlayerIntel(token.playerNum)?.racePluralName;
	$: tokenState = token && battle.getTokenForPhase(token.num, phase);
	$: armor = design?.spec.armor ?? 0;
	$: shields = design?.spec.shields ?? 0;
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
				on:pointerdown|preventDefault={(e) => onShipDesignTooltip(e, design)}
			>
				<div class="flex flex-col">
					<div>
						{design?.name}
						{#if (token.quantity ?? 0) > 1}
							x{token.quantity}
						{/if}
						<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
					</div>
					<div class="flex flex-row justify-center">
						<div
							class="border"
							style={`border-color: ${playerFinder.getPlayerColor(design.playerNum)};`}
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
				Armor: {armor}dp
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
			Shields: {shields ?? 'none'}
		</div>
		<div>
			Tactic: {startCase(token.tactic)}
		</div>
		<div>
			Primary Target: {startCase(token.primaryTarget)}
		</div>
		<div>
			Secondary Target: {startCase(token.secondaryTarget)}
		</div>
	{/if}
</div>

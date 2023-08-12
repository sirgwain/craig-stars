<script lang="ts">
	import { designFinderKey, playerFinderKey } from '$lib/services/Contexts';
	import type { DesignFinder, PlayerFinder } from '$lib/services/Universe';
	import type { Battle, PhaseToken } from '$lib/types/Battle';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { startCase } from 'lodash-es';
	import { getContext } from 'svelte';
	import { onShipDesignTooltip } from '../tooltips/ShipDesignTooltip.svelte';

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
		<div class="text-primary">
			<button
				type="button"
				class="w-full h-full cursor-help text-left"
				on:pointerdown|preventDefault={(e) => onShipDesignTooltip(e, design)}
			>
				{design?.name}
				{#if (token.quantity ?? 0) > 1}
					x{token.quantity}
				{/if}
				<Icon src={QuestionMarkCircle} size="16" class=" cursor-help inline-block" />
			</button>
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

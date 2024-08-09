<script lang="ts">
	import { designFinderKey, playerFinderKey } from '$lib/services/GameContext';
	import type { DesignFinder, PlayerFinder } from '$lib/services/Universe';
	import { Battle, TokenActionType, type TokenAction } from '$lib/types/Battle';
	import { getContext } from 'svelte';

	const designFinder = getContext<DesignFinder>(designFinderKey);
	const playerFinder = getContext<PlayerFinder>(playerFinderKey);

	export let battle: Battle;
	export let action: TokenAction | undefined;
	export let phase: number;

	function getTokenDescription(tokenNum?: number): string {
		if (!tokenNum) {
			return '';
		}
		const token = battle.getTokenForPhase(tokenNum, phase);
		if (token) {
			const design = designFinder.getDesign(token.playerNum, token.designNum);
			const raceName = playerFinder.getPlayerIntel(token.playerNum);
			if (design && raceName) {
				return `${raceName.racePluralName} ${design.name}`;
			}
		}
		return '';
	}
</script>

{#if action}
	{#if action.type === TokenActionType.Move}
		{`${getTokenDescription(action.tokenNum)} moved from ${action.from.x}, ${action.from.y} to ${
			action.to.x
		},${action.to.y}`}
	{:else if action.type === TokenActionType.RanAway}
		{`${getTokenDescription(action.tokenNum)} ran away`}
	{:else if action.type === TokenActionType.BeamFire}
		{`The ${getTokenDescription(action.tokenNum)} at (${action.from.x}, ${
			action.from.y
		}) attacks the ${getTokenDescription(action.targetNum)} at (${action.to.x}, ${action.to.y})`}
		{#if action.damageDoneArmor && action.damageDoneShields}
			{`doing ${action.damageDoneArmor} damage to armor and ${action.damageDoneShields} damage to shields`}
		{:else if action.damageDoneArmor}
			{`doing ${action.damageDoneArmor} damage to armor`}
		{:else}
			{`doing ${action.damageDoneShields} damage to shields`}
		{/if}
	{:else if action.type === TokenActionType.TorpedoFire}
		{`The ${getTokenDescription(action.tokenNum)} attacks the ${getTokenDescription(
			action.targetNum
		)} at (${action.to.x}, ${action.to.y})`}
		{#if action.damageDoneArmor && action.damageDoneShields}
			{`doing ${action.damageDoneArmor} damage to armor and ${action.damageDoneShields} damage to shields with ${action.torpedoHits} hits`}
		{:else if action.damageDoneArmor}
			{`doing ${action.damageDoneArmor} damage to armor`}
		{:else if action.damageDoneArmor}
			{`doing ${action.damageDoneShields} damage to shields`}
		{/if}
		{#if action.torpedoMisses}
			{`but ${action.torpedoMisses ?? 0} torpedo${
				(action.torpedoMisses ?? 0) > 1 ? 's' : ''
			} missed`}
		{/if}
	{/if}

	{#if action.tokensDestroyed && action.tokensDestroyed > 0}
		destroying {action.tokensDestroyed} ship{action.tokensDestroyed > 1 ? 's' : ''}
	{/if}
{/if}
<div />

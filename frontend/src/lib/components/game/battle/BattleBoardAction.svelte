<script lang="ts">
	import { designFinderKey, playerFinderKey } from '#lib/services/GameContext.js';
	import type { DesignFinder, PlayerFinder } from '#lib/services/Universe.js';
	import { Battle } from '#lib/types/Battle.js';
	import { type BattleRecordTokenActionJson } from '#lib/types/cs-proto.js';
	import { getContext } from 'svelte';

	const designFinder = getContext<DesignFinder>(designFinderKey);
	const playerFinder = getContext<PlayerFinder>(playerFinderKey);

	type Props = {
		battle: Battle;
		action: BattleRecordTokenActionJson | undefined;
		phase: number;
	};

	let { battle, action, phase }: Props = $props();

	function getTokenDescription(tokenNum: number | undefined = undefined): string {
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

	// the component in the shooter's firing slot, when its design is known
	function getWeaponDescription(): string {
		if (!action?.slot) {
			return '';
		}
		const token = battle.getToken(action.tokenNum ?? 0);
		const design = token && designFinder.getDesign(token.playerNum, token.designNum);
		const name = design?.slots.find((s) => s.hullSlotIndex === action.slot)?.hullComponent;
		return name ? ` with ${name}` : '';
	}
</script>

{#if action}
	{#if action.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_MOVE'}
		{`${getTokenDescription(action.tokenNum)} moved from ${action.from?.x ?? 0}, ${action.from?.y ?? 0} to ${
			action.to?.x ?? 0
		},${action.to?.y ?? 0}`}
	{:else if action.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_RAN_AWAY'}
		{`${getTokenDescription(action.tokenNum)} ran away`}
	{:else if action.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_BEAM_FIRE'}
		{`The ${getTokenDescription(action.tokenNum)} at (${action.from?.x ?? 0}, ${
			action.from?.y ?? 0
		}) attacks the ${getTokenDescription(action.targetNum)} at (${action.to?.x ?? 0}, ${action.to?.y ?? 0})${getWeaponDescription()}`}
		{#if action.damageDoneArmor && action.damageDoneShields}
			{`doing ${action.damageDoneArmor} damage to armor and ${action.damageDoneShields} damage to shields`}
		{:else if action.damageDoneArmor}
			{`doing ${action.damageDoneArmor} damage to armor`}
		{:else}
			{`doing ${action.damageDoneShields} damage to shields`}
		{/if}
	{:else if action.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_TORPEDO_FIRE'}
		{`The ${getTokenDescription(action.tokenNum)} attacks the ${getTokenDescription(
			action.targetNum
		)} at (${action.to?.x ?? 0}, ${action.to?.y ?? 0})${getWeaponDescription()}`}
		{#if action.damageDoneArmor && action.damageDoneShields}
			{`doing ${action.damageDoneArmor} damage to armor and ${action.damageDoneShields} damage to shields with ${action.torpedoHits} hits`}
		{:else if action.damageDoneArmor}
			{`doing ${action.damageDoneArmor} damage to armor`}
		{:else if action.damageDoneShields}
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
<div></div>

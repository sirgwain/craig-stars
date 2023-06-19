<script lang="ts">
	import { TokenActionType, type BattleRecord, type TokenAction, Battle } from '$lib/types/Battle';
	import type { Player } from '$lib/types/Player';

	export let battle: Battle;
	export let player: Player;
	export let action: TokenAction | undefined;

	function getTokenDescription(tokenNum?: number): string {
		if (!tokenNum) {
			return '';
		}
		const token = battle.getToken(tokenNum);
		if (token) {
			const design = player.getDesign(token.playerNum, token.designNum);
			const raceName = player.getPlayerIntel(token.playerNum);
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
		{`The ${getTokenDescription(action.tokenNum)} at (${action.from.x}, ${action.from.y}) attacks the ${getTokenDescription(
			action.targetNum
		)} at (${action.to.x}, ${action.to.y})`}
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
		{:else}
			{`but ${action.torpedoMisses ?? 0} torpedo${
				(action.torpedoMisses ?? 0) > 1 ? 's' : ''
			} missed`}
		{/if}
	{:else}
		<pre>${JSON.stringify(action, undefined, ' ')}</pre>
	{/if}
{/if}
<div />
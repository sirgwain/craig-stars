<script lang="ts">
	import type { Battle, PhaseToken } from '$lib/types/Battle';
	import type { Player } from '$lib/types/Player';
	import { startCase } from 'lodash-es';

	export let battle: Battle;
	export let player: Player;
	export let phase: number;
	export let token: PhaseToken | undefined;

	$: design = token && player.getDesign(token.playerNum, token.designNum);
	$: raceName = token && player.getPlayerIntel(token.playerNum)?.racePluralName;
	$: tokenState = token && battle.getTokenForPhase(token.num, phase);
	$: armor =
		design && ('spec' in design ? design.spec.armor : 'armor' in design ? design.armor : 0);
	$: shields =
		design && ('spec' in design ? design.spec.shield : 'armor' in design ? design.shields : 0);
</script>

<div class="w-full">
	{#if token && design && tokenState}
		<div>
			Selection: ({token.x}, {token.y})
		</div>
		<div>
			The {raceName ?? ''}
		</div>
		<div class="text-primary">
			{design?.name}
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
			{#if tokenState.damage && armor}
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

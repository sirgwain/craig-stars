<script lang="ts">
	import type { Battle, PhaseToken } from '$lib/types/Battle';
	import type { Player } from '$lib/types/Player';
	import { startCase } from 'lodash-es';

	export let battle: Battle;
	export let player: Player;
	export let token: PhaseToken | undefined;

	$: design = token && player.getDesign(token.playerNum, token.token.designNum);
	$: raceName = token && player.getPlayerIntel(token.playerNum)?.racePluralName;
</script>

<div class="w-full">
	{#if token && design}
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
				Armor: {design.spec?.armor ?? design.armor ?? 0}dp
			</div>
			<!-- <div>
				Damage: tbd - we have to figure this out per phase
			</div> -->
		</div>
		<div>
			Shields: {design.spec?.shields ?? design.shields ?? 'none'}
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

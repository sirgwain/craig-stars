<script lang="ts">
	import type { Battle, PhaseToken, TokenAction } from '$lib/types/Battle';
	import type { Player } from '$lib/types/Player';
	import type { Vector } from '$lib/types/Vector';
	import CommandTile from '../../../../routes/(user)/games/(game)/[id]/(main)/command/CommandTile.svelte';
	import BattleBoardAction from './BattleBoardAction.svelte';
	import BattleBoardAttack from './BattleBoardAttack.svelte';
	import BattleBoardPhaseControls from './BattleBoardPhaseControls.svelte';
	import BattleBoardSelectedToken from './BattleBoardSelectedToken.svelte';
	import BattleBoardSquare from './BattleBoardSquare.svelte';

	export let battle: Battle;
	export let player: Player;
	export let phase: number = 0;

	let selectedToken: PhaseToken | undefined;

	$: action = battle.getActionForPhase(phase ?? 0);
</script>

<div class="flex w-full">
	<div class="mx-auto">
		<div class="flex flex-row flex-wrap">
			<!-- the grid of the board -->
			<div class="flex flex-col">
				<div
					class="w-[690px] grid grid-cols-10 grid-cols-max grid-rows-max border-2 border-secondary rounded-md gap-0 h-[690px]"
				>
					<BattleBoardAttack {battle} {phase} />

					{#each [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] as y}
						{#each [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] as x}
							<BattleBoardSquare
								{player}
								tokens={battle.getTokensAtLocation(phase, x, y)}
								selected={selectedToken?.x === x && selectedToken?.y === y}
								on:selected={(e) => {
									selectedToken = e.detail;
								}}
							/>
						{/each}
					{/each}
				</div>
				<div class="mx-auto">
					<BattleBoardPhaseControls
						{battle}
						bind:phase
						on:phaseupdated={(e) => {
							action = battle.getActionForPhase(e.detail);
							selectedToken = action?.tokenNum
								? battle.getTokenForPhase(action.tokenNum, phase)
								: selectedToken;
						}}
					/>
				</div>
			</div>

			<!-- the right pane with descriptions -->
			<div class="pl-2 w-64">
				{#if phase}
					<div class="text-xl font-semibold text-center">
						Round {action?.round ?? 0} of {battle.totalRounds}
					</div>
					<div class="w-full card bg-base-200 shadow rounded-sm border-2 border-base-300 mb-2">
						<div class="card-body p-3 gap-0">
							<h2 class="text-lg font-semibold text-center mb-1 text-secondary">
								{`Phase ${phase} of ${battle.totalPhases}`}
							</h2>
							<BattleBoardAction {battle} {action} {player} {phase} />
						</div>
					</div>
				{/if}
				{#if selectedToken}
					<div class="w-full card bg-base-200 shadow rounded-sm border-2 border-base-300">
						<div class="card-body p-3 gap-0">
							<h2 class="text-lg font-semibold text-center mb-1 text-secondary">Selection</h2>
							<BattleBoardSelectedToken {battle} {player} token={selectedToken} {phase} />
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>

<script lang="ts">
	import { type Battle, type PhaseToken } from '$lib/types/Battle';
	import BattleBoardAction from './BattleBoardAction.svelte';
	import BattleBoardAttack from './BattleBoardAttack.svelte';
	import BattleBoardPhaseControls from './BattleBoardPhaseControls.svelte';
	import BattleBoardSquare from './BattleBoardSquare.svelte';
	import BattleBoardTokenDetails from './BattleBoardTokenDetails.svelte';

	type Props = {
		battle: Battle;
		phase?: number;
	};

	let { battle, phase = $bindable(0) }: Props = $props();

	let action = $derived(battle.getActionForPhase(phase ?? 0));
	let selectedToken: PhaseToken | undefined = $state();
	let actionToken: PhaseToken | undefined = $state();
	let target: PhaseToken | undefined = $derived(battle.getTargetForPhase(phase));
</script>

<div class="flex w-full">
	<div class="mx-auto">
		<div class="flex flex-row flex-wrap">
			<!-- the grid of the board -->
			<div class="flex flex-col md:mt-7">
				<div
					class="w-[690px] h-[690px] relative grid grid-cols-10 grid-cols-max grid-rows-max border-2 border-secondary rounded-md gap-0"
				>
					<BattleBoardAttack {battle} {phase} />

					{#each [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] as y (y)}
						{#each [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] as x (x)}
							<BattleBoardSquare
								{phase}
								{selectedToken}
								tokens={battle.getTokensAtLocation(phase, x, y)}
								selected={selectedToken?.x === x && selectedToken?.y === y}
								onSelected={(token) => {
									selectedToken = token;
								}}
							/>
						{/each}
					{/each}
				</div>
				<div class="mx-auto">
					<BattleBoardPhaseControls
						{battle}
						bind:phase
						onPhaseUpdated={(updatedPhase) => {
							phase = updatedPhase;
							const newAction = battle.getActionForPhase(phase);
							selectedToken = newAction?.tokenNum
								? battle.getTokenForPhase(newAction.tokenNum, phase)
								: selectedToken;
							actionToken = selectedToken;
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
							<BattleBoardAction {battle} {action} {phase} />
						</div>
					</div>
				{:else}
					<div class="text-xl font-semibold text-center">&nbsp</div>
				{/if}
				{#if selectedToken}
					<div class="w-full card bg-base-200 shadow rounded-sm border-2 border-base-300 mb-2">
						<div class="card-body p-3 gap-0">
							<h2 class="text-lg font-semibold text-center mb-1 text-secondary">
								{#if selectedToken.action?.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_BEAM_FIRE' || selectedToken.action?.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_TORPEDO_FIRE'}
									Attacker
								{:else}
									Selection
								{/if}
							</h2>
							<BattleBoardTokenDetails {battle} token={selectedToken} {phase} />
						</div>
					</div>
				{/if}
				{#if target && selectedToken?.num === actionToken?.num}
					<div class="w-full card bg-base-200 shadow rounded-sm border-2 border-base-300">
						<div class="card-body p-3 gap-0">
							<h2 class="text-lg font-semibold text-center mb-1 text-secondary">Target</h2>
							<BattleBoardTokenDetails {battle} token={target} {phase} />
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>

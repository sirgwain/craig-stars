<script lang="ts">
	import type { Battle, PhaseToken } from '$lib/types/Battle';
	import type { Player } from '$lib/types/Player';
	import type { Vector } from '$lib/types/Vector';
	import BattleBoardAction from './BattleBoardAction.svelte';
	import BattleBoardAttack from './BattleBoardAttack.svelte';
	import BattleBoardPhaseControls from './BattleBoardPhaseControls.svelte';
	import BattleBoardSelectedToken from './BattleBoardSelectedToken.svelte';
	import BattleBoardSquare from './BattleBoardSquare.svelte';

	export let battle: Battle;
	export let player: Player;
	export let phase: number = 0;

	let selectedSquare: Vector = { x: 0, y: 0 };
	let selectedToken: PhaseToken | undefined;

	$: action = battle.getActionForPhase(phase ?? 0);
</script>

<div class="flex justify-center">
	<div class="flex flex-row">
		<!-- the grid of the board -->
		<div class="flex flex-col">
			<div
				class="grid grid-cols-10 border-2 border-secondary rounded-md gap-0 min-w-[690px] w-[690px] overflow-auto"
			>
				<BattleBoardAttack {battle} {phase} />

				{#each [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] as y}
					{#each [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] as x}
						<BattleBoardSquare
							{player}
							{x}
							{y}
							tokens={battle.getTokensAtLocation(phase, x, y)}
							selected={selectedSquare.x === x && selectedSquare.y === y}
							on:selected={(e) => {
								selectedToken = e.detail;
								selectedSquare = { x, y };
							}}
						/>
					{/each}
				{/each}
			</div>
			<div class="mx-auto">
				<BattleBoardPhaseControls {battle} bind:phase />
			</div>
		</div>

		<!-- the right pane with descriptions -->
		<div class="pl-2 w-64">
			<div>Phase {phase} of {battle.totalPhases}</div>
			<div>Round {action?.round ?? 0} of {battle.totalRounds}</div>
			<div><BattleBoardAction {battle} {action} {player} /></div>
			<div><BattleBoardSelectedToken {battle} {player} token={selectedToken} /></div>
		</div>
	</div>
</div>

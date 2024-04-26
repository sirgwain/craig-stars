<script lang="ts">
	import TorpedoHit from '$lib/components/icons/TorpedoHit.svelte';
	import { Battle, TokenActionType } from '$lib/types/Battle';
	import { emptyVector, minus } from '$lib/types/Vector';
	import { tweened } from 'svelte/motion';

	export let battle: Battle;
	export let phase: number;

	// let tweenedX = tweened(0);
	// let tweenedY = tweened(0);

	$: actionToken = battle.getActionToken(phase ?? 0);
	$: action = battle.getActionForPhase(phase ?? 0);
	$: targetVector = action && actionToken && minus(action.to, actionToken);

	// if we are doing a same square attack
	$: targetVector && targetVector.x === 0 && targetVector.y === 0
		? (targetVector = { x: 0.5, y: 0.5 })
		: undefined;

	// $: {
	// 	if (actionToken && action) {
	// 		// $tweenedX = actionToken.x * 66 + 32;
	// 		// $tweenedY = actionToken.y * 66 + 32;
	// 		targetVector = minus(action.to, actionToken);
	// 	} else {
	// 		targetVector = emptyVector;
	// 	}
	// }
	// $: {
	// 	if (actionToken && targetVector != emptyVector) {
	// 		$tweenedX = actionToken.x * 66 + 32 + targetVector.x * 66 + 32;
	// 		$tweenedY = actionToken.y * 66 + 32 + targetVector.y * 66 + 32;
	// 	}
	// }
</script>

<div class="absolute w-full h-full z-30 pointer-events-none">
	{#if actionToken?.action?.type == TokenActionType.BeamFire && targetVector}
		<div class="relative left-0 top-0 w-full h-full">
			<svg class="w-full h-full">
				<path
					d={`M${actionToken.x * 66 + 32}, ${actionToken.y * 66 + 32} l${targetVector.x * 66},${
						targetVector.y * 66
					}`}
					class="beam-line"
				/>
			</svg>
		</div>
	{:else if action?.type == TokenActionType.TorpedoFire}
		<div class="relative left-0 top-0 w-full h-full">
			<TorpedoHit
				class="w-8 h-8 fill-transparent"
				fill={'#FF0000'}
				style={`transform: translate(${(action.to.x ?? 0) * 66 + 32 - 16}px, ${(action.to.y ?? 1) * 68 + 32 - 16}px)`}
			/>
		</div>
	{/if}
</div>

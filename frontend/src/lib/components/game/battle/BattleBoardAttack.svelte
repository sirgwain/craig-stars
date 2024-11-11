<script lang="ts">
	import TorpedoHit from '$lib/components/icons/TorpedoHit.svelte';
	import { Battle, TokenActionType } from '$lib/types/Battle';
	import { subtract } from '$lib/types/Vector';

	type Props = {
		battle: Battle;
		phase: number;
	};

	let { battle, phase }: Props = $props();

	let actionToken = $derived(battle.getActionToken(phase ?? 0));
	let action = $derived(battle.getActionForPhase(phase ?? 0));
	let targetVector = $derived.by(() => {
		if (action && actionToken) {
			const target = subtract(action.to, actionToken);

			if (target.x === 0 && target.y === 0) {
				target.x = 0.5;
				target.y = 0.5;
			}

			return target;
		}

		return undefined;
	});
</script>

<div class="absolute w-full h-full z-30 pointer-events-none">
	{#if actionToken?.action?.type === TokenActionType.BeamFire && targetVector}
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
	{:else if action?.type === TokenActionType.TorpedoFire}
		<div class="relative left-0 top-0 w-full h-full">
			<TorpedoHit
				class="w-8 h-8 fill-transparent"
				fill={'#FF0000'}
				style={`transform: translate(${action.to.x * 66 + 32 - 16}px, ${action.to.y * 68 + 32 - 16}px)`}
			/>
		</div>
	{/if}
</div>

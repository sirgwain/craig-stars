<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { type Battle } from '$lib/types/Battle';
	import { BattleRecordTokenActionType } from '$lib/types/cs-proto';
	import {
		ArrowLongLeft,
		ArrowLongRight,
		ChevronDoubleLeft,
		ChevronDoubleRight
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	type Props = {
		phase: number;
		battle: Battle;
		onPhaseUpdated?: (phase: number) => void;
	};

	let { phase = $bindable(), battle, onPhaseUpdated: onPhaseUpdated }: Props = $props();

	const previous = () => {
		phase--;
		onPhaseUpdated?.(phase);
	};
	const next = () => {
		phase++;
		onPhaseUpdated?.(phase);
	};
	const nextAttack = () => {
		const nextPhase = battle.actions.findIndex(
			(a, index) =>
				index > phase - 1 &&
				(a.type === BattleRecordTokenActionType.BEAM_FIRE ||
					a.type === BattleRecordTokenActionType.TORPEDO_FIRE)
		);
		if (nextPhase != -1) {
			phase = nextPhase + 1;
		}
		onPhaseUpdated?.(phase);
	};

	const begin = () => {
		phase = 0;
		onPhaseUpdated?.(phase);
	};
	const end = () => {
		phase = battle.totalPhases;
		onPhaseUpdated?.(phase);
	};
</script>

<div class="flex">
	<div>
		<button
			onclick={begin}
			disabled={phase === 0}
			class="btn btn-outline btn-sm normal-case btn-secondary"
			title="begin"
			><Icon src={ChevronDoubleLeft} size="16" class="hover:stroke-accent inline" /></button
		>
	</div>

	<div>
		<button
			onclick={previous}
			disabled={phase === 0}
			class="btn btn-outline btn-sm normal-case btn-secondary"
			title="previous"
			><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /></button
		>
	</div>
	<div>
		<input
			type="number"
			class="input input-sm input-bordered hide-spinner"
			onchange={(e) => (phase = clamp(parseInt(e.currentTarget.value), 0, battle.totalPhases))}
			onclick={(e) => e.currentTarget.select()}
			min={0}
			max={battle.totalPhases}
			value={phase}
		/>
	</div>
	<div>
		<button
			onclick={nextAttack}
			disabled={phase === battle.totalPhases}
			class="btn btn-outline btn-sm normal-case btn-secondary"
			title="next attack"
			>Next Attack<Icon src={ArrowLongRight} size="16" class="hover:stroke-accent inline" /></button
		>
	</div>
	<div>
		<button
			onclick={next}
			disabled={phase === battle.totalPhases}
			class="btn btn-outline btn-sm normal-case btn-secondary"
			title="next"
			><Icon src={ArrowLongRight} size="16" class="hover:stroke-accent inline" /></button
		>
	</div>
	<div>
		<button
			onclick={end}
			disabled={phase === battle.totalPhases}
			class="btn btn-outline btn-sm normal-case btn-secondary"
			title="end"
			><Icon src={ChevronDoubleRight} size="16" class="hover:stroke-accent inline" /></button
		>
	</div>
</div>

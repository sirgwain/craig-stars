<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { TokenActionType, type Battle } from '$lib/types/Battle';
	import {
		ArrowLongLeft,
		ArrowLongRight,
		ChevronDoubleLeft,
		ChevronDoubleRight
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let phase: number;
	export let battle: Battle;

	const previous = () => {
		phase--;
		dispatch('phaseupdated', phase);
	};
	const next = () => {
		phase++;
		dispatch('phaseupdated', phase);
	};
	const nextAttack = () => {
		const nextPhase = battle.actions.findIndex(
			(a, index) =>
				index > phase - 1 &&
				(a.type == TokenActionType.BeamFire || a.type == TokenActionType.TorpedoFire)
		);
		if (nextPhase != -1) {
			phase = nextPhase + 1;
		}
		dispatch('phaseupdated', phase);
	};

	const begin = () => {
		phase = 0;
		dispatch('phaseupdated', phase);
	};
	const end = () => {
		phase = battle.totalPhases;
		dispatch('phaseupdated', phase);
	};
</script>

<div class="flex">
	<div>
		<button
			on:click={begin}
			disabled={phase === 0}
			class="btn btn-outline btn-sm normal-case btn-secondary"
			title="begin"
			><Icon src={ChevronDoubleLeft} size="16" class="hover:stroke-accent inline" /></button
		>
	</div>

	<div>
		<button
			on:click={previous}
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
			on:change={(e) =>
				(phase = clamp(parseInt(e.currentTarget.value) ?? 0, 0, battle.totalPhases))}
			on:click={(e) => e.currentTarget.select()}
			min={0}
			max={battle.totalPhases}
			value={phase}
		/>
	</div>
	<div>
		<button
			on:click={nextAttack}
			disabled={phase === battle.totalPhases}
			class="btn btn-outline btn-sm normal-case btn-secondary"
			title="next attack"
			>Next Attack<Icon src={ArrowLongRight} size="16" class="hover:stroke-accent inline" /></button
		>
	</div>
	<div>
		<button
			on:click={next}
			disabled={phase === battle.totalPhases}
			class="btn btn-outline btn-sm normal-case btn-secondary"
			title="next"
			><Icon src={ArrowLongRight} size="16" class="hover:stroke-accent inline" /></button
		>
	</div>
	<div>
		<button
			on:click={end}
			disabled={phase === battle.totalPhases}
			class="btn btn-outline btn-sm normal-case btn-secondary"
			title="end"
			><Icon src={ChevronDoubleRight} size="16" class="hover:stroke-accent inline" /></button
		>
	</div>
</div>

<script lang="ts">
	import { clamp } from '#lib/services/Math.js';
	import { type Battle } from '#lib/types/Battle.js';
	import {
		ArrowLongLeft,
		ArrowLongRight,
		ChevronDoubleLeft,
		ChevronDoubleRight,
		Pause,
		Play
	} from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { SPEEDS, type BattlePlayback } from './animation/BattlePlayback.svelte';

	type Props = {
		phase: number;
		battle: Battle;
		playback: BattlePlayback;
	};

	let { phase, battle, playback }: Props = $props();

	let atEnd = $derived(phase === battle.totalPhases && !playback.animating);

	const btn = 'btn btn-sm normal-case btn-secondary';
	const iconBtn = `${btn} btn-outline btn-square`;
</script>

{#snippet phaseInput()}
	<input
		type="number"
		class="input input-sm hide-spinner w-14"
		aria-label="phase"
		onchange={(e) =>
			playback.seek(clamp(parseInt(e.currentTarget.value) || 0, 0, battle.totalPhases))}
		onclick={(e) => e.currentTarget.select()}
		min={0}
		max={battle.totalPhases}
		value={phase}
	/>
{/snippet}

<div class="battle-controls flex flex-col gap-1 w-full mt-1">
	<div class="flex flex-row flex-wrap items-center justify-center gap-x-2 gap-y-1">
		<div class="flex flex-row items-center gap-1">
			<button
				onclick={() => playback.begin()}
				disabled={phase === 0 && !playback.animating}
				class={iconBtn}
				title="begin (Home)"
				aria-label="begin"
				><Icon src={ChevronDoubleLeft} size="16" class="hover:stroke-accent inline" /></button
			>
			<button
				onclick={() => playback.stepBack()}
				disabled={phase === 0 && !playback.animating}
				class={iconBtn}
				title="previous (Left arrow)"
				aria-label="previous"
				><Icon src={ArrowLongLeft} size="16" class="hover:stroke-accent inline" /></button
			>
			<button
				onclick={() => playback.togglePlay()}
				class={`${btn} w-20 flex-nowrap`}
				title={playback.playing ? 'pause (Space)' : 'play (Space)'}
			>
				{#if playback.playing}
					<Icon src={Pause} size="16" class="inline" /> Pause
				{:else}
					<Icon src={Play} size="16" class="inline" /> {atEnd ? 'Replay' : 'Play'}
				{/if}
			</button>
			<button
				onclick={() => playback.stepForward()}
				disabled={atEnd}
				class={iconBtn}
				title="next (Right arrow)"
				aria-label="next"
				><Icon src={ArrowLongRight} size="16" class="hover:stroke-accent inline" /></button
			>
			<button
				onclick={() => playback.end()}
				disabled={atEnd}
				class={iconBtn}
				title="end (End)"
				aria-label="end"
				><Icon src={ChevronDoubleRight} size="16" class="hover:stroke-accent inline" /></button
			>
		</div>

		<div class="flex flex-row items-center gap-1">
			{@render phaseInput()}
			<button
				onclick={() => playback.nextAttack()}
				disabled={atEnd}
				class={`${btn} btn-outline`}
				title="next attack"
				>Next Attack<Icon
					src={ArrowLongRight}
					size="16"
					class="hover:stroke-accent inline"
				/></button
			>
		</div>

		<div class="flex flex-row items-center gap-1">
			<div class="join" role="group" aria-label="playback speed">
				{#each SPEEDS as speed (speed)}
					<button
						class="join-item btn btn-sm px-2 normal-case"
						class:btn-secondary={playback.speed === speed}
						class:btn-outline={playback.speed !== speed}
						onclick={() => playback.setSpeed(speed)}
						title={`playback speed ${speed}x`}
						aria-pressed={playback.speed === speed}>{speed}x</button
					>
				{/each}
			</div>
		</div>
		<div class="flex flex-col gap-1">
			<label
				class="label cursor-pointer gap-1 text-sm"
				title="fly torpedoes to their targets; when off, attacks are drawn right away"
			>
				<input
					type="checkbox"
					class="toggle toggle-sm toggle-secondary"
					checked={playback.animations}
					onchange={(e) => playback.setAnimations(e.currentTarget.checked)}
				/>
				Animations
			</label>
			<label
				class="label cursor-pointer gap-1 text-sm"
				title="mark the damage an attack does on its targets"
			>
				<input
					type="checkbox"
					class="toggle toggle-sm toggle-secondary"
					checked={playback.damage}
					onchange={(e) => playback.setDamage(e.currentTarget.checked)}
				/>
				Damage
			</label>
		</div>
	</div>
	<input
		type="range"
		class="range range-xs w-full"
		min="0"
		max={battle.totalPhases}
		step="1"
		aria-label="battle timeline"
		value={phase}
		oninput={(e) => playback.seek(parseInt(e.currentTarget.value) || 0)}
	/>
</div>

<style>
	/* bigger controls on phones, for fingers (everywhere the desk layout doesn't apply) */
	@media (max-width: 47.99rem), (max-height: 29.99rem) {
		.battle-controls .btn {
			--size: calc(var(--size-field, 0.25rem) * 10);
		}
		.battle-controls .input,
		.battle-controls .select {
			--size: calc(var(--size-field, 0.25rem) * 10);
		}
		.battle-controls .range {
			--range-thumb-size: calc(var(--size-selector, 0.25rem) * 6);
		}
	}
</style>

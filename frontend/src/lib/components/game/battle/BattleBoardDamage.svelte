<script lang="ts">
	import { Battle } from '#lib/types/Battle.js';
	import { attackAtPhase } from './animation/original';
	import { ARMOR_ICON_PATH, damageByTarget, SHIELD_ICON_PATH } from './damage';

	type Props = {
		battle: Battle;
		phase: number;
		// size of a board square in px
		cell: number;
		// draw the damage marks this much bigger, to stay readable on a scaled down board
		markScale?: number;
	};

	let { battle, phase, cell, markScale = 1 }: Props = $props();

	// what the attack on the board did to its targets, shown on their squares
	let targets = $derived(damageByTarget(attackAtPhase(battle.actions, phase)));
</script>

<!-- over the canvas, so a bang on the target does not cover its damage -->
<div class="absolute size-full z-40 pointer-events-none">
	{#each targets as target (`${target.x},${target.y}`)}
		<!-- damage on the target: shields in blue with a bubble, armor in orange with a plate -->
		<div
			class="battle-damage absolute flex flex-col items-center"
			class:battle-damage-armor={target.marks.some((m) => m.kind === 'armor')}
			data-testid="battle-damage"
			style={`left: ${target.x * cell}px; top: ${target.y * cell}px; width: ${cell}px; height: ${cell}px;`}
		>
			<div
				class="flex flex-col items-center gap-px origin-top"
				style={markScale !== 1 ? `transform: scale(${markScale});` : ''}
			>
				{#each target.marks as mark (mark.kind)}
					<div
						class="battle-damage-mark flex flex-row items-center gap-0.5 rounded-sm px-1 text-xs leading-4 font-bold"
						style={`color: ${mark.color};`}
						data-kind={mark.kind}
						title={mark.label}
						aria-label={mark.label}
					>
						<svg viewBox="0 0 24 24" class="size-3 shrink-0" aria-hidden="true">
							{#if mark.kind === 'shield'}
								<path d={SHIELD_ICON_PATH} fill="none" stroke="currentColor" stroke-width="3" />
							{:else}
								<path d={ARMOR_ICON_PATH} fill="currentColor" />
							{/if}
						</svg>
						{mark.text}
					</div>
				{/each}
			</div>
		</div>
	{/each}
</div>

<style>
	.battle-damage {
		justify-content: flex-start;
		padding-top: 2px;
		border-radius: 0.25rem;
		box-shadow: inset 0 0 0 2px rgb(143 227 255 / 0.8);
	}

	.battle-damage-armor {
		box-shadow: inset 0 0 0 2px rgb(255 176 80 / 0.9);
	}

	.battle-damage-mark {
		background: rgb(0 0 0 / 0.75);
		white-space: nowrap;
	}
</style>

<script lang="ts">
	import TorpedoHit from '#lib/components/icons/TorpedoHit.svelte';
	import { Battle } from '#lib/types/Battle.js';
	import { ARMOR_ICON_PATH, damageMarks, SHIELD_ICON_PATH } from './damage';

	type Props = {
		battle: Battle;
		phase: number;
		// size of a board square in px
		cell?: number;
		// draw the damage marks this much bigger, to stay readable on a scaled down board
		markScale?: number;
	};

	let { battle, phase, cell = 66, markScale = 1 }: Props = $props();
	let half = $derived(cell / 2);

	let actionToken = $derived(battle.getActionToken(phase));
	let action = $derived(battle.getActionForPhase(phase));
	let targetVector = $derived.by(() => {
		if (action?.to && actionToken) {
			const target = {
				x: (action.to.x ?? 0) - actionToken.x,
				y: (action.to.y ?? 0) - actionToken.y
			};

			if (target.x === 0 && target.y === 0) {
				target.x = 0.5;
				target.y = 0.5;
			}

			return target;
		}

		return undefined;
	});

	// what the shot did to its target, shown on the target's square
	let marks = $derived(
		action?.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_BEAM_FIRE' ||
			action?.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_TORPEDO_FIRE'
			? damageMarks(action)
			: []
	);
	let armorHit = $derived(marks.some((m) => m.kind === 'armor'));
</script>

<div class="absolute size-full z-30 pointer-events-none">
	{#if actionToken?.action?.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_BEAM_FIRE' && targetVector}
		<div class="relative left-0 top-0 size-full">
			<svg class="size-full">
				<path
					d={`M${actionToken.x * cell + half}, ${actionToken.y * cell + half} l${targetVector.x * cell},${
						targetVector.y * cell
					}`}
					class="beam-line"
				/>
			</svg>
		</div>
	{:else if action?.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_TORPEDO_FIRE' && action.to}
		<div class="relative left-0 top-0 size-full">
			<TorpedoHit
				class="size-8  fill-transparent"
				fill="#FF0000"
				style={`transform: translate(${(action.to.x ?? 0) * cell + half - 16}px, ${(action.to.y ?? 0) * cell + half - 16}px)`}
			/>
		</div>
	{/if}
	{#if marks.length && action?.to}
		<!-- damage on the target: shields in blue, armor in orange -->
		<div
			class="battle-damage absolute flex flex-col items-center"
			class:battle-damage-armor={armorHit}
			data-testid="battle-damage"
			style={`left: ${(action.to.x ?? 0) * cell}px; top: ${(action.to.y ?? 0) * cell}px; width: ${cell}px; height: ${cell}px;`}
		>
			<div
				class="flex flex-col items-center gap-px origin-top"
				style={markScale !== 1 ? `transform: scale(${markScale});` : ''}
			>
				{#each marks as mark (mark.kind)}
					<div
						class="battle-damage-mark flex flex-row items-center gap-0.5 rounded-sm px-1 text-xs leading-4 font-bold"
						style={`color: ${mark.color};`}
						data-kind={mark.kind}
						title={mark.label}
						aria-label={mark.label}
					>
						{#if mark.kind !== 'destroyed'}
							<svg viewBox="0 0 24 24" class="size-3 shrink-0" aria-hidden="true">
								<path
									d={mark.kind === 'shield' ? SHIELD_ICON_PATH : ARMOR_ICON_PATH}
									fill={mark.kind === 'shield' ? 'none' : 'currentColor'}
									stroke="currentColor"
									stroke-width={mark.kind === 'shield' ? 3 : 1}
									stroke-linejoin="round"
								/>
							</svg>
						{/if}
						{mark.text}
					</div>
				{/each}
			</div>
		</div>
	{/if}
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

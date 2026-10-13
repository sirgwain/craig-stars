<script lang="ts">
	import { designFinderKey, gameKey, type GameContext } from '#lib/services/GameContext.js';
	import { techs } from '#lib/services/Stores.js';
	import type { DesignFinder } from '#lib/services/Universe.js';
	import { type Battle, type PhaseToken } from '#lib/types/Battle.js';
	import { getDisplayColor } from '#lib/utils/colorUtils.js';
	import { Play } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { getContext, onDestroy, onMount, untrack } from 'svelte';
	import { MediaQuery } from 'svelte/reactivity';
	import { readable } from 'svelte/store';
	import { BattlePlayback } from './animation/BattlePlayback.svelte';
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

	// the board can be shown without a game (the battle test page), so fall back to fixed colors
	const gameContext = getContext<GameContext | undefined>(gameKey);
	const player = gameContext?.player ?? readable(undefined);
	const universe = gameContext?.universe ?? readable(undefined);
	const settings = gameContext?.settings ?? readable(undefined);
	const designFinder = getContext<DesignFinder | undefined>(designFinderKey);
	const fallbackColors = ['#0000FF', '#FF0000', '#FFFF00', '#00C000', '#FF00FF', '#00FFFF'];

	function colorOf(playerNum: number | undefined): string {
		if ($player && $universe && $settings) {
			return getDisplayColor(playerNum, $player, $universe, $settings);
		}
		return fallbackColors[(Math.max(1, playerNum ?? 1) - 1) % fallbackColors.length];
	}

	let action = $derived(battle.getActionForPhase(phase));
	let selectedToken: PhaseToken | undefined = $state();
	let actionToken: PhaseToken | undefined = $state();
	let target: PhaseToken | undefined = $derived(battle.getTargetForPhase(phase));

	// a new phase was committed (by playback or the controls): follow the acting token
	function setPhase(updatedPhase: number) {
		phase = updatedPhase;
		const newAction = battle.getActionForPhase(phase);
		selectedToken = newAction?.tokenNum
			? battle.getTokenForPhase(newAction.tokenNum, phase)
			: selectedToken;
		actionToken = selectedToken;
	}

	const playback = new BattlePlayback(
		untrack(() => battle),
		{
			getPhase: () => phase,
			setPhase,
			weaponPower
		}
	);
	onMount(() => playback.autoplay());
	onDestroy(() => playback.destroy());

	// damage per part of the weapon in a token's slot (the original draws strong beams blue)
	function weaponPower(tokenNum: number, slot: number): number | undefined {
		const token = battle.getToken(tokenNum);
		const design = token && designFinder?.getDesign(token.playerNum, token.designNum);
		const name = design?.slots.find((s) => s.hullSlotIndex === slot)?.hullComponent;
		return name ? $techs.getHullComponent(name)?.power : undefined;
	}

	$effect(() => {
		const b = battle;
		untrack(() => playback.setBattle(b));
	});

	// The board is laid out at 690px and scaled down to fit a phone: to the width when upright,
	// to the height when on its side (with the controls beside it).
	const boardSize = 690;
	const phoneLandscape = new MediaQuery('(orientation: landscape) and (max-height: 29.99rem)');
	let availableWidth = $state(boardSize);
	let innerHeight = $state(0);
	let boardTop = $state(0);
	let boardArea: HTMLDivElement | undefined = $state();
	let scale = $derived.by(() => {
		let size = Math.min(boardSize, availableWidth);
		if (phoneLandscape.current && innerHeight) {
			size = Math.min(size, Math.max(200, innerHeight - boardTop - 8));
		}
		return size / boardSize;
	});
	let scaledSize = $derived(Math.floor(boardSize * scale));

	$effect(() => {
		// where the board starts on the page, to fit it under the header on a sideways phone
		void innerHeight;
		void phoneLandscape.current;
		if (boardArea) {
			boardTop = boardArea.getBoundingClientRect().top + window.scrollY;
		}
	});

	// board geometry, measured so the attacks line up with the squares
	let boardWidth = $state(686);
	let boardHeight = $state(686);
	let cell = $derived(boardWidth / 10);
	let canvas: HTMLCanvasElement | undefined = $state();

	$effect(() => {
		const [c, w, h] = [canvas, boardWidth, boardHeight];
		untrack(() => playback.setCanvas(c, w, h));
	});

	// keyboard control of the original playback (the classic viewer has none)
	function onKeyDown(e: KeyboardEvent) {
		const el = e.target as HTMLElement | null;
		if (
			!playback.animate ||
			e.ctrlKey ||
			e.metaKey ||
			e.altKey ||
			(el &&
				(el.tagName === 'INPUT' ||
					el.tagName === 'TEXTAREA' ||
					el.tagName === 'SELECT' ||
					el.isContentEditable))
		) {
			return;
		}
		switch (e.key) {
			case ' ':
				playback.togglePlay();
				break;
			case 'ArrowRight':
				playback.stepForward();
				break;
			case 'ArrowLeft':
				playback.stepBack();
				break;
			case 'Home':
				playback.begin();
				break;
			case 'End':
				playback.end();
				break;
			default:
				return;
		}
		e.preventDefault();
	}

	// a focused button would otherwise also "click" on space
	function onKeyUp(e: KeyboardEvent) {
		if (
			playback.animate &&
			e.key === ' ' &&
			(e.target as HTMLElement | null)?.tagName === 'BUTTON'
		) {
			e.preventDefault();
		}
	}
</script>

<svelte:window onkeydown={onKeyDown} onkeyup={onKeyUp} bind:innerHeight />

<div class="flex w-full" bind:clientWidth={availableWidth}>
	<div class="mx-auto">
		<div class="flex flex-row flex-wrap">
			<!-- the grid of the board -->
			<div
				class="flex flex-col md:mt-7 phone-landscape:mt-0 phone-landscape:flex-row"
				style={`--board-size: ${scaledSize}px;`}
			>
				<div class="relative size-(--board-size)" bind:this={boardArea}>
					<div
						class="size-[690px] relative grid grid-cols-10 grid-cols-max grid-rows-max border-2 border-secondary rounded-md gap-0 overflow-hidden origin-top-left"
						style={scale !== 1 ? `transform: scale(${scale});` : ''}
						bind:clientWidth={boardWidth}
						bind:clientHeight={boardHeight}
						data-testid="battle-board"
					>
						<!-- the original style draws its own attacks -->
						{#if playback.style === 'classic'}
							<BattleBoardAttack {battle} {phase} {cell} markScale={Math.min(2, 1 / scale)} />
						{/if}

						{#each [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] as y (y)}
							{#each [0, 1, 2, 3, 4, 5, 6, 7, 8, 9] as x (x)}
								<BattleBoardSquare
									{phase}
									{selectedToken}
									{colorOf}
									tokens={battle.getTokensAtLocation(phase, x, y)}
									selected={selectedToken?.x === x && selectedToken.y === y}
									onSelected={(token) => {
										selectedToken = token;
									}}
								/>
							{/each}
						{/each}

						<!-- the original style's beams, torpedoes and bangs -->
						<canvas
							bind:this={canvas}
							class="absolute left-0 top-0 size-full z-30 pointer-events-none"
							aria-hidden="true"
						></canvas>
					</div>

					{#if playback.animate && phase === 0 && !playback.playing && !playback.animating && battle.totalPhases > 0}
						<div class="absolute inset-0 z-40 flex items-center justify-center pointer-events-none">
							<button
								class="btn btn-secondary btn-lg shadow-lg pointer-events-auto opacity-90"
								onclick={() => playback.play()}
								title="play the battle (Space)"
							>
								<Icon src={Play} size="24" class="inline" /> Play battle
							</button>
						</div>
					{/if}
				</div>
				<div
					class="mx-auto w-(--board-size) phone-landscape:mx-0 phone-landscape:w-72 phone-landscape:pl-2"
				>
					<BattleBoardPhaseControls {battle} {phase} {playback} />
				</div>
			</div>

			<!-- the right pane with descriptions -->
			<div class="pl-2 w-64">
				{#if phase}
					<div class="text-xl font-semibold text-center">
						Round {action?.round ?? 0} of {battle.totalRounds}
					</div>
					<div class="w-full card bg-base-200 shadow-sm rounded-xs border-2 border-base-300 mb-2">
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
					<div class="w-full card bg-base-200 shadow-sm rounded-xs border-2 border-base-300 mb-2">
						<div class="card-body p-3 gap-0">
							<h2 class="text-lg font-semibold text-center mb-1 text-secondary">
								{#if selectedToken.action?.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_BEAM_FIRE' || selectedToken.action?.type === 'BATTLE_RECORD_TOKEN_ACTION_TYPE_TORPEDO_FIRE'}
									Attacker
								{:else}
									Selection
								{/if}
							</h2>
							<BattleBoardTokenDetails {battle} token={selectedToken} {phase} {colorOf} />
						</div>
					</div>
				{/if}
				{#if target && selectedToken?.num === actionToken?.num}
					<div class="w-full card bg-base-200 shadow-sm rounded-xs border-2 border-base-300">
						<div class="card-body p-3 gap-0">
							<h2 class="text-lg font-semibold text-center mb-1 text-secondary">Target</h2>
							<BattleBoardTokenDetails {battle} token={target} {phase} {colorOf} />
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>

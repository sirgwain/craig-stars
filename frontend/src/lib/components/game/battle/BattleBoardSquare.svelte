<script lang="ts">
	import { designFinderKey, getGameContext } from '$lib/services/GameContext';
	import type { DesignFinder } from '$lib/services/Universe';
	import { getHullIcon } from '$lib/techicon';
	import type { PhaseToken } from '$lib/types/Battle';
	import { getDisplayColor } from '$lib/utils/colorUtils';
	import { getContext } from 'svelte';

	const { player, universe, settings } = getGameContext();

	const designFinder = getContext<DesignFinder>(designFinderKey);

	type Props = {
		tokens?: PhaseToken[] | undefined;
		phase: number;
		selectedToken: PhaseToken | undefined;
		selected?: boolean;
		onSelected?: (token: PhaseToken) => void;
	};

	let {
		tokens = undefined,
		phase,
		selectedToken,
		selected = false,
		onSelected: onSelected
	}: Props = $props();

	let targetTokenIndex = $derived(tokens?.findIndex((t) => t.target));
	let selectedTokenIndex = $derived(
		selectedToken && tokens?.findIndex((t) => t.num === selectedToken.num)
	);
	let tokenIndex = $derived(
		targetTokenIndex && targetTokenIndex != -1
			? targetTokenIndex
			: selectedTokenIndex && selectedTokenIndex != -1
				? selectedTokenIndex
				: 0
	);

	let topToken = $derived(tokens && tokens[tokenIndex]);

	const icon = (tokens: PhaseToken[] | undefined, tokenIndex: number) => {
		if (
			tokens &&
			tokens.filter((t) => !(t.ranAway || (t.destroyedPhase && phase > t.destroyedPhase))).length >
				0
		) {
			tokenIndex = tokenIndex % tokens.length;
			const token = tokens[tokenIndex];
			const design = designFinder.getDesign(token.playerNum, token.designNum);
			if (design) {
				return getHullIcon(design);
			}
		}
		return '';
	};
</script>

<div
	class={`bg-black tech-avatar border-2 box-content relative ${icon(tokens, tokenIndex)}`}
	class:border-neutral={!selected && !tokens}
	class:border-accent={selected}
	style={!selected && tokens && tokens[tokenIndex]
		? `border-color: ${getDisplayColor(tokens[tokenIndex].playerNum, $player, $universe, $settings)};`
		: ''}
>
	{#if tokens}
		<div class="absolute flex flex-row justify-between w-full text-[11px] -top-0.5">
			<div>
				{#if tokens.length > 1}
					{tokens.length}
				{/if}
			</div>
			<div>
				{#if (topToken?.quantity ?? 0) > 1}
					{topToken?.quantity}
				{/if}
			</div>
		</div>
		<button
			type="button"
			class="w-full h-full cursor-pointer"
			aria-label="Selects the token on the board"
			onclick={() => {
				const newTokenIndex = selected ? (tokenIndex + 1) % tokens.length : tokenIndex;
				onSelected?.(tokens[newTokenIndex]);
			}}
		></button>
	{/if}
</div>

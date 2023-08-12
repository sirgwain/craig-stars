<script lang="ts">
	import { designFinderKey, playerFinderKey } from '$lib/services/Contexts';
	import type { DesignFinder, PlayerFinder } from '$lib/services/Universe';
	import type { PhaseToken } from '$lib/types/Battle';
	import { kebabCase } from 'lodash-es';
	import { createEventDispatcher, getContext } from 'svelte';

	const designFinder = getContext<DesignFinder>(designFinderKey);
	const playerFinder = getContext<PlayerFinder>(playerFinderKey);

	const dispatch = createEventDispatcher();

	export let tokens: PhaseToken[] | undefined = undefined;
	export let phase: number;
	export let selectedToken: PhaseToken | undefined;
	export let selected = false;

	$: targetTokenIndex = tokens?.findIndex((t) => t.target);
	$: selectedTokenIndex = selectedToken && tokens?.indexOf(selectedToken);
	$: tokenIndex =
		targetTokenIndex && targetTokenIndex != -1
			? targetTokenIndex
			: selectedTokenIndex && selectedTokenIndex != -1
			? selectedTokenIndex
			: 0;

	$: topToken = tokens && tokens[tokenIndex];

	const icon = (tokens: PhaseToken[] | undefined, tokenIndex: number) => {
		if (
			tokens &&
			tokens.filter((t) => !(t.ranAway || (t.destroyedPhase && phase > t.destroyedPhase))).length >
				0
		) {
			tokenIndex = tokenIndex % (tokens?.length ?? 0);
			const token = tokens[tokenIndex];
			if (token) {
				const design = designFinder.getDesign(token.playerNum, token.designNum);
				if (design) {
					const name = kebabCase(design.hull.replace("'", '').replace(' ', '').replace('±', ''));
					return `hull-${name}-${design.hullSetNumber ?? 0}`;
				}
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
		? `border-color: ${playerFinder.getPlayerColor(tokens[tokenIndex].playerNum)};`
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
			on:click={() => {
				if (tokens) {
					if (selected) {
						tokenIndex = (tokenIndex + 1) % (tokens?.length ?? 0);
					}
					dispatch('selected', tokens[tokenIndex]);
				}
			}}
		/>
	{/if}
</div>

<script lang="ts">
	import type { PhaseToken } from '$lib/types/Battle';
	import type { Player } from '$lib/types/Player';
	import { kebabCase } from 'lodash-es';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let player: Player;
	export let x: number;
	export let y: number;
	export let tokens: PhaseToken[] | undefined = undefined;
	export let selected = false;

	$: actionToken = tokens?.find((t) => t.action);
	$: actionTokenIndex = tokens?.findIndex((t) => t.action) ?? -1;
	$: activeTokenIndex = actionTokenIndex === -1 ? 0 : actionTokenIndex;

	const icon = (tokens: PhaseToken[] | undefined, tokenIndex: number) => {
		if (tokens && tokens.filter((t) => !(t.ranAway || t.destroyed)).length > 0) {
			const token = tokens[tokenIndex];
			if (token) {
				const design = player.getDesign(token.playerNum, token.designNum);
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
	class={`bg-black tech-avatar border-2 box-content ${icon(tokens, activeTokenIndex)}`}
	class:border-secondary={!selected && !actionToken}
	class:border-accent={!selected && !!actionToken}
	class:border-primary={selected}
	class:z-20={!!actionToken}
>
	{#if tokens}
		<button
			type="button"
			class="w-full h-full cursor-pointer"
			on:click={() => {
				if (tokens) {
					dispatch('selected', tokens[activeTokenIndex]);
					if (selected) {
						activeTokenIndex = (activeTokenIndex + 1) % (tokens?.length ?? 0);
					}
				}
			}}
		/>
	{/if}
</div>

<script lang="ts">
	import type { PhaseToken } from '$lib/types/Battle';
	import type { Player } from '$lib/types/Player';
	import { kebabCase } from 'lodash-es';
	import { afterUpdate, createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let player: Player;
	export let tokens: PhaseToken[] | undefined = undefined;
	export let selected = false;

	let tokenIndex = 0;

	const icon = (tokens: PhaseToken[] | undefined, tokenIndex: number) => {
		if (tokens && tokens.filter((t) => !(t.ranAway || t.destroyed)).length > 0) {
			tokenIndex = tokenIndex % (tokens?.length ?? 0);
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
	class={`bg-black tech-avatar border-2 box-content ${icon(tokens, tokenIndex)}`}
	class:border-neutral={!selected}
	class:border-accent={selected}
>
	{#if tokens}
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

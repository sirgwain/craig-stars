<script lang="ts">
	import { getTechIcon } from '$lib/techicon';
	import type { TechHull } from '$lib/types/cs-proto';
	import { isHull, type TechLike } from '$lib/types/Tech';
	import { onTechHullTooltip } from '../game/tooltips/TechHullTooltip.svelte';
	import { onTechTooltip } from '../game/tooltips/TechTooltip.svelte';

	type Props = {
		tech: TechLike | undefined;
		hullSetNumber?: number;
		hullTooltip?: boolean;
	};

	let { tech, hullSetNumber = 0, hullTooltip = false }: Props = $props();
	let hull = $derived(isHull(tech) && (tech as TechHull));
</script>

<div
	class="tech-avatar {getTechIcon(tech, hullSetNumber)}"
	role="link"
	tabindex="-1"
	oncontextmenu={(e) => {
		e.preventDefault();
		onTechTooltip(e, tech);
	}}
>
	{#if hullTooltip && hull}
		<button
			type="button"
			aria-label="Brings up information on technology"
			class="w-full h-full"
			onpointerdown={(e) => {
				onTechHullTooltip(e, hull);
			}}
		></button>
	{/if}
</div>

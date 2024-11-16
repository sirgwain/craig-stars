<script lang="ts">
	import { isHull, TechCategory, type Tech, type TechHull } from '$lib/types/Tech';
	import { kebabCase } from 'lodash-es';
	import { onTechHullTooltip } from '../game/tooltips/TechHullTooltip.svelte';
	import { onTechTooltip } from '../game/tooltips/TechTooltip.svelte';
	import { getTechIcon } from '$lib/techicon';

	type Props = {
		tech: Tech | undefined;
		hullSetNumber?: number;
		hullTooltip: boolean;
	};

	let { tech = undefined, hullSetNumber = 0, hullTooltip = false }: Props = $props();
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
				e.preventDefault();
				onTechHullTooltip(e, hull);
			}}
		></button>
	{/if}
</div>

<script lang="ts">
	import { TechCategory, type Tech, type TechHull } from '$lib/types/Tech';
	import { kebabCase } from 'lodash-es';
	import { onTechHullTooltip } from '../game/tooltips/TechHullTooltip.svelte';
	import { onTechTooltip } from '../game/tooltips/TechTooltip.svelte';

	interface Props
	{
		tech: Tech | undefined,
		hullSetNumber?: number,
		hullTooltip: boolean
	}

	let {
		tech = undefined,
		hullSetNumber = 0,
		hullTooltip = false
	}: Props = $props();

	let hull = $derived(tech &&
		[TechCategory.ShipHull, TechCategory.StarbaseHull].includes(tech.category)
		? tech as TechHull : undefined);

	const icon = (tech: Tech | undefined, hullSetNumber: number) => {
		const name = kebabCase(tech?.name.replace("'", '').replace(' ', '').replace('±', ''));
		if (hull) {
			return `hull-${name}-${hullSetNumber ?? 0}`;
		} else {
			return name;
		}
	};
</script>

<div
	class="tech-avatar {icon(tech, hullSetNumber)}"
	role="link"
	tabindex="-1"
	oncontextmenu={e => { e.preventDefault(); onTechTooltip(e, tech); } }
>
	{#if hullTooltip && hull}
		<button
			type="button"
			aria-label="Brings up information on technology"
			class="w-full h-full"
			onpointerdown={e => {e.preventDefault(); onTechHullTooltip(e, hull); } }
		></button>
	{/if}
</div>

<script lang="ts">
	import { TechCategory, type Tech, type TechHull } from '$lib/types/Tech';
	import { kebabCase } from 'lodash-es';
	import { onTechHullTooltip } from '../game/tooltips/TechHullTooltip.svelte';
	import { onTechTooltip } from '../game/tooltips/TechTooltip.svelte';

	export let tech: Tech | undefined = undefined;
	export let hullSetNumber = 0;
	export let hullTooltip = false;

	let hull: TechHull;

	const icon = (hullSetNumber: number) => {
		const name = kebabCase(tech?.name.replace("'", '').replace(' ', '').replace('±', ''));
		if (hull) {
			return `hull-${name}-${hullSetNumber ?? 0}`;
		} else {
			return name;
		}
	};

	$: tech &&
		(tech.category == TechCategory.ShipHull || tech.category == TechCategory.StarbaseHull) &&
		(hull = tech as TechHull);
</script>

<div
	class="avatar tech-avatar {icon(hullSetNumber)}"
	on:contextmenu|preventDefault={(e) => onTechTooltip(e, tech)}
>
	{#if hullTooltip && hull}
		<button
			type="button"
			class="w-full h-full"
			on:pointerdown|preventDefault={(e) => onTechHullTooltip(e, hull)}
		/>
	{/if}
</div>

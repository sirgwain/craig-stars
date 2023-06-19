<script lang="ts">
	import { TechCategory, type Tech, type TechHull } from '$lib/types/Tech';
	import { kebabCase } from 'lodash-es';
	import { onTechTooltip } from '../game/tooltips/TechTooltip.svelte';

	export let tech: Tech | undefined = undefined;
	export let hullSetNumber = 0;
	export let numHullSets = 4;
	export let hullSetChangeable = false;

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

<div class="avatar tech-avatar {icon(hullSetNumber)}">
	<button type="button" />
</div>

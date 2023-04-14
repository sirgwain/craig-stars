<script lang="ts">
	import { TechCategory, type Tech, type TechHull } from '$lib/types/Tech';
	import { kebabCase } from 'lodash-es';

	export let tech: Tech | undefined = undefined;
	export let hullSetNumber = 1;
	export let numHullSets = 4;
	export let hullSetChangeable = false;

	let hull: TechHull;

	const icon = (hullSetNumber: number) => {
		const name = kebabCase(tech?.name.replace("'", '').replace(' ', '').replace('±', ''));
		if (hull) {
			return `hull-${name}-${hullSetNumber}`;
		} else {
			return name;
		}
	};

	$: tech &&
		(tech.category == TechCategory.ShipHull || tech.category == TechCategory.StarbaseHull) &&
		(hull = tech as TechHull);
</script>

<div
	on:click={() => {
		if (hullSetChangeable) {
			hullSetNumber += 1;
			hullSetNumber %= numHullSets;
		}
	}}
	class="avatar border border-secondary tech-avatar {icon(hullSetNumber)}"
/>

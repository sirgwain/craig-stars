<script lang="ts">
	
	import type { Tech, TechHullComponent } from '$lib/types/Tech';
	import { startCase } from 'lodash-es';
	import { onMount } from 'svelte';

	export let tech: Tech;

	let warnings: string[] = [];

	onMount(() => {
		if ('hullSlotType' in tech) {
			const hullComponent = tech as TechHullComponent;
			if (hullComponent) {
				if (hullComponent.radiating) {
					warnings.push(
						`This ${startCase(
							hullComponent.category
						).toLowerCase()} creates powerful waves of radiation and will kill some of your colonists if the midpoint of your race's Radiation band isn't at least 85mR.`
					);
				}
			}
		}
		warnings = warnings;
	});
</script>

<div class="flex flex-col p-1">
	{#each warnings as warning}
		<div class="text-warning text-base">{warning}</div>
	{/each}
</div>

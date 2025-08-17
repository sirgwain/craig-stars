<script lang="ts">
	import { TechCategory, type TechHullComponent } from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import type { TechLike } from '$lib/types/Tech';
	import { onMount } from 'svelte';

	type Props = {
		tech: TechLike;
	};

	let { tech }: Props = $props();

	let warnings: string[] = $state([]);

	onMount(() => {
		if ('hullSlotType' in tech) {
			const hullComponent = tech as TechHullComponent;
			if (hullComponent) {
				if (hullComponent.radiating) {
					warnings.push(
						`This ${enumToString(
							TechCategory,
							hullComponent.tech?.category ?? TechCategory.UNSPECIFIED
						).toLowerCase()} creates powerful waves of radiation and will kill some of your colonists if the midpoint of your race's Radiation band isn't at least 85mR.`
					);
				}
			}
		}
		warnings = warnings;
	});
</script>

<div class="flex flex-col p-1">
	{#each warnings as warning (warning)}
		<div class="text-warning text-base">{warning}</div>
	{/each}
</div>

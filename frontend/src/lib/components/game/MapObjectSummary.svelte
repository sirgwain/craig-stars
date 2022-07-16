<script lang="ts">
	import { selectedMapObject, commandedPlanet, game, player } from '$lib/services/Context';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { PaperAirplane } from '@steeze-ui/heroicons';
	import PlanetSummary from './PlanetSummary.svelte';
	import type { Planet } from '$lib/types/Planet';

	let title = '';

	$: {
		if ($selectedMapObject) {
			title = $selectedMapObject.name;
		}
	}

	let selectedPlanet: Planet | undefined;
	$: selectedPlanet = $selectedMapObject as Planet;
</script>

<div class="card bg-base-200 shadow-xl rounded-sm border-2 border-base-300">
	<div class="card-body p-2 gap-0">
		<div class="flex flex-row items-center">
			<div class="flex-1 text-center text-lg font-semibold text-secondary">
				{title}
			</div>
			<Icon src={PaperAirplane} size="16" class="hover:stroke-accent" />
		</div>
		{#if selectedPlanet}
			<PlanetSummary />
		{/if}
	</div>
</div>

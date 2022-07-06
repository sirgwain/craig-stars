<script lang="ts">
	import { commandedPlanet, game, player } from '$lib/services/Context';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { PaperAirplane } from '@steeze-ui/heroicons';
	import PlanetSummary from './PlanetSummary.svelte';

	let title = '';

	$: {
		if ($commandedPlanet) {
			title = $commandedPlanet.name;
		}
	}
</script>

<div class="card bg-base-200 shadow-xl w-[31rem] max-h-fit min-h-fit rounded-sm border-2 border-base-300">
	<div class="card-body p-2 text-base gap-0">
		<div class="flex flex-row items-center">
			<div class="flex-1 text-center text-xl font-semibold text-secondary">
				{title}
			</div>
			<Icon src={PaperAirplane} size="16" class="hover:stroke-accent" />
		</div>
		{#if $commandedPlanet}
			<PlanetSummary />
		{/if}
	</div>
</div>

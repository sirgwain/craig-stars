<script lang="ts">
	import { clamp } from '$lib/services/Math';
	import { PRT, type Race } from '$lib/types/Race';
	import type { FocusEventHandler } from 'svelte/elements';

	export let race: Race;
</script>

{#if race.prt === PRT.AR}
	<p class="mb-1">
		Annual Resources = Planet Value * sqrt(Population * Energy Tech /
		<span class="text-accent">{race.popEfficiency}</span>)
	</p>
{:else}
	<p class="mb-1">
		One resource is generated each year for every
		<span class="text-accent">{race.popEfficiency}</span>
		colonists.
	</p>
	<p class="mb-1">
		Every 10 factories produce
		<span class="text-accent">{race.factoryOutput}</span>
		resources each year.
	</p>
	<p class="mb-1">
		Factories require
		<span class="text-accent">{race.factoryCost}</span>
		resources to build.
	</p>
	<p class="mb-1">
		Every 10,000 colonists may operate up to
		<span class="text-accent">{race.numFactories}</span>
		factories.
	</p>

	{#if race.factoriesCostLess}
		<p class="mb-1">Factories cost <span class="text-accent">1kT</span> less of Germanium to build</p>
	{/if}
	<p class="mb-1">
		Every 10 mines produce up to
		<span class="text-accent">{race.mineOutput}</span>
		kT of each mineral every year.
	</p>
	<p class="mb-1">
		Mines require
		<span class="text-accent">{race.mineCost}</span>
		resources to build.
	</p>
	<p class="mb-1">
		Every 10,000 colonists may operate up to
		<span class="text-accent">{race.numMines}</span>
		mines.
	</p>
{/if}

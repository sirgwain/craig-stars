<script lang="ts" context="module">
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { Player } from '$lib/types/Player';
	export type PopulationTooltipProps = {
		game: FullGame;
		player: Player;
		planet: Planet;
	};
</script>

<script lang="ts">
	import { owned, ownedBy } from '$lib/types/MapObject';
	import type { FullGame } from '$lib/services/FullGame';

	export let game: FullGame;
	export let player: Player;
	export let planet: Planet;

</script>

<div class="flex flex-col sm:w-[26rem] m-auto">
	<div>
		{#if ownedBy(planet, player.num) && planet.spec.population}
			{#if planet.spec.habitability ?? 0 > 0}
				Your population on <span class="font-semibold">{planet.name}</span> is
				<span class="font-semibold">{planet.spec.population.toLocaleString()}</span>.
				<span class="font-semibold">{planet.name}</span> will support a population of up to
				<span class="font-semibold">{planet.spec.maxPopulation.toLocaleString()}</span>
				of your colonists.
			{:else}
				Your population on <span class="font-semibold">{planet.name}</span> is
				<span class="font-semibold">{planet.spec.population.toLocaleString()}</span>." + $"<span
					class="font-semibold">{planet.name}</span
				> has a hostile environment and will no support any of your colonists.
			{/if}
			{#if planet.spec.growthAmount > 0}
				Your population on <span class="font-semibold">{planet.name}</span> will grow by {planet.spec.growthAmount.toLocaleString()}
				to {(planet.spec.population + planet.spec.growthAmount).toLocaleString()} next year.
			{:else if planet.spec.growthAmount === 0}
				Your population on <span class="font-semibold">{planet.name}</span> will not grow next year.
			{:else if planet.spec.growthAmount < 0}
				Approximately {Math.abs(planet.spec.growthAmount).toLocaleString()} of your colonists will die
				next year.
			{/if}
		{:else if !owned(planet) && planet.reportAge != Unexplored}
			<span class="font-semibold">{planet.name} is uninhabited. {planet.spec.habitability}</span>
			{#if planet.spec.habitability && planet.spec.habitability > 0}
				If you were to colonize <span class="font-semibold">{planet.name}</span>, it would support
				up to <span class="font-semibold">{planet.spec.maxPopulation?.toLocaleString()}</span>
				of your colonists
			{:else}
				<span class="font-semibold">{planet.name}</span> will kill off approximately
				<span class="font-semibold"
					>{(Math.abs(planet.spec.habitability ?? 0) / 10).toFixed(1)}%</span
				> of all colonists you settle on it every turn.
			{/if}
		{:else if owned(planet) && planet.reportAge != Unexplored}
			<span class="font-semibold">{planet.name}</span> is currently occupied by the
			<span class="font-semibold">{game.getPlayerName(planet.playerNum)}</span>.
			{#if planet.spec.habitability ?? 0 > 0}
				If you were to colonize <span class="font-semibold">{planet.name}</span>, it would support
				up to <span class="font-semibold">{planet.spec.maxPopulation?.toLocaleString()}</span>
				of your colonists
			{:else}
				<span class="font-semibold">{planet.name}</span> will kill off approximately
				<span class="font-semibold"
					>{(Math.abs(planet.spec.habitability ?? 0) / 10).toFixed(1)}%</span
				> of all colonists you settle on it every turn.
			{/if}
		{:else}
			<span class="font-semibold">{planet.name}</span> is unexplored. Send a scout ship to this planet
			to determine its habitability.
		{/if}
	</div>
</div>

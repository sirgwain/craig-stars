<script lang="ts" context="module">
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import type { Player } from '$lib/types/Player';
	export type PopulationTooltipProps = {
		playerFinder: PlayerFinder;
		player: Player;
		planet: Planet;
	};
</script>

<script lang="ts">
	import type { PlayerFinder } from '$lib/services/Universe';
	import { owned, ownedBy } from '$lib/types/MapObject';

	export let playerFinder: PlayerFinder;
	export let player: Player;
	export let planet: Planet;
</script>

<div class="flex flex-col sm:w-[26rem] m-auto">
	<div>
		{#if ownedBy(planet, player.num) && planet.spec.population}
			{#if (planet.spec.habitability ?? 0) > 0 || player.race.spec?.livesOnStarbases}
				Your population on <span class="font-semibold">{planet.name}</span> is
				<span class="font-semibold">{planet.spec.population.toLocaleString()}</span> ({(
					planet.spec.populationDensity * 100
				).toFixed()}% of capacity).
				<span class="font-semibold">{planet.name}</span> will support a population of up to
				<span class="font-semibold">{planet.spec.maxPopulation?.toLocaleString() ?? 0}</span>
				of your colonists.
			{:else}
				Your population on <span class="font-semibold">{planet.name}</span> is
				<span class="font-semibold">{planet.spec.population.toLocaleString()}</span>.
				<span class="font-semibold">{planet.name}</span> has a hostile environment and will only
				support up to
				<span class="font-semibold">{planet.spec.maxPopulation?.toLocaleString() ?? 0}</span>
				of your colonists.
			{/if}

			{#if planet.spec.growthAmount > 0}
				Your population on <span class="font-semibold">{planet.name}</span> will grow by
				{planet.spec.growthAmount.toLocaleString()}
				to {(planet.spec.population + planet.spec.growthAmount).toLocaleString()} next year.
			{:else if planet.spec.growthAmount === 0}
				Your population on <span class="font-semibold">{planet.name}</span> will not grow next year.
			{:else if planet.spec.growthAmount < 0}
				{#if planet.spec.populationDensity > 1}
					<span class="font-semibold">{planet.name}</span> is overcrowded.
				{/if}
				Approximately {Math.abs(planet.spec.growthAmount).toLocaleString()} of your colonists will die
				next year.
			{/if}
		{:else if !owned(planet) && planet.reportAge != Unexplored}
			<span class="font-semibold">{planet.name} is uninhabited. </span>
			{#if planet.spec.habitability && planet.spec.habitability > 0}
				If you were to colonize <span class="font-semibold">{planet.name}</span>, it would support
				up to <span class="font-semibold">{planet.spec.maxPopulation?.toLocaleString()}</span>
				of your colonists.
			{:else}
				<span class="font-semibold">{planet.name}</span> will kill off approximately
				<span class="font-semibold"
					>{(Math.abs(planet.spec.habitability ?? 0) / 10).toFixed(1)}%</span
				> of all colonists you settle on it every turn.
			{/if}
		{:else if owned(planet) && planet.reportAge != Unexplored}
			The <span class="font-semibold">{playerFinder.getPlayerName(planet.playerNum)}</span>
			population on
			<span class="font-semibold">{planet.name}</span> is approximately
			<span class="font-semibold">{(planet.spec.population ?? 0).toLocaleString()}</span>.
			{#if (planet.spec.habitability ?? 0) > 0}
				If you were to colonize <span class="font-semibold">{planet.name}</span>, it would support
				up to <span class="font-semibold">{planet.spec.maxPopulation?.toLocaleString()}</span>
				of your colonists.
			{:else}
				<span class="font-semibold">{planet.name}</span> will kill off approximately
				<span class="font-semibold"
					>{(Math.abs(planet.spec.habitability ?? 0) / 10).toFixed(1)}%</span
				> of all colonists you settle on it every turn.
			{/if}

			{#if (planet.spec.defenseCoverage ?? 0) == 0}
				<span class="font-semibold">{planet.name}</span> appears to have no planetary defenses.
			{:else}
				<span class="font-semibold">{planet.name}</span> appears to have planetary defenses with
				approximately {Math.round((planet.spec.defenseCoverage ?? 0) * 100)}% coverage.
			{/if}
		{:else}
			<span class="font-semibold">{planet.name}</span> is unexplored. Send a scout ship to this planet
			to determine its habitability.
		{/if}
	</div>
</div>

<script lang="ts" context="module">
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import { roundToNearest100 } from '$lib/services/Math'
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
			<p>Your population on <span class="font-semibold">{planet.name}</span> is
			<span class="font-semibold">{planet.spec.population.toLocaleString()}</span> ({(
			planet.spec.populationDensity * 100
			).toFixed()}% of capacity).</p>
			{#if (planet.spec.habitability ?? 0) > 0 || player.race.spec?.livesOnStarbases}
				<p><span class="font-semibold">{planet.name}</span> will support a population of up to
				<span class="font-semibold">{planet.spec.maxPopulation?.toLocaleString() ?? 0}</span>
				of your colonists.</p>
			{:else}
				<p><span class="font-semibold">{planet.name}</span> has a hostile environment and will only
				support up to
				<span class="font-semibold">{planet.spec.maxPopulation?.toLocaleString() ?? 0}</span>
				of your colonists.</p>
			{/if}

			{#if planet.spec.growthAmount > 0}
				<p>
				Your population on <span class="font-semibold">{planet.name}</span> will grow by
				<span class="font-semibold">{planet.spec.growthAmount.toLocaleString()}</span>
				to {(planet.spec.population + planet.spec.growthAmount).toLocaleString()}
				next year.</p>
			{:else if planet.spec.growthAmount === 0}
				<p>Your population on <span class="font-semibold">{planet.name}</span> will not grow next year.</p>
			{:else if planet.spec.growthAmount < 0}
				{#if planet.spec.populationDensity > 1}
					<p><span class="font-semibold">{planet.name}</span> is overcrowded.</p>
				{/if}
				<p>Approximately
				<span class="font-semibold">{Math.abs(planet.spec.growthAmount).toLocaleString()}</span>
				of your colonists will die next year.</p>
			{/if}
		{:else if !owned(planet) && planet.reportAge != Unexplored}
			<p><span class="font-semibold">{planet.name} is uninhabited.</span></p>

			{#if planet.spec.habitability && planet.spec.habitability > 0}
				<p>If you were to colonize <span class="font-semibold">{planet.name}</span>, it would support
				up to <span class="font-semibold">{planet.spec.maxPopulation?.toLocaleString()}</span>
				of your colonists.</p>
			{:else}
				<p><span class="font-semibold">{planet.name}</span> will kill off approximately
				<span class="font-semibold"
					>{(Math.abs(planet.spec.habitability ?? 0) / 10).toFixed(1)}%</span
				> of all colonists you settle on it every turn.</p>
			{/if}
		{:else if owned(planet) && planet.reportAge != Unexplored}
			<p>The <span class="font-semibold">{playerFinder.getPlayerName(planet.playerNum)}</span>
			population on
			<span class="font-semibold">{planet.name}</span> is approximately
			<span class="font-semibold">{roundToNearest100((planet.spec.population ?? 0)).toLocaleString()}</span>.</p>
			{#if (planet.spec.habitability ?? 0) > 0}
				<p>If you were to colonize <span class="font-semibold">{planet.name}</span>, it would support
				up to <span class="font-semibold">{planet.spec.maxPopulation?.toLocaleString()}</span>
				of your colonists.</p>
			{:else}
				<p><span class="font-semibold">{planet.name}</span> will kill off approximately
				<span class="font-semibold"
					>{(Math.abs(planet.spec.habitability ?? 0) / 10).toFixed(1)}%</span
				> of all colonists you settle on it every turn.</p>
			{/if}

			{#if (planet.spec.defenseCoverage ?? 0) == 0}
				<p>
				<span class="font-semibold">{planet.name}</span> appears to have no planetary defenses.</p>
			{:else}
				<p><span class="font-semibold">{planet.name}</span> appears to have planetary defenses with
				approximately {Math.round((planet.spec.defenseCoverage ?? 0) * 100)}% coverage.</p>
			{/if}
		{:else}
			<p>
			<span class="font-semibold">{planet.name}</span> is unexplored. Send a scout ship to this planet
			to determine its habitability.</p>
		{/if}
	</div>
</div>

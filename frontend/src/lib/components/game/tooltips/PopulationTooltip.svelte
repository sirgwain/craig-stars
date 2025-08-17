<script lang="ts" module>
	import { ReportAgeUnexplored } from '$lib/types/Consts';
	import type { CommandedPlayer } from '$lib/types/Player';
	export type PopulationTooltipProps = {
		playerFinder: PlayerFinder;
		player: CommandedPlayer;
		planet: AnyPlanet;
	};
</script>

<script lang="ts">
	import type { AnyPlanet, PlayerFinder } from '$lib/services/Universe';
	import { owned, ownedBy } from '$lib/types/MapObject';
	import { population } from '$lib/types/Cargo';
	import { getGrowth } from '$lib/types/Planet';

	let { playerFinder, player, planet }: PopulationTooltipProps = $props();

	let reportAge = $derived('reportAge' in planet ? (planet.reportAge ?? 0) : 0);
	let habitability = $derived(planet.spec?.habitability ?? 0);
	let pop = $derived(population(planet.cargo));
	let growthAmount = $derived(getGrowth(planet));
</script>

<div class="flex flex-col sm:w-[26rem] m-auto">
	<div>
		{#if ownedBy(planet, player.num) && pop}
			<p>
				Your population on <span class="font-semibold">{planet.mapObject?.name}</span> is
				<span class="font-semibold">{pop.toLocaleString()}</span> ({(
					(planet.spec?.populationDensity ?? 0) * 100
				).toFixed()}% of capacity).
			</p>
			{#if (planet.spec?.habitability ?? 0) > 0 || player.race.spec.livesOnStarbases}
				<p>
					<span class="font-semibold">{planet.mapObject?.name}</span> will support a population of
					up to
					<span class="font-semibold">{planet.spec?.maxPopulation?.toLocaleString() ?? 0}</span>
					of your colonists.
				</p>
			{:else}
				<p>
					<span class="font-semibold">{planet.mapObject?.name}</span> has a hostile environment and
					will only support up to
					<span class="font-semibold">{planet.spec?.maxPopulation?.toLocaleString() ?? 0}</span>
					of your colonists.
				</p>
			{/if}

			{#if growthAmount > 0}
				<p>
					Your population on <span class="font-semibold">{planet.mapObject?.name}</span> will grow
					by
					<span class="font-semibold">{growthAmount.toLocaleString()}</span>
					to {(pop + growthAmount).toLocaleString()}
					next year.
				</p>
			{:else if growthAmount < 0}
				{#if (planet.spec?.populationDensity ?? 0) > 1}
					<p><span class="font-semibold">{planet.mapObject?.name}</span> is overcrowded.</p>
				{/if}
				<p>
					Approximately
					<span class="font-semibold">{Math.abs(growthAmount).toLocaleString()}</span>
					of your colonists will die next year.
				</p>
			{:else}
				<p>
					Your population on <span class="font-semibold">{planet.mapObject?.name}</span> will not grow
					next year.
				</p>
			{/if}
		{:else if !owned(planet) && reportAge !== ReportAgeUnexplored}
			<p><span class="font-semibold">{planet.mapObject?.name} is uninhabited.</span></p>

			{#if habitability > 0}
				<p>
					If you were to colonize <span class="font-semibold">{planet.mapObject?.name}</span>, it
					would support up to
					<span class="font-semibold">{planet.spec?.maxPopulation?.toLocaleString()}</span>
					of your colonists.
				</p>
			{:else}
				<p>
					<span class="font-semibold">{planet.mapObject?.name}</span> will kill off approximately
					<span class="font-semibold">{(Math.abs(habitability) / 10).toFixed(1)}%</span> of all colonists
					you settle on it every turn.
				</p>
			{/if}
		{:else if owned(planet) && reportAge != ReportAgeUnexplored}
			<p>
				The <span class="font-semibold"
					>{playerFinder.getPlayerName(planet.mapObject?.playerNum)}</span
				>
				population on
				<span class="font-semibold">{planet.mapObject?.name}</span> is approximately
				<span class="font-semibold">{pop.toLocaleString()}</span>. <!-- Rounded to 100 in backend -->
			</p>
			{#if habitability > 0}
				<p>
					If you were to colonize <span class="font-semibold">{planet.mapObject?.name}</span>, it
					would support up to
					<span class="font-semibold">{planet.spec?.maxPopulation?.toLocaleString()}</span>
					of your colonists.
				</p>
			{:else}
				<p>
					<span class="font-semibold">{planet.mapObject?.name}</span> will kill off approximately
					<span class="font-semibold">{(Math.abs(habitability) / 10).toFixed(1)}%</span> of all colonists
					you settle on it every turn.
				</p>
			{/if}

			{#if (planet.spec?.defenseCoverage ?? 0) == 0}
				<p>
					<span class="font-semibold">{planet.mapObject?.name}</span> appears to have no planetary defenses.
				</p>
			{:else}
				<p>
					<span class="font-semibold">{planet.mapObject?.name}</span> appears to have planetary
					defenses with approximately {Math.round((planet.spec?.defenseCoverage ?? 0) * 100)}%
					coverage.
				</p>
			{/if}
		{:else}
			<p>
				<span class="font-semibold">{planet.mapObject?.name}</span> is unexplored. Send a scout ship
				to this planet to determine its habitability.
			</p>
		{/if}
	</div>
</div>

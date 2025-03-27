<script lang="ts">
	import PlanetHabBars from '$lib/components/game/PlanetHabBars.svelte';
	import PlanetHabValue from '$lib/components/game/PlanetHabValue.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { population } from '$lib/types/Cargo';
	import { ReportAgeUnexplored, type MapObject } from '$lib/types/cs';
	import { getUnderlyingMapObject, ownedBy } from '$lib/types/MapObject';
	import MapObjectIcon from './MapObjectIcon.svelte';

	const { player } = getGameContext();

	type Props = {
		mapObject: MapObject | undefined;
	};

	let { mapObject }: Props = $props();

	let { planet, fleet } = $derived(getUnderlyingMapObject(mapObject));
</script>

<div class="flex flex-row justify-start gap-3 text-sm">
	<MapObjectIcon {mapObject} />
	{#if planet}
		<div class="flex flex-col w-full">
			{#if ('reportAge' in planet && planet.reportAge !== ReportAgeUnexplored) || !('reportAge' in planet)}
				<div class="flex flex-row justify-between">
					<div>
						Value: <PlanetHabValue {planet} />
					</div>
					<div>
						{#if population(planet.cargo)}
							Pop: {population(planet.cargo).toLocaleString()}
						{/if}
					</div>
				</div>
			{:else}
				Unexplored
			{/if}

			{#if ownedBy(planet, $player.num)}
				<div class="mt-1">
					<div class="flex justify-between">
						<div class="text-tile-item-title text-ironium">Ironium</div>
						<div>
							{(planet.cargo.ironium ?? 0).toLocaleString()}kT ({planet.mineralConcentration
								?.ironium})
						</div>
					</div>
					<div class="flex justify-between">
						<div class="text-tile-item-title text-boranium">Boranium</div>
						<div>
							{(planet.cargo.boranium ?? 0).toLocaleString()}kT ({planet.mineralConcentration
								?.boranium})
						</div>
					</div>
					<div class="flex justify-between">
						<div class="text-tile-item-title text-germanium">Germanium</div>
						<div>
							{(planet.cargo.germanium ?? 0).toLocaleString()}kT ({planet.mineralConcentration
								?.germanium})
						</div>
					</div>
				</div>
			{:else if 'reportAge' in planet && planet.reportAge !== ReportAgeUnexplored}
				<div class="my-auto">
					<PlanetHabBars {planet} player={$player} />
				</div>
			{/if}
		</div>
	{:else if fleet}
		<!-- else content here -->
	{:else}
		<!-- else content here -->
	{/if}
</div>

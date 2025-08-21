<script lang="ts">
	import PlanetHabBars from '$lib/components/game/PlanetHabBars.svelte';
	import PlanetHabValue from '$lib/components/game/PlanetHabValue.svelte';
	import type { HabTooltipProps } from '$lib/components/game/tooltips/HabTooltip.svelte';
	import HabTooltip from '$lib/components/game/tooltips/HabTooltip.svelte';
	import PopulationTooltip, {
		type PopulationTooltipProps
	} from '$lib/components/game/tooltips/PopulationTooltip.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { showTooltip } from '$lib/services/Stores';
	import { population } from '$lib/types/Cargo';
	import { None, ReportAgeUnexplored } from '$lib/types/Consts';
	import type { Planet } from '$lib/types/cs-proto';
	import { Grav, Rad, Temp } from '$lib/types/Hab';
	import { ownedBy } from '$lib/types/MapObject';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import MapObjectIcon from './MapObjectIcon.svelte';
	import PlanetMineralsGraph from './PlanetMineralsGraph.svelte';

	const { player, universe } = getGameContext();

	type Props = {
		planet: Planet;
	};

	let { planet }: Props = $props();

	function onPopulationTooltip(e: PointerEvent) {
		e.preventDefault();
		showTooltip<PopulationTooltipProps>(e.x, e.y, PopulationTooltip, {
			playerFinder: $universe,
			player: $player,
			planet
		});
	}

	function onGravityTooltip(e: PointerEvent) {
		e.preventDefault();
		showTooltip<HabTooltipProps>(e.x, e.y, HabTooltip, {
			player: $player,
			planet,
			habType: Grav
		});
	}

	function onTemperatureTooltip(e: PointerEvent) {
		e.preventDefault();
		showTooltip<HabTooltipProps>(e.x, e.y, HabTooltip, {
			player: $player,
			planet,
			habType: Temp
		});
	}

	function onRadiationTooltip(e: PointerEvent) {
		e.preventDefault();
		showTooltip<HabTooltipProps>(e.x, e.y, HabTooltip, {
			player: $player,
			planet,
			habType: Rad
		});
	}
</script>

<div class="flex flex-col md:min-h-[11rem] select-none w-full">
	{#if planet.mapObject?.reportAge === ReportAgeUnexplored}
		<div class="relative w-full m-auto">
			<!-- Icon on the left -->
			<div class="absolute top-1/2 -translate-y-1/2">
				<!-- Your icon here -->
				<MapObjectIcon mapObject={planet} />
			</div>

			<!-- Centered content -->
			<div>
				<Icon src={QuestionMarkCircle} size="64" class="hover:stroke-accent m-auto" />
			</div>
		</div>
	{:else}
		<div class="flex justify-between cursor-help" onpointerdown={onPopulationTooltip}>
			<div class="ml-[5.5rem]">
				Value: <PlanetHabValue {planet} />
			</div>
			{#if population(planet.cargo)}
				<div>Population: {population(planet.cargo).toLocaleString()}</div>
			{/if}
		</div>
		<div class="flex justify-between">
			<div class="ml-[5.5rem]">
				<div>
					{#if ownedBy(planet, $player.num)}
						Report is current
					{:else if planet.mapObject?.reportAge === 0}
						Report is current
					{:else if planet.mapObject?.reportAge === 1}
						Report is 1 year old
					{:else}
						Report is {planet.mapObject?.reportAge} years old
					{/if}
				</div>
			</div>
			<div>
				{#if planet.mapObject?.reportAge !== ReportAgeUnexplored && (planet.mapObject?.playerNum ?? None) != $player.num && (planet.mapObject?.playerNum ?? None) != None}
					<span style={`color: ${$universe.getPlayerColor(planet.mapObject?.playerNum)}`}
						>{$universe.getPlayerPluralName(planet.mapObject?.playerNum)}</span
					>
				{/if}
			</div>
		</div>

		<PlanetHabBars
			{planet}
			player={$player}
			{onGravityTooltip}
			{onTemperatureTooltip}
			{onRadiationTooltip}
		/>

		<div class="mb-1"></div>

		<PlanetMineralsGraph {planet} />
	{/if}
</div>

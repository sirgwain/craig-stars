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
	import type { AnyPlanet } from '$lib/services/Universe';
	import { population } from '$lib/types/Cargo';
	import { Grav, None, Rad, ReportAgeUnexplored, Temp } from '$lib/types/cs';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import MapObjectIcon from './MapObjectIcon.svelte';
	import PlanetMineralsGraph from './PlanetMineralsGraph.svelte';

	const { player, universe } = getGameContext();

	type Props = {
		planet: AnyPlanet;
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
	{#if 'reportAge' in planet && planet.reportAge === ReportAgeUnexplored}
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
					{#if 'reportAge' in planet}
						{#if (planet.reportAge ?? 0) == 0}
							Report is current
						{:else if planet.reportAge == 1}
							Report is 1 year old
						{:else}
							Report is {planet.reportAge} years old
						{/if}
					{:else}
						Report is current
					{/if}
				</div>
			</div>
			<div>
				{#if 'reportAge' in planet && planet.reportAge !== ReportAgeUnexplored && planet.playerNum != $player.num && planet.playerNum != None}
					<span style={`color: ${$universe.getPlayerColor(planet.playerNum)}`}
						>{$universe.getPlayerPluralName(planet.playerNum)}</span
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

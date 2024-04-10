<script lang="ts">
	import PlanetHabPoint from '$lib/components/game/PlanetHabPoint.svelte';
	import PlanetHabTerraformLine from '$lib/components/game/PlanetHabTerraformLine.svelte';
	import type { HabTooltipProps } from '$lib/components/game/tooltips/HabTooltip.svelte';
	import HabTooltip from '$lib/components/game/tooltips/HabTooltip.svelte';
	import PopulationTooltip, {
		type PopulationTooltipProps
	} from '$lib/components/game/tooltips/PopulationTooltip.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { clamp } from '$lib/services/Math';
	import { showTooltip } from '$lib/services/Stores';
	import { HabTypes, add, getGravString, getRadString, getTempString } from '$lib/types/Hab';
	import { None } from '$lib/types/MapObject';
	import { Unexplored, type Planet } from '$lib/types/Planet';
	import { QuestionMarkCircle } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import PlanetMineralsGraph from './PlanetMineralsGraph.svelte';

	const { game, player, universe } = getGameContext();

	export let planet: Planet;

	$: habLow = $player.race.habLow;
	$: habHigh = $player.race.habHigh;
	$: habPoint = planet.hab ?? {};
	$: terraformHabPoint = add(habPoint, planet.spec.terraformAmount ?? {});
	$: habWidth = {
		grav: (habHigh.grav ?? 0) - (habLow.grav ?? 0),
		temp: (habHigh.temp ?? 0) - (habLow.temp ?? 0),
		rad: (habHigh.rad ?? 0) - (habLow.rad ?? 0)
	};

	$: habPointPercent = {
		grav: clamp(habPoint.grav ? (habPoint.grav / 100) * 100 : 0, 0, 100),
		temp: clamp(habPoint.temp ? (habPoint.temp / 100) * 100 : 0, 0, 100),
		rad: clamp(habPoint.rad ? (habPoint.rad / 100) * 100 : 0, 0, 100)
	};
	$: terraformHabPointPercent = {
		grav: clamp(terraformHabPoint.grav ? (terraformHabPoint.grav / 100) * 100 : 0, 0, 100),
		temp: clamp(terraformHabPoint.temp ? (terraformHabPoint.temp / 100) * 100 : 0, 0, 100),
		rad: clamp(terraformHabPoint.rad ? (terraformHabPoint.rad / 100) * 100 : 0, 0, 100)
	};
	$: habLowPercent = {
		grav: clamp(habLow.grav ? (habLow.grav / 100) * 100 : 0, 0, 100),
		temp: clamp(habLow.temp ? (habLow.temp / 100) * 100 : 0, 0, 100),
		rad: clamp(habLow.rad ? (habLow.rad / 100) * 100 : 0, 0, 100)
	};
	$: habWidthPercent = {
		grav: clamp(habWidth.grav ? (habWidth.grav / 100) * 100 : 0, 0, 100),
		temp: clamp(habWidth.temp ? (habWidth.temp / 100) * 100 : 0, 0, 100),
		rad: clamp(habWidth.rad ? (habWidth.rad / 100) * 100 : 0, 0, 100)
	};

	function onPopulationTooltip(e: PointerEvent) {
		showTooltip<PopulationTooltipProps>(e.x, e.y, PopulationTooltip, {
			playerFinder: $universe,
			player: $player,
			planet
		});
	}

	function onGravityTooltip(e: PointerEvent) {
		showTooltip<HabTooltipProps>(e.x, e.y, HabTooltip, {
			player: $player,
			planet,
			habType: HabTypes.Gravity
		});
	}

	function onTemperatureTooltip(e: PointerEvent) {
		showTooltip<HabTooltipProps>(e.x, e.y, HabTooltip, {
			player: $player,
			planet,
			habType: HabTypes.Temperature
		});
	}

	function onRadiationTooltip(e: PointerEvent) {
		showTooltip<HabTooltipProps>(e.x, e.y, HabTooltip, {
			player: $player,
			planet,
			habType: HabTypes.Radiation
		});
	}
</script>

<div class="flex flex-col min-h-[11rem] select-none">
	{#if planet.reportAge == Unexplored}
		<div class="m-auto">
			<Icon src={QuestionMarkCircle} size="64" class="hover:stroke-accent" />
		</div>
	{:else}
		<div
			class="flex justify-between cursor-help"
			on:pointerdown|preventDefault={onPopulationTooltip}
		>
			<div class="ml-[5.5rem]">
				Value: <span
					class:text-habitable={(planet.spec.habitability ?? 0) > 0}
					class:text-uninhabitable={(planet.spec.habitability ?? 0) < 0}
					class:text-terraformable={(planet.spec.habitability ?? 0) < 0 &&
						(planet.spec.terraformedHabitability ?? 0) > 0}
					>{planet.spec.habitability}%{planet.spec.terraformedHabitability &&
					planet.spec.terraformedHabitability !== planet.spec.habitability
						? ` (${planet.spec.terraformedHabitability}%)`
						: ''}</span
				>
			</div>
			{#if planet?.spec.population}
				<div>Population: {planet.spec.population.toLocaleString()}</div>
			{/if}
		</div>
		<div class="flex justify-between">
			<div class="ml-[5.5rem]">
				<div>
					{#if (planet.reportAge ?? 0) == 0}
						Report is current
					{:else if planet.reportAge == 1}
						Report is 1 year old
					{:else}
						Report is {planet.reportAge} years old
					{/if}
				</div>
			</div>
			<div>
				{#if planet.reportAge != Unexplored && planet.playerNum != $player.num && planet.playerNum != None}
					<span style={`color: ${$universe.getPlayerColor(planet.playerNum)}`}
						>{$universe.getPlayerName(planet.playerNum)}</span
					>
				{/if}
			</div>
		</div>

		<div class="flex flex-row cursor-help" on:pointerdown|preventDefault={onGravityTooltip}>
			<div class="text-right w-[5.5rem] text-tile-item-title">Gravity</div>
			<div class="grow border-b border-base-300 bg-black mx-1 overflow-hidden">
				<div class="h-full relative">
					{#if !$player.race.immuneGrav}
						<div
							style={`left: ${habLowPercent.grav.toFixed()}%; width: ${habWidthPercent.grav?.toFixed()}%`}
							class="absolute grav-bar h-full"
						/>
					{/if}
					<PlanetHabPoint
						style={`left: ${habPointPercent.grav.toFixed()}%;`}
						class="absolute grav-point h-full -translate-x-1/2"
					/>
					<!-- Terraform line -->
					<div class="absolute h-full w-full">
						<svg width="100%" height="100%" viewBox="0 0 100 100" preserveAspectRatio="none">
							<line
								x1={habPointPercent.grav}
								y1="50"
								x2={terraformHabPointPercent.grav}
								y2="50"
								vector-effect="non-scaling-stroke"
								stroke-width="1"
								class="grav-point"
							/>
						</svg>
					</div>
				</div>
			</div>
			<div class="w-[3rem]">{getGravString(planet.hab?.grav ?? 0)}</div>
		</div>
		<div class="flex flex-row cursor-help" on:pointerdown|preventDefault={onTemperatureTooltip}>
			<div class="text-right w-[5.5rem] text-tile-item-title">Temperature</div>
			<div class="grow border-b border-base-300 bg-black mx-1 overflow-hidden">
				<div class="h-full relative">
					{#if !$player.race.immuneTemp}
						<div
							style={`left: ${habLowPercent.temp.toFixed()}%; width: ${habWidthPercent.temp?.toFixed()}%`}
							class="absolute temp-bar h-full"
						/>
					{/if}
					<PlanetHabPoint
						style={`left: ${habPointPercent.temp.toFixed()}%;`}
						class="absolute temp-point h-full -translate-x-1/2"
					/>
					<!-- Terraform line -->
					<PlanetHabTerraformLine
						x1={habPointPercent.temp}
						x2={terraformHabPointPercent.temp}
						class="temp-point"
					/>
				</div>
			</div>
			<div class="w-[3rem]">{getTempString(planet.hab?.temp ?? 0)}</div>
		</div>
		<div class="flex flex-row cursor-help" on:pointerdown|preventDefault={onRadiationTooltip}>
			<div class="text-right w-[5.5rem] text-tile-item-title">Radiation</div>
			<div class="grow bg-black mx-1 overflow-hidden">
				<div class="h-full relative">
					{#if !$player.race.immuneRad}
						<div
							style={`left: ${habLowPercent.rad.toFixed()}%; width: ${habWidthPercent.rad?.toFixed()}%`}
							class="absolute rad-bar h-full"
						/>
					{/if}
					<PlanetHabPoint
						style={`left: ${habPointPercent.rad.toFixed()}%;`}
						class="absolute rad-point h-full -translate-x-1/2"
					/>
					<!-- Terraform line -->
					<PlanetHabTerraformLine
						x1={habPointPercent.rad}
						x2={terraformHabPointPercent.rad}
						class="rad-point"
					/>
				</div>
			</div>
			<div class="w-[3rem]">{getRadString(planet.hab?.rad ?? 0)}</div>
		</div>

		<div class="mb-1" />

		<PlanetMineralsGraph {planet} />
	{/if}
</div>

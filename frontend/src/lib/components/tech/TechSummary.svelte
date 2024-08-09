<script lang="ts">
	import type { Player } from '$lib/types/Player';

	import {
		isHullComponent,
		TechCategory,
		type Tech,
		type TechDefense,
		type TechEngine,
		type TechHull,
		type TechHullComponent
	} from '$lib/types/Tech';
	import Cost from '../game/Cost.svelte';
	import TechDescription from './TechDescription.svelte';
	import TechEngineGraph from './TechEngineGraph.svelte';
	import TechLevelRequirements from './TechLevelRequirements.svelte';
	import TechTraitRequirements from './TechTraitRequirements.svelte';

	import { levelsAbove } from '$lib/types/TechLevel';
	import { kebabCase } from 'lodash-es';
	import TechAvatar from './TechAvatar.svelte';
	import TechDefenseGraph from './TechDefenseGraph.svelte';
	import { PlayerService } from '$lib/services/PlayerService';
	import type { FullGame } from '$lib/services/FullGame';
	import TechWarnings from './TechWarnings.svelte';

	export let tech: Tech;
	export let player: Player | undefined = undefined;
	export let game: FullGame | undefined = undefined;
	export let showResearchCost = false;
	export let hideGraph = false;

	let defense: TechDefense;
	let hullComponent: TechHullComponent;
	let hull: TechHull;
	let engine: TechEngine;
	let researchCost = 0;

	$: tech && isHullComponent(tech.category) && (hullComponent = tech as TechHullComponent);
	$: tech && tech.category == TechCategory.ShipHull && (hull = tech as TechHull);
	$: tech && tech.category == TechCategory.Engine && (engine = tech as TechEngine);
	$: tech && tech.category == TechCategory.PlanetaryDefense && (defense = tech as TechDefense);
	$: above = player?.hasTech(tech) ? levelsAbove(tech.requirements, player.techLevels) : 0;

	$: {
		if (showResearchCost && player && game) {
			PlayerService.getResearchCost(game.id, tech.requirements).then(
				(result) => (researchCost = result.resources)
			);
		}
	}
</script>

{#if tech}
	<div
		class="card bg-base-200 shadow rounded-sm border-2 border-base-300 max-h-fit min-h-fit w-full h-full"
	>
		<div class="card-body p-3 gap-0">
			<h2 class="text-lg font-semibold text-center mb-1 text-secondary">
				<div class="indicator w-full">
					{#if player?.hasTech(tech)}
						<span class:hidden={!player || above !== 0} class="indicator-item badge badge-accent"
							>new
						</span>
					{/if}
					<div class="w-full">
						{#if player}
							<a href={`/games/${player.gameId}/techs/${kebabCase(tech.name.replaceAll("'", ''))}`}
								>{tech.name}</a
							>
						{:else}
							<a href="/techs/{kebabCase(tech.name.replaceAll("'", ''))}">{tech.name}</a>
						{/if}
					</div>
				</div>
			</h2>

			<div class="flex flex-row gap-2">
				<div class="flex flex-col flex-initial min-w-[6rem]">
					<!-- icon and tech requirements row-->
					<TechAvatar {tech} hullTooltip={true} />
					<TechLevelRequirements {tech} {player} />
					{#if showResearchCost && researchCost}
						<div
							class="flex flex-row justify-between gap-1"
							class:text-error={player?.techLevels &&
								(player.techLevels.energy ?? 0) < (tech.requirements.energy ?? 0)}
						>
							<div>Research Cost:</div>
							<div>{researchCost}</div>
						</div>
					{/if}
				</div>

				<div class="flex flex-col flex-1">
					<div class="flex flex-row gap-2">
						<!-- cost -->
						{#if player}
							<Cost cost={player.getTechCost(tech)} />
						{:else}
							<Cost cost={tech.cost} />
						{/if}

						{#if hullComponent}
							<div class="flex justify-between gap-2">
								<div>Mass:</div>
								<div>{hullComponent.mass ?? 0}kT</div>
							</div>
						{:else if hull}
							<div class="flex justify-between gap-2">
								<div>Mass:</div>
								<div>{hull.mass ?? 0}kT</div>
							</div>
						{/if}
					</div>
					{#if engine}
						<div class="grow min-h-[14rem]">
							{#if !hideGraph}
								<TechEngineGraph {engine} />
							{/if}
						</div>
					{:else if defense}
						<div class="grow min-h-[14rem]">
							{#if !hideGraph}
								<TechDefenseGraph {defense} />
							{/if}
						</div>
					{:else}
						<div class="border border-base-300 bg-base-100 grow min-h-[14rem]">
							<TechDescription {tech} />
							<TechWarnings {tech} />
						</div>
					{/if}
				</div>
			</div>
		</div>
		<div class="flex flex-col min-h-[2rem]">
			<TechWarnings {tech} />
			<TechTraitRequirements {tech} {player} />
		</div>
	</div>
{/if}

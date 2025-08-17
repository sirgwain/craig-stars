<script lang="ts">
	import type { CommandedPlayer } from '$lib/types/Player';
	import { isHullComponent, type TechLike } from '$lib/types/Tech';
	import Cost from '../game/Cost.svelte';
	import TechDescription from './TechDescription.svelte';
	import TechEngineGraph from './TechEngineGraph.svelte';
	import TechLevelRequirements from './TechLevelRequirements.svelte';
	import TechTraitRequirements from './TechTraitRequirements.svelte';

	import {
		TechCategory,
		type TechDefense,
		type TechHull,
		type TechHullComponent
	} from '$lib/types/cs-proto';
	import { levelsAbove } from '$lib/types/TechLevel';
	import type { CS } from '$lib/wasm';
	import { kebabCase } from 'lodash-es';
	import TechAvatar from './TechAvatar.svelte';
	import TechDefenseGraph from './TechDefenseGraph.svelte';
	import TechWarnings from './TechWarnings.svelte';

	type Props = {
		tech: TechLike;
		player?: CommandedPlayer | undefined;
		cs?: CS | undefined;
		showResearchCost?: boolean;
		hideGraph?: boolean;
	};

	let { tech: techLike, player, cs, showResearchCost = false, hideGraph = false }: Props = $props();

	let defense = $derived(
		techLike.tech?.category == TechCategory.PLANETARY_DEFENSE
			? (techLike as TechDefense)
			: undefined
	);
	let hullComponent = $derived(
		isHullComponent(techLike.tech?.category) ? (techLike as TechHullComponent) : undefined
	);
	let hull = $derived(
		techLike.tech?.category == TechCategory.SHIP_HULL ? (techLike as TechHull) : undefined
	);
	let engine = $derived(
		techLike.tech?.category == TechCategory.ENGINE ? (techLike as TechHullComponent) : undefined
	);
	let researchCost = $state(0);
	let above = $derived(
		techLike.tech && player?.hasTech(techLike)
			? levelsAbove(techLike.tech?.requirements?.techLevel, player.techLevels)
			: 0
	);

	let cost = $state(techLike.tech?.cost);
	$effect(() => {
		if (!(cs && techLike.tech && player)) {
			return;
		}
		cs.wasmService
			.getTechCost({ tech: techLike.tech })
			.then((resp) => (cost = resp.cost ?? techLike.tech?.cost));
	});
	$effect(() => {
		if (!(techLike.tech?.requirements && showResearchCost && player && cs)) {
			return;
		}
		cs.wasmService
			.getResearchCost({ techLevel: techLike.tech.requirements?.techLevel })
			.then((resp) => {
				researchCost = resp.resources;
			});
	});
</script>

{#if techLike.tech}
	{@const tech = techLike.tech}
	<div
		class="card bg-base-200 shadow rounded-sm border-2 border-base-300 max-h-fit min-h-fit w-full h-full"
	>
		<div class="card-body p-3 gap-0">
			<div class="text-lg font-semibold text-center mb-1 text-secondary">
				<div class="indicator w-full">
					{#if player?.hasTech(techLike)}
						<span class:hidden={!player || above !== 0} class="indicator-item badge badge-accent"
							>new
						</span>
					{/if}
					<div class="w-full">
						{#if player}
							<a
								href={`/games/${player.gameDbObject.gameId}/techs/${kebabCase(tech.name?.replaceAll("'", ''))}`}
								>{tech.name}</a
							>
						{:else}
							<a href="/techs/{kebabCase(tech.name?.replaceAll("'", ''))}">{tech.name}</a>
						{/if}
					</div>
				</div>
			</div>

			<div class="flex flex-row gap-2">
				<div class="flex flex-col flex-initial min-w-[6rem]">
					<!-- icon and tech requirements row-->
					<TechAvatar tech={techLike} hullTooltip={true} />
					<TechLevelRequirements tech={techLike} {player} />
					{#if showResearchCost && researchCost}
						<div
							class="flex flex-row justify-between gap-1"
							class:text-error={player?.techLevels &&
								(player.techLevels.energy ?? 0) < (tech.requirements?.techLevel?.energy ?? 0)}
						>
							<div>Research Cost:</div>
							<div>{researchCost}</div>
						</div>
					{/if}
				</div>

				<div class="flex flex-col flex-1">
					<div class="flex flex-row gap-2">
						<Cost {cost} />

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
							<TechDescription tech={techLike} />
							<TechWarnings tech={techLike} />
						</div>
					{/if}
				</div>
			</div>
		</div>
		<div class="flex flex-col min-h-[2rem]">
			<TechWarnings tech={techLike} />
			<TechTraitRequirements tech={techLike} {player} />
		</div>
	</div>
{/if}

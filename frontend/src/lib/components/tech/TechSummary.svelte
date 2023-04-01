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

	import { kebabCase } from 'lodash-es';
	import TechDefenseGraph from './TechDefenseGraph.svelte';

	export let tech: Tech;
	export let player: Player | undefined = undefined;

	let hull: TechHull;
	let defense: TechDefense;
	let hullComponent: TechHullComponent;
	let engine: TechEngine;
	let hullSetNumber = 0;
	const numHullSets = 4;

	$: tech && isHullComponent(tech.category) && (hullComponent = tech as TechHullComponent);
	$: tech && tech.category == TechCategory.Engine && (engine = tech as TechEngine);
	$: tech && tech.category == TechCategory.PlanetaryDefense && (defense = tech as TechDefense);
	$: tech &&
		(tech.category == TechCategory.ShipHull || tech.category == TechCategory.StarbaseHull) &&
		(hull = tech as TechHull);

	const icon = (hullSetNumber: number) => {
		const name = kebabCase(tech.name.replace("'", '').replace(' ', '').replace('±', ''));
		if (hull) {
			return `hull-${name}-${hullSetNumber}`;
		} else {
			return name;
		}
	};
</script>

{#if tech}
	<div
		class="card bg-base-200 shadow-xl w-full max-h-fit min-h-fit rounded-sm border-2 border-base-300"
	>
		<div class="card-body p-3 gap-0">
			<h2 class="text-lg font-semibold text-center mb-1 text-secondary"><a href="/techs/{kebabCase(tech.name)}">{tech.name}</a></h2>

			<div class="flex flex-row gap-2">
				<div class="flex flex-col flex-initial min-w-[6rem]">
					<!-- icon and tech requirements row-->
					<div
						on:click={() => {
							hullSetNumber += 1;
							hullSetNumber %= numHullSets;
						}}
						class="avatar border border-secondary tech-avatar {icon(hullSetNumber)}"
					/>
					<TechLevelRequirements {tech} {player} />
				</div>

				<div class="flex flex-col flex-1">
					<div class="flex flex-row gap-2">
						<!-- cost -->
						<Cost cost={tech.cost} />

						{#if hullComponent}
							<div class="flex justify-between gap-2">
								<div>Mass:</div>
								<div>{hullComponent.mass ?? 0}kT</div>
							</div>
						{/if}
					</div>
					{#if engine}
						<div class="grow min-h-[14rem]">
							<TechEngineGraph {engine} />
						</div>
					{:else if defense}
						<div class="grow min-h-[14rem]">
							<TechDefenseGraph {defense} />
						</div>
					{:else}
						<div class="border border-base-300 bg-base-100 grow min-h-[14rem]">
							<TechDescription {tech} />
						</div>
					{/if}
				</div>
			</div>
		</div>
		<div class="flex flex-row min-h-[2rem]">
			<TechTraitRequirements {tech} {player} />
		</div>
	</div>
{/if}

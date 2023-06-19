<script lang="ts">
	import type { Player } from '$lib/types/Player';

	import type { Tech, TechEngine, TechHullComponent, TechRequirements } from '$lib/types/Tech';
	import Cost from '../game/Cost.svelte';
	import TechDescription from './TechDescription.svelte';
	import TechEngineGraph from './TechEngineGraph.svelte';
	import TechLevelRequirements from './TechLevelRequirements.svelte';
	import TechTraitRequirements from './TechTraitRequirements.svelte';

	import { kebabCase, replace } from 'lodash-es';

	export let tech: Tech;
	export let player: Player | undefined = undefined;

	let hullComponent: TechHullComponent;
	let engine: TechEngine;

	$: tech && 'mass' in tech && (hullComponent = tech as TechHullComponent);
	$: tech && 'idealSpeed' in tech && (engine = tech as TechEngine);
</script>

{#if tech}
	<div
		class="card bg-base-200 shadow-xl w-[27rem] max-h-fit min-h-fit rounded-sm border-2 border-base-300"
	>
		<div class="card-body p-3 gap-0">
			<h2 class="text-lg font-semibold text-center mb-1 text-secondary">{tech.name}</h2>

			<div class="flex flex-row gap-2">
				<div class="flex flex-col flex-none min-w-[6rem]">
					<!-- icon and tech requirements row-->
					<div
						class="avatar border tech-avatar {kebabCase(
							tech.name.replace("'", '').replace(' ', '').replace('±', '')
						)}"
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
						<div class="border border-base-300 bg-base-100 grow min-h-[14rem]">
							<TechEngineGraph {engine} />
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

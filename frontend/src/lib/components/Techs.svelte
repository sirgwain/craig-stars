<script lang="ts">
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import techjson from '$lib/ssr/techs.json';
	import { TechCategory, type Tech, type TechStore } from '$lib/types/Tech';
	import { kebabCase, startCase } from 'lodash-es';
	import { onMount } from 'svelte';
	import { $enum as eu } from 'ts-enum-util';
	import SectionHeader from './SectionHeader.svelte';
	import TableSearchInput from './TableSearchInput.svelte';
	import { game } from '$lib/services/Stores';
	import { canLearnTech, hasRequiredLevels } from '$lib/types/Player';

	// for ssr, we start with techs from a json file
	export let techStore: TechStore = techjson as TechStore;
	export let techs: Tech[] = [
		...techStore.engines,
		...techStore.planetaryScanners,
		...techStore.defenses,
		...techStore.hullComponents,
		...techStore.hulls,
		...techStore.terraforms
	];

	let filter = '';
	let showAll = !$game?.player ?? false;

	let techsByCategory: Record<TechCategory, Tech[]> = {
		Armor: [],
		BeamWeapon: [],
		Bomb: [],
		Electrical: [],
		Engine: [],
		Mechanical: [],
		MineLayer: [],
		MineRobot: [],
		Orbital: [],
		PlanetaryScanner: [],
		PlanetaryDefense: [],
		Scanner: [],
		Shield: [],
		ShipHull: [],
		StarbaseHull: [],
		Terraforming: [],
		Torpedo: []
	};

	function clearTechsByCategory() {
		techsByCategory = {
			Armor: [],
			BeamWeapon: [],
			Bomb: [],
			Electrical: [],
			Engine: [],
			Mechanical: [],
			MineLayer: [],
			MineRobot: [],
			Orbital: [],
			PlanetaryScanner: [],
			PlanetaryDefense: [],
			Scanner: [],
			Shield: [],
			ShipHull: [],
			StarbaseHull: [],
			Terraforming: [],
			Torpedo: []
		};
	}

	$: filteredTechs = techs.filter(
		(t) =>
			t.name.toLocaleLowerCase().indexOf(filter.toLocaleLowerCase()) != -1 ||
			t.category.toLocaleLowerCase().indexOf(filter.toLocaleLowerCase()) != -1
	);

	$: {
		clearTechsByCategory();
		filteredTechs.forEach((tech) => {
			techsByCategory[tech.category].push(tech);
		});
	}

	onMount(async () => {
		const response = await fetch(`/api/techs`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			techStore = (await response.json()) as TechStore;
			techs = [];
			techs = techs.concat(techStore.engines);
			techs = techs.concat(techStore.planetaryScanners);
			techs = techs.concat(techStore.defenses);
			techs = techs.concat(techStore.hullComponents);
			techs = techs.concat(techStore.hulls);
			techs = techs.concat(techStore.terraforms);
		} else {
			console.error(response);
		}
	});
</script>

<div class="flex justify-between">
	<div><TableSearchInput bind:value={filter} /></div>
	<div class="form-control" class:hidden={!$game?.player}>
		<label class="label cursor-pointer">
			<span class="label-text mr-1">Show All</span>
			<input type="checkbox" class="toggle" bind:checked={showAll} />
		</label>
	</div>
</div>
{#each eu(TechCategory).getKeys() as category}
	{#if techsByCategory[category].length > 0}
		<a id={kebabCase(category)} href={`#${kebabCase(category)}`}
			><SectionHeader title={startCase(category)} /></a
		>
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
			{#each techsByCategory[category] as tech (tech.name)}
				{#if showAll || ($game?.player && canLearnTech($game?.player, tech) && hasRequiredLevels($game?.player.techLevels, tech.requirements))}
					<div>
						<TechSummary {tech} />
					</div>
				{/if}
			{/each}
		</div>
	{/if}
{/each}

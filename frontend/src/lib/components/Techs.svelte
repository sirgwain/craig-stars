<script lang="ts">
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import techjson from '$lib/ssr/techs.json';
	import { type Tech, type TechCategory, type TechStore } from '$lib/types/cs';
	import { CommandedPlayer, canLearnTech } from '$lib/types/Player';
	import { TechCategories } from '$lib/types/Tech';
	import { hasRequiredLevels, levelsAbove } from '$lib/types/TechLevel';
	import type { CS } from '$lib/wasm';
	import { kebabCase, sortBy, startCase } from 'lodash-es';
	import { onMount } from 'svelte';
	import ItemTitle from './ItemTitle.svelte';
	import SectionHeader from './SectionHeader.svelte';

	type Props = {
		player?: CommandedPlayer | undefined;
		cs?: CS | undefined;
	};

	let { player, cs }: Props = $props();

	let filter = $state('');
	let showAll = $state(player === undefined);

	// for ssr, we start with techs from a json file
	let techStore = $state(techjson as unknown as TechStore);
	let techs = $derived([
		...techStore.engines,
		...techStore.planetaryScanners,
		...techStore.defenses,
		...techStore.planetaries,
		...techStore.hullComponents,
		...techStore.hulls,
		...techStore.terraforms
	]);

	let filteredTechs = $derived(
		techs.filter(
			(t) =>
				t.name.toLocaleLowerCase().indexOf(filter.toLocaleLowerCase()) != -1 ||
				t.category.toLocaleLowerCase().indexOf(filter.toLocaleLowerCase()) != -1
		)
	);

	let techsByCategory: Record<TechCategory, Tech[]> = $derived.by(() => {
		const techsByCategory: Record<TechCategory, Tech[]> = {
			Armor: [],
			BeamWeapon: [],
			Bomb: [],
			Electrical: [],
			Engine: [],
			Mechanical: [],
			MineLayer: [],
			MineRobot: [],
			Orbital: [],
			Planetary: [],
			PlanetaryScanner: [],
			PlanetaryDefense: [],
			Scanner: [],
			Shield: [],
			ShipHull: [],
			StarbaseHull: [],
			Terraforming: [],
			Torpedo: []
		};
		filteredTechs.forEach((tech) => {
			techsByCategory[tech.category].push(tech);
		});
		return techsByCategory;
	});

	onMount(async () => {
		const response = await fetch(`/api/techs`, {
			method: 'GET',
			headers: {
				accept: 'application/json'
			}
		});

		if (response.ok) {
			techStore = (await response.json()) as TechStore;
		} else {
			console.error(response);
		}
	});

	let newTechs = $derived(
		player &&
			techs.filter((t) => player?.hasTech(t) && levelsAbove(t.requirements, player.techLevels) == 0)
	);

</script>

<div class="flex justify-between">
	<div><TableSearchInput bind:value={filter} /></div>
	<div class="form-control" class:hidden={!player}>
		<label class="label cursor-pointer">
			<span class="label-text mr-1">Show All</span>
			<input type="checkbox" class="toggle" class:toggle-accent={showAll} bind:checked={showAll} />
		</label>
	</div>
</div>

{#if player && newTechs && filter === ''}
	<ItemTitle>Recently Learned Techs</ItemTitle>
	<div class="grid grid-cols-1 md:grid-cols-2 gap-2">
		{#each newTechs as tech (tech.name)}
			<div class="mx-3">
				<TechSummary {tech} {player} />
			</div>
		{/each}
	</div>
{/if}

{#each TechCategories as category (category)}
	{#if techsByCategory[category]?.length > 0}
		<a id={kebabCase(category)} href={`#${kebabCase(category)}`}
			><SectionHeader title={startCase(category)} /></a
		>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-2">
			{#each sortBy(techsByCategory[category], 'rank') as tech (tech.name)}
				{#if showAll || (player && canLearnTech(player, tech) && hasRequiredLevels(player.techLevels, tech.requirements))}
					<div class="mx-3" data-type="tech-card" data-name={tech.name}>
						<!-- Hide the graph on safari until svelte5 -->
						<TechSummary
							{tech}
							{player}
							{cs}
							showResearchCost={!!(
								player &&
								cs &&
								canLearnTech(player, tech) &&
								!hasRequiredLevels(player.techLevels, tech.requirements)
							)}
						/>
					</div>
				{/if}
			{/each}
		</div>
	{/if}
{/each}

<script lang="ts">
	import TableSearchInput from '$lib/components/table/TableSearchInput.svelte';
	import TechSummary from '$lib/components/tech/TechSummary.svelte';
	import { techClient } from '$lib/services/connect';
	import techjson from '$lib/ssr/techs.json';
	import {
		GetTechsResponseSchema,
		TechCategory,
		type GetTechsResponseJson
	} from '$lib/types/cs-proto';
	import { enumToString } from '$lib/types/Enums';
	import { CommandedPlayer, canLearnTech } from '$lib/types/Player';
	import { TechCategories, type TechLike } from '$lib/types/Tech';
	import { hasRequiredLevels, levelsAbove } from '$lib/types/TechLevel';
	import type { CS } from '$lib/wasm';
	import { fromJson } from '@bufbuild/protobuf';
	import { kebabCase, sortBy } from 'lodash-es';
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
	let techStore = $state(
		fromJson(GetTechsResponseSchema, techjson as unknown as GetTechsResponseJson)
	);
	let techs = $derived([
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
				t.tech?.name?.toLocaleLowerCase().indexOf(filter.toLocaleLowerCase()) != -1 ||
				enumToString(TechCategory, t.tech?.category)
					.toLocaleLowerCase()
					.indexOf(filter.toLocaleLowerCase()) != -1
		)
	);

	let techsByCategory: Record<TechCategory, TechLike[]> = $derived.by(() => {
		const techsByCategory: Record<TechCategory, TechLike[]> = {
			[TechCategory.UNSPECIFIED]: [],
			[TechCategory.ARMOR]: [],
			[TechCategory.BEAM_WEAPON]: [],
			[TechCategory.BOMB]: [],
			[TechCategory.ELECTRICAL]: [],
			[TechCategory.ENGINE]: [],
			[TechCategory.MECHANICAL]: [],
			[TechCategory.MINE_LAYER]: [],
			[TechCategory.MINE_ROBOT]: [],
			[TechCategory.ORBITAL]: [],
			[TechCategory.PLANETARY]: [],
			[TechCategory.PLANETARY_SCANNER]: [],
			[TechCategory.PLANETARY_DEFENSE]: [],
			[TechCategory.SCANNER]: [],
			[TechCategory.SHIELD]: [],
			[TechCategory.SHIP_HULL]: [],
			[TechCategory.STARBASE_HULL]: [],
			[TechCategory.TERRAFORMING]: [],
			[TechCategory.TORPEDO]: []
		};
		filteredTechs.forEach((t) => {
			techsByCategory[t.tech?.category ?? TechCategory.UNSPECIFIED]?.push(t);
		});
		return techsByCategory;
	});

	onMount(async () => {
		const resp = await techClient.getTechs({});
		techStore = { ...resp };
	});

	let newTechs = $derived(
		player &&
			techs.filter(
				(t) =>
					player?.hasTech(t) && levelsAbove(t.tech?.requirements?.techLevel, player.techLevels) == 0
			)
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
		{#each newTechs as tech (tech.tech?.name)}
			<div class="mx-3">
				<TechSummary {tech} {player} />
			</div>
		{/each}
	</div>
{/if}

{#each TechCategories as category (category)}
	{#if techsByCategory[category]?.length > 0}
		<a
			id={kebabCase(enumToString(TechCategory, category))}
			href={`#${kebabCase(enumToString(TechCategory, category))}`}
			><SectionHeader title={enumToString(TechCategory, category)} /></a
		>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-2">
			{#each sortBy(techsByCategory[category], 'rank') as tech (tech.tech?.name)}
				{#if showAll || (player && canLearnTech(player, tech) && hasRequiredLevels(player.techLevels, tech.tech?.requirements?.techLevel))}
					<div class="mx-3" data-type="tech-card" data-name={tech.tech?.name}>
						<!-- Hide the graph on safari until svelte5 -->
						<TechSummary
							{tech}
							{player}
							{cs}
							showResearchCost={!!(
								player &&
								cs &&
								canLearnTech(player, tech) &&
								!hasRequiredLevels(player.techLevels, tech.tech?.requirements?.techLevel)
							)}
						/>
					</div>
				{/if}
			{/each}
		</div>
	{/if}
{/each}

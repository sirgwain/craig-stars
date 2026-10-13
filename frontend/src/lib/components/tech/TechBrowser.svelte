<script lang="ts">
	import TableSearchInput from '../table/TableSearchInput.svelte';
	import { TechCategory } from '#lib/types/cs-proto.js';
	import { TechCategories, type TechLike } from '#lib/types/Tech.js';
	import { enumToString } from '#lib/types/Enums.js';
	import { sortBy } from 'lodash-es';
	import TechBrowseCard from './TechBrowseCard.svelte';
	let { techs }: { techs: TechLike[] } = $props();
	let search = $state('');
	let selected = $state('all');
	let sort = $state('rank');
	let categories = $derived(
		TechCategories.filter((c) => techs.some((t) => t.tech?.category === c))
	);
	let filtered = $derived(
		techs.filter(
			(t) =>
				t.tech &&
				(selected === 'all' || t.tech.category === Number(selected)) &&
				`${t.tech.name} ${enumToString(TechCategory, t.tech.category)}`
					.toLowerCase()
					.includes(search.toLowerCase())
		)
	);
	let visibleCategories = $derived(
		categories.filter((c) => filtered.some((t) => t.tech?.category === c))
	);
</script>

<div class="space-y-4">
	<div class="flex flex-col sm:flex-row gap-2">
		<div class="grow">
			<TableSearchInput
				bind:value={search}
				aria-label="Search technologies"
				placeholder="Search technologies…"
				class="input input-sm w-full"
			/>
		</div>
		<select class="select select-sm" bind:value={sort} aria-label="Sort technologies"
			><option value="rank">Technology order</option><option value="name">Name A–Z</option></select
		>
	</div>
	<div class="lg:hidden">
		<select class="select w-full" aria-label="Technology category" bind:value={selected}
			><option value="all">All technologies</option>{#each categories as category (category)}<option
					value={String(category)}>{enumToString(TechCategory, category)}</option
				>{/each}</select
		>
	</div>
	<div class="grid grid-cols-1 lg:grid-cols-[12rem_minmax(0,1fr)] gap-4 items-start">
		<aside
			class="hidden lg:block sticky top-20 bg-base-200 border border-base-300 rounded-box p-2 max-h-[calc(100dvh-6rem)] overflow-y-auto"
			aria-label="Technology categories"
		>
			<h2 class="text-secondary font-semibold p-2">Categories</h2>
			<ul class="menu w-full p-0">
				<li>
					<button class:menu-active={selected === 'all'} onclick={() => (selected = 'all')}
						>All technologies</button
					>
				</li>
				{#each categories as category (category)}<li>
						<button
							class:menu-active={selected === String(category)}
							onclick={() => (selected = String(category))}
							>{enumToString(TechCategory, category)}</button
						>
					</li>{/each}
			</ul>
		</aside>
		<div class="min-w-0 space-y-5">
			<p class="text-base-content/70" aria-live="polite">{filtered.length} technologies</p>
			{#each visibleCategories as category (category)}
				<section aria-label={enumToString(TechCategory, category)}>
					<h2 class="text-xl font-semibold text-secondary border-b border-base-300 pb-2 mb-3">
						{enumToString(TechCategory, category)}
					</h2>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						{#each sortBy( filtered.filter((t) => t.tech?.category === category), sort === 'name' ? 'tech.name' : 'tech.ranking' ) as tech (tech.tech?.name)}<div
								class="min-w-0 h-full"
								data-type="tech-card"
								data-name={tech.tech?.name}
							>
								<TechBrowseCard {tech} />
							</div>{/each}
					</div>
				</section>
			{:else}<div class="card bg-base-200 border border-base-300">
					<div class="card-body">
						<h2 class="text-2xl">No matching technologies</h2>
						<p>Try another search or category.</p>
					</div>
				</div>{/each}
		</div>
	</div>
</div>

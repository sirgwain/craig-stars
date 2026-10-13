<script lang="ts">
	import { isHull, isHullComponent, type TechLike } from '#lib/types/Tech.js';
	import {
		TechCategory,
		type TechDefense,
		type TechHull,
		type TechHullComponent
	} from '#lib/types/cs-proto.js';
	import { kebabCase } from 'lodash-es';
	import TechAvatar from './TechAvatar.svelte';
	import TechLevelRequirements from './TechLevelRequirements.svelte';
	import TechDefenseGraph from './TechDefenseGraph.svelte';
	import TechDescription from './TechDescription.svelte';
	import TechEngineGraph from './TechEngineGraph.svelte';
	import TechTraitRequirements from './TechTraitRequirements.svelte';
	import TechWarnings from './TechWarnings.svelte';
	import Cost from '../game/Cost.svelte';
	let { tech }: { tech: TechLike } = $props();
	let mass = $derived(
		isHull(tech)
			? (tech as TechHull).mass
			: isHullComponent(tech.tech?.category)
				? (tech as TechHullComponent).mass
				: undefined
	);
	let href = $derived(`/techs/${kebabCase(tech.tech?.name.replaceAll("'", '') ?? '')}`);

	// engines and planetary defenses can show their graph in place, it's only rendered on request
	let engine = $derived(
		tech.tech?.category === TechCategory.ENGINE ? (tech as TechHullComponent) : undefined
	);
	let defense = $derived(
		tech.tech?.category === TechCategory.PLANETARY_DEFENSE ? (tech as TechDefense) : undefined
	);
	let graphLabel = $derived(engine ? 'usage' : defense ? 'coverage' : undefined);
	let showGraph = $state(false);

	// specs taller than a few lines are clamped behind a "more…" toggle
	const clampLines = 3;
	// the spec components pad themselves, so allow for that before deciding the specs are long
	const specPadding = 20;
	let specs: HTMLElement | undefined = $state();
	let specsHeight = $state(0);
	let expanded = $state(false);
	let lineHeight = $derived(specs ? parseFloat(getComputedStyle(specs).lineHeight) || 18 : 18);
	let hasSpecs = $derived(specsHeight > specPadding);
	let long = $derived(specsHeight > clampLines * lineHeight + specPadding);
</script>

{#if tech.tech}
	<article class="card bg-base-200 border border-base-300 h-full min-w-0">
		<div class="card-body p-3 gap-3">
			<div class="flex items-start gap-3">
				<div class="shrink-0"><TechAvatar {tech} hullTooltip={true} /></div>
				<div class="min-w-0">
					<h3 class="text-xl text-secondary font-semibold">
						<a class="hover:text-accent" {href}>{tech.tech.name}</a>
					</h3>
					<div class="mt-2"><TechLevelRequirements {tech} /></div>
				</div>
			</div>
			<div class="border-t border-base-300 pt-3">
				<Cost cost={tech.tech.cost} layout="horizontal" />
				{#if mass !== undefined}<dl class="flex gap-3 mt-3">
						<dt class="text-base-content/70">Mass</dt>
						<dd class="tabular-nums">{mass}kT</dd>
					</dl>{/if}
			</div>
			<div class:border-t={hasSpecs} class:pt-2={hasSpecs} class="border-base-300">
				<div
					class="overflow-hidden"
					class:specs-clamped={long && !expanded}
					style:max-height={long && !expanded ? `${clampLines * lineHeight + 4}px` : undefined}
				>
					<div bind:this={specs} bind:clientHeight={specsHeight}>
						<TechDescription {tech} />
						<TechTraitRequirements {tech} />
					</div>
				</div>
				{#if long}<button
						class="cs-link px-1"
						aria-expanded={expanded}
						onclick={() => (expanded = !expanded)}>{expanded ? 'less' : 'more…'}</button
					>{/if}
			</div>
			{#if graphLabel}<div>
					<button
						class="cs-link px-1"
						aria-expanded={showGraph}
						onclick={() => (showGraph = !showGraph)}
						>{showGraph ? `hide ${graphLabel}` : `${graphLabel}…`}</button
					>
					{#if showGraph}<div class="h-56 pb-5" data-type="tech-graph">
							{#if engine}<TechEngineGraph {engine} />{:else if defense}<TechDefenseGraph
									{defense}
								/>{/if}
						</div>{/if}
				</div>{/if}
			{#if 'radiating' in tech && tech.radiating}<TechWarnings {tech} />{/if}
			{#if isHull(tech)}<div class="card-actions justify-end mt-auto">
					<a class="btn btn-primary" {href}>Hull makeup →</a>
				</div>{/if}
		</div>
	</article>
{/if}

<style>
	/* fade the last visible line so a clamped spec doesn't end mid-line */
	.specs-clamped {
		mask-image: linear-gradient(to bottom, black 60%, transparent);
	}
</style>

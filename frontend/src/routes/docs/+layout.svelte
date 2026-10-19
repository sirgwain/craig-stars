<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import type { Snippet } from 'svelte';
	import DarkModeToggler from '#lib/components/DarkModeToggler.svelte';
	import { docPages } from '#lib/docs/navigation.js';

	let { children }: { children: Snippet } = $props();
	const index = $derived(docPages.findIndex((p) => p.href === page.url.pathname));
	const current = $derived(docPages.find((p) => p.href === page.url.pathname));
	const previous = $derived(index > 0 ? docPages[index - 1] : undefined);
	const next = $derived(
		index >= 0 && index < docPages.length - 1 ? docPages[index + 1] : undefined
	);
</script>

<svelte:head>
	<title>{current?.title ?? 'Documentation'} · craig-stars</title>
</svelte:head>

<div class="min-h-dvh bg-base-100 text-base-content">
	<header class="sticky top-0 z-10 border-b border-base-300 bg-base-100">
		<div class="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-2">
			<a class="text-xl font-bold text-accent" href="/">craig-stars</a>
			<a class="link" href="/docs">Player guide</a>
			<div class="flex items-center gap-4">
				<a class="link" href="/">Play</a>
				<DarkModeToggler />
			</div>
		</div>
	</header>

	<div class="mx-auto grid max-w-7xl gap-6 px-4 py-6 md:grid-cols-[12rem_minmax(0,1fr)] lg:gap-10">
		<aside>
			<nav aria-label="Documentation" class="sticky top-20 hidden md:block">
				<p class="mb-2 font-semibold uppercase tracking-wide text-secondary">Player guide</p>
				<ul class="menu w-full p-0">
					{#each docPages as item (item.href)}
						<li>
							<a
								href={item.href}
								aria-current={item === current ? 'page' : undefined}
								class:menu-active={item === current}>{item.title}</a
							>
						</li>
					{/each}
				</ul>
				{#if current?.sections.length}
					<p class="mb-2 mt-6 font-semibold text-secondary">On this page</p>
					<ul class="menu w-full p-0">
						{#each current.sections as [id, title] (id)}
							<li><a href={`#${id}`}>{title}</a></li>
						{/each}
					</ul>
				{/if}
			</nav>
			<div class="md:hidden">
				<label class="mb-2 block font-semibold" for="doc-topic">Player guide</label>
				<select
					id="doc-topic"
					class="select w-full"
					value={page.url.pathname}
					onchange={(event) => goto(event.currentTarget.value)}
				>
					{#each docPages as item (item.href)}
						<option value={item.href}>{item.title}</option>
					{/each}
				</select>
				{#if current?.sections.length}
					<details class="mt-3 rounded-box border border-base-300 p-3">
						<summary class="cursor-pointer font-semibold">On this page</summary>
						<ul class="menu p-0 pt-2">
							{#each current.sections as [id, title] (id)}
								<li><a href={`#${id}`}>{title}</a></li>
							{/each}
						</ul>
					</details>
				{/if}
			</div>
		</aside>

		<main class="min-w-0 max-w-4xl">
			<nav aria-label="Breadcrumb" class="breadcrumbs mb-4 text-sm text-secondary">
				<ul>
					<li><a href="/docs">Documentation</a></li>
					{#if current && index > 0}<li>{current.title}</li>{/if}
				</ul>
			</nav>
			<article class="docs-article text-lg/relaxed">
				{@render children()}
			</article>
			<nav
				aria-label="Guide pages"
				class="mt-8 flex justify-between gap-4 border-t border-base-300 pt-4"
			>
				<div>
					{#if previous}<a class="link" href={previous.href}>← {previous.title}</a>{/if}
				</div>
				{#if next}<a class="link" href={next.href}>{next.title} →</a>{/if}
			</nav>
		</main>
	</div>
</div>

<style>
	:global(.docs-article > h1) {
		font-size: 2rem;
		font-weight: 700;
		line-height: 1.2;
		margin-bottom: 1rem;
	}
	:global(.docs-article > h2) {
		font-size: 1.4rem;
		font-weight: 600;
		margin-top: 2rem;
		margin-bottom: 0.75rem;
		scroll-margin-top: 6rem;
	}
	:global(.docs-article > h3) {
		font-size: 1.1rem;
		font-weight: 600;
		margin-top: 1.5rem;
		margin-bottom: 0.5rem;
		scroll-margin-top: 6rem;
	}
	:global(.docs-article > p),
	:global(.docs-article > ul),
	:global(.docs-article > ol) {
		margin-bottom: 1rem;
	}
	:global(.docs-article > ul) {
		list-style: disc;
		padding-left: 1.5rem;
	}
	:global(.docs-article > ol) {
		list-style: decimal;
		padding-left: 1.5rem;
	}
	:global(.docs-article > :is(ul, ol) > li) {
		margin-bottom: 0.5rem;
	}
	:global(.docs-article > p a),
	:global(.docs-article > :is(ul, ol) a) {
		text-decoration: underline;
		text-underline-offset: 3px;
		color: var(--color-accent);
	}
	:global(.docs-article > table) {
		display: block;
		overflow-x: auto;
		max-width: 100%;
		border-collapse: collapse;
		margin-bottom: 1rem;
	}
	:global(.docs-article > table :is(th, td)) {
		border-bottom: 1px solid var(--color-base-300);
		padding: 0.75rem;
		text-align: left;
		vertical-align: top;
	}
	:global(.docs-article > pre) {
		overflow-x: auto;
		background: var(--color-base-200);
		padding: 1rem;
		margin-bottom: 1rem;
	}
</style>

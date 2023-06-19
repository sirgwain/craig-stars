<script lang="ts">
	import { page } from '$app/stores';
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
	import { game } from '$lib/services/Context';
	import { FullGame } from '$lib/services/FullGame';
	import { getContext, onMount } from 'svelte';
	import TestBreadcrumb from './TestBreadcrumb.svelte';

	onMount(() => {
		// setup the quantityModifier
		bindQuantityModifier();

		return () => {
			unbindQuantityModifier();
		};
	});

	let title = getContext<string>('title');

	$: title = getContext('title') ?? $page.routeId?.replace('tests/', '') ?? '';

	// create a test player
	game.update(() => new FullGame());
</script>

<TestBreadcrumb {title} />

<slot />

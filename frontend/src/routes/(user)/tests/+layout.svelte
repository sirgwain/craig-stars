<script lang="ts">
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
	import { page } from '$app/stores';
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
</script>

<TestBreadcrumb {title} />

<slot />

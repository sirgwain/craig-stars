<script lang="ts">
	import { page } from '$app/stores';
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
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

	$: title = getContext('title') ?? $page.route.id?.replace('tests/', '') ?? '';
</script>

<TestBreadcrumb {title} />

<slot />

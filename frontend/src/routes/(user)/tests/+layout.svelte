<script lang="ts">
	import { bindQuantityModifier, unbindQuantityModifier } from '$lib/quantityModifier';
	import { page } from '$app/stores';
	import { getContext, onMount } from 'svelte';
	import TestBreadcrumb from './TestBreadcrumb.svelte';
	import { humanoid } from '$lib/types/Race';
	import { player } from '$lib/services/Context';

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
	player.update(() => {
		return {
			gameId: 0,
			userId: 0,
			num: 0,
			race: humanoid,
			techLevels: {
				energy: 3,
				weapons: 3,
				propulsion: 3,
				construction: 3,
				electronics: 3,
				biotechnology: 3
			},
			techLevelsSpent: {},
			messages: [],
			designs: [],
			planets: [],
			fleets: [],
			planetIntels: [],
			fleetIntels: [],
			color: '#0000FF'
		};
	});
</script>

<TestBreadcrumb {title} />

<slot />

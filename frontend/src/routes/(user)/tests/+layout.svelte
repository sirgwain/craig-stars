<script lang="ts">
	import { page } from '$app/stores';
	import { getContext } from 'svelte';
	import TestBreadcrumb from './TestBreadcrumb.svelte';

	let title = getContext<string>('title');

	$: title = getContext('title') ?? $page.route.id?.replace('tests/', '') ?? '';
</script>

<TestBreadcrumb {title} />

<slot />

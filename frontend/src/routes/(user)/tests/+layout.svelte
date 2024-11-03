<script lang="ts">
	import { run } from 'svelte/legacy';

	import { page } from '$app/stores';
	import { getContext } from 'svelte';
	import TestBreadcrumb from './TestBreadcrumb.svelte';
	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();

	let title = $state(getContext<string>('title'));

	run(() => {
		title = getContext('title') ?? $page.route.id?.replace('tests/', '') ?? '';
	});
</script>

<TestBreadcrumb {title} />

{@render children?.()}

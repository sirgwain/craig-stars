<script lang="ts">
	import { page } from '$app/state';
	import { getContext, type Snippet } from 'svelte';
	import TestBreadcrumb from './TestBreadcrumb.svelte';
	type Props = {
		children?: Snippet;
	};

	let { children }: Props = $props();

	let title = $derived(getContext<string>('title'));

	$effect(() => {
		title = getContext('title') || page.route.id?.replace('tests/', '') || '';
	});
</script>

<TestBreadcrumb {title} />

{@render children?.()}

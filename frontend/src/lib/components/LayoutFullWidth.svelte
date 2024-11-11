<script lang="ts">
	import ErrorToast from '$lib/components/ErrorToast.svelte';
	import Menu from '$lib/components/Menu.svelte';
	import NotificationToast from '$lib/components/NotificationToast.svelte';
	import { errors } from '$lib/services/Errors';
	import { me } from '$lib/services/Stores';
	import { onMount } from 'svelte';
	type Props = {
		children?: import('svelte').Snippet;
	};

	let { children }: Props = $props();

	onMount(() => {
		$errors = [];
	});
</script>

<main class="p-3 flex flex-col">
	<div class="sticky top-0 z-10">
		<Menu user={$me} />
	</div>
	<div class="w-full mx-auto">
		{#if children}{@render children()}{:else}This is the main content{/if}
	</div>
	<ErrorToast />
	<NotificationToast />
</main>

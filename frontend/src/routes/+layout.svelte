<script lang="ts">
	import { page, updated } from '$app/state';
	import { authGuard } from '#lib/authGuard.js';
	import HomePage from '#lib/components/HomePage.svelte';
	import VersionToast from '#lib/components/VersionToast.svelte';
	import { me } from '#lib/services/Stores.js';
	import type { Snippet } from 'svelte';
	import { onMount } from 'svelte';

	import '../css/app.css';
	import '../css/hulls.css';
	import '../css/mapobjects.css';
	import '../css/planets.css';
	import '../css/techs.css';
	import { UserStatuses } from '#lib/types/User.js';

	type Props = { children?: Snippet };

	let { children }: Props = $props();
	const publicRoute = $derived(
		page.url.pathname.startsWith('/auth') || page.url.pathname.startsWith('/docs')
	);
	const wasmExecUrl = new URL('#lib/wasm/wasm_exec.js', import.meta.url).href;

	// check the user
	onMount(() => {
		if (!publicRoute) {
			authGuard();
		}
	});
</script>

<svelte:head>
	<script src={wasmExecUrl}></script>
</svelte:head>

<!-- Show the main content if we've logged in, otherwise show the login page -->
{#if $me.status == UserStatuses.LoggedIn || publicRoute}
	{#if children}{@render children()}{:else}This is the main content{/if}
{:else if $me.status == UserStatuses.NotFound}
	<HomePage />
{/if}

<VersionToast frontendUpdated={updated.current} />

<script lang="ts">
	import { page } from '$app/stores';
	import { authGuard } from '$lib/authGuard';
	import HomePage from '$lib/components/HomePage.svelte';
	import { me } from '$lib/services/Stores';
	import { UserStatus } from '$lib/types/User';
	import { onMount } from 'svelte';

	import '../css/app.css';
	import '../css/hulls.css';
	import '../css/mapobjects.css';
	import '../css/planets.css';
	import '../css/techs.css';

	const loggingIn = $page.url.pathname.startsWith('/auth');
	const wasmExecUrl = new URL('$lib/wasm/wasm_exec.js', import.meta.url).href;

	// check the user
	onMount(() => {
		if (!loggingIn) {
			authGuard();
		}
	});
</script>

<svelte:head>
	<script src={wasmExecUrl}></script>
</svelte:head>

<!-- Show the main content if we've logged in, otherwise show the login page -->
{#if $me.status == UserStatus.LoggedIn || loggingIn}
	<slot>This is the main content</slot>
{:else if $me.status == UserStatus.NotFound}
	<HomePage />
{/if}

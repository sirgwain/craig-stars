<script lang="ts">
	import { authGuard } from '$lib/authGuard';
	import type { User } from '$lib/types/User';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import '../css/app.css';
	import '../css/planets.css';
	import '../css/techs.css';
	import '../css/hulls.css';
	import Login from '$lib/components/Login.svelte';

	let mounted = false;
	let showLogin = false;
	// verify the user, redirect otherwise
	onMount(async () => {
		mounted = true;
		if (!$page.routeId?.startsWith('techs')) {
			user = await authGuard();
			if (!user) {
				showLogin = true;
			}
		} else {
		}
	});

	$: {
		if (mounted && !user && !$page.routeId?.startsWith('techs')) {
			(async () => (user = await authGuard()))();
		}
	}

	let user: User | undefined;
</script>

<!-- Show the main content if we've logged in, otherwise show the login page -->
{#if user}
	<slot>This is the main content</slot>
{:else}
	<Login />
{/if}

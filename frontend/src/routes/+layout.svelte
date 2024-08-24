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

	const loggingIn = $page.url.pathname.startsWith('/auth')

	// check the user
	onMount(() => {
		if (!loggingIn) {
			authGuard();
		}
	});
</script>

<!-- Show the main content if we've logged in, otherwise show the login page -->
{#if $me.status == UserStatus.LoggedIn || loggingIn}
	<slot>This is the main content</slot>
{:else if $me.status == UserStatus.NotFound}
	<HomePage />
{/if}

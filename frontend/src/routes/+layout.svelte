<script lang="ts">
	import { authGuard } from '$lib/authGuard';
	import Login from '$lib/components/auth/Login.svelte';
	import { me } from '$lib/services/Context';
	import { UserStatus } from '$lib/types/User';
	import { onMount } from 'svelte';
	import '../css/app.css';
	import '../css/hulls.css';
	import '../css/planets.css';
	import '../css/mapobjects.css';
	import '../css/techs.css';
	import { page } from '$app/stores';

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
	<Login />
{/if}

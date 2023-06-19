<script lang="ts">
	import { authGuard } from '$lib/authGuard';
	import { UserStatus, type User } from '$lib/types/User';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { me } from '$lib/services/Context';
	import '../css/app.css';
	import '../css/planets.css';
	import '../css/techs.css';
	import '../css/hulls.css';
	import Login from '$lib/components/Login.svelte';

	// check the user
	onMount(() => {
		authGuard();
	});
</script>

<!-- Show the main content if we've logged in, otherwise show the login page -->
{#if $me.status == UserStatus.LoggedIn}
	<slot>This is the main content</slot>
{:else if $me.status == UserStatus.NotFound}
	<Login />
{/if}

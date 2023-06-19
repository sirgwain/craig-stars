<script lang="ts">
	import { authGuard } from '$lib/authGuard';
	import Menu from '$lib/components/Menu.svelte';
	import type { User } from '$lib/types/User';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import '../css/app.css';
	import '../css/planets.css';
	import '../css/techs.css';
	import '../css/hulls.css';

	// verify the user, redirect otherwise
	onMount(async () => {
		if (!$page.routeId?.startsWith('techs')) {
			user = await authGuard();
		}
	});

	$: {
		if (!user && !$page.routeId?.startsWith('techs')) {
			(async () => (user = await authGuard()))();
		}
	}

	let user: User | undefined;
</script>

<main class="p-3 flex flex-col h-screen">
	<div class="flex-initial">
		<Menu {user} />
	</div>
	<div class="flex-1 h-full">
		<slot>This is the main content</slot>
	</div>
</main>

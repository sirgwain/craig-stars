<script lang="ts">
	import { authGuard } from '$lib/authGuard';
	import Menu from '$lib/components/Menu.svelte';
	import type { User } from '$lib/types/User';
	import { onMount } from 'svelte';
	import '../app.css';

	// verify the user, redirect otherwise
	onMount(async () => (user = await authGuard()));

	let user: User | undefined;
</script>

<main class="p-5 flex flex-col h-screen">
	<div class="flex-none">
		{#if user}
			<Menu {user} />
		{/if}
	</div>
	<div class="flex-1 h-full">
		<slot>This is the main content</slot>
	</div>
</main>

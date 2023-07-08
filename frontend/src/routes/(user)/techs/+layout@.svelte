<script lang="ts">
	import Menu from '$lib/components/Menu.svelte';
	import { me } from '$lib/services/Stores';
</script>

<main class="p-3 flex flex-col">
	<div class="sticky top-0 z-10">
		<Menu user={$me} />
	</div>
	<div class="relative">
		<div class="w-full mx-auto md:max-w-3xl">
			<slot />
		</div>
	</div>
</main>

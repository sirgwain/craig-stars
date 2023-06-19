<script lang="ts">
	import { me } from '$lib/services/Context';
	import Menu from './Menu.svelte';

	export let title = 'Resource not found';
</script>

<main class="p-3 flex flex-col">
	<div class="sticky top-0 z-10">
		<Menu user={$me} />
	</div>
	<div class="w-full mx-auto md:max-w-2xl">
		{title}
	</div>
</main>

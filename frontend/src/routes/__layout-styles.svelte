<script lang="ts">
	import DarkModeToggler from '$lib/components/DarkModeToggler.svelte';
</script>

<main class="m-5">
	<div class="flex justify-end">
		<DarkModeToggler />
	</div>
	<slot />
</main>

<script lang="ts">
	import { page } from '$app/stores';
	import { designs } from '$lib/services/Context';
	import { DesignService } from '$lib/services/DesignService';
	import { onMount } from 'svelte';

	let gameId = parseInt($page.params.id);

	onMount(async () => {
		try {
			DesignService.load(gameId).then((items) => designs.update(() => items));
		} catch (error) {
			// TODO: handle error
		}
	});
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<slot />
</div>

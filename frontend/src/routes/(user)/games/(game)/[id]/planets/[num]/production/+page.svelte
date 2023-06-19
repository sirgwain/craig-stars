<script lang="ts">
	import { page } from '$app/stores';
	import { player } from '$lib/services/Context';

	let num = parseInt($page.params.num);

	$: planet = $player?.planets.find((p) => p.num === num);
</script>

{#if planet}
	{planet.name} Production
{:else}
	Planet not found
{/if}

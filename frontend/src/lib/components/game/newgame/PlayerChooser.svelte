<script lang="ts">
	import ColorInput from '$lib/components/ColorInput.svelte';
	import { RaceService } from '$lib/services/RaceService';
	import type { Race } from '$lib/types/Race';
	import { onMount } from 'svelte';

	// races for the host
	let races: Race[];
	export let raceId: number;

	onMount(async () => {
		races = await RaceService.load();
	});
</script>

{#if races}
	<label class="label" for="hostRace">Race</label>
	<select name="hostRaceId" class="select select-bordered" bind:value={raceId}>
		{#each races as race}
			<option value={race.id}>{race.name}</option>
		{/each}
	</select>
{/if}

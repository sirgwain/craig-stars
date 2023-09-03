<script lang="ts">
	import { RaceService } from '$lib/services/RaceService';
	import { humanoid, type Race } from '$lib/types/Race';
	import { onMount } from 'svelte';

	// races for the host
	let races: Race[] = [humanoid];
	export let race: Race;

	onMount(async () => {
		const userRaces = await RaceService.load();
		if (userRaces.length > 0) {
			races = userRaces;
		}
	});

	function raceChanged(id: number) {
		const newRace = races.find((r) => r.id == id);
		if (newRace) {
			race = newRace;
		}
	}
</script>

{#if races}
	<label class="label" for="hostRace">Race</label>
	<select
		class="select select-bordered"
		on:change={(e) => raceChanged(parseInt(e.currentTarget.value))}
	>
		{#each races as race}
			<option value={race.id}>{race.name}</option>
		{/each}
	</select>
{/if}

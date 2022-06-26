<script lang="ts">
	import { RaceService } from '$lib/services/RaceService';
	import type { NewGamePlayer } from '$lib/types/Game';
	import type { Race } from '$lib/types/Race';
	import { onMount } from 'svelte';

	// races for the host
	let hostRaces: Race[];
	export let player: NewGamePlayer;

	const raceService = new RaceService();

	onMount(async () => {
		hostRaces = await raceService.loadRaces();
	});
</script>

{#if hostRaces}
	<label class="label" for="hostRace">Host Race</label>
	<select name="hostRaceId" class="select select-bordered" bind:value={player.raceId}>
		{#each hostRaces as race}
			<option value={race.id}>{race.name}</option>
		{/each}
	</select>
{/if}

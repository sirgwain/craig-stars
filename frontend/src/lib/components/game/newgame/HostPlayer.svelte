<script lang="ts">
	import ColorInput from '$lib/components/ColorInput.svelte';
	import Select from '$lib/components/Select.svelte';
	import { RaceService } from '$lib/services/RaceService';
	import type { NewGamePlayer } from '$lib/types/Game';
	import { humanoid, type Race } from '$lib/types/Race';
	import { onMount } from 'svelte';

	// races for the host
	let hostRaces: Race[] = [humanoid];

	export let player: NewGamePlayer;

	onMount(async () => {
		player.race = hostRaces[0];
		const races = await RaceService.load();
		if (races.length > 0) {
			hostRaces = races;
			player.race = hostRaces[0];
		}
	});

	function raceChanged(id: number) {
		const newRace = hostRaces.find((r) => r.id == id);
		if (newRace) {
			player.race = newRace;
		}
	}
</script>

{#if hostRaces}
	<Select
		values={hostRaces.map((r) => {
			return { value: r.id, title: r.pluralName };
		})}
		name="Host"
		value={player.race?.id ?? 0}
		on:change={(e) => raceChanged(e.detail)}
	/>

	<!-- <ColorInput bind:value={player.color} name="color" /> -->
{/if}

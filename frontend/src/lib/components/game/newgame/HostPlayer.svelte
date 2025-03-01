<script lang="ts">
	import Select from '$lib/components/Select.svelte';
	import { RaceService } from '$lib/services/RaceService';
	import type { NewGamePlayer } from '$lib/types/cs';
	import { humanoid } from '$lib/types/Race';
	import { type Race } from '$lib/types/cs';
	import { onMount } from 'svelte';

	// races for the host
	let hostRaces: Race[] = $state([humanoid()]);

	type Props = {
		player: NewGamePlayer;
	};

	let { player = $bindable() }: Props = $props();

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
		onchange={(e) => raceChanged(parseInt(e.currentTarget.value))}
	/>

	<!-- <ColorInput bind:value={player.color} name="color" /> -->
{/if}

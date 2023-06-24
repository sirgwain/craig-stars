<script lang="ts">
	import ColorInput from '$lib/components/ColorInput.svelte';
	import Select from '$lib/components/Select.svelte';
	import { RaceService } from '$lib/services/RaceService';
	import type { NewGamePlayer } from '$lib/types/Game';
	import type { Race } from '$lib/types/Race';
	import { onMount } from 'svelte';

	// races for the host
	let hostRaces: Race[];

	export let player: NewGamePlayer;

	onMount(async () => {
		hostRaces = await RaceService.load();
	});
</script>

{#if hostRaces}
	<Select
		values={hostRaces.map((r) => {
			return { value: r.id, title: r.pluralName };
		})}
		name="race"
		bind:value={player.raceId}
	/>

	<ColorInput bind:value={player.color} name="color" />
{/if}

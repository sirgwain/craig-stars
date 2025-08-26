<script lang="ts">
	import Select from '$lib/components/Select.svelte';
	import { humanoid } from '$lib/types/Race';
	import { onMount } from 'svelte';
	import { raceClient } from '$lib/services/connect';
	import type { Race } from '$lib/types/cs-proto';
	import type { NewGamePlayer } from '$lib/types/cs-proto';

	// races for the host
	let hostRaces: Race[] = $state([humanoid()]);

	type Props = {
		player: NewGamePlayer;
	};

	let { player = $bindable() }: Props = $props();

	onMount(async () => {
		player.race = hostRaces[0];
		const { races } = await raceClient.getRaces({});
		if (races.length > 0) {
			hostRaces = races;
			player.race = hostRaces[0];
		}
	});

	function raceChanged(id: bigint) {
		const newRace = hostRaces.find((r) => r.id === id);
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
		onchange={(e) => raceChanged(BigInt(e.currentTarget.value))}
	/>

	<!-- <ColorInput bind:value={player.color} name="color" /> -->
{/if}

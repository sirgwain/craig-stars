<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import { humanoid } from '$lib/types/Race';
	import { onMount } from 'svelte';
	import RaceEditor from '../../../../routes/(user)/races/[id]/RaceEditor.svelte';
	import RacePoints from '../../../../routes/(user)/races/[id]/RacePoints.svelte';
	import type { Race } from '$lib/types/cs-proto';
	import { raceClient } from '$lib/services/connect';
	import { loadWasm, type CS } from '$lib/wasm';

	type Props = {
		raceUpdated?: (race: Race, valid: boolean) => void;
	};

	let { raceUpdated }: Props = $props();

	// races for the host
	let races: Race[] = $state([]);
	let race = $state(humanoid());
	let cs: CS | undefined = $state();

	onMount(async () => {
		// load wasm for the points calculator
		loadWasm().then((resp) => (cs = resp));

		// load the user's races
		const { races: userRaces } = await raceClient.getRaces({});
		if (userRaces?.length > 0) {
			races = userRaces;
			raceUpdated?.(races[0], true);
		}
	});

	function raceChanged(id: bigint) {
		const newRace = races.find((r) => r.id === id);
		if (newRace) {
			raceUpdated?.(newRace, true);
		}
	}
</script>

{#if races.length > 0}
	<label class="label" for="hostRace">Race</label>
	<select
		class="select select-bordered"
		onchange={(e) => raceChanged(BigInt(e.currentTarget.value))}
	>
		{#each races as race (race.id)}
			<option value={race.id}>{race.name}</option>
		{/each}
	</select>
{:else}
	<ItemTitle>Your Race</ItemTitle>
	{#if cs}
		<RacePoints
			wasmClient={cs.wasmService}
			{race}
			onPointsUpdated={(points) => raceUpdated?.(race, points >= 0)}
		/>
	{/if}
	<RaceEditor bind:race />
{/if}

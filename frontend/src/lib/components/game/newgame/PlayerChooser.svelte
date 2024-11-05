<script lang="ts">
	import { run } from 'svelte/legacy';

	import { RaceService } from '$lib/services/RaceService';
	import { humanoid, type Race } from '$lib/types/Race';
	import { onMount } from 'svelte';
	import RaceEditor from '../../../../routes/(user)/races/[id]/RaceEditor.svelte';
	import RacePoints from '../../../../routes/(user)/races/[id]/RacePoints.svelte';
	import ItemTitle from '$lib/components/ItemTitle.svelte';

	// races for the host
	let races: Race[] = $state([]);
	interface Props {
		race?: Race;
		valid?: boolean;
	}

	let { race = $bindable(humanoid()), valid = $bindable(true) }: Props = $props();

	onMount(async () => {
		const userRaces = await RaceService.load();
		if (userRaces.length > 0) {
			races = userRaces;
			race = races[0];
		}
	});

	function raceChanged(id: number) {
		const newRace = races.find((r) => r.id == id);
		if (newRace) {
			race = newRace;
		}
	}

	let points = $state(0);
	run(() => {
		valid = points >= 0;
	});
</script>

{#if races.length > 0}
	<label class="label" for="hostRace">Race</label>
	<select
		class="select select-bordered"
		onchange={(e) => raceChanged(parseInt(e.currentTarget.value))}
	>
		{#each races as race}
			<option value={race.id}>{race.name}</option>
		{/each}
	</select>
{:else}
	<ItemTitle>Your Race</ItemTitle>
	<RacePoints bind:points {race} />
	<RaceEditor bind:race />
{/if}

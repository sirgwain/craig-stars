<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import { RaceService } from '$lib/services/RaceService';
	import { humanoid } from '$lib/types/Race';
	import { type Race } from '$lib/types/cs';
	import { onMount } from 'svelte';
	import RaceEditor from '../../../../routes/(user)/races/[id]/RaceEditor.svelte';
	import RacePoints from '../../../../routes/(user)/races/[id]/RacePoints.svelte';

	// races for the host
	let races: Race[] = $state([]);
	let race = $state(humanoid());
	type Props = {
		raceUpdated?: (race: Race, valid: boolean) => void;
	};

	let { raceUpdated }: Props = $props();

	onMount(async () => {
		const userRaces = await RaceService.load();
		if (userRaces.length > 0) {
			races = userRaces;
			raceUpdated?.(races[0], true);
		}
	});

	function raceChanged(id: number) {
		const newRace = races.find((r) => r.id == id);
		if (newRace) {
			raceUpdated?.(newRace, true);
		}
	}
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
	<RacePoints {race} onPointsUpdated={(points) => raceUpdated?.(race, points >= 0)} />
	<RaceEditor bind:race />
{/if}

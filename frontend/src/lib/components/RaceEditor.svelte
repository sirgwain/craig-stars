<script lang="ts">
	import { RaceService } from '$lib/services/RaceService';
	import type { Race } from '$lib/types/Race';
	import { onMount } from 'svelte';

	export let id: number;
	let race: Race;
	let raceService: RaceService = new RaceService();

	onMount(async () => {
		// load the race on mount
		race = await raceService.loadRace(id);
	});

	async function onSubmit() {}
</script>

{#if race}
	<div class="prose"><h2>{race.name}</h2></div>
	<form on:submit|preventDefault={onSubmit}>
		<fieldset class="form-control">
			<label class="label" for="name">Name</label>
			<input name="name" class="input input-bordered" bind:value={race.name} />
		</fieldset>

		<button class="btn btn-primary">Save</button>
	</form>
{/if}

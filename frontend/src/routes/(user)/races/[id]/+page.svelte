<script lang="ts">
	import { page } from '$app/stores';

	import { goto } from '$app/navigation';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import { RaceService } from '$lib/services/RaceService';
	import { Service } from '$lib/services/Service';
	import { humanoid, type Race } from '$lib/types/Race';
	import { onMount } from 'svelte';
	import RaceEditor from './RaceEditor.svelte';
	import RacePoints from './RacePoints.svelte';
	import { notify } from '$lib/services/Notifications';

	let id = $page.params.id;
	let race: Race;

	onMount(async () => {
		if (id !== 'new') {
			try {
				race = await RaceService.get(id);
			} catch (err) {
				// TODO: show error
			}
		} else {
			// create a new humanoid
			race = Object.assign({}, humanoid());
		}
	});

	const onSubmit = async () => {
		const body = JSON.stringify(race);
		const create = race?.id ? false : true;
		const response = await fetch(`/api/races${race?.id ? '/' + race.id : ''}`, {
			method: create ? 'POST' : 'PUT',
			headers: {
				accept: 'application/json'
			},
			body
		});

		if (!response.ok) {
			await Service.throwError(response);
		}

		race = (await response.json()) as Race;
		// redirect to page with id
		if (create) {
			await goto(`/races/${race.id}`);
		}

		notify('Saved ' + race.pluralName);
	};

	let points = 0;
</script>

{#if race}
	<form on:submit|preventDefault={onSubmit}>
		<div class="w-full flex justify-end gap-2">
			<button class="btn btn-success" type="submit" disabled={points < 0}>Save</button>
		</div>

		<ItemTitle>{race.name}</ItemTitle>
		<RacePoints bind:points {race} />
		<RaceEditor bind:race />
	</form>
{/if}

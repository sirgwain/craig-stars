<script lang="ts">
	import { page } from '$app/stores';

	import { goto } from '$app/navigation';
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import { humanoid, PRT, type Race } from '$lib/types/Race';
	import { onMount } from 'svelte';

	let id = $page.params.id;
	let race: Race;

	onMount(async () => {
		if (id !== 'new') {
			const response = await fetch(`/api/races/${id}`, {
				method: 'GET',
				headers: {
					accept: 'application/json'
				}
			});

			if (response.ok) {
				race = (await response.json()) as Race;
			} else {
				console.error(response);
			}
		} else {
			// create a new humanoid
			race = Object.assign({}, humanoid);
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

		race = (await response.json()) as Race;
		// redirect to page with id
		if (create) {
			await goto(`/races/${race.id}`);
		}
	};

	$: points = 0;
	$: {
		if (race) {
			computeRacePoints();
		}
	}

	// update points from the server anytime things change
	const computeRacePoints = async () => {
		const body = JSON.stringify(race);
		const response = await fetch(`/api/races/points`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body
		});

		const result = (await response.json()) as { points: number };
		points = result.points;
	};
</script>

{#if race}
	<form on:submit|preventDefault={onSubmit}>
		<div class="w-full flex justify-end gap-2">
			<button class="btn btn-success" type="submit">Save</button>
		</div>

		<ItemTitle>{race.name}</ItemTitle>

		Points: {points}
		<TextInput name="name" bind:value={race.name} />
		<TextInput name="pluralName" bind:value={race.pluralName} />
		<EnumSelect name="prt" enumType={PRT} typeTitle={(type) => type} bind:value={race.prt} />
	</form>
{/if}

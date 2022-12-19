<script lang="ts">
	import { page } from '$app/stores';

	import { goto } from '$app/navigation';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import { onMount } from 'svelte';
	import { setContext } from 'svelte';
	import { humanoid, LRT, PRT, type Race } from '$lib/types/Race';
	import TextInput from '$lib/components/TextInput.svelte';
	import EnumSelect from '$lib/components/EnumSelect.svelte';

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

	// update race calculations
	$: race && true; // (race.spec = computeRaceSpec(race));
</script>

{#if race}
	<form on:submit|preventDefault={onSubmit}>
		<div class="w-full flex justify-end gap-2">
			<button class="btn btn-success" type="submit">Save</button>
		</div>

		<ItemTitle>{race.name}</ItemTitle>

		<TextInput name="name" bind:value={race.name} />
		<TextInput name="pluralName" bind:value={race.pluralName} />
		<EnumSelect name="prt" enumType={PRT} typeTitle={(type) => type} bind:value={race.prt} />
	</form>
{/if}

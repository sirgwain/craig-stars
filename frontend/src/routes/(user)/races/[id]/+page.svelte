<script lang="ts">
	import { page } from '$app/stores';

	import { goto } from '$app/navigation';
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import {
		getLabelForPRT,
		humanoid,
		LRT,
		PRT,
		SpendLeftoverPointsOn,
		type Race
	} from '$lib/types/Race';
	import { onMount } from 'svelte';
	import PRTDescription from './PRTDescription.svelte';
	import LRTs from './LRTs.svelte';
	import PlanetaryProduction from './PlanetaryProduction.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import Habitability from './Habitability.svelte';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { User } from '@steeze-ui/heroicons';
	import Research from './Research.svelte';
	import { RaceService } from '$lib/services/RaceService';

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
			<button class="btn btn-success" type="submit" disabled={points < 0}>Save</button>
		</div>

		<ItemTitle>{race.name}</ItemTitle>

		<div class="sticky top-[9rem] z-10">
			<div class="flex justify-end">
				<div class="stats stats-horizontal shadow border border-base-200">
					<div class="stat place-items-center">
						<div class="stat-title">Points</div>
						<div class="stat-figure"><Icon class="w-8 h-8" src={User} /></div>
						<div class="stat-value" class:text-error={points < 0} class:text-success={points >= 0}>
							{points}
						</div>
						<div class="stat-desc pt-1" />
					</div>
				</div>
			</div>
		</div>

		<TextInput name="name" bind:value={race.name} />
		<TextInput name="pluralName" bind:value={race.pluralName} />
		<EnumSelect
			name="spendLeftoverPointsOn"
			enumType={SpendLeftoverPointsOn}
			typeTitle={(type) => type}
			bind:value={race.spendLeftoverPointsOn}
		/>

		<SectionHeader>Primary Racial Trait</SectionHeader>
		<EnumSelect
			name="prt"
			title="Primary Racial Trait"
			enumType={PRT}
			typeTitle={(prt) => getLabelForPRT(prt)}
			bind:value={race.prt}
		/>
		<div class="card bg-base-200 shadow">
			<div class="card-body">
				<PRTDescription prt={race.prt} />
			</div>
		</div>

		<SectionHeader>Lesser Racial Traits</SectionHeader>
		<LRTs bind:race />

		<SectionHeader>Habitability</SectionHeader>
		<Habitability bind:race />

		<SectionHeader>Planetary Production</SectionHeader>
		<PlanetaryProduction bind:race />

		<SectionHeader>Research</SectionHeader>
		<Research bind:race />
	</form>
{/if}

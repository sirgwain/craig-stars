<script lang="ts">
	import TableSearchInput from '#lib/components/table/TableSearchInput.svelte';
	import { raceClient } from '#lib/services/connect.js';
	import { addError } from '#lib/services/Errors.js';
	import { type Race } from '#lib/types/cs-proto.js';
	import { getAbbreviationForLRT, getLabelForPRT, getLabelForLRT, lrts } from '#lib/types/Race.js';
	import { Grav, Temp, Rad } from '#lib/types/Hab.js';
	import type { ConnectError } from '@connectrpc/connect';
	import { onMount } from 'svelte';
	import HabBar from '../games/(game)/[id]/race/HabBar.svelte';
	let races: Race[] = $state([]);
	let search = $state('');
	let trait = $state('all');
	let sort = $state('name');
	let loading = $state(true);
	let failed = $state(false);
	let filteredRaces = $derived(
		races
			.filter(
				(r) =>
					`${r.name} ${r.pluralName}`.toLowerCase().includes(search.toLowerCase()) &&
					(trait === 'all' || r.prt === Number(trait))
			)
			.toSorted((a, b) =>
				sort === 'created'
					? Number(b.createdAt?.seconds ?? 0) - Number(a.createdAt?.seconds ?? 0)
					: a.pluralName.localeCompare(b.pluralName)
			)
	);
	let traits = $derived(
		[...new Set(races.map((r) => r.prt))].sort((a, b) =>
			getLabelForPRT(a).localeCompare(getLabelForPRT(b))
		)
	);
	async function loadRaces() {
		loading = true;
		failed = false;
		try {
			races = (await raceClient.getRaces({})).races;
		} catch (e) {
			failed = true;
			addError(e as ConnectError);
		} finally {
			loading = false;
		}
	}
	async function removeItem(race: Race) {
		if (confirm(`Are you sure you want to delete ${race.name}?`)) {
			try {
				await raceClient.deleteRace({ raceId: race.id });
				races = races.filter((r) => r.id !== race.id);
			} catch (e) {
				addError(e as ConnectError);
			}
		}
	}
	onMount(loadRaces);
</script>

<svelte:head><title>Races</title></svelte:head>
<div class="py-4 space-y-5">
	<header class="flex items-center justify-between gap-3">
		<div>
			<h1 class="text-3xl font-bold">Races</h1>
			<p class="text-base-content/70 mt-1">Your race designs.</p>
		</div>
		<a href="/races/new" class="btn btn-primary shrink-0">Create race</a>
	</header>
	<div class="flex flex-col sm:flex-row gap-2">
		<div class="grow">
			<TableSearchInput
				bind:value={search}
				aria-label="Search races"
				placeholder="Search races…"
				class="input input-sm w-full"
			/>
		</div>
		<select class="select select-sm min-w-0" aria-label="Filter primary trait" bind:value={trait}
			><option value="all">All primary traits</option>{#each traits as prt (prt)}<option
					value={String(prt)}>{getLabelForPRT(prt)}</option
				>{/each}</select
		>
		<select class="select select-sm" aria-label="Sort races" bind:value={sort}
			><option value="name">Name A–Z</option><option value="created">Newest first</option></select
		>
	</div>
	{#if loading}<div class="flex justify-center p-10">
			<span class="loading loading-spinner" aria-label="Loading races"></span>
		</div>
	{:else if failed}<div class="alert alert-error">
			<span>Unable to load races.</span><button class="btn btn-sm" onclick={loadRaces}>Retry</button
			>
		</div>
	{:else}
		<p class="text-base-content/70">{filteredRaces.length} of {races.length} race designs</p>
		<div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
			{#each filteredRaces as race (race.id)}
				<article class="card bg-base-200 border border-base-300 min-w-0">
					<div class="card-body p-4 gap-3">
						<div class="flex items-start justify-between gap-2">
							<div class="min-w-0">
								<h2 class="text-2xl font-semibold">
									<a class="hover:text-accent" href="/races/{race.id}">{race.pluralName}</a>
								</h2>
								{#if race.description}<p
										class="text-base-content/70 mt-1"
										data-type="race-description"
									>
										{race.description}
									</p>{/if}
								<span class="badge badge-secondary h-auto mt-2">{getLabelForPRT(race.prt)}</span>
							</div>
							<details class="dropdown dropdown-end">
								<summary class="btn btn-ghost" aria-label={`Manage ${race.pluralName}`}>•••</summary
								>
								<ul
									class="dropdown-content menu bg-base-100 border border-base-300 rounded-box shadow-lg z-20 w-36 p-2"
								>
									<li>
										<button
											class="text-error"
											onclick={() => removeItem(race)}
											data-type="delete-button"
											data-id={race.id}>Delete race</button
										>
									</li>
								</ul>
							</details>
						</div>
						<div class="flex flex-wrap items-center gap-1">
							{#each lrts.filter((lrt) => (race.lrts & lrt) !== 0) as lrt (lrt)}<abbr
									class="badge badge-outline badge-sm no-underline cursor-help"
									title={getLabelForLRT(lrt)}>{getAbbreviationForLRT(lrt)}</abbr
								>{:else}<span class="text-base-content/70">No lesser traits</span>{/each}
						</div>
						<div class="space-y-1 race-hab-bars">
							<HabBar
								compact
								habType={Grav}
								habLow={race.habLow?.grav}
								habHigh={race.habHigh?.grav}
								immune={race.immuneGrav}
							/>
							<HabBar
								compact
								habType={Temp}
								habLow={race.habLow?.temp}
								habHigh={race.habHigh?.temp}
								immune={race.immuneTemp}
							/>
							<HabBar
								compact
								habType={Rad}
								habLow={race.habLow?.rad}
								habHigh={race.habHigh?.rad}
								immune={race.immuneRad}
							/>
						</div>
						<div
							class="flex flex-wrap items-center justify-between gap-2 mt-auto pt-3 border-t border-base-300"
						>
							<span class="text-base-content/70">Growth {race.growthRate}% / year</span><a
								class="btn btn-primary"
								href="/races/{race.id}">Edit race →</a
							>
						</div>
					</div>
				</article>
			{:else}<div class="card bg-base-200 border border-base-300 md:col-span-2 xl:col-span-3">
					<div class="card-body text-center">
						<h2 class="text-2xl">
							{races.length ? 'No matching races' : 'Design your first race'}
						</h2>
						<p>
							{races.length
								? 'Try another search or primary trait.'
								: 'Choose your traits and habitability for your next galaxy.'}
						</p>
					</div>
				</div>{/each}
		</div>{/if}
</div>

<script lang="ts">
	import { page } from '$app/state';

	import { goto } from '$app/navigation';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import { raceClient } from '$lib/services/connect';
	import { addError } from '$lib/services/Errors';
	import { notify } from '$lib/services/Notifications';
	import type { Race } from '$lib/types/cs-proto';
	import { humanoid } from '$lib/types/Race';
	import { loadWasm, type CS } from '$lib/wasm';
	import { ConnectError } from '@connectrpc/connect';
	import { onMount } from 'svelte';
	import RaceEditor from './RaceEditor.svelte';
	import RacePoints from './RacePoints.svelte';

	let id = page.params.id;
	let race: Race = $state(humanoid());
	let cs: CS | undefined = $state();

	onMount(async () => {
		loadWasm().then((res) => (cs = res));
		if (id !== 'new') {
			try {
				const resp = await raceClient.getRace({ raceId: BigInt(id) });
				if (resp.race) {
					race = resp.race;
				}
			} catch (e) {
				addError(e as ConnectError);
			}
		} else {
			// create a new humanoid
			race = humanoid();
		}
	});

	const onSubmit = async () => {
		if (id === 'new') {
			const { race: created } = await raceClient.createRace({ race });
			await goto(`/races/${created?.id}`);
		} else {
			await raceClient.updateRace({ race });
		}

		notify('Saved ' + race.pluralName);
	};

	let saveDisabled = $state(false);
</script>

{#if race}
	<form
		onsubmit={(e) => {
			e.preventDefault();
			onSubmit();
		}}
	>
		<div class="w-full flex justify-end gap-2">
			<button class="btn btn-success" type="submit" disabled={saveDisabled}>Save</button>
		</div>

		<ItemTitle>{race.name}</ItemTitle>
		{#if cs}
			<RacePoints
				wasmClient={cs.wasmService}
				{race}
				onPointsUpdated={(points) => (saveDisabled = points < 0)}
			/>
		{/if}
		<RaceEditor bind:race />
	</form>
{/if}

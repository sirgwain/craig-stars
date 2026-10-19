<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount, untrack } from 'svelte';
	import LoadingSpinner from '#lib/components/LoadingSpinner.svelte';
	import { gameClient, raceClient } from '#lib/services/connect.js';
	import { addError } from '#lib/services/Errors.js';
	import { me } from '#lib/services/Stores.js';
	import { NewGamePlayerType, type GameSettings, type Race } from '#lib/types/cs-proto.js';
	import { humanoid } from '#lib/types/Race.js';
	import type { ConnectError } from '@connectrpc/connect';
	import GameCreationForm from './GameCreationForm.svelte';
	import { allPlayerTypes, singlePlayerSettings } from './defaults';

	type Props = { initialSettings?: GameSettings };
	let { initialSettings }: Props = $props();
	let settings = $state(untrack(() => initialSettings ?? singlePlayerSettings()));
	let races: Race[] = $state([humanoid()]);
	let loaded = $state(false);
	let submitting = $state(false);
	let allowedPlayerTypes = $derived($me.isGuest() ? [NewGamePlayerType.AI] : allPlayerTypes);

	onMount(async () => {
		try {
			const response = await raceClient.getRaces({});
			if (response.races.length) races = response.races;
		} catch (e) {
			addError(e as ConnectError);
		} finally {
			loaded = true;
		}
	});

	async function onSubmit() {
		if (submitting) return;
		submitting = true;
		try {
			const response = await gameClient.createGame({ settings });
			if (response.game?.game) await goto(`/games/${response.game.game.id}`);
		} catch (e) {
			addError(e as ConnectError);
		} finally {
			submitting = false;
		}
	}
</script>

{#if loaded}
	<GameCreationForm bind:settings {races} {allowedPlayerTypes} {submitting} {onSubmit} />
{:else}
	<LoadingSpinner />
{/if}

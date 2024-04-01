<script lang="ts" context="module">
	export type SubmitTurnEvent = {
		'submit-turn': void;
	};
</script>

<script lang="ts">
	import { page } from '$app/stores';
	import ErrorToast from '$lib/components/ErrorToast.svelte';
	import LoadingModal from '$lib/components/LoadingModal.svelte';
	import Popup from '$lib/components/game/tooltips/Popup.svelte';
	import Tooltip from '$lib/components/game/tooltips/Tooltip.svelte';
	import { loadingModalText } from '$lib/services/Stores';
	import { createEventDispatcher } from 'svelte';
	import GameMenu from './GameMenu.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	let id = parseInt($page.params.id);
	const { game } = getGameContext();
	const dispatch = createEventDispatcher<SubmitTurnEvent>();

	// bubble the submit turn event up
	function onSubmitTurn() {
		dispatch('submit-turn');
	}
</script>

<main class="flex flex-col h-[100dvh]">
	<header class="flex-none z-50">
		<GameMenu on:submit-turn={onSubmitTurn} />
	</header>
	<ErrorToast />
	<LoadingModal text={$loadingModalText} />
	<slot>Game</slot>
</main>
<Tooltip />
<Popup />

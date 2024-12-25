<script lang="ts">
	import ErrorToast from '$lib/components/ErrorToast.svelte';
	import LoadingModal from '$lib/components/LoadingModal.svelte';
	import NotificationToast from '$lib/components/NotificationToast.svelte';
	import Popup from '$lib/components/game/tooltips/Popup.svelte';
	import Tooltip from '$lib/components/game/tooltips/Tooltip.svelte';
	import { loadingModalText } from '$lib/services/Stores';
	import { type Snippet } from 'svelte';
	import GameMenu from './GameMenu.svelte';

	type Props = {
		children?: Snippet;
		onSubmitTurn?: () => void;
	};

	let { children, onSubmitTurn }: Props = $props();
</script>

<main class="flex flex-col h-[100dvh]">
	<header class="flex-none z-50">
		<GameMenu {onSubmitTurn} />
	</header>
	<ErrorToast />
	<NotificationToast />
	<LoadingModal text={$loadingModalText} />
	{#if children}{@render children()}{:else}Game{/if}
</main>
<Tooltip />
<Popup />

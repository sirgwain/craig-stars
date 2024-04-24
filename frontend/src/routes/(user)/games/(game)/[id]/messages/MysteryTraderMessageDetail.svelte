<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { MessageType, type Message } from '$lib/types/Message';
	import type { MysteryTrader } from '$lib/types/MysteryTrader';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;
	export let mysteryTrader: MysteryTrader;
</script>

{#if message.type === MessageType.MysteryTraderDiscovered}
	A mysterious trading vessel broadcasting a proposal has been detected entering known space.
{:else}
	<!-- Fallback for unknown message types -->
	<FallbackMessageDetail {message} />
{/if}

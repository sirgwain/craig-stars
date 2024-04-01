<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { MessageType, type Message } from '$lib/types/Message';
	import { $enum as eu } from 'ts-enum-util';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;
</script>

{#if message.type === MessageType.PlayerNoPlanets}
	All your planets have been overrun.
	{#if (message.spec.amount ?? 0) > 0}
		You still have colonists on a freighter, you can still recover from this this setback!
	{:else}
		You have no colonists on any of your remaining ships. You may remain in the game and harry your
		opponents with your rogue fleets, but you have lost.
	{/if}
{:else if message.type === MessageType.PlayerDead}
	You are dead. All your planets have been overrun and your spaceships defeated.
{:else}
	<FallbackMessageDetail {message} />
{/if}

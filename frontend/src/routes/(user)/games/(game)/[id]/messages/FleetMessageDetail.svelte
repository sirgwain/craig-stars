<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { MessageType, type Message } from '$lib/types/Message';
	import type { Fleet } from '$lib/types/Fleet';
	import type { PlayerIntel } from '$lib/types/Player';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;
	export let fleet: Fleet;
	export let owner: PlayerIntel | undefined;
</script>

{#if message.type === MessageType.FleetColonistDieoff}
	Engine radiation has killed {message.spec.amount} colonists traveling in {fleet.name}.
{:else}
	{message.text}
{/if}

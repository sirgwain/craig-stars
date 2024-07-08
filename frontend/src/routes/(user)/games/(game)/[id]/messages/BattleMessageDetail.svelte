<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { getBattleRecordDetails } from '$lib/types/Battle';
	import { MessageType, type Message } from '$lib/types/Message';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;

	type Details = {
		present: boolean;
		location: string;
		ours: number;
		theirs: number;
		ourDead: number;
		theirDead: number;
		oursLeft: number;
		theirsLeft: number;
	};

	function getBattleMessageDetails(message: Message): Details | undefined {
		const battle = $universe.getBattle(message.battleNum);
		if (battle) {
			return getBattleRecordDetails(battle, $player, $universe);
		}
	}

	$: details = getBattleMessageDetails(message);
</script>

{#if message.text}
	{message.text}
{:else if details}
	{#if message.type === MessageType.Battle}
		A battle took place at {details.location}.
		{#if details.ourDead === 0 && details.theirDead === 0}
			No ships were lost on either side.
		{:else if details.ourDead === 0 && details.theirDead === details.theirs}
			All enemy forces were destroyed. You did not suffer a single casualty.
		{:else if details.ourDead === details.ours && details.theirDead === 0}
			Your forces were annihilated by the enemy, and they lost no ships.
		{:else if details.ourDead > 0 && details.ours && details.theirDead > 0}
			Both you and the enemy suffered losses during the exchange.
		{/if}
	{:else if message.type === MessageType.BattleAlly}
		Your ally was involved in a battle at {details.location}.
		{#if details.ourDead === 0 && details.theirDead === 0}
			No ships were lost on either side.
		{:else if details.ourDead === 0 && details.theirDead === details.theirs}
			All enemy forces were destroyed. Your allies did not suffer a single casualty.
		{:else if details.ourDead === details.ours && details.theirDead === 0}
			Your ally's forces were annihilated by the enemy, who lost no ships.
		{:else if details.ourDead > 0 && details.ours && details.theirDead > 0}
			Both your allies and the enemy suffered losses during the exchange.
		{/if}
	{:else}
		A battle took place at an unknown location.
	{/if}
{:else}
	<FallbackMessageDetail {message} />
{/if}

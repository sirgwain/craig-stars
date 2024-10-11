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
			{#if details.theirs > 1}
				All {details.theirs ?? 0} enemy forces were destroyed
			{:else}
				A single enemy ship was destroyed
			{/if} by your fleet of {details.ours ?? 0} {details.ours === 1 ? 'ship' : 'ships'}, which did
			not suffer a single casualty.
		{:else if details.ourDead === details.ours && details.theirDead === 0}
			{#if details.ours > 1}
				All {details.ours ?? 0} of your ships were annihilated
			{:else}
				A single ship was annihilated
			{/if}
			by the enemy's {details.theirs ?? 0} {details.theirs === 1 ? 'vessel' : 'vessels'}, which
			suffered no casualties.
		{:else if details.ourDead > 0 && details.ours && details.theirDead > 0}
			Both you and the enemy suffered losses during the exchange. You lost {details.ourDead ?? 0} out
			of {details.ours ?? 0} {details.ours === 1 ? 'ship' : 'ships'}, while they lost {details.theirDead ??
				0} out of {details.theirs ?? 0} ships.
		{/if}
	{:else if message.type === MessageType.BattleAlly}
		Your ally was involved in a battle at {details.location}.
		{#if details.ourDead === 0 && details.theirDead === 0}
			No ships were lost on either side.
		{:else if details.ourDead === 0 && details.theirDead === details.theirs}
			{#if details.theirs > 1}
				All {details.theirs ?? 0} enemy forces were destroyed
			{:else}
				A single enemy ship was destroyed
			{/if} by their fleet of {details.ours ?? 0} {details.ours === 1 ? 'ship' : 'ships'}, which
			did not suffer a single casualty.
		{:else if details.ourDead === details.ours && details.theirDead === 0}
			{#if details.ours > 1}
				All {details.ours ?? 0} of your ally's ships were annihilated
			{:else}
				A single ally ship was annihilated
			{/if}
			by the enemy's {details.theirs ?? 0} {details.theirs === 1 ? 'vessel' : 'vessels'}, which
			suffered no casualties.
		{:else if details.ourDead > 0 && details.ours && details.theirDead > 0}
			Both your ally and the enemy suffered losses during the exchange. Your allies lost {details.ourDead ??
				0} out of {details.ours ?? 0} {details.ours === 1 ? 'ship' : 'ships'}, while they lost {details.theirDead ??
				0} out of {details.theirs ?? 0} ships.
		{/if}
	{:else}
		A battle took place at an unknown location.
	{/if}
{:else}
	<FallbackMessageDetail {message} />
{/if}

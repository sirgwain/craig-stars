<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { getBattleRecordDetails } from '$lib/types/Battle';
	import { PlayerMessageType, type PlayerMessage } from '$lib/types/cs-proto';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	const { player, universe } = getGameContext();

	type Props = {
		message: PlayerMessage;
	};

	let { message }: Props = $props();

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

	function getBattleMessageDetails(message: PlayerMessage): Details | undefined {
		const battle = $universe.getBattle(message.battleNum);
		if (battle) {
			return getBattleRecordDetails(battle, $player, $universe);
		}
	}

	let details = $derived(getBattleMessageDetails(message));
</script>

{#if message.text}
	{message.text}
{:else if details}
	{#if message.type === PlayerMessageType.BATTLE}
		A battle took place at {details.location}.
		{#if details.ourDead === 0 && details.theirDead === 0}
			No ships were lost by either side.
		{:else if details.ourDead === 0 && details.theirs === details.theirDead}
			Your fleet of {details.ours}
			{details.ours === 1 ? 'ship' : 'ships'}
			decimated the enemy's {details.theirs}
			{details.theirs === 1 ? 'vessel' : 'vessels'}
			without suffering a single casualty.
		{:else if details.ours === details.ourDead && details.theirDead === 0}
			Your fleet of {details.ours}
			{details.ours === 1 ? 'ship' : 'ships'}
			was annihilated by the enemy's {details.theirs}
			{details.theirs === 1 ? 'vessel' : 'vessels'}, which suffered no casualties.
		{:else if details.ours === details.ourDead && details.theirs === details.theirDead}
			Your {details.ours}
			{details.ours === 1 ? 'ship' : 'ships'}
			and the enemy's {details.theirs}
			{details.theirs === 1 ? 'ships' : 'ships'}
			completely destroyed each other. No survivors were left on either side.
		{:else}
			Both you and the enemy suffered losses during the exchange.
			{#if details.ours === details.ourDead && details.theirs > details.theirDead}
				Your {details.ours}
				{details.ours === 1 ? 'ship was' : 'ships were'}
				able to defeat {details.theirDead} out of {details.theirs}
				{details.theirs === 1 ? 'enemy ship' : 'enemy ships'} before being wiped out.
			{:else if details.theirs === details.theirDead && details.ours > details.ourDead}
				Your fleet of {details.ours}
				{details.ours === 1 ? 'ship' : 'ships'}
				defeated the enemy's entire armada of {details.theirs}
				{details.theirs === 1 ? 'ship' : 'ships'}, but lost {details.ourDead}
				{details.ourDead === 1 ? 'ship' : 'ships'} in the process.
			{:else}
				You lost {details.ourDead} of {details.ours}
				{details.ours === 1 ? 'ship' : 'ships'}, while they lost {details.theirDead} of {details.theirs}
				ships.
			{/if}
		{/if}
	{:else if message.type === PlayerMessageType.BATTLE_ALLY}
		Your ally was involved in a battle at {details.location}.
		{#if details.ourDead === 0 && details.theirDead === 0}
			No ships were lost by either side.
		{:else if details.ourDead === 0 && details.theirs === details.theirDead}
			Your ally's {details.ours}
			{details.ours === 1 ? 'ship' : 'ships'}
			decimated the enemy's {details.theirs}
			{details.theirs === 1 ? 'vessel' : 'vessels'}
			without suffering a single casualty.
		{:else if details.ours === details.ourDead && details.theirDead === 0}
			Your ally's fleet of {details.ours}
			{details.ours === 1 ? 'ship' : 'ships'}
			was annihilated by the enemy's {details.theirs}
			{details.theirs === 1 ? 'vessel' : 'vessels'}, which suffered no casualties.
		{:else if details.ours === details.ourDead && details.theirs === details.theirDead}
			Your ally's {details.ours}
			{details.ours === 1 ? 'ship' : 'ships'}
			and the enemy's {details.theirs}
			{details.theirs === 1 ? 'vessels' : 'vessels'}
			completely destroyed each other. No survivors were left on either side.
		{:else}
			Both your ally and the enemy suffered losses during the exchange.
			{#if details.ours === details.ourDead && details.theirs > details.theirDead}
				Your ally's {details.ours}
				{details.ours === 1 ? 'ship was' : 'ships were'}
				able to defeat {details.theirDead} out of {details.theirs}
				{details.theirs === 1 ? 'enemy ship' : ' enemy ships'} before being wiped out.
			{:else if details.theirs === details.theirDead && details.ours > details.ourDead}
				Your ally's fleet of {details.ours}
				{details.ours === 1 ? 'ship' : 'ships'}
				defeated the enemy's fleet of {details.theirs}
				{details.theirs === 1 ? 'ship' : 'ships'}, but lost {details.ourDead}
				{details.ourDead === 1 ? 'ship' : 'ships'} in the process.
			{:else}
				Your ally lost {details.ourDead} out of {details.ours}
				{details.ours === 1 ? 'ship' : 'ships'}, while the enemy lost {details.theirDead} out of
				{details.theirs} ships.
			{/if}
		{/if}
	{:else}
		A battle took place at an unknown location.
	{/if}
{:else}
	<FallbackMessageDetail {message} />
{/if}

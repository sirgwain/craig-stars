<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { MessageType, type Message } from '$lib/types/Message';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	const { game, player, universe, settings } = getGameContext();

	export let message: Message;
</script>

{#if message.text}
	{message.text}
{:else if message.type === MessageType.BattleReports}
	{#if $universe.battles?.length === 1}
		You have received a battle recording this year.
	{:else}
		You have received {$universe.battles?.length} battle recordings this year.
	{/if}
{:else if message.type === MessageType.PlayerNoPlanets}
	All your planets have been overrun.
	{#if (message.spec.amount ?? 0) > 0}
		You still have colonists on a freighter, you can still recover from this this setback!
	{:else}
		You have no colonists on any of your remaining ships. You may remain in the game and harry your
		opponents with your rogue fleets, but you have lost.
	{/if}
{:else if message.type === MessageType.PlayerDead}
	{#if message.targetPlayerNum == $player.num}
		You are dead. All your planets have been overrun and your spaceships defeated.
	{:else}
		All traces of the {$universe.getPlayerPluralName(message.targetPlayerNum)} have been eliminated from
		the galaxy. May they rest in peace.
	{/if}
{:else if message.type === MessageType.PlayerTechLevelGainedBattle}
	{@const battle = $universe.getBattle(message.battleNum)}
	Wreckage from the battle that occurred at ({battle?.position.x ?? 0}, {battle?.position.y ?? 0})
	has boosted your research in {message.spec.field} by 1 level.
{:else}
	<FallbackMessageDetail {message} />
{/if}

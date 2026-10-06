<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { PlayerMessageType, TechCategory, type PlayerMessage } from '$lib/types/cs-proto';
	import { techs } from '$lib/services/Stores';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	const { game, player, universe } = getGameContext();

	type Props = {
		message: PlayerMessage;
	};

	let { message }: Props = $props();
</script>

{#if message.text}
	{message.text}
{:else if message.type === PlayerMessageType.ERROR}
	Something went wrong on the server. Please contact the administrator, {message.spec?.error}
{:else if message.type === PlayerMessageType.PLAYER_DISCOVERY}
	You have discovered a new species, the {message.spec?.name ||
		$universe.getPlayerPluralName(message.target?.targetPlayerNum)}. You are not alone in this
	universe!
{:else if message.type === PlayerMessageType.PLAYER_GAIN_TECH_LEVEL}
	<!-- Amount is the research level reached. -->
	Your scientists have completed research into Tech Level {message.spec?.amount ?? 0} for {message
		.spec?.field}. They will continue their efforts in the {message.spec?.nextField} field.
{:else if message.type === PlayerMessageType.PLAYER_TECH_GAINED}
	{@const category = $techs.getTech(message.spec?.techGained ?? '')?.tech?.category}
	Your recent breakthrough in {message.spec?.field}
	{#if category === TechCategory.SHIP_HULL || category === TechCategory.STARBASE_HULL}
		has also given you the {message.spec?.techGained}
		{category === TechCategory.SHIP_HULL ? 'ship' : 'starbase'} hull. To create a design with this hull,
		go to <a class="link" href="/games/{$game.id}/designer">Ship Designs</a> and select "Create".
	{:else if category === TechCategory.PLANETARY_DEFENSE}
		has also taught you how to build {message.spec?.techGained} defenses. All existing planetary defenses
		have been upgraded to the new technology.
	{:else if category === TechCategory.PLANETARY_SCANNER}
		has also taught you how to build the {message.spec?.techGained} scanner. All existing planetary scanners
		have been upgraded to the new technology.
	{:else}
		has also given you the {message.spec?.techGained} benefit.
	{/if}
{:else if message.type === PlayerMessageType.PLAYER_VICTOR}
	{#if message.target?.targetPlayerNum === $player.num}
		You have been declared the winner of this grand game. You may continue to play though, if you
		wish to really rub your nose in everyone else's faces.
	{:else}
		The {message.spec?.name || $universe.getPlayerName(message.target?.targetPlayerNum)} have been declared
		the winner of this game. You are advised to accept their supremacy, though you may continue the fight
		regardless.
	{/if}
{:else if message.type === PlayerMessageType.BATTLE_REPORTS}
	{#if $universe.battleRecords.length === 1}
		You have received a battle recording this year.
	{:else}
		You have received {$universe.battleRecords.length} battle recordings this year.
	{/if}
{:else if message.type === PlayerMessageType.PLAYER_NO_PLANETS}
	All your planets have been overrun.
	{#if (message.spec?.amount ?? 0) > 0}
		You still have colonists on a freighter, you can still recover from this this setback!
	{:else}
		You have no colonists on any of your remaining ships. You may remain in the game and harry your
		opponents with your rogue fleets, but you have lost.
	{/if}
{:else if message.type === PlayerMessageType.PLAYER_DEAD}
	{#if message.target?.targetPlayerNum == $player.num}
		You are dead. All your planets have been overrun and your spaceships defeated.
	{:else}
		All traces of the {$universe.getPlayerPluralName(message.target?.targetPlayerNum)} have been eliminated
		from the galaxy. May they rest in peace.
	{/if}
{:else if message.type === PlayerMessageType.PLAYER_TECH_LEVEL_GAINED_BATTLE}
	{@const battle = $universe.getBattle(message.battleNum)}
	Wreckage from the battle that occurred at ({battle?.position?.x ?? 0}, {battle?.position?.y ?? 0})
	has boosted your research in {message.spec?.field} by 1 level.
{:else}
	<FallbackMessageDetail {message} />
{/if}

<script lang="ts">
	import {
		MysteryTraderRewardType,
		PlayerMessageType,
		type PlayerMessage
	} from '$lib/types/cs-proto';
	import { isHullComponent } from '$lib/types/MysteryTrader';
	import { sum } from '$lib/types/TechLevel';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	type Props = {
		message: PlayerMessage;
	};

	let { message }: Props = $props();
</script>

{#if message.type === PlayerMessageType.MYSTERY_TRADER_DISCOVERED}
	A mysterious trading vessel broadcasting a proposal has been detected entering known space.
{:else if message.type === PlayerMessageType.MYSTERY_TRADER_AGAIN}
	The Mystery Trader has decided to make another pass through known space to increase access to its
	wares.
{:else if message.type === PlayerMessageType.MYSTERY_TRADER_CHANGED_COURSE}
	The Mystery Trader has unexplicably changed course and/or speed. Perhaps something startled him?
{:else if message.type === PlayerMessageType.MYSTERY_TRADER_ALREADY_REWARDED}
	The Mystery Trader eyes the captain of {message.spec?.mapObjectTarget?.targetName} suspiciously and
	suggests that he is still recovering from the last transaction with you.
{:else if message.type === PlayerMessageType.MYSTERY_TRADER_MET_WITHOUT_REWARD}
	{@const detail = message.spec?.mysteryTrader}
	{#if detail?.ship}
		<!-- This will occur if the player has a design with the same name that isn't flagged as a MysteryTrader design -->
		The Mystery Trader tried to give you an auxillary ship called the {detail.ship}, but you were
		unable to learn the design.
	{:else}
		The Mystery Trader has refused to give the captain of {message.spec?.mapObjectTarget
			?.targetName} an audience. It may be due to an insufficient quantity of minerals carried by your
		fleet.
	{/if}
{:else if message.type === PlayerMessageType.MYSTERY_TRADER_MET_WITH_REWARD}
	{@const detail = message.spec?.mysteryTrader}
	{message.spec?.mapObjectTarget?.targetName} has been absorbed by the Mystery Trader.
	{#if detail}
		{#if detail.type === MysteryTraderRewardType.RESEARCH}
			The trader has given you {sum(detail.techLevels)} technology advances.
		{:else if isHullComponent(detail.type)}
			You have been given the plans for a unique part to place on your ships. The trader suggests
			you visit other traders.
		{:else if detail.type === MysteryTraderRewardType.SHIP_HULL}
			In return, you have been given the plans for a new ship hull. The trader suggests you visit
			other traders.
		{:else if detail.type === MysteryTraderRewardType.LIFEBOAT}
			In return, you have been given {detail.shipCount} of the Trader's auxillary ships for your own
			use.
		{:else if detail.type === MysteryTraderRewardType.GENESIS}
			In return, you have been given the plans for a powerful planetary device. The trader suggests
			you visit other traders.
		{:else if detail.type === MysteryTraderRewardType.UNSPECIFIED}
			However, the trader was unable to teach you anything new.
		{:else}
			The trader has given you a boon, but this paltry web client can't tell what it is.
		{/if}
	{:else}
		The trader has given you a boon, but this paltry web client can't tell what it is.
	{/if}
	The trader suggests that you visit other traders as they may carry different items of interest.
{:else}
	<!-- Fallback for unknown message types -->
	<FallbackMessageDetail {message} />
{/if}

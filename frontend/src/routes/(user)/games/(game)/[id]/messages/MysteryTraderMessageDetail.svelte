<script lang="ts">
	import {
		MysteryTraderRewardGenesis,
		MysteryTraderRewardLifeboat,
		MysteryTraderRewardNone,
		MysteryTraderRewardResearch,
		MysteryTraderRewardShipHull,
		PlayerMessageMysteryTraderAgain,
		PlayerMessageMysteryTraderAlreadyRewarded,
		PlayerMessageMysteryTraderChangedCourse,
		PlayerMessageMysteryTraderDiscovered,
		PlayerMessageMysteryTraderMetWithoutReward,
		PlayerMessageMysteryTraderMetWithReward,
		type PlayerMessage
	} from '$lib/types/cs';
	import { isHullComponent } from '$lib/types/MysteryTrader';
	import { sum } from '$lib/types/TechLevel';
	import FallbackMessageDetail from './FallbackMessageDetail.svelte';

	type Props = {
		message: PlayerMessage;
	};

	let { message }: Props = $props();
</script>

{#if message.type === PlayerMessageMysteryTraderDiscovered}
	A mysterious trading vessel broadcasting a proposal has been detected entering known space.
{:else if message.type === PlayerMessageMysteryTraderAgain}
	The Mystery Trader has decided to make another pass through known space to increase access to its
	wares.
{:else if message.type === PlayerMessageMysteryTraderChangedCourse}
	The Mystery Trader has unexplicably changed course and/or speed. Perhaps something startled him?
{:else if message.type === PlayerMessageMysteryTraderAlreadyRewarded}
	The Mystery Trader eyes the captain of {message.spec.targetName} suspiciously and suggests that he
	is still recovering from the last transaction with you.
{:else if message.type === PlayerMessageMysteryTraderMetWithoutReward}
	{@const detail = message.spec.mysteryTrader}
	{#if detail?.ship}
		<!-- This will occur if the player has a design with the same name that isn't flagged as a MysteryTrader design -->
		The Mystery Trader tried to give you an auxillary ship called the {detail.ship}, but you were
		unable to learn the design.
	{:else}
		The Mystery Trader has refused to give the captain of {message.spec.targetName} an audience. It may
		be due to an insufficient quantity of minerals carried by your fleet.
	{/if}
{:else if message.type === PlayerMessageMysteryTraderMetWithReward}
	{@const detail = message.spec.mysteryTrader}
	{message.spec.targetName} has been absorbed by the Mystery Trader.
	{#if detail}
		{#if detail.type === MysteryTraderRewardResearch}
			The trader has given you {sum(detail.techLevels)} technology advances.
		{:else if isHullComponent(detail.type)}
			You have been given the plans for a unique part to place on your ships. The trader suggests
			you visit other traders.
		{:else if detail.type === MysteryTraderRewardShipHull}
			In return, you have been given the plans for a new ship hull. The trader suggests you visit
			other traders.
		{:else if detail.type === MysteryTraderRewardLifeboat}
			In return, you have been given {detail.shipCount ?? 0} of the Trader's auxillary ships for your
			own use.
		{:else if detail.type === MysteryTraderRewardGenesis}
			In return, you have been given the plans for a powerful planetary device. The trader suggests
			you visit other traders.
		{:else if detail.type === MysteryTraderRewardNone}
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

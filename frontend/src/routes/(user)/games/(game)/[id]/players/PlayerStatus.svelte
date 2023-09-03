<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { GameState } from '$lib/types/Game';
	import type { PlayerStatus } from '$lib/types/Player';
	import { CheckBadge, XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { game, player, universe } = getGameContext();

	export let playerStatus: PlayerStatus;
</script>

{#if $game.state === GameState.Setup}
	{#if playerStatus.ready}
		<div class="flex flex-row">
			<div class="w-20">Ready</div>
			<Icon src={CheckBadge} size="24" class="stroke-success" />
		</div>
	{:else}
		<div class="flex flex-row">
			<div class="w-20">Waiting</div>
			<Icon src={XMark} size="24" class="stroke-error" />
		</div>
	{/if}
{:else if playerStatus.submittedTurn}
	<div class="flex flex-row">
		<div class="w-20">Submitted</div>
		<Icon src={CheckBadge} size="24" class="stroke-success" />
	</div>
{:else}
	<div class="flex flex-row">
		<div class="w-20">Playing</div>
		<Icon src={XMark} size="24" class="stroke-error" />
	</div>
{/if}

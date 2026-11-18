<script lang="ts">
	import EnumSelect from '#lib/components/EnumSelect.svelte';
	import { me } from '#lib/services/Stores.js';
	import { NewGamePlayerType, type NewGamePlayer } from '#lib/types/cs-proto.js';
	import { XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import AiPlayer from './AIPlayer.svelte';
	import HostPlayer from './HostPlayer.svelte';

	type Props = {
		player: NewGamePlayer;
		index: number;
		onRemove: () => void;
	};

	let { player = $bindable(), index, onRemove }: Props = $props();
</script>

{#if player}
	<div class="border border-base-300 bg-base-200 p-3" data-type="player-slot" data-player={index}>
		<div class="mb-2 flex items-center justify-between gap-2">
			<h4 class="font-semibold">Player {index}{index === 1 ? ' · You' : ''}</h4>
			{#if index !== 1}
				<button
					data-type="delete-button"
					data-id={`Player ${index}`}
					aria-label={`Remove slot ${index}`}
					onclick={onRemove}
					type="button"
					class="btn btn-ghost btn-sm"><Icon size="16" src={XMark} /></button
				>
			{/if}
		</div>
		<div class="grid gap-3 sm:grid-cols-2">
			{#if index !== 1 && !$me.isGuest()}
				<EnumSelect
					name={`player-${index}-type`}
					stacked
					enumType={NewGamePlayerType}
					bind:value={player.type}
					title={`Player ${index}`}
				/>
			{/if}
			{#if player.type === NewGamePlayerType.AI}
				<AiPlayer bind:player />
			{:else if player.type === NewGamePlayerType.HOST}
				<HostPlayer bind:player />
			{:else}
				<p class="my-auto text-sm text-base-content/70">
					Open slot · A player can join and choose their race.
				</p>
			{/if}
		</div>
	</div>
{/if}

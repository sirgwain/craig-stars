<script lang="ts">
	import EnumSelect from '#lib/components/EnumSelect.svelte';
	import { NewGamePlayerType, type NewGamePlayer, type Race } from '#lib/types/cs-proto.js';
	import { XMark } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import AiPlayer from './AIPlayer.svelte';
	import HostPlayer from './HostPlayer.svelte';

	type Props = {
		player: NewGamePlayer;
		index: number;
		onRemove: () => void;
		races: Race[];
		allowedPlayerTypes: NewGamePlayerType[];
	};

	let { player = $bindable(), index, onRemove, races, allowedPlayerTypes }: Props = $props();
</script>

{#if player}
	<div class="block">
		{#if index !== 1}
			<div class="flex flex-row justify-end">
				<div class="grow">
					{#if allowedPlayerTypes.length > 1}
						<div class="grow">
							<EnumSelect
								name="type"
								enumType={NewGamePlayerType}
								typeFilter={(type) => allowedPlayerTypes.includes(type)}
								bind:value={player.type}
								title={`Player ${index}`}
							/>
						</div>
					{:else}
						<div class="text-xl mr-2 my-auto">AI Player {index}</div>
					{/if}
					{#if player.type === NewGamePlayerType.AI}
						<AiPlayer bind:player />
					{:else if player.type === NewGamePlayerType.HOST}
						<!-- multiple host players make a hot seat game -->
						<HostPlayer bind:player {races} />
					{/if}
				</div>
				<div class="my-auto mx-1">
					<button
						data-type="delete-button"
						data-id={`Player ${index}`}
						onclick={onRemove}
						type="button"
						class="btn btn-outline btn-sm my-1 normal-case"><Icon size="16" src={XMark} /></button
					>
				</div>
			</div>
		{:else}
			<HostPlayer bind:player {races} />
		{/if}
	</div>

	<div></div>
{/if}

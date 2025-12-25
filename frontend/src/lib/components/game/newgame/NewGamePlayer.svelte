<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import { me } from '$lib/services/Stores';
	import { NewGamePlayerType, type NewGamePlayer } from '$lib/types/cs-proto';
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
	<div class="block">
		{#if index !== 1}
			<div class="flex flex-row justify-end">
				<div class="grow">
					{#if !$me.isGuest()}
						<div class="grow">
							<EnumSelect
								name="type"
								enumType={NewGamePlayerType}
								bind:value={player.type}
								title={`Player ${index}`}
							/>
						</div>
					{:else}
						<div class="text-xl mr-2 my-auto">AI Player {index}</div>
					{/if}
					{#if player.type === NewGamePlayerType.AI}
						<AiPlayer bind:player />
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
			<HostPlayer bind:player />
		{/if}
	</div>

	<div></div>
{/if}

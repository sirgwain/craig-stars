<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import type { NewGamePlayer } from '$lib/types/Game';
	import { Icon } from '@steeze-ui/svelte-icon';
	import AiPlayer from './AIPlayer.svelte';
	import HostPlayer from './HostPlayer.svelte';
	import { XMark } from '@steeze-ui/heroicons';
	import { createEventDispatcher } from 'svelte';
	import { me } from '$lib/services/Stores';

	const dispatch = createEventDispatcher();

	enum NewGamePlayerChooseType {
		Open = 'Open',
		Gueset = 'Guest',
		AI = 'AI'
	}

	export let player: NewGamePlayer;
	export let index: number;
</script>

{#if player}
	<div class="block">
		{#if index != 1}
			<div class="flex flex-row justify-end">
				{#if !$me.isGuest()}
					<div class="grow">
						<EnumSelect
							enumType={NewGamePlayerChooseType}
							name="type"
							bind:value={player.type}
							title={`Player ${index}`}
						/>
					</div>
				{:else}
					<div class="text-xl mr-2 my-auto">AI Player {index}</div>
				{/if}
				<div class="my-auto">
					<button
						on:click={() => dispatch('remove')}
						type="button"
						class="btn btn-outline btn-sm my-1 normal-case"><Icon size="16" src={XMark} /></button
					>
				</div>
			</div>
			<AiPlayer {player} />
		{:else}
			<HostPlayer {player} />
		{/if}
	</div>
	<div />
{/if}

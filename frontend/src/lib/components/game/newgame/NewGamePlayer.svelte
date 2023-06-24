<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import type { NewGamePlayer } from '$lib/types/Game';
	import AiPlayer from './AIPlayer.svelte';
	import HostPlayer from './HostPlayer.svelte';

	enum NewGamePlayerChooseType {
		Open = 'Open',
		AI = 'AI'
	}

	export let player: NewGamePlayer;
	export let index: number;
</script>

{#if player}
	<div class="block">
		{#if index != 1}
			<EnumSelect
				enumType={NewGamePlayerChooseType}
				name="type"
				bind:value={player.type}
				title={`Player ${index}`}
			/>
			<AiPlayer {player} />
		{:else}
			<HostPlayer {player} />
		{/if}
	</div>
{/if}

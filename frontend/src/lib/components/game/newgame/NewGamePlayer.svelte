<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import { NewGamePlayerType, type NewGamePlayer } from '$lib/types/Game';
	import AiPlayer from './AIPlayer.svelte';
	import HostPlayer from './HostPlayer.svelte';

	export let player: NewGamePlayer;
	export let index: number;
</script>

{#if player}
	<div class="block">
		<EnumSelect
			enumType={NewGamePlayerType}
			name="type"
			bind:value={player.type}
			title={`Player ${index}`}
		/>

		{#if player.type === NewGamePlayerType.Host}
			<HostPlayer {player} />
		{:else if player.type == NewGamePlayerType.AI}
			<AiPlayer {player} />
		{/if}
	</div>
{/if}

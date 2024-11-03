<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import { AIDifficulty, type NewGamePlayer } from '$lib/types/Game';

	interface Props {
		player: NewGamePlayer;
	}

	let { player = $bindable() }: Props = $props();
</script>

<div class="w-3/4 ml-auto grow">
	<EnumSelect
		enumType={AIDifficulty}
		name="type"
		bind:value={player.aiDifficulty}
		title={`Difficulty`}
	/>
</div>

<!-- <ColorInput bind:value={player.color} name="color" /> -->

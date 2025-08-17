<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import { AiDifficulty, type NewGamePlayer } from '$lib/types/cs-proto';

	type Props = {
		player: NewGamePlayer;
	};

	let { player = $bindable() }: Props = $props();
</script>

<div class="w-3/4 ml-auto grow">
	<EnumSelect
		name="type"
		bind:value={player.aiDifficulty}
		enumType={AiDifficulty}
		typeFilter={(t) => t === AiDifficulty.CHEATER || t === AiDifficulty.NORMAL}
		title="Difficulty"
	/>
</div>

<!-- <ColorInput bind:value={player.color} name="color" /> -->

<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import { AIDifficulties } from '$lib/types/Game';
	import { type NewGamePlayer } from '$lib/types/cs';

	type Props = {
		player: NewGamePlayer;
	};

	let { player = $bindable() }: Props = $props();
</script>

<div class="w-3/4 ml-auto grow">
	<EnumSelect
		name="type"
		bind:value={player.aiDifficulty}
		options={AIDifficulties}
		title="Difficulty"
	/>
</div>

<!-- <ColorInput bind:value={player.color} name="color" /> -->

<script lang="ts">
	import ColorInput from '$lib/components/ColorInput.svelte';
	import type { NewGamePlayer } from '$lib/types/Game';

	export let player: NewGamePlayer;
</script>

<!-- <ColorInput bind:value={player.color} name="color" /> -->

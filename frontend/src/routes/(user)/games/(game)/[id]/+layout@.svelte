<script lang="ts">
	import GameLayout from './GameLayout.svelte';
</script>

<GameLayout>
	<!-- We want our main game view to only fill the screen (minus the toolbar) -->
	<div class="grow viewport">
		<slot>Game</slot>
	</div>
</GameLayout>

<style>
	.viewport {
		max-height: calc(100vh-4rem);
		max-height: calc((var(--vh, 1vh) * 100)-4rem);
	}
</style>

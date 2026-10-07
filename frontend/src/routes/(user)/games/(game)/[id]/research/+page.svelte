<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import Research from './Research.svelte';

	const { updatePlayerOrders } = getGameContext();
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<Research onUpdatePlayer={updatePlayerOrders} />
</div>

<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import Relations from './Relations.svelte';

	const { player, updatePlayerRelationships } = getGameContext();
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<Relations player={$player} onUpdatePlayerRelationships={updatePlayerRelationships} />
</div>

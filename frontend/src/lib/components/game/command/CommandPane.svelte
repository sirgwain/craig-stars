<script lang="ts">
	import MineralsOnHandTile from './MineralsOnHandTile.svelte';
	import PlanetProductionTile from './PlanetProductionTile.svelte';
	import PlanetStatusTile from './PlanetStatusTile.svelte';
	import { game, player } from '$lib/services/Context';
</script>

{#if $game && $player.planets?.length}
	<div class="grid grid-cols-1 lg:grid-cols-2 gap-2">
		<MineralsOnHandTile />
		<PlanetStatusTile />
		<PlanetProductionTile />
	</div>
{/if}

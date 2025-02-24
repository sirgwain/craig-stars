<script lang="ts">
	import { goto } from '$app/navigation';
	import { getGameContext } from '$lib/services/GameContext';
	import { ownedBy } from '$lib/types/MapObject';
	import { type MapObject } from '$lib/types/cs';
	import SearchResults from './SearchResults.svelte';

	const { game, player, commandMapObject, zoomToMapObject, selectMapObject } = getGameContext();

	function selectSearchResult(mo: MapObject | undefined) {
		if (mo) {
			if (ownedBy(mo, $player.num)) {
				commandMapObject(mo);
			}
			selectMapObject(mo);
			zoomToMapObject(mo);
			goto(`/games/${$game.id}`);
		}
	}
</script>

<div class="w-full mx-auto md:max-w-2xl h-full">
	<SearchResults onOk={(e) => selectSearchResult(e)} />
</div>

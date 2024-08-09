<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { equal, type MapObject } from '$lib/types/MapObject';
	import { distance } from '$lib/types/Vector';

	const { highlightedMapObject, selectedMapObject, commandedMapObject } = getGameContext();

	let dist = 0;
	let from: MapObject | undefined;
	let to: MapObject | undefined;

	$: {
		if ($highlightedMapObject) {
			to = $highlightedMapObject;
			from = equal($selectedMapObject, $highlightedMapObject)
				? $commandedMapObject
				: $selectedMapObject;
		} else {
			to = $selectedMapObject;
			from = $commandedMapObject;
		}
	}
	$: dist = from && to ? distance(from.position, to?.position) : 0;
</script>

<div class="flex flex-row justify-start gap-3 h-4 text-sm">
	{#if to && dist}
		<div class="w-10">
			ID: {to.num}
		</div>
		<div class="w-20">
			X: {to.position.x}, Y: {to.position.y}
		</div>
		<div>
			{to.name}
		</div>
		{#if from && dist}
			<div>
				{dist.toFixed(1)} ly from {from.name}
			</div>
		{/if}
	{/if}
</div>

<script lang="ts">
	import { run } from 'svelte/legacy';

	import { getGameContext } from '$lib/services/GameContext';
	import { equal, getMapObjectName, type MapObject } from '$lib/types/MapObject';
	import { distance } from '$lib/types/Vector';

	const { highlightedMapObject, selectedMapObject, commandedMapObject } = getGameContext();

	let dist = $state(0);
	let from: MapObject | undefined = $state();
	let to: MapObject | undefined = $state();

	run(() => {
		if ($highlightedMapObject) {
			to = $highlightedMapObject;
			from = equal($selectedMapObject, $highlightedMapObject)
				? $commandedMapObject
				: $selectedMapObject;
		} else {
			to = $selectedMapObject;
			from = $commandedMapObject;
		}
	});
	run(() => {
		dist = from && to ? distance(from.position, to?.position) : 0;
	});
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
			{getMapObjectName(to)}
		</div>
		{#if from && dist}
			<div>
				{dist.toFixed(1)} ly from {getMapObjectName(from)}
			</div>
		{/if}
	{/if}
</div>

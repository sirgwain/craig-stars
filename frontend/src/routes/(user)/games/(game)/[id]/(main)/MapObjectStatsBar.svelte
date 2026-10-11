<script lang="ts">
	import { getGameContext } from '#lib/services/GameContext.js';
	import { equal, getMapObjectName, type MapObjectLike } from '#lib/types/MapObject.js';
	import { distance, emptyVector } from '#lib/types/Vector.js';

	const { highlightedMapObject, selectedMapObject, commandedMapObject } = getGameContext();

	type Props = {
		// bigger text, for where the bar is also something to tap
		large?: boolean;
	};

	let { large = false }: Props = $props();

	let to: MapObjectLike | undefined = $derived(
		$highlightedMapObject ? $highlightedMapObject : $selectedMapObject
	);
	let from: MapObjectLike | undefined = $derived(
		$highlightedMapObject
			? equal($selectedMapObject, $highlightedMapObject)
				? $commandedMapObject
				: $selectedMapObject
			: $commandedMapObject
	);
	const posOf = (m: MapObjectLike | undefined) => m?.mapObject?.position ?? emptyVector();
	const numOf = (m: MapObjectLike | undefined) => m?.mapObject?.num ?? 0;
	let dist = $derived(from && to ? distance(posOf(from), posOf(to)) : 0);
</script>

<div
	class="flex flex-row justify-start gap-3 whitespace-nowrap {large
		? 'h-6 text-base'
		: 'h-4 text-sm'}"
>
	{#if to}
		<div class={large ? 'w-14' : 'w-10'}>
			ID: {numOf(to)}
		</div>
		<div class={large ? 'w-28' : 'w-20'}>
			X: {posOf(to).x}, Y: {posOf(to).y}
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

<script lang="ts">
	import { highlightedMapObject, selectedMapObject } from '$lib/services/Context';
	import { distance } from '$lib/types/Vector';

	let dist = 0;

	$: $selectedMapObject &&
		$highlightedMapObject &&
		(dist = distance($selectedMapObject.position, $highlightedMapObject.position));
</script>

<div class="flex flex-row justify-start gap-3 h-4 text-sm">
	{#if $highlightedMapObject}
		<div class="w-10">
			ID: {$highlightedMapObject.num}
		</div>
		<div class="w-20">
			X: {$highlightedMapObject.position.x}, Y: {$highlightedMapObject.position.y}
		</div>
		<div>
			{$highlightedMapObject.name}
		</div>
		{#if dist && $selectedMapObject}
			<div>
				{dist.toFixed(1)} ly from {$selectedMapObject.name}
			</div>
		{/if}
	{/if}
</div>

<script lang="ts">
	import { highlightedMapObject } from '$lib/services/Context';
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
	{/if}
</div>

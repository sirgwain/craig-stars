<script lang="ts">
	import type { Mineral } from '$lib/types/Mineral';

	export let mineral: Mineral | undefined;
</script>

{#if mineral}
	<div class="tracking-wider text-center">
		<span class="text-ironium">{mineral.ironium ?? 0}</span>
		<span class="text-boranium">{mineral.boranium ?? 0}</span>
		<span class="text-germanium">{mineral.germanium ?? 0}</span>
	</div>
{/if}

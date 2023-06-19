<script lang="ts">
	import type { Mineral } from '$lib/types/Mineral';

	export let mineral: Mineral | undefined;
	export let showUnits = false;
</script>

{#if mineral}
	<div class="tracking-wider text-center">
		<span class="text-ironium">{Math.floor((mineral.ironium ?? 0)).toFixed()}</span>{showUnits ? ' kT' : ''}
		<span class="text-boranium">{Math.floor((mineral.boranium ?? 0)).toFixed()}</span>{showUnits ? ' kT' : ''}
		<span class="text-germanium">{Math.floor((mineral.germanium ?? 0)).toFixed()}</span>{showUnits ? ' kT' : ''}
	</div>
{/if}

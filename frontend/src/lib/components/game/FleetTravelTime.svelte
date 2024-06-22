<script lang="ts">
	import type { Waypoint } from '$lib/types/Fleet';
	import { StargateWarpSpeed } from '$lib/types/MapObject';
	import { distance } from '$lib/types/Vector';

	export let wp1: Waypoint | undefined;
	export let wp2: Waypoint | undefined;

	$: dist = wp1 && wp2 ? distance(wp1.position, wp2.position) : 0;
</script>

{#if wp1 && wp2}
	{#if wp2.warpSpeed === StargateWarpSpeed}
		1 year
	{:else}
		{Math.ceil(dist / (wp2.warpSpeed * wp2.warpSpeed))} years
	{/if}
{:else}
	--
{/if}

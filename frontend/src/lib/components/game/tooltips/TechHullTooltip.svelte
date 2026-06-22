<script module lang="ts">
	import type { TechHullTooltipProps } from './TechHullTooltip';
</script>

<script lang="ts">
	import Hull from '../design/Hull.svelte';

	let { hull }: TechHullTooltipProps = $props();
</script>

<div class="flex flex-row justify-center">
	<Hull {hull} cargoCapacity={hull.cargoCapacity} />
</div>

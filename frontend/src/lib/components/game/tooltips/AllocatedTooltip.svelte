<script module lang="ts">
	import type { AllocatedTooltipProps } from './AllocatedTooltip';
</script>

<script lang="ts">
	import CostMini from '../CostMini.svelte';

	let { cost }: AllocatedTooltipProps = $props();
</script>

<div>
	Allocated <CostMini {cost} />
</div>

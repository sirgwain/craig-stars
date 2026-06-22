<script module lang="ts">
	import type { ShipDesignTooltipProps } from './ShipDesignTooltip';
</script>

<script lang="ts">
	import Design from '$lib/components/game/design/Design.svelte';

	let { design }: ShipDesignTooltipProps = $props();
</script>

<div>
	<Design {design} />
</div>

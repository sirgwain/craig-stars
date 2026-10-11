<script module lang="ts">
	import type { ShipDesignTooltipProps } from './ShipDesignTooltip';
</script>

<script lang="ts">
	import Design from '#lib/components/game/design/Design.svelte';

	let { design, damage }: ShipDesignTooltipProps = $props();
</script>

<div>
	<Design {design} {damage} />
</div>

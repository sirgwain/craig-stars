<script lang="ts">
	import MineralTooltip from '$lib/components/game/tooltips/MineralTooltip.svelte';
	import Tooltip from '$lib/components/game/tooltips/Tooltip.svelte';
</script>

<h1 class="text-xl">Mineral Tooltip</h1>
<Tooltip>
	<MineralTooltip mineralType="Ironium" surfaceAmount={1248} concentration={114} miningRate={52} />
</Tooltip>

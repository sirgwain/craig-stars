<script lang="ts">
	import FactoriesTooltip from '$lib/components/game/tooltips/FactoriesTooltip.svelte';
	import MineralTooltip from '$lib/components/game/tooltips/MineralTooltip.svelte';
	import Tooltip from '$lib/components/game/tooltips/Tooltip.svelte';
</script>

<h1 class="text-xl">Mineral Tooltip</h1>
<MineralTooltip mineralType="Ironium" surfaceAmount={1248} concentration={114} miningRate={52} />

<h1 class="text-xl">Factories Tooltip</h1>
<FactoriesTooltip
	planetName="Brin"
	factories={10}
	maxFactories={1000}
	maxPossibleFactories={1200}
/>

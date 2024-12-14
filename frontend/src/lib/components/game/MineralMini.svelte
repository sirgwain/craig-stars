<script lang="ts">
	import type { Mineral } from '$lib/types/Mineral';

	type Props = {
		mineral: Mineral | undefined;
		showUnits?: boolean;
	};

	let { mineral, showUnits = false }: Props = $props();
</script>

{#if mineral}
	<div class="tracking-wider text-center">
		<span class="text-ironium">{Math.floor(mineral.ironium ?? 0).toFixed()}</span>{showUnits
			? ' kT'
			: ''}
		<span class="text-boranium">{Math.floor(mineral.boranium ?? 0).toFixed()}</span>{showUnits
			? ' kT'
			: ''}
		<span class="text-germanium">{Math.floor(mineral.germanium ?? 0).toFixed()}</span>{showUnits
			? ' kT'
			: ''}
	</div>
{/if}

<script lang="ts">
	import type { Cargo } from '$lib/types/Cargo';

	export let cargo: Cargo | undefined;
</script>

{#if cargo}
	<div class="tracking-wider text-center">
		<span class="text-ironium">{cargo.ironium ?? 0}</span>
		<span class="text-boranium">{cargo.boranium ?? 0}</span>
		<span class="text-germanium">{cargo.germanium ?? 0}</span>
		<span class="text-colonists">{cargo.colonists ?? 0}</span>
	</div>
{/if}

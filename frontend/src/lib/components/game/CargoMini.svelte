<script lang="ts">
	import type { Cargo } from '$lib/types/cs-proto';

	type Props = {
		cargo: Cargo | undefined;
	};

	let { cargo }: Props = $props();
</script>

{#if cargo}
	<div class="tracking-wider text-center">
		<span class="text-ironium">{cargo.ironium}</span>
		<span class="text-boranium">{cargo.boranium}</span>
		<span class="text-germanium">{cargo.germanium}</span>
		<span class="text-colonists">{cargo.colonists}</span>
	</div>
{/if}

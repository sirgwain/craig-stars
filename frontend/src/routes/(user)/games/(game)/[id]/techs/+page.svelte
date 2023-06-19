<script lang="ts">
	import Techs from '$lib/components/Techs.svelte';
	import { game } from '$lib/services/Stores';
</script>

{#if game}
	<div class="w-full">
		<Techs techStore={$game?.rules.techs} />
	</div>
{/if}
<script lang="ts">
	import Techs from '$lib/components/Techs.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
</script>

<Breadcrumb>
	{#snippet crumbs()}
	
			<li>Techs</li>
		
	{/snippet}
</Breadcrumb>

<div class="mt-2">
	<Techs />
</div>

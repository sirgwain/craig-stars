<script lang="ts">
	import Techs from '$lib/components/Techs.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li>Techs</li>
	</svelte:fragment>
</Breadcrumb>

<div class="mt-2">
	<Techs />
</div>

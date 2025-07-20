<script lang="ts">
	import Techs from '$lib/components/Techs.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/GameContext';

	const { player, cs } = getGameContext();
</script>

<Breadcrumb>
	{#snippet crumbs()}
		<li>Techs</li>
	{/snippet}
</Breadcrumb>

<Techs player={$player} {cs} />

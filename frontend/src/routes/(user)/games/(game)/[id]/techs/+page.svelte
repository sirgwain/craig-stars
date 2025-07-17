<script lang="ts">
	import Techs from '$lib/components/Techs.svelte';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { techs } from '$lib/services/Stores';

	const { player, cs } = getGameContext();
</script>

<Breadcrumb>
	{#snippet crumbs()}
		<li>Techs</li>
	{/snippet}
</Breadcrumb>

<Techs techStore={$techs} player={$player} {cs} />

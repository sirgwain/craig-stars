<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import BattleView from '$lib/components/game/battle/BattleView.svelte';
	import { player } from '$lib/services/Context';
	import { GameService } from '$lib/services/GameService';
	import { onMount } from 'svelte';

	let id = $page.params.id;
	let num = parseInt($page.params.num);

	onMount(async () => {
		const p = await GameService.loadFullPlayer(id);
		player.update(() => p);
	});

	$: battle = $player?.getBattle(num);
</script>

{#if $player && battle}
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a class="cs-link" href={`/games/${id}/battles`}>Battles</a></li>
			<li>{$player?.getBattleLocation(battle)}</li>
		</svelte:fragment>
	</Breadcrumb>

	<div class="grow px-1">
		<BattleView battleRecord={battle} player={$player} />
	</div>
{/if}

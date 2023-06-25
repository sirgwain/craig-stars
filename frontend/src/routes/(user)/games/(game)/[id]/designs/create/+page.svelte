<script lang="ts">
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import { techs } from '$lib/services/Stores';
	import { canLearnTech } from '$lib/types/Player';
	import { hasRequiredLevels } from '$lib/types/TechLevel';
	import { kebabCase } from 'lodash-es';

	const { game, player, universe } = getGameContext();
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li><a class="cs-link" href={`/games/${$game.id}/designs`}>Ship Designs</a></li>
		<li>Choose Hull</li>
	</svelte:fragment>
</Breadcrumb>
<ul class="px-1">
	{#each $techs.hulls as hull}
		{#if $player && canLearnTech($game.player, hull) && hasRequiredLevels($player.techLevels, hull.requirements)}
			<li>
				<a class="cs-link" href={`/games/${$game.id}/designs/create/${kebabCase(hull.name)}`}>
					<div class="flex flex-row place-items-center">
						<div class="mr-2 mb-2 border border-secondary bg-black p-1">
							<TechAvatar tech={hull} />
						</div>
						<div>
							{hull.name}
						</div>
					</div>
				</a>
			</li>
		{/if}
	{/each}
</ul>

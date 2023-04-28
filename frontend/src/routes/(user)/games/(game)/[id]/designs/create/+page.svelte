<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { player, techs } from '$lib/services/Context';
	import { canLearnTech, hasRequiredLevels } from '$lib/types/Player';
	import { kebabCase } from 'lodash-es';

	let gameId = $page.params.id;
</script>

<Breadcrumb>
	<svelte:fragment slot="crumbs">
		<li><a class="cs-link" href={`/games/${gameId}/designs`}>Designs</a></li>
		<li>Choose Hull</li>
	</svelte:fragment>
</Breadcrumb>
<ul class="px-1">
	{#each $techs.hulls as hull}
		{#if $player && canLearnTech($player, hull) && hasRequiredLevels($player.techLevels, hull.requirements)}
			<li>
				<a class="cs-link" href={`/games/${gameId}/designs/create/${kebabCase(hull.name)}`}>
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

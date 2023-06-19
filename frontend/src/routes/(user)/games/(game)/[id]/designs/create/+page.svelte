<script lang="ts">
	import { page } from '$app/stores';
	import TechAvatar from '$lib/components/tech/TechAvatar.svelte';
	import { player, techs } from '$lib/services/Context';
	import { canLearnTech, hasRequiredLevels } from '$lib/types/Player';
	import { kebabCase } from 'lodash-es';

	let gameId = $page.params.id;
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<div class="w-full flex justify-between gap-2 border-primary border-b-2 mb-2">
		<div>Hulls</div>
	</div>
	<ul class="px-1">
		{#each $techs.hulls as hull}
			{#if $player && canLearnTech($player, hull) && hasRequiredLevels($player.techLevels, hull.requirements)}
				<li>
					<div class="flex flex-row place-items-center">
						<div class="mr-2 mb-2">
							<TechAvatar tech={hull} />
						</div>
						<div>
							<a class="cs-link" href={`/games/${gameId}/designs/create/${kebabCase(hull.name)}`}
								>{hull.name}</a
							>
						</div>
					</div>
				</li>
			{/if}
		{/each}
	</ul>
</div>

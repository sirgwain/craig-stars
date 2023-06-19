<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/components/game/Breadcrumb.svelte';
	import Design from '$lib/components/game/design/Design.svelte';
	import { game } from '$lib/services/Context';
	import type { ErrorResponse } from '$lib/types/ErrorResponse';
	import type { ShipDesign } from '$lib/types/ShipDesign';

	let gameId = $page.params.id;
	let num = parseInt($page.params.num);

	$: design = $game && ($game.player.getDesign($game.player.num, num) as ShipDesign);

	let error: ErrorResponse | undefined;
</script>

{#if error?.error}
	{error.error}
{:else if design}
	<Breadcrumb>
		<svelte:fragment slot="crumbs">
			<li><a class="cs-link" href={`/games/${gameId}/designs`}>Designs</a></li>
			<li>{design?.name}</li>
			{#if !design.spec?.numInstances}
				<li><a class="cs-link" href={`/games/${gameId}/designs/${design.num}/edit`}>Edit</a></li>
			{/if}
		</svelte:fragment>
	</Breadcrumb>

	<div class="grow h-full px-1 md:p-0">
		<Design {design} />
	</div>
{/if}

<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import LoadingSpinner from './LoadingSpinner.svelte';

	const { game, player, universe } = getGameContext();

	export let text: string | undefined;
</script>

<div class="modal" class:modal-open={text !== undefined}>
	<div class="modal-box w-64 h-40 bg-base-100 rounded-sm border-2 border-secondary">
		<div class="flex flex-col justify-center h-full">
			<div class="mx-auto">
				<LoadingSpinner />
			</div>
			<div class="mx-auto mt-2">
				{text ?? ''}
			</div>
		</div>
	</div>
</div>

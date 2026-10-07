<script lang="ts">
	import InfoToast from '#lib/components/InfoToast.svelte';
	import { getGameContext } from '#lib/services/GameContext.js';
	import { Square2Stack } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';

	const { game } = getGameContext();

	let copiedText = $state('');
	let link = $derived(`${window.location.origin}/join-private-game/${$game.hash}`);
</script>

<InfoToast bind:text={copiedText} />

<div class="flex grow">
	<div class="w-full grow">
		<div class="cs-form-control">
			<label class="cs-form-label"
				><span class="cs-label-text w-32 text-right">Invite Link</span>

				<input class="input ml-2 grow" readonly value={link} />
			</label>
		</div>
	</div>
	<div>
		<div class="tooltip" data-tip="Copy Invite Link">
			<button
				onclick={() => {
					navigator.clipboard.writeText(link);
					copiedText = 'Copied invite link to clipboard';
				}}
				type="button"
				class="btn btn-outline my-2 normal-case"
				><Icon src={Square2Stack} size="24" class="stroke-success" /></button
			>
		</div>
	</div>
</div>

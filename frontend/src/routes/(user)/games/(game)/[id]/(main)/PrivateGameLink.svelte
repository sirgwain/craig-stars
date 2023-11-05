<script lang="ts">
	import InfoToast from '$lib/components/InfoToast.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import type { PlayerStatus } from '$lib/types/Player';
	import type { SessionUser } from '$lib/types/User';
	import { Square2Stack } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';

	const { game } = getGameContext();

	let copiedText = '';

	onMount(async () => {});

	$: link = `${window.location.origin}/join-private-game/${$game.hash}`;
</script>

<InfoToast bind:text={copiedText} />

<div class="flex flex-grow">
	<div class="w-full flex-grow">
		<div class="form-control">
			<label class="label"
				><span class="label-text w-32 text-right">Invite Link</span>

				<input class="input input-bordered ml-2 flex-grow" readonly value={link} />
			</label>
		</div>
	</div>
	<div>
		<div class="tooltip" data-tip="Copy Invite Link">
			<button
				on:click={() => {
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

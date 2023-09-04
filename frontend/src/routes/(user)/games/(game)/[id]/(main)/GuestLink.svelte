<script lang="ts">
	import InfoToast from '$lib/components/InfoToast.svelte';
	import { getGameContext } from '$lib/services/Contexts';
	import type { PlayerStatus } from '$lib/types/Player';
	import type { User } from '$lib/types/User';
	import { Square2Stack } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';

	const { game } = getGameContext();

	export let player: PlayerStatus;
	export let hideText = false;

	let guest: User | undefined;
	let copiedText = '';

	onMount(async () => {
		console.log(player);
		if (player.guest) {
			guest = await $game.loadGuest(player.num);
		}
	});

	$: link = `${window.location.origin}/auth/guest/${guest?.password}`;
</script>

{#if guest}
	<InfoToast bind:text={copiedText} />
	<div class="flex flex-row">
		<div class="my-auto grow" class:hidden={hideText}>
			<input class="input input-sm input-bordered w-full" readonly value={link} />
		</div>
		<div>
			<div
				class="tooltip"
				data-tip="Copy Invite Link"
			>
				<button
					on:click={() => {
						navigator.clipboard.writeText(link);
						copiedText = 'Copied invite link to clipboard';
					}}
					type="button"
					class="btn btn-outline btn-sm my-1 normal-case"
					><Icon src={Square2Stack} size="24" class="stroke-success" /></button
				>
			</div>
		</div>
	</div>
{/if}

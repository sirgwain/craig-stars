<script lang="ts">
	import InfoToast from '$lib/components/InfoToast.svelte';
	import type { GuestUser } from '$lib/types/cs-proto';
	import type { PlayerStatus } from '$lib/types/cs-proto';
	import { gameClient } from '$lib/services/connect';
	import { getGameContext } from '$lib/services/GameContext';
	import { Square2Stack } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';
	import { addError } from '$lib/services/Errors';
	import type { ConnectError } from '@connectrpc/connect';

	const { game } = getGameContext();

	type Props = {
		player: PlayerStatus;
		hideText?: boolean;
	};

	let { player, hideText = false }: Props = $props();

	let guest: GuestUser | undefined = $state();
	let copiedText = $state('');

	onMount(async () => {
		try {
			if (player.guest) {
				const resp = await gameClient.getGuestUser({ gameId: $game.id, playerNum: player.num });
				guest = resp.user;
			}
		} catch (e) {
			addError(e as ConnectError);
		}
	});

	let link = $derived(`${window.location.origin}/auth/guest/${guest?.hash}`);
</script>

{#if guest}
	<InfoToast bind:text={copiedText} />
	<div class="flex flex-row">
		<div class="my-auto grow" class:hidden={hideText}>
			<input class="input input-sm input-bordered w-full" readonly value={link} />
		</div>
		<div>
			<div class="tooltip" data-tip="Copy Invite Link">
				<button
					onclick={() => {
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

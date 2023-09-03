<script lang="ts">
	import { getGameContext } from '$lib/services/Contexts';
	import { me } from '$lib/services/Stores';
	import type { PlayerStatus } from '$lib/types/Player';
	import type { User } from '$lib/types/User';
	import { Square2Stack } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';

	const { game } = getGameContext();

	export let player: PlayerStatus;

	let guest: User | undefined;

	onMount(async () => {
		console.log(player);
		if (player.guest) {
			guest = await $game.loadGuest(player.num);
		}
	});

	$: link = `${window.location.origin}/auth/guest/${guest?.password}`;
</script>

{#if guest}
	<button
		on:click={() => navigator.clipboard.writeText(link)}
		type="button"
		class="btn btn-outline btn-sm my-1 normal-case"
		><Icon src={Square2Stack} size="24" class="stroke-success" /></button
	>
{/if}

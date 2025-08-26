<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import Select from '$lib/components/Select.svelte';
	import { adminClient } from '$lib/services/connect';
	import { UserRole, type GameWithPlayers, type User } from '$lib/types/cs-proto';
	import { onMount } from 'svelte';

	let users: User[] = $state([]);
	let games: GameWithPlayers[] = $state([]);
	let id = page.params.id;
	let guestUser: User | undefined = $state();
	let targetUserId: bigint | undefined = $state();

	onMount(async () => {
		try {
			const usersResp = await adminClient.getUsers({});
			users = usersResp.users;

			guestUser = usersResp.users.find((u) => u.id == BigInt(id));

			// load the games for this user
			if (guestUser) {
				const gamesResp = await adminClient.getUserGames({ userId: guestUser.id });
				games = gamesResp.games;
			}
		} catch (_err) {
			// TODO: show error
		}
	});

	async function onSubmit() {
		if (!targetUserId || !guestUser) {
			return;
		}
		await adminClient.convertGuestUser({ guestUserId: guestUser.id, userId: targetUserId });
		await goto(`/admin/users`);
	}
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<form
		onsubmit={(e) => {
			e.preventDefault();
			onSubmit();
		}}
	>
		<div class="w-full flex justify-end gap-2">
			<button class="btn btn-success" type="submit">Convert</button>
		</div>

		{#if guestUser}
			<ItemTitle>{guestUser.username}</ItemTitle>

			<Select
				values={users
					.filter((u) => u.role !== UserRole.GUEST)
					.map((u) => {
						return { value: u.id, title: u.username };
					})}
				name="Target User"
				bind:value={targetUserId}
			/>

			{#if games}
				<SectionHeader>Guest User Games</SectionHeader>
				<ul>
					{#each games as game (game.game?.id)}
						<li>{game.game?.name} - {game.players.length} players</li>
					{/each}
				</ul>
			{/if}
		{/if}
	</form>
</div>

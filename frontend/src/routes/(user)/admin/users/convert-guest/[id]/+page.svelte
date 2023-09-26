<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import SectionHeader from '$lib/components/SectionHeader.svelte';
	import Select from '$lib/components/Select.svelte';
	import { AdminService } from '$lib/services/AdminService';
	import { Service } from '$lib/services/Service';
	import type { Game } from '$lib/types/Game';
	import type { User } from '$lib/types/User';
	import { onMount } from 'svelte';

	let users: User[];
	let games: Game[];
	let id = $page.params.id;
	let guestUser: User | undefined;
	let targetUserId: number | undefined;

	onMount(async () => {
		try {
			users = await AdminService.loadUsers();

			guestUser = users.find((u) => u.id == parseInt(id));

			// load the games for this user
			if (guestUser) {
				games = await AdminService.loadUserGames(guestUser.id);
			}
		} catch (err) {
			// TODO: show error
		}
	});

	async function onSubmit() {
		const body = JSON.stringify({ userId: targetUserId });
		const response = await fetch(`/api/admin/users/${id}/convert-guest-user`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			},
			body
		});

		if (!response.ok) {
			await Service.throwError(response);
		}

		await goto(`/admin/users`);
	}
</script>

<div class="w-full mx-auto md:max-w-2xl">
	<form on:submit|preventDefault={onSubmit}>
		<div class="w-full flex justify-end gap-2">
			<button class="btn btn-success" type="submit">Convert</button>
		</div>

		{#if guestUser}
			<ItemTitle>{guestUser?.username}</ItemTitle>

			<Select
				values={users
					.filter((u) => !u.isGuest())
					.map((u) => {
						return { value: u.id, title: u.username };
					})}
				name="Target User"
				bind:value={targetUserId}
			/>

			{#if games}
				<SectionHeader>Guest User Games</SectionHeader>
				<ul>
					{#each games as game}
						<li>{game.name} - {game.numPlayers} players</li>
					{/each}
				</ul>
			{/if}
		{/if}
	</form>
</div>

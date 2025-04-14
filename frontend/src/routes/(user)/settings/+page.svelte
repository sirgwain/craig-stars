<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import UserAvatar from '$lib/components/UserAvatar.svelte';
	import { notify } from '$lib/services/Notifications';
	import { Service } from '$lib/services/Service';
	import { me } from '$lib/services/Stores';
	import { UserService } from '$lib/services/UserService';
	import type { User } from '$lib/types/cs';
	import { onMount } from 'svelte';

	let user: User | undefined = $state();
	let changedWebhook = $state(false);

	onMount(async () => {
		user = await UserService.get($me?.id);
	});

	const onSubmit = async () => {
		if (!$me) {
			return;
		}
		const body = JSON.stringify({
			userSettings: {
				discordWebhookUrl: user?.discordWebhookUrl
			}
		});
		const response = await fetch(`/api/users/${$me?.id}`, {
			method: 'PUT',
			headers: {
				accept: 'application/json'
			},
			body
		});

		if (!response.ok) {
			await Service.throwError(response);
		}

		notify(`Saved ${user?.username} settings`);
		changedWebhook = false;
	};

	const testWebhook = async () => {
		if (!user || !user?.discordWebhookUrl) {
			return;
		}

		const response = await fetch(`/api/users/${user.id}/test-discord-webhook`, {
			method: 'POST',
			headers: {
				accept: 'application/json'
			}
		});

		if (!response.ok) {
			await Service.throwError(response);
		}

		notify(`Sent discord test`);
	};
</script>

<div>
	{#if user}
		<form
			onsubmit={(e) => {
				e.preventDefault();
				onSubmit();
			}}
		>
			<div class="w-full flex justify-end gap-2">
				<button class="btn btn-success" type="submit">Save</button>
			</div>

			<ItemTitle>{user.username}</ItemTitle>
			<div class="flex justify-center">
				<UserAvatar {user} />
			</div>
			<TextInput
				name="discordWebhookUrl"
				bind:value={user.discordWebhookUrl}
				onchange={() => (changedWebhook = true)}
			/>
			<div class="flex justify-end">
				<button
					class="btn btn-secondary"
					disabled={!user.discordWebhookUrl || changedWebhook}
					type="button"
					onclick={testWebhook}>Test</button
				>
			</div>
		</form>
	{/if}
</div>

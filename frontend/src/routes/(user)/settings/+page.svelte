<script lang="ts">
	import ItemTitle from '$lib/components/ItemTitle.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import UserAvatar from '$lib/components/UserAvatar.svelte';
	import { userClient } from '$lib/services/connect';
	import { addError } from '$lib/services/Errors';
	import { notify } from '$lib/services/Notifications';
	import { me } from '$lib/services/Stores';
	import type { User } from '$lib/types/cs-proto';
	import { ConnectError } from '@connectrpc/connect';
	import { onMount } from 'svelte';

	let user: User | undefined = $state();
	let changedWebhook = $state(false);

	onMount(async () => {
		const resp = await userClient.getUser({ userId: $me.id });
		user = resp.user;
	});

	const onSubmit = async () => {
		if (!$me) {
			return;
		}

		try {
			await userClient.updateUserSettings({
				userSettings: {
					discordWebhookUrl: user?.discordWebhookUrl
				}
			});

			notify(`Saved ${user?.username} settings`);
			changedWebhook = false;
		} catch (err) {
			addError(err as ConnectError);
		}
	};

	const testWebhook = async () => {
		if (!user || !user?.discordWebhookUrl) {
			return;
		}

		try {
			await userClient.testDiscordWebhook({});
		} catch (err) {
			addError(err as ConnectError);
		}
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

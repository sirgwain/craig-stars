<script lang="ts">
	type Props = {
		user: { discordAvatar?: string; discordId?: string; username?: string };
	};

	let { user }: Props = $props();
</script>

{#if user.discordAvatar && user.discordAvatar != ''}
	<div class="rounded-full">
		{#if user.discordId && user.discordAvatar}
			<img
				src={`https://cdn.discordapp.com/avatars/${user.discordId}/${user.discordAvatar}`}
				alt="user avatar"
			/>
		{/if}
	</div>
{:else}
	<div class="avatar placeholder select-none">
		<div class="bg-neutral-focus text-neutral-content rounded-full border">
			<span class="text-2xl uppercase">{user.username ? user.username[0] : '?'}</span>
		</div>
	</div>
{/if}

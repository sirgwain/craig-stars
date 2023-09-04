<script lang="ts">
	import type { User } from '$lib/types/User';

	export let user: User;
</script>

{#if user.discordAvatar && user.discordAvatar != ''}
	<div class="rounded-full">
		<img
			src={`https://cdn.discordapp.com/avatars/${user.discordId}/${user.discordAvatar}`}
			alt="user avatar"
		/>
	</div>
{:else}
	<div class="avatar placeholder select-none">
		<div class="bg-neutral-focus text-neutral-content rounded-full border">
			<span class="text-2xl uppercase">{user.username[0]}</span>
		</div>
	</div>
{/if}

<script lang="ts">
	import type { User } from '$lib/types/User';
	import { Bars3 } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import DarkModeToggler from './DarkModeToggler.svelte';
	import UserAvatar from './UserAvatar.svelte';
	import Discord from './icons/Discord.svelte';
	import GitHub from './icons/GitHub.svelte';

	export let user: User | undefined;
</script>

<!-- svelte-ignore a11y-no-noninteractive-tabindex -->
<!-- svelte-ignore a11y-label-has-associated-control -->
<div class="navbar bg-base-100 flex flex-row">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl text-accent" href="/"
			><span class="hidden md:block">craig-stars</span><span class="md:hidden">cs</span></a
		>
	</div>
	<div class="flex-none">
		<div class="hidden md:block flex-none">
			<ul class="menu menu-horizontal p-0">
				{#if user}
					<li><a href="/games">Games</a></li>
					<li><a href="/races">Races</a></li>
				{/if}
				<li><a href="/techs">Techs</a></li>
			</ul>
		</div>

		{#if user}
			<div class="dropdown dropdown-end">
				<label class="avatar w-8 h-8 place-content-center">
					<UserAvatar {user} />
				</label>
			</div>
		{/if}

		<div class="dropdown dropdown-end">
			<label for="menu" tabindex="0" class="btn btn-ghost">
				<div id="menu">
					<Icon src={Bars3} size="24" />
				</div>
			</label>
			<ul
				tabindex="0"
				class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-30"
			>
				{#if user}
					<li>{user.username}</li>
				{/if}
				<li>
					<DarkModeToggler />
				</li>
				{#if user}
					<li class="md:hidden"><a href="/games">Games</a></li>
					<li class="md:hidden"><a href="/races">Races</a></li>
					<li class="md:hidden"><a href="/techs">Techs</a></li>
					{#if user.isAdmin()}
						<li><div class="divider" /></li>
						<li>
							<a href={`/admin/games`} class="justify-between">All Games</a>
						</li>
						<li>
							<a href={`/admin/users`} class="justify-between">Users</a>
						</li>
					{/if}
					<li><div class="divider" /></li>
					<li>
						<a href="https://discord.gg/Ctdx7h6UZS" target="_blank"
							><Discord class="fill-base-content w-5 h-5" /> Discord</a
						>
					</li>
					<li>
						<a href="https://github.com/sirgwain/craig-stars" target="_blank"
							><GitHub class="fill-base-content w-5 h-5" /> GitHub</a
						>
					</li>
					<li><div class="divider" /></li>
					<li><a href="/auth/logout">Logout, {user.username}</a></li>
					<li><div class="divider" /></li>
					<li class="text-center">version {PKG.version}</li>
					{:else}
					<li class="md:hidden"><a href="/techs">Techs</a></li>
				{/if}
			</ul>
		</div>
	</div>
</div>

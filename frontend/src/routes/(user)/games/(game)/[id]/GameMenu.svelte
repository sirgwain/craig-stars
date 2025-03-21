<script lang="ts">
	import { page } from '$app/stores';
	import DarkModeToggler from '$lib/components/DarkModeToggler.svelte';
	import UserAvatar from '$lib/components/UserAvatar.svelte';
	import { getGameContext } from '$lib/services/GameContext';
	import { me } from '$lib/services/Stores';
	import { GameStateWaitingForPlayers } from '$lib/types/cs';
	import { ArrowUpTray, Bars3 } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { onMount } from 'svelte';

	const { game, player } = getGameContext();

	type Props = {
		onSubmitTurn?: () => void;
	};

	let { onSubmitTurn }: Props = $props();

	const updateTitle = () => (document.title = `${$game.name} - ${$game.year}`);

	onMount(() => {
		// update the title of the page every time the game updates
		return game.subscribe(updateTitle);
	});
</script>

<div class="navbar bg-base-100 flex flex-row w-full">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl text-primary" href="/">cs</a>
		<div class="md:block">
			<a class="btn btn-ghost text-lg text-accent" href={`/games/${$game.id}`}
				>{$game.name} - {$game.year}</a
			>
		</div>
	</div>
	<div class="flex-initial">
		{#if $page.url.pathname === `/games/${$game.id}` && !$player.submittedTurn && $game.state === GameStateWaitingForPlayers}
			<button type="button" onclick={onSubmitTurn} class="btn btn-primary" title="submit turn"
				><span class="hidden md:inline-block mr-1">Submit Turn</span><Icon
					src={ArrowUpTray}
					size="16"
				/></button
			>
		{/if}

		{#if !$player.submittedTurn}
			<div class="hidden md:inline-block">
				<div class="dropdown">
					<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
					<label for="reports" tabindex="0" class="btn btn-ghost w-40">Commands</label>
					<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
					<ul
						id="commands"
						tabindex="0"
						class=" menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-300"
					>
						<li><a href={`/games/${$game.id}/research`}>Research</a></li>
						<li><a href={`/games/${$game.id}/designer`}>Ship Designer</a></li>
						<li><a href={`/games/${$game.id}/relations`}>Relations</a></li>
						<li><a href={`/games/${$game.id}/battle-plans`}>Battle Plans</a></li>
						<li><a href={`/games/${$game.id}/production-plans`}>Production Plans</a></li>
						<li><a href={`/games/${$game.id}/transport-plans`}>Transport Plans</a></li>
					</ul>
				</div>

				<div class="dropdown">
					<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
					<label for="reports" tabindex="0" class="btn btn-ghost">Reports</label>
					<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
					<ul
						id="reports"
						tabindex="0"
						class=" menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-300"
					>
						<li><a href={`/games/${$game.id}/players`}>Players</a></li>
						<li><a href={`/games/${$game.id}/planets`}>Planets</a></li>
						<li><a href={`/games/${$game.id}/fleets`}>Fleets</a></li>
						<li><a href={`/games/${$game.id}/designs`}>Designs</a></li>
						<li><a href={`/games/${$game.id}/messages`}>Messages</a></li>
						<li><a href={`/games/${$game.id}/battles`}>Battles</a></li>
					</ul>
				</div>
			</div>
		{/if}

		<div class="dropdown">
			<label class="avatar w-8 h-8 place-content-center mx-1">
				<UserAvatar user={$me} />
			</label>
		</div>

		<div class="dropdown dropdown-end">
			<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
			<label for="menu" tabindex="0" class="btn btn-ghost">
				<div id="menu">
					<Icon src={Bars3} size="24" />
				</div>
			</label>
			<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
			<div
				tabindex="0"
				class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-300 w-[22rem] md:w-auto"
			>
				<div class="flex flex-row justify-between">
					<ul class="mt-11">
						{#if !$player.submittedTurn}
							<li class="md:hidden menu-title">
								<span>Commands</span>
							</li>
							<li class="md:hidden"><a href={`/games/${$game.id}/research`}>Research</a></li>
							<li class="md:hidden"><a href={`/games/${$game.id}/designer`}>Ship Designer</a></li>
							<li class="md:hidden"><a href={`/games/${$game.id}/relations`}>Relations</a></li>
							<li class="md:hidden">
								<a href={`/games/${$game.id}/battle-plans`}>Battle Plans</a>
							</li>
							<li class="md:hidden">
								<a href={`/games/${$game.id}/production-plans`}>Production Plans</a>
							</li>
							<li class="md:hidden">
								<a href={`/games/${$game.id}/transport-plans`}>Transport Plans</a>
							</li>

							<li class="md:hidden menu-title">
								<span>Reports</span>
							</li>
							<li class="md:hidden"><a href={`/games/${$game.id}/players`}>Players</a></li>
							<li class="md:hidden"><a href={`/games/${$game.id}/planets`}>Planets</a></li>
							<li class="md:hidden"><a href={`/games/${$game.id}/fleets`}>Fleets</a></li>
							<li class="md:hidden"><a href={`/games/${$game.id}/designs`}>Designs</a></li>
							<li class="md:hidden"><a href={`/games/${$game.id}/messages`}>Messages</a></li>
							<li class="md:hidden"><a href={`/games/${$game.id}/battles`}>Battles</a></li>
						{/if}
					</ul>
					<ul>
						<li>
							<DarkModeToggler />
						</li>
						<li class="md:hidden menu-title">
							<span>Game</span>
						</li>
						<li>
							<a href={`/games/${$game.id}/race`} class="justify-between">Race</a>
						</li>
						<li>
							<a href={`/games/${$game.id}/techs`} class="justify-between">Techs</a>
						</li>
						{#if $me.isAdmin()}
							<li><div class="divider"></div></li>
							<li>
								<a href="/admin/games" class="justify-between">All Games</a>
							</li>
							<li>
								<a href="/admin/users" class="justify-between">Users</a>
							</li>
						{/if}
						<li><div class="divider"></div></li>
						<li><a href="/auth/logout">Logout, {$me.username}</a></li>
						<li><div class="divider"></div></li>
						<li class="text-center">version {PKG.version}</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
</div>

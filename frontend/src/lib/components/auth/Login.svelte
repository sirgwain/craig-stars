<script lang="ts">
	import DiscordLink from './DiscordLink.svelte';

	const onSubmit = async () => {
		const data = JSON.stringify({ user, passwd });

		const response = await fetch(`/api/auth/local/login?session=1`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				accept: 'application/json'
			},
			body: data
		});

		if (response.ok) {
			document.location = '/';
		} else {
			const resolvedResponse = await response?.json();
			loginError = resolvedResponse.error;
			console.error(loginError);
		}
	};

	let user = '';
	let passwd = '';
	$: loginError = '';
	let showAdmin = false;
</script>

<div class="flex flex-col p-2">
	{#if !showAdmin}
		<div>
			<DiscordLink />
		</div>
	{/if}
	<div>
		<button class="cs-link" on:click={() => (showAdmin = !showAdmin)}>
			{#if showAdmin}
				I'm not an admin
			{:else}
				I'm an admin
			{/if}
		</button>
	</div>
	{#if showAdmin}
		<div class="text-left mx-auto">
			<!-- content here -->
			<form on:submit|preventDefault={onSubmit}>
				<label class="label block">
					<span class="label-text">Username</span>
					<input bind:value={user} required type="text" name="user" class="input input-bordered" autocapitalize="off" />
				</label>

				<label class="label block">
					<span class="label-text">Password</span>
					<input
						bind:value={passwd}
						required
						type="password"
						name="passwd"
						class="input input-bordered"
					/>
				</label>
				<button class="btn btn-primary" type="submit">Submit</button>
			</form>
			{#if loginError}
				<div class="text-red-600">{loginError}</div>
			{/if}
		</div>
	{/if}
</div>

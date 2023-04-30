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
</script>

<div class="flex items-center justify-center min-h-screen card">
	<div class="px-8 py-6 mt-4 bg-base-200 shadow-xl">
		<h2 class="text-2xl card-title">Login</h2>
		<div class="card-body">
			<form on:submit|preventDefault={onSubmit}>
				<label class="label block">
					<div class="label-text">Username</div>
					<input bind:value={user} required type="text" name="user" class="input input-bordered" />
				</label>

				<label class="label block">
					<div class="label-text">Password</div>
					<input
						bind:value={passwd}
						required
						type="password"
						name="passwd"
						class="input input-bordered"
					/>
				</label>
				<button class="btn btn-primary" type="submit">Submit</button>
				<DiscordLink />
			</form>
			{#if loginError}
				<div class="text-red-600">{loginError}</div>
			{/if}
		</div>
	</div>
</div>

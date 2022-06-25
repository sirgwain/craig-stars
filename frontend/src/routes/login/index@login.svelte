<script lang="ts">
	const onSubmit = async () => {
		const data = JSON.stringify({ username, password });

		const response = await fetch(`/api/login`, {
			method: 'post',
			headers: {
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

	let username = '';
	let password = '';
	$: loginError = '';
</script>

<div class="flex items-center justify-center min-h-screen ">
	<div
		class="px-8 py-6 mt-4 text-left shadow-lg rounded-md bg-white dark:bg-slate-800 border dark:border dark:border-slate-600 "
	>
		<h3 class="text-2xl font-bold text-center">Login</h3>

		<form on:submit|preventDefault={onSubmit}>
			<label class="block">
				<span class="label-text">Username</span>
				<input bind:value={username} autofocus required type="text" name="username" class="input" />
			</label>

			<label class="block">
				<span class="label-text">Password</span>
				<input bind:value={password} required type="password" name="password" class="input" />
			</label>
			<button class="submit-button" type="submit">Submit</button>
		</form>
		{#if loginError}
			<div class="text-red-600">{loginError}</div>
		{/if}
	</div>
</div>

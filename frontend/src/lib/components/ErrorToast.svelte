<script lang="ts">
	import { CSError, errors } from '$lib/services/Errors';
	import { fade } from 'svelte/transition';

	function onFadeOut(err: CSError) {
		$errors = $errors.filter((e) => e !== err);
	}
</script>

<div class="toast toast-top toast-center z-50 w-full md:max-w-2xl">
	{#each $errors as err}
		<div>
			<div
				class="alert alert-error"
				in:fade
				out:fade={{ delay: 3000 }}
				on:introend={() => onFadeOut(err)}
			>
				<div>
					<span>{err.error}</span>
				</div>
			</div>
		</div>
	{/each}
</div>

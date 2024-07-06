<script lang="ts">
	import { notification } from '$lib/services/Notifications';
	import { fade } from 'svelte/transition';

	function onFadeOut() {
		$notification = '';
	}
</script>

<div class="toast toast-top toast-center z-50 w-full md:max-w-2xl">
	{#if $notification}
		<div>
			<div class="alert alert-info" in:fade out:fade={{ delay: 1000 }} on:introend={onFadeOut}>
				<div>
					<span>{$notification}</span>
				</div>
			</div>
		</div>
	{/if}
</div>

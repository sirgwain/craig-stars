<script lang="ts">
	import { fade } from 'svelte/transition';

	export let text = '';

	function onFadeOut() {
		text = '';
	}
</script>

<div class="toast toast-top toast-center z-50 w-full md:max-w-2xl">
	{#if text}
		<div>
			<div
				class="alert alert-info"
				in:fade
				out:fade={{ delay: 3000 }}
				on:introend={() => onFadeOut()}
			>
				<div>
					<span>{text}</span>
				</div>
			</div>
		</div>
	{/if}
</div>

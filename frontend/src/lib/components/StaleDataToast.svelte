<script lang="ts">
	import { staleData } from '#lib/services/Stale.js';
</script>

{#if $staleData}
	<div class="toast toast-bottom toast-center z-50 w-full md:max-w-2xl">
		<div class="alert alert-warning" role="alert">
			<div>
				<p>This game has newer changes from another device or tab.</p>
				<p class="text-sm">
					Changes made here won't be saved until you reload to get the latest orders.
				</p>
			</div>
			<button class="btn btn-sm btn-primary" onclick={() => location.reload()}>Reload</button>
		</div>
	</div>
{/if}

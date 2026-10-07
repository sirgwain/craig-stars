<script lang="ts">
	import { versionUpdate } from '#lib/services/Version.js';

	let { frontendUpdated = false }: { frontendUpdated?: boolean } = $props();

	$effect(() => {
		if (frontendUpdated) versionUpdate.observeFrontendUpdate();
	});
</script>

{#if $versionUpdate.available && !$versionUpdate.dismissed}
	<div class="toast toast-bottom toast-center z-50 w-full md:max-w-2xl">
		<div class="alert alert-info" role="status">
			<div>
				<p>
					A new version of CraigStars! is available{#if $versionUpdate.serverVersion}
						({$versionUpdate.serverVersion}){/if}.
				</p>
				<p class="text-sm">Save your work before reloading.</p>
			</div>
			<div class="flex gap-2">
				<button class="btn btn-sm btn-primary" onclick={() => location.reload()}>Reload</button>
				<button class="btn btn-sm btn-ghost" onclick={() => versionUpdate.dismiss()}>Later</button>
			</div>
		</div>
	</div>
{/if}

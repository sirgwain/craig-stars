<script lang="ts">
	export let title = '';
</script>

<div class="card bg-base-200 shadow-xl mb-1 w-[16rem] max-h-fit min-h-fit">
	<div class="card-body p-4 text-sm gap-0">
		<h2 class="text-xl font-semibold text-center mb-1 text-secondary">{title}</h2>
		<slot />
	</div>
</div>

<script lang="ts">
	export let title = '';
</script>

<div class="card bg-base-200 shadow-xl w-[14rem] max-h-fit rounded-sm border-2 border-base-300">
	<div class="card-body p-3 gap-0">
		<h2 class="text-lg font-semibold text-center mb-1 text-secondary">{title}</h2>
		<slot />
	</div>
</div>

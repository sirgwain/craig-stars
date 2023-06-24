<script lang="ts">
	
	export let error: string;
</script>

<div class="mx-auto z-50 w-full md:max-w-2xl">
	<div>
		<div class="alert alert-error">
			<div>
				<span>{error}</span>
			</div>
		</div>
	</div>
</div>

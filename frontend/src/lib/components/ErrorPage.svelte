<script lang="ts">
	interface Props {
		error: string;
	}

	let { error }: Props = $props();
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

<script lang="ts">
</script>

<div class="w-full flex flex-row justify-between border-primary border-b-2 mb-2">
	<div class="breadcrumbs ml-2">
		<ul>
			<slot name="crumbs" />
		</ul>
	</div>
	<slot name="end" />
</div>

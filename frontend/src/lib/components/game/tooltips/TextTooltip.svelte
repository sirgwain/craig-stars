<script lang="ts" context="module">
	export type TextTooltipProps = {
		text: string;
		clazz?: string;
	};
</script>

<script lang="ts">
	export let text: string;
	export let clazz: string;
</script>

<div class="flex flex-col sm:w-[26rem] m-auto ${clazz}">
	<div>
		{text}
	</div>
</div>

<script lang="ts" module>
	export type TextTooltipProps = {
		text: string;
		clazz?: string;
	};
</script>

<script lang="ts">
	interface Props {
		text: string;
		clazz: string;
	}

	let { text, clazz }: Props = $props();
</script>

<div class="flex flex-col sm:w-[26rem] m-auto ${clazz}">
	<div>
		{text}
	</div>
</div>

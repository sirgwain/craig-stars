<script lang="ts" module>
	export type TextTooltipProps = {
		text: string;
	};
</script>

<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';

	type Props = TextTooltipProps & HTMLAttributes<HTMLDivElement>;

	let { text, ...props }: Props = $props();
</script>

<div class="flex flex-col sm:w-[26rem] m-auto ${props.class ?? ''}">
	<div>
		{text}
	</div>
</div>

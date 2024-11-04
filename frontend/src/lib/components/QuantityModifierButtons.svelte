<script lang="ts">
	import { quantityModifier } from '$lib/quantityModifier';

	interface Props {
		modifier?: number;
	}

	let { modifier = $bindable(1) }: Props = $props();

	let buttonModifer = $state(1);

	function updateModifier(value: number) {
		modifier = value;
		buttonModifer = modifier;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.metaKey || e.shiftKey || e.ctrlKey) {
			modifier = quantityModifier(e);
		}
	}

	function handleKeyup(e: KeyboardEvent) {
		const keyModifer = quantityModifier(e);

		if (keyModifer === 1) {
			modifier = buttonModifer;
		} else {
			modifier = keyModifer;
		}
	}
</script>

<!-- watch for key events to account for quantityModifier changes -->
<svelte:window onkeydown={handleKeydown} onkeyup={handleKeyup} />

<button
	class:btn-primary={modifier == 1}
	class="btn btn-xs border-secondary normal-case rounded-full"
	onclick={() => updateModifier(1)}
	>x1
</button>
<button
	class:btn-primary={modifier == 10 || modifier == 1000}
	class="btn btn-xs border-secondary normal-case rounded-full"
	onclick={() => updateModifier(10)}
	>x10
</button>
<button
	class:btn-primary={modifier == 100 || modifier == 1000}
	class="btn btn-xs border-secondary normal-case rounded-full"
	onclick={() => updateModifier(100)}
	>x100
</button>

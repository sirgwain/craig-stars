<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import { BattleAttackWhos, BattleTactics, BattleTargets } from '$lib/types/Battle';
	import type { BattlePlan } from '$lib/types/cs';

	type Props = {
		plan: BattlePlan;
	};

	let { plan = $bindable() }: Props = $props();
</script>

<TextInput name="name" bind:value={plan.name} required />
<EnumSelect name="primaryTarget" options={BattleTargets} bind:value={plan.primaryTarget} />
<EnumSelect name="secondaryTarget" options={BattleTargets} bind:value={plan.secondaryTarget} />
<EnumSelect name="tactic" options={BattleTactics} bind:value={plan.tactic} />
<EnumSelect name="attackWho" options={BattleAttackWhos} bind:value={plan.attackWho} />

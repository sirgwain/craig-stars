<script lang="ts">
	import EnumSelect from '$lib/components/EnumSelect.svelte';
	import TextInput from '$lib/components/TextInput.svelte';
	import { BattleAttackWho, BattleTactic, BattleTarget } from '$lib/types/Battle';
	import type { BattlePlan } from '$lib/types/Player';

	export let plan: BattlePlan;
</script>

<TextInput name="name" bind:value={plan.name} required />
<EnumSelect name="primaryTarget" enumType={BattleTarget} bind:value={plan.primaryTarget} />
<EnumSelect name="secondaryTarget" enumType={BattleTarget} bind:value={plan.secondaryTarget} />
<EnumSelect name="tactic" enumType={BattleTactic} bind:value={plan.tactic} />
<EnumSelect name="attackWho" enumType={BattleAttackWho} bind:value={plan.attackWho} />

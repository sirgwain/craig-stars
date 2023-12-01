import { HabTypes, type Hab, type HabType, getHabValue } from '$lib/types/Hab';
import type { Player } from '$lib/types/Player';
import { TerraformHabTypes, type TerraformHabType, type TechStore } from '$lib/types/Tech';

export function fromHabType(habType: HabType): TerraformHabType {
	switch (habType) {
		case HabTypes.Gravity:
			return TerraformHabTypes.Gravity;
		case HabTypes.Temperature:
			return TerraformHabTypes.Temperature;
		case HabTypes.Radiation:
			return TerraformHabTypes.Radiation;
	}
}

// getTerraformAmount returns the total amount we can terraform this planet
export function getTerraformAmount(
	techStore: TechStore,
	hab: Hab,
	baseHab: Hab,
	player: Player,
	terraformer?: Player
): Hab {
	const terraformAmount: [number, number, number] = [0, 0, 0];

	// default the terraformer to the player
	if (!terraformer) {
		terraformer = player;
	}

	const isEnemy = terraformer.isEnemy(player.num);

	const ta = terraformer.getTerraformAbility(techStore);
	const terraformAbility: [number, number, number] = [ta.grav ?? 0, ta.temp ?? 0, ta.rad ?? 0];
	const habCenter: [number, number, number] = [
		player.race.spec?.habCenter?.grav ?? 0,
		player.race.spec?.habCenter?.temp ?? 0,
		player.race.spec?.habCenter?.rad ?? 0
	];

	const immune: [boolean, boolean, boolean] = [
		player.race.immuneGrav ?? false,
		player.race.immuneTemp ?? false,
		player.race.immuneRad ?? false
	];

	Object.values(HabTypes).forEach((habType) => {
		if (immune[habType]) {
			return;
		}

		const ability = terraformAbility[habType];
		if (ability === 0) {
			return;
		}

		// The distance from the starting hab of this planet
		const fromIdealBase = habCenter[habType] - getHabValue(baseHab, habType);

		// the distance from the current hab of this planet
		const fromIdeal = habCenter[habType] - getHabValue(hab, habType);

		// if we have any left to terraform
		if (fromIdeal > 0) {
			// i.e. our ideal is 50 and the planet hab is 47
			if (isEnemy) {
				const alreadyTerraformed = Math.abs(fromIdealBase - fromIdeal);
				terraformAmount[habType] = -(ability - alreadyTerraformed);
			} else {
				// we can either terrform up to our full ability, or however much
				// we have left to terraform on this
				const alreadyTerraformed = fromIdeal - fromIdealBase;
				terraformAmount[habType] = Math.min(ability - alreadyTerraformed, fromIdeal);
			}
		} else if (fromIdeal < 0) {
			if (isEnemy) {
				const alreadyTerraformed = Math.abs(fromIdealBase - fromIdeal);
				terraformAmount[habType] = ability - alreadyTerraformed;
			} else {
				// i.e. our ideal is 50 and the planet hab is 53
				const alreadyTerraformed = fromIdeal - fromIdealBase;
				terraformAmount[habType] = Math.max(-(ability - alreadyTerraformed), fromIdeal);
			}
		} else if (isEnemy) {
			// the terrformer is enemies with the player, terraform away from ideal
			const alreadyTerraformed = Math.abs(fromIdealBase - fromIdeal);
			terraformAmount[habType] = ability - alreadyTerraformed;
		}
	});

	return {
		grav: terraformAmount[0],
		temp: terraformAmount[1],
		rad: terraformAmount[2]
	};
}


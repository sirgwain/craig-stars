import { HabTypes, type Hab, type HabType, getHabValue } from '$lib/types/Hab';
import type { Player } from '$lib/types/Player';
import { getPlanetHabitability } from '$lib/types/Race';
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
				const alreadyTerraformed = fromIdealBase - fromIdeal;
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

// getMinTerraformAmount gets the minimum amount we need to terraform this planet to make it habitable (if we can terraform it at all)
export function getMinTerraformAmount(
	techStore: TechStore,
	hab: Hab,
	baseHab: Hab,
	player: Player
): Hab {
	const terraformAmount: [number, number, number] = [0, 0, 0];

	// can't terraform, return an empty hab
	if (!player) {
		return {};
	}

	// no min terraform, this planet is habitable
	const habValue = getPlanetHabitability(player.race, hab);
	if (habValue >= 0) {
		return {};
	}

	const ta = player.getTerraformAbility(techStore);
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

	const planetHab: [number, number, number] = [hab.grav ?? 0, hab.temp ?? 0, hab.rad ?? 0];
	const planetBaseHab: [number, number, number] = [
		baseHab.grav ?? 0,
		baseHab.temp ?? 0,
		baseHab.rad ?? 0
	];

	Object.values(HabTypes).forEach((habType) => {
		if (immune[habType]) {
			return;
		}

		const ability = terraformAbility[habType];
		if (ability === 0) {
			return;
		}

		// the distance from the current hab of this planet to our minimum hab threshold
		// If this is positive, it means we need to terraform a certain percent to get it in range
		let fromHabitableDistance = 0;
		const planetHabValue = planetHab[habType];
		const playerHabIdeal = habCenter[habType];
		if (planetHabValue > playerHabIdeal) {
			// this planet is higher that we want, check the upper bound distance
			// if the player's high temp is 85, and this planet is 83, we are already in range
			// and don't need to min-terraform. If the planet is 87, we need to drop it 2 to be in range.
			fromHabitableDistance = planetHabValue - getHabValue(player.race.habHigh, habType);
		} else {
			// this planet is lower than we want, check the lower bound distance
			// if the player's low temp is 15, and this planet is 17, we are already in range
			// and don't need to min-terraform. If the planet is 13, we need to increase it 2 to be in range.
			fromHabitableDistance = getHabValue(player.race.habLow, habType) - planetHabValue;
		}

		// if we are already in range, set this to 0 because we don't want to terraform anymore
		if (fromHabitableDistance < 0) {
			fromHabitableDistance = 0;
		}

		// if we have any left to terraform
		if (fromHabitableDistance > 0) {
			// the distance from the current hab of this planet
			// the distance from the current hab of this planet
			const fromIdeal = habCenter[habType] - planetHabValue;
			const fromIdealDistance = Math.abs(fromIdeal);

			// The distance from the starting hab of this planet
			const fromIdealBaseDistance = Math.abs(playerHabIdeal - planetBaseHab[habType]);

			// we can either terrform up to our full ability, or however much
			// we have left to terraform on this
			const alreadyTerraformed = fromIdealBaseDistance - fromIdealDistance;
			const terraformAmountPossible = Math.min(ability - alreadyTerraformed, fromIdealDistance);

			// if we are in range for this hab type, we won't terraform at all, otherwise return the max possible terraforming
			// left.
			terraformAmount[habType] = Math.min(fromHabitableDistance, terraformAmountPossible);
		}
	});

	return {
		grav: terraformAmount[0],
		temp: terraformAmount[1],
		rad: terraformAmount[2]
	};
}

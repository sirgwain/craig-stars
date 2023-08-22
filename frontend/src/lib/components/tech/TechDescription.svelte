<script lang="ts">
	import { defaultRules, type Rules } from '$lib/types/Game';

	import {
		getCloakPercentForCloakUnits,
		InfinteGate,
		ScanWithZeroRange,
		TechCategory,
		TerraformHabType,
		type Tech,
		type TechHull,
		type TechHullComponent,
		type TechPlanetaryScanner,
		type TechTerraform
	} from '$lib/types/Tech';
	import { onMount } from 'svelte';

	export let tech: Tech;
	export let rules: Rules = defaultRules;

	type Stat = {
		label: string;
		text: string;
	};

	let stats: Stat[] = [];
	let descriptions: string[] = [];
	let warnings: string[] = [];

	onMount(() => {
		if (tech.category == TechCategory.ShipHull || tech.category == TechCategory.StarbaseHull) {
			const hull = tech as TechHull;
			if (hull) {
				if (hull.fuelCapacity && hull.fuelCapacity > 0) {
					stats.push({ label: 'Fuel Capacity', text: `${hull.fuelCapacity}mg` });
				}
				if (hull.cargoCapacity && hull.cargoCapacity > 0) {
					stats.push({ label: 'Cargo Capacity', text: `${hull.cargoCapacity}kT` });
				}
				stats.push({ label: 'Armor Strength', text: hull.armor.toString() });
				hull.initiative && stats.push({ label: 'Initiative', text: hull.initiative.toString() });

				if (hull.repairBonus && hull.repairBonus > 0) {
					if (hull.starbase) {
						descriptions.push(
							`This starbase will repair fleets in orbit ${hull.repairBonus * 100}% faster.`
						);
					} else {
						descriptions.push(
							`This hull will help the fleet repair ${hull.repairBonus * 100}% faster.`
						);
					}
				}
			}
		}

		if (tech.category == TechCategory.PlanetaryScanner) {
			const planetaryScanner = tech as TechPlanetaryScanner;

			if (planetaryScanner.scanRange > 0) {
				descriptions.push(
					`Enemy fleets not orbiting a planet can be detected up to ${planetaryScanner.scanRange} light years away.`
				);
			}

			if (planetaryScanner.scanRangePen && planetaryScanner.scanRangePen > 0) {
				descriptions.push(
					`This scanner can determine a planet's basic stats from a distance up to ${planetaryScanner.scanRangePen} light years. The scanner will also spot enemy fleets attempting to hide behind planets within range.`
				);
			}
		}

		if (tech.category == TechCategory.Terraforming) {
			const terraform = tech as TechTerraform;
			if (terraform.habType == TerraformHabType.All) {
				descriptions.push(
					`Allows you to modify any of a planet's three environment variables up to ${terraform.ability}% from its original value`
				);
			} else {
				descriptions.push(
					`Allows you to modify a planet's ${terraform.habType} by up to ${terraform.ability}% from its original value`
				);
			}
		}

		if ('hullSlotType' in tech) {
			const hullComponent = tech as TechHullComponent;
			if (hullComponent) {
				if (hullComponent.category == TechCategory.MineLayer && hullComponent.mineFieldType) {
					const mineFieldStats = rules.mineFieldStatsByType[hullComponent.mineFieldType];
					stats.push({ label: 'Mines laid per year', text: `${hullComponent.mineLayingRate}` });
					stats.push({ label: 'Maximum safe speed', text: `${mineFieldStats.maxSpeed}` });
					stats.push({
						label: 'Chance/l.y. of a hit',
						text: `${(mineFieldStats.chanceOfHit * 100).toFixed(1)}%`
					});
					stats.push({
						label: 'Dmg done to each ship',
						text: `${mineFieldStats.damagePerEngine} (${mineFieldStats.damagePerEngineRS}) / engine`
					});
					stats.push({
						label: 'Min damage done to fleet',
						text: `${mineFieldStats.minDamagePerFleet} (${mineFieldStats.minDamagePerFleetRS})`
					});
					descriptions.push(
						'Numbers in parenthesis are for fleets containing a ship with ram scoop engines. Note that the chance of hitting a mine goes up the % listed for EACH warp you exceed the safe speed.'
					);
				}

				if (hullComponent.category == TechCategory.Shield && (hullComponent.armor ?? 0) > 0) {
					// if this is a shield with armor, it sounds cooler to make the armor a description
					descriptions.push(
						`This shield also contains an armor component which will absorb ${hullComponent.armor} damage points.`
					);
				} else if ((hullComponent.armor ?? 0) > 0) {
					stats.push({
						label: 'Armor Strength',
						text: hullComponent.armor ? hullComponent.armor.toString() : '0'
					});
				}

				if ((hullComponent.category == TechCategory.Armor && hullComponent.shield) ?? 0 > 0) {
					// if this is an armor with a shield, it sounds cooler to make the shield a description
					descriptions.push(
						`This armor also acts as part shield which will absorb ${hullComponent.shield} damage points.`
					);
				} else if ((hullComponent.shield ?? 0) > 0) {
					stats.push({
						label: 'Shield Strength',
						text: hullComponent.shield ? hullComponent.shield.toString() : ''
					});
				}

				if (hullComponent.power) {
					stats.push({ label: 'Power', text: `${hullComponent.power}` });
				}
				if (hullComponent.range || hullComponent.category == TechCategory.BeamWeapon) {
					stats.push({ label: 'Range', text: `${hullComponent.range ?? 0}` });
				}
				if (hullComponent.initiative) {
					stats.push({ label: 'Initiative', text: `${hullComponent.initiative}` });
				}
				if (hullComponent.accuracy) {
					stats.push({ label: 'Accuracy', text: `${hullComponent.accuracy}%` });
				}
				if (hullComponent.hitsAllTargets) {
					descriptions.push(`This weapon hits all targets in range each time it is fired.`);
				}
				if (hullComponent.gattling) {
					descriptions.push(
						`This weapon also makes an excellent mine sweeper, capable of sweeping ${
							(hullComponent.power ?? 0) * Math.pow(hullComponent.range ?? 0, 4)
						} mines per year.`
					);
				}
				if (hullComponent.damageShieldsOnly) {
					descriptions.push(`This weapon will only damage shields, it has no effect on armor.`);
				}

				if ((hullComponent.killRate ?? 0) > 0 && !hullComponent.orbitalConstructionModule) {
					// we have special text for orbital construction modules.
					descriptions.push(
						`This bomb will kill approimately ${hullComponent.killRate}% of a planet's populatation each year.`
					);
					if ((hullComponent.minKillRate ?? 0) > 0) {
						descriptions.push(
							`If a planet has no defenses, this bomb is guaranteed to kill at least ${hullComponent.minKillRate} colonists.`
						);
					}
					if (!hullComponent.structureDestroyRate) {
						descriptions.push("This bomb will not damage a planet's mines or factories.");
					}
				}

				if ((hullComponent.miningRate ?? 0) > 0) {
					descriptions.push(
						`This module contains robots capable of mining up to ${hullComponent.miningRate}kT of each mineral (depending on concentration) from an uninhabited planet the ship is orbiting. The fleet must have orders set to 'Remote Mining'.`
					);
				}

				if ((hullComponent.terraformRate ?? 0) > 0) {
					descriptions.push(
						`This modified mining robot terraforms inhabited planets by ${hullComponent.terraformRate} per year. It has a positive effect on friendly planets, a negative effect on neutral and enemy planets.`
					);

					if ((hullComponent.cloakUnits ?? 0) > 0) {
						descriptions.push(`It also provides ${hullComponent.cloakUnits}% cloaking.`);
					}
				}

				if ((hullComponent.structureDestroyRate ?? 0) > 0) {
					descriptions.push(
						`This bomb will destroy approximately ${hullComponent.structureDestroyRate} of a planet's mines, factories, and/or defenses each year.`
					);
				}

				if ((hullComponent.unterraformRate ?? 0) > 0) {
					descriptions.push(
						`This bomb does not kill colonists or destroy installations. This bomb 'unterraforms' planets toward their original state up to ${hullComponent.unterraformRate}% per variable per bombing run. Planetary defenses have no effect on this bomb.`
					);
				}

				if (
					hullComponent.cloakUnits &&
					hullComponent.cloakUnits > 0 &&
					!hullComponent.terraformRate
				) {
					if (hullComponent.cloakUnarmedOnly) {
						descriptions.push(
							`Cloaks unarmed hulls, reducing the range at which scanners detect it by up to ${getCloakPercentForCloakUnits(
								hullComponent.cloakUnits
							).toFixed()}%.`
						);
					} else {
						descriptions.push(
							`Cloaks any ship, reducing the range at which scanners detect it by up to ${getCloakPercentForCloakUnits(
								hullComponent.cloakUnits
							).toFixed()}%.`
						);
					}
				}

				if ((hullComponent.fuelBonus ?? 0) > 0) {
					descriptions.push(`This part acts as a ${hullComponent.fuelBonus}mg fuel tank.`);
				}

				if ((hullComponent.fuelRegenerationRate ?? 0) > 0) {
					descriptions.push(
						`This part generates ${hullComponent.fuelRegenerationRate}mg of fuel every year.`
					);
				}

				if (hullComponent.colonizationModule) {
					descriptions.push(
						'This pod allows a ship to colonize a planet and will dismantle the ship upon arrival and convert it into supplies for the colonists.'
					);
				}

				if (hullComponent.orbitalConstructionModule) {
					descriptions.push(
						'This module contains an empty orbital hull which can be deployed in orbit of an uninhabited planet.'
					);
					if ((hullComponent.minKillRate ?? 0) > 0) {
						descriptions.push(
							`This pod also contains viral weapons capable of killing ${hullComponent.minKillRate} enemy colonists per attack.`
						);
					}
				}

				if ((hullComponent.cargoBonus ?? 0) > 0) {
					descriptions.push(
						`This pod increases the cargo capacity of the ship by ${hullComponent.cargoBonus}kT`
					);
				}

				if ((hullComponent.movementBonus ?? 0) > 0) {
					descriptions.push(
						`Increases speed in battle by ${hullComponent.movementBonus} square of movement.`
					);
				}

				if (hullComponent.beamDefense && hullComponent.beamDefense > 0) {
					descriptions.push(
						`The deflector decreases damage done by beam weapons to this ship by up to ${(
							hullComponent.beamDefense * 100
						).toFixed()}%`
					);
				}

				if ((hullComponent.torpedoBonus ?? 0) > 0 || (hullComponent.initiativeBonus ?? 0) > 0) {
					if ((hullComponent.torpedoBonus ?? 0) > 0 && (hullComponent.initiativeBonus ?? 0) > 0) {
						descriptions.push(
							`This module increases the accuracy of your torpedos by ${
								(hullComponent.torpedoBonus ?? 0) * 100
							}% and increases your initiative by ${
								hullComponent.initiativeBonus
							}. If an enemy ship has jammers the computer acts to offset their effects.`
						);
					} else if ((hullComponent.initiativeBonus ?? 0) > 0) {
						descriptions.push(
							`This module increases your initiative by ${hullComponent.initiativeBonus}.`
						);
					} else if ((hullComponent.torpedoBonus ?? 0) > 0) {
						descriptions.push(
							`This module increases the accuracy of your torpedos by ${
								(hullComponent.torpedoBonus ?? 0) * 100
							}%. If an enemy ship has jammers the computer acts to offset their effects.`
						);
					}
				}

				if (hullComponent.torpedoJamming && hullComponent.torpedoJamming > 0) {
					descriptions.push(
						`Has a ${
							hullComponent.torpedoJamming * 100
						}% chance of deflecting incoming torpedos. Deflected torpedoes will still reduce shields (in any by 1/8 the damage value).`
					);
				}

				if (hullComponent.beamBonus && hullComponent.beamBonus > 0) {
					descriptions.push(
						`Increases the damage done by all beam weapons on their ship by ${
							hullComponent.beamBonus * 100
						}%.`
					);
				}

				if ((hullComponent.reduceMovement ?? 0) > 0) {
					descriptions.push(
						`Slows all ships in combat by ${hullComponent.reduceMovement} square of movement.`
					);
				}

				if (hullComponent.reduceCloaking) {
					descriptions.push(
						`Reduces the effectiveness of other players cloaks by ${rules.tachyonCloakReduction}%.`
					);
				}

				if ((hullComponent.safeRange ?? 0) > 0) {
					descriptions.push(
						'Allows fleets without cargo to jump to any other planet with a Stargate in a single year.'
					);
					stats.push({
						label: 'Safe hull mass',
						text:
							hullComponent.safeHullMass == InfinteGate
								? 'Unlimited'
								: `${hullComponent.safeHullMass}kT`
					});
					stats.push({
						label: 'Safe range',
						text:
							hullComponent.safeRange == InfinteGate
								? 'Unlimited'
								: `${hullComponent.safeRange} light years`
					});

					if (hullComponent.maxHullMass != InfinteGate && hullComponent.maxRange != InfinteGate) {
						warnings.push(
							`Warning: Ships up to ${hullComponent.maxHullMass}kT might be successfully gated up to ${hullComponent.maxRange} l.y. but exceeding the stated limits will cause damage to the fleet.`
						);
					} else if (hullComponent.maxHullMass != InfinteGate) {
						warnings.push(
							`Warning: Ships up to ${hullComponent.maxHullMass}kT can be successfully gated up but exceeding the stated limits will cause damage to the fleet.`
						);
					} else if (hullComponent.maxRange != InfinteGate) {
						warnings.push(
							`Warning: Ships can be successfully gated up to ${hullComponent.maxRange} l.y. but exceeding the stated limits will cause damage to the fleet.`
						);
					}
				}

				if ((hullComponent.packetSpeed ?? 0) > 0) {
					stats.push({ label: 'Warp', text: `${hullComponent.packetSpeed}` });
					descriptions.push('Allows planets to fling mineral packets at other planets.');
					warnings.push(
						'Warning: The receiving planet must have a mass driver at least as capable or it will take damage.'
					);
				}

				if ((hullComponent.scanRange ?? 0) > 0) {
					if ((hullComponent.scanRange ?? 0) > 0) {
						if (hullComponent.scanRange == ScanWithZeroRange) {
							// special case for bat scanner
							descriptions.push(
								'Enemy fleets cannot be detected boy this scanner unless they are at the same location as the scanner.'
							);
						} else {
							descriptions.push(
								`Enemy fleets not orbiting a planet can be detected up to ${hullComponent.scanRange} light years away.`
							);
						}

						if (!hullComponent.scanRangePen) {
							// we have no pen scan, but we are a normal scanner, we can still scan planets we orbit
							descriptions.push(
								"This scanner is capable of determining a planet's environment and composition while in orbit of the planet."
							);
						}
					}

					if ((hullComponent.scanRangePen ?? 0) > 0) {
						descriptions.push(
							`This scanner can determine a planet's basic stats from a distance up to ${hullComponent.scanRangePen} light years. The scanner will also spot enemy fleets attempting to hide behind planets within range.`
						);
					}

					if (hullComponent.canStealFleetCargo || hullComponent.canStealPlanetCargo) {
						let target = '';
						if (hullComponent.canStealFleetCargo && !hullComponent.canStealPlanetCargo) {
							target = 'fleets';
						} else if (!hullComponent.canStealFleetCargo && hullComponent.canStealPlanetCargo) {
							target = 'planets';
						} else {
							target = 'fleets and planets';
						}
						descriptions.push(
							`This scanner is capable of penetrating the defenses of enemy ${target} allowing you to steal their cargo.`
						);
					}
				}
			}
		}

		stats = stats;
		descriptions = descriptions;
		warnings = warnings;
	});
</script>

<div class="flex flex-col p-1">
	{#each stats as stat}
		<div class="flex">
			<div class="w-1/2 text-right font-semibold">{stat.label}:</div>
			<div class="w-1/2 text-left ml-2">{stat.text}</div>
		</div>
	{/each}

	<div class="mt-1" />
	{#each descriptions as description}
		<div>{description}</div>
	{/each}
	{#each warnings as warning}
		<div class="text-warning">{warning}</div>
	{/each}
</div>

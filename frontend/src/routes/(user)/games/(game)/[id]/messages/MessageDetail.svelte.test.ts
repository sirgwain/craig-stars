import { create } from '@bufbuild/protobuf';
import { page } from 'vitest/browser';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import {
	MapObjectType,
	PlayerMessageSchema,
	PlayerMessageTargetType,
	PlayerMessageType as Type,
	ResourceType,
	TechCategory,
	TerraformHabType
} from '#lib/types/cs-proto.js';
import MessageDetail from './MessageDetail.svelte';

const state = vi.hoisted(() => ({ category: 0 }));
vi.mock('#lib/services/GameContext.js', async () => {
	const { writable } = await import('svelte/store');
	const context = {
		game: writable({ id: 1n }),
		player: writable({
			num: 1,
			getAllies: () => [1],
			race: { growthRate: 15, spec: { growthFactor: 1 } }
		}),
		universe: writable({
			// Merged fleets and landed packets are absent from the current universe.
			getMapObject: () => undefined,
			getFleet: () => undefined,
			getPlayerIntel: () => undefined,
			getPlayerName: () => 'Private race notes',
			getBattle: () => undefined,
			getBattleLocation: () => 'Earth',
			getPlayerPluralName: () => 'Visitors',
			battleRecords: []
		})
	};
	return { getGameContext: () => context };
});
vi.mock('#lib/services/Stores.js', async () => {
	const { writable } = await import('svelte/store');
	return { techs: writable({ getTech: () => ({ tech: { category: state.category } }) }) };
});
const fleetTarget = {
	targetType: PlayerMessageTargetType.FLEET,
	targetPlayerNum: 1,
	targetNum: 7,
	targetName: 'Voyager'
};
const planetTarget = {
	targetType: PlayerMessageTargetType.PLANET,
	targetNum: 3,
	targetName: 'Earth'
};
const references = {
	mapObjectTarget: { targetType: MapObjectType.PLANET, targetName: 'Earth', targetNum: 3 },
	routeTarget: { targetType: MapObjectType.PLANET, targetName: 'Mars', targetNum: 4 }
};
beforeEach(() => {
	state.category = TechCategory.UNSPECIFIED;
});

describe('structured fleet messages', () => {
	it.each([
		[Type.FLEET_ORDERS_COMPLETE, 'completed its assigned orders'],
		[Type.FLEET_ENGINE_FAILURE, 'unable to engage its engines'],
		[Type.FLEET_MERGED, 'Original has been merged into Voyager.'],
		[Type.FLEET_MERGE_INVALID_NOT_FLEET, "waypoint destination wasn't a fleet"],
		[Type.FLEET_MERGE_INVALID_UNOWNED, "destination fleet wasn't one of yours"],
		[Type.FLEET_ROUTE_INVALID_NOT_PLANET, 'not orbiting a planet'],
		[Type.FLEET_ROUTE_INVALID_NOT_FRIENDLY_PLANET, 'inhabitants of Earth'],
		[Type.FLEET_ROUTE_INVALID_NO_ROUTE_TARGET, 'at Earth as the planet has no route set'],
		[Type.FLEET_ROUTE, 'routed by the citizens of Earth to Mars'],
		[Type.FLEET_COLONIZE_INVALID_NOT_PLANET, 'not currently orbiting a planet'],
		[Type.FLEET_COLONIZE_INVALID_OWNED_PLANET, 'Earth is already populated'],
		[Type.FLEET_COLONIZE_INVALID_NO_MODULE, 'without a colonization module'],
		[Type.FLEET_COLONIZE_INVALID_NO_COLONISTS, 'failed to bring along any colonists'],
		[Type.FLEET_LAY_MINES_INVALID_NO_MINE_LAYERS, 'has no mine layers'],
		[
			Type.FLEET_REMOTE_MINE_INVALID_NO_MINERS,
			"remote mine Earth, but the fleet doesn't have any remote mining modules"
		],
		[
			Type.FLEET_REMOTE_MINE_INVALID_INHABITED,
			'remote mine Earth, but the planet is already inhabited'
		],
		[Type.FLEET_REMOTE_MINE_INVALID_DEEP_SPACE, 'remote mine in deep space'],
		[Type.FLEET_STARGATE_INVALID_SOURCE, 'no stargate exists there'],
		[Type.FLEET_STARGATE_INVALID_SOURCE_OWNER, 'starbase is not owned by you or your allies'],
		[Type.FLEET_STARGATE_INVALID_DEST, 'no stargate could be detected at the destination'],
		[
			Type.FLEET_STARGATE_INVALID_DEST_OWNER,
			'destination starbase is not owned by you or your allies'
		],
		[Type.FLEET_STARGATE_INVALID_MASS, "far too massive for the gate's limits"],
		[Type.FLEET_STARGATE_INVALID_COLONISTS, "carrying colonists and can't drop them off"],
		[Type.FLEET_STARGATE_DESTROYED, 'The fleet never arrived.']
	] as const)('renders type %s without a live fleet', async (type, expected) => {
		render(MessageDetail, {
			message: create(PlayerMessageSchema, {
				type,
				target: fleetTarget,
				spec: { ...references, name: 'Original' }
			})
		});
		await expect.element(page.getByText(expected, { exact: false })).toBeInTheDocument();
	});
	it('preserves warp zero and fractional distances', async () => {
		const view = render(MessageDetail, {
			message: create(PlayerMessageSchema, {
				type: Type.FLEET_OUT_OF_FUEL,
				target: fleetTarget,
				spec: { amount: 0 }
			})
		});
		await expect.element(page.getByText('Warp 0.', { exact: false })).toBeInTheDocument();
		await view.rerender({
			message: create(PlayerMessageSchema, {
				type: Type.FLEET_STARGATE_INVALID_RANGE,
				target: fleetTarget,
				spec: { ...references, distance: 123.45 }
			})
		});
		await expect
			.element(page.getByText('distance of 123.5 Ly.', { exact: false }))
			.toBeInTheDocument();
	});
	it.each([
		[ResourceType.COLONISTS, -12, 'beamed 1,200 colonists from Earth'],
		[ResourceType.COLONISTS, 12, 'beamed 1,200 colonists to Earth'],
		[ResourceType.FUEL, -12, 'loaded 12mg of Fuel from Earth'],
		[ResourceType.IRONIUM, 12, 'unloaded 12kt of Ironium to Earth']
	] as const)('formats cargo %s with signed amount %s', async (cargoType, transfered, expected) => {
		render(MessageDetail, {
			message: create(PlayerMessageSchema, {
				type: Type.FLEET_TRANSFERRED_CARGO,
				target: fleetTarget,
				spec: { ...references, cargoTransfer: { cargoType, transfered } }
			})
		});
		await expect.element(page.getByText(expected, { exact: false })).toBeInTheDocument();
	});
	it.each([
		[0, 'suffering 25 dp of damage'],
		[1, 'only 1 ship to the treacherous void'],
		[4, 'only 4 ships to the treacherous void'],
		[5, '5 ships to the unforgiving void'],
		[10, '10 ships to the unforgiving void'],
		[11, '11 ships to the great unknown'],
		[50, '50 ships to the great unknown'],
		[51, 'unbelievable 51 ships to the cosmic ocean']
	] as const)('renders stargate loss boundary %s', async (amount2, expected) => {
		render(MessageDetail, {
			message: create(PlayerMessageSchema, {
				type: Type.FLEET_STARGATE_DAMAGED,
				target: fleetTarget,
				spec: { ...references, amount: 25, amount2 }
			})
		});
		await expect.element(page.getByText(expected, { exact: false })).toBeInTheDocument();
	});
	it.each([
		[{ colonists: 12 }, '1,200 colonists'],
		[{ ironium: 3, boranium: 4 }, '7kt of minerals'],
		[{ colonists: 12, ironium: 3 }, '1,200 colonists and 3kt of minerals']
	])('renders unloaded stargate cargo %j', async (cargo, expected) => {
		render(MessageDetail, {
			message: create(PlayerMessageSchema, {
				type: Type.FLEET_DUMPED_CARGO,
				target: fleetTarget,
				spec: { ...references, cargo }
			})
		});
		await expect.element(page.getByText(expected, { exact: false })).toBeInTheDocument();
	});
	it.each([MapObjectType.FLEET, MapObjectType.WORMHOLE])(
		'renders lost target %s',
		async (lostTargetType) => {
			render(MessageDetail, {
				message: create(PlayerMessageSchema, {
					type: Type.FLEET_TARGET_LOST,
					target: fleetTarget,
					spec: { name: 'Lost Object', lostTargetType }
				})
			});
			await expect
				.element(
					page.getByText(
						lostTargetType === MapObjectType.FLEET
							? 'outrun the range of your scanners'
							: 'appears to have disappeared',
						{ exact: false }
					)
				)
				.toBeInTheDocument();
			await expect.element(page.getByText('Lost Object', { exact: false })).toBeInTheDocument();
		}
	);
});

describe('structured planet messages', () => {
	it.each([
		[Type.PLANET_COLONIZED, 'now in control of Earth'],
		[Type.PLANET_INVADE_INVALID_EMPTY, 'the planet is uninhabited'],
		[Type.PLANET_INVADE_INVALID_STARBASE, 'the planet is protected by a starbase'],
		[Type.PLANET_PACKET_LANDED, '100kT of minerals has arrived at Earth'],
		[Type.PLANET_PACKET_CAUGHT, 'captured a packet containing 100kT']
	] as const)('dispatches type %s without a live planet', async (type, expected) => {
		render(MessageDetail, {
			message: create(PlayerMessageSchema, {
				type,
				target: planetTarget,
				spec: { amount: 100, mapObjectTarget: { targetName: 'Voyager' } }
			})
		});
		await expect.element(page.getByText(expected, { exact: false })).toBeInTheDocument();
	});
	it.each([
		[true, false, 0, 'Unable to completely slow the packet'],
		[true, true, 3, 'Unfortunately, 1,200 of your colonists and 3 of your defenses'],
		[false, true, 0, 'Earth was annihilated by a mineral packet'],
		[false, false, 0, '1,200 of your colonists were killed'],
		[false, false, 3, '1,200 of your colonists and 3 of your defenses']
	] as const)(
		'renders packet flags %s/%s, defenses %s',
		async (hasMassDriver, planetEmptied, defensesDestroyed, expected) => {
			render(MessageDetail, {
				message: create(PlayerMessageSchema, {
					type: Type.PLANET_PACKET_DAMAGE,
					target: planetTarget,
					spec: {
						amount: 100,
						hasMassDriver,
						planetEmptied,
						mineralPacketDamage: { killed: 1200, defensesDestroyed }
					}
				})
			});
			await expect.element(page.getByText(expected, { exact: false })).toBeInTheDocument();
		}
	);
	for (const type of [
		Type.PLANET_BUILT_TERRAFORM,
		Type.PLANET_PERMAFORM,
		Type.PLANET_PACKET_TERRAFORM,
		Type.PLANET_PACKET_PERMAFORM
	]) {
		it.each([
			[TerraformHabType.GRAV, 0, '0.12g'],
			[TerraformHabType.TEMP, 60, '40°C'],
			[TerraformHabType.RAD, 70, '70mR']
		] as const)(
			`renders axis %s for type ${type} from its snapshot`,
			async (habType, amount2, expected) => {
				const view = render(MessageDetail, {
					message: create(PlayerMessageSchema, {
						type,
						target: planetTarget,
						spec: { habType, amount: 1, amount2 }
					})
				});
				await expect
					.element(page.getByText(`to ${expected}`, { exact: false }))
					.toBeInTheDocument();
				await expect.element(page.getByText('increased', { exact: false })).toBeInTheDocument();
				await view.rerender({
					message: create(PlayerMessageSchema, {
						type,
						target: planetTarget,
						spec: { habType, amount: -1, amount2 }
					})
				});
				await expect.element(page.getByText('decreased', { exact: false })).toBeInTheDocument();
			}
		);
	}
});

describe('player and legacy messages', () => {
	it.each([
		[TechCategory.SHIP_HULL, 'ship hull'],
		[TechCategory.STARBASE_HULL, 'starbase hull'],
		[TechCategory.PLANETARY_DEFENSE, 'defenses'],
		[TechCategory.PLANETARY_SCANNER, 'scanner'],
		[TechCategory.ARMOR, 'benefit']
	] as const)('renders research unlock category %s', async (category, expected) => {
		state.category = category;
		render(MessageDetail, {
			message: create(PlayerMessageSchema, {
				type: Type.PLAYER_TECH_GAINED,
				spec: { field: 'Energy', techGained: 'New Technology' }
			})
		});
		await expect
			.element(page.getByText(`New Technology ${expected}`, { exact: false }))
			.toBeInTheDocument();
	});
	it('renders discovery, research progress, and error details', async () => {
		const view = render(MessageDetail, {
			message: create(PlayerMessageSchema, {
				type: Type.PLAYER_DISCOVERY,
				target: { targetPlayerNum: 2 },
				spec: { name: 'Visitors' }
			})
		});
		await expect
			.element(page.getByText('new species, the Visitors', { exact: false }))
			.toBeInTheDocument();
		await view.rerender({
			message: create(PlayerMessageSchema, {
				type: Type.PLAYER_GAIN_TECH_LEVEL,
				spec: { field: 'Energy', nextField: 'Weapons', amount: 4 }
			})
		});
		await expect
			.element(page.getByText('Tech Level 4 for Energy', { exact: false }))
			.toBeInTheDocument();
		await expect.element(page.getByText('Weapons field.', { exact: false })).toBeInTheDocument();
		await view.rerender({
			message: create(PlayerMessageSchema, { type: Type.ERROR, spec: { error: 'Missing object' } })
		});
		await expect
			.element(page.getByText('Please contact the administrator, Missing object', { exact: false }))
			.toBeInTheDocument();
	});
	it.each([1, 2])('renders victory for player %s', async (targetPlayerNum) => {
		render(MessageDetail, {
			message: create(PlayerMessageSchema, {
				type: Type.PLAYER_VICTOR,
				target: { targetPlayerNum },
				spec: { name: 'Visitors' }
			})
		});
		await expect
			.element(
				page.getByText(
					targetPlayerNum === 1
						? 'You have been declared the winner'
						: 'The Visitors have been declared the winner',
					{ exact: false }
				)
			)
			.toBeInTheDocument();
	});
	it.each([fleetTarget, planetTarget, undefined])(
		'preserves legacy text for target %j',
		async (target) => {
			render(MessageDetail, {
				message: create(PlayerMessageSchema, {
					type: Type.INVALID,
					target,
					text: 'Saved backend message.'
				})
			});
			await expect
				.element(page.getByText('Saved backend message.', { exact: true }))
				.toBeInTheDocument();
		}
	);
});

describe('invasion reports', () => {
	it.each([
		[Type.FLEET_INVADED_PLANET, true],
		[Type.FLEET_INVADED_PLANET, false],
		[Type.PLANET_INVADED, true],
		[Type.PLANET_INVADED, false]
	] as const)('renders invasion %s, successful %s', async (type, successful) => {
		const view = render(MessageDetail, {
			message: create(PlayerMessageSchema, {
				type,
				target: planetTarget,
				spec: {
					amount: 12000,
					amount2: 10000,
					invasion: {
						fleetName: 'Invaders',
						attackerPlayerNum: 2,
						defenderPlayerNum: 3,
						attackersKilled: 8700,
						defendersKilled: 10000,
						successful
					}
				}
			})
		});
		await expect
			.element(
				page.getByText('12,000 attacking colonists and 10,000 defending colonists', {
					exact: false
				})
			)
			.toBeInTheDocument();
		await expect.element(page.getByText('Visitors', { exact: false })).toBeInTheDocument();
		await expect
			.element(page.getByText('Private race notes', { exact: false }))
			.not.toBeInTheDocument();
		if (type === Type.PLANET_INVADED) {
			await expect
				.element(page.getByText('Your troops beaming down', { exact: false }))
				.not.toBeInTheDocument();
		}
		await view.rerender({
			message: create(PlayerMessageSchema, {
				type,
				target: planetTarget,
				spec: { invasion: { successful, attackersKilled: 8700, defendersKilled: 10000 } }
			})
		});
		await expect
			.element(page.getByText('The invasion began', { exact: false }))
			.not.toBeInTheDocument();
	});
});

describe('production and orbital adjustment reports', () => {
	it.each([
		[Type.PLANET_PRODUCTION_QUEUE_EMPTY, 'production queue on Earth is empty'],
		[
			Type.PLANET_PRODUCTION_QUEUE_COMPLETE,
			'Earth has completed its orders. The production queue is empty.'
		]
	] as const)('renders production status %s', async (type, expected) => {
		render(MessageDetail, { message: create(PlayerMessageSchema, { type, target: planetTarget }) });
		await expect.element(page.getByText(expected, { exact: false })).toBeInTheDocument();
	});
	it.each([
		[1, 60, 80, 1, 'Your fleet Adjuster has improved Earth from a value of 60% to 80%.'],
		[2, 80, 60, -1, 'A Visitors fleet Adjuster has degraded Earth from a value of 80% to 60%.'],
		[1, 60, 60, 1, 'is currently unable to improve the value of Earth beyond 60%.'],
		[2, -10, -10, -1, 'is currently unable to degrade the value of Earth beyond -10%.'],
		[1, -10, -5, 1, 'has improved Earth from a value of -10% to -5%.']
	] as const)(
		'renders original-style planet value for owner %s, %s to %s',
		async (sourcePlayerNum, prevAmount, amount, amount2, expected) => {
			render(MessageDetail, {
				message: create(PlayerMessageSchema, {
					type: Type.PLANET_REMOTE_TERRAFORM,
					target: planetTarget,
					spec: {
						sourcePlayerNum,
						prevAmount,
						amount,
						amount2,
						terraformAmount: { grav: -2, temp: 1 },
						mapObjectTarget: { targetName: 'Adjuster' }
					}
				})
			});
			await expect.element(page.getByText(expected, { exact: false })).toBeInTheDocument();
			await expect
				.element(page.getByText('Gravity has decreased by 2%.', { exact: false }))
				.toBeInTheDocument();
			await expect
				.element(page.getByText('Temperature has increased by 1%.', { exact: false }))
				.toBeInTheDocument();
			await expect
				.element(page.getByText('Radiation has', { exact: false }))
				.not.toBeInTheDocument();
		}
	);

	it.each([
		[1, 1, 'Your fleet Adjuster', 'increasing'],
		[2, -1, 'A Visitors fleet Adjuster', 'decreasing']
	] as const)(
		'renders remote terraforming by %s',
		async (sourcePlayerNum, amount, owner, direction) => {
			render(MessageDetail, {
				message: create(PlayerMessageSchema, {
					type: Type.PLANET_REMOTE_TERRAFORM,
					target: planetTarget,
					spec: {
						sourcePlayerNum,
						amount,
						amount2: 60,
						habType: TerraformHabType.TEMP,
						mapObjectTarget: { targetName: 'Adjuster' }
					}
				})
			});
			await expect.element(page.getByText(owner, { exact: false })).toBeInTheDocument();
			await expect
				.element(page.getByText(`${direction} its Temperature by 1% to 40°C`, { exact: false }))
				.toBeInTheDocument();
		}
	);
});

describe('battle loss reports', () => {
	for (const type of [Type.BATTLE, Type.BATTLE_ALLY]) {
		it.each([
			[0, 0, 'No ships were lost by either side.'],
			[1, 0, type === Type.BATTLE ? 'Only you suffered losses' : 'Only your ally suffered losses'],
			[0, 1, 'Only the enemy suffered losses'],
			[
				1,
				1,
				type === Type.BATTLE
					? 'Both you and the enemy suffered losses'
					: 'Both your ally and the enemy suffered losses'
			],
			[2, 0, 'was annihilated'],
			[0, 6, 'without suffering a single casualty'],
			[2, 6, 'completely destroyed each other']
		] as const)(
			`renders battle type ${type} with losses %s/%s`,
			async (ourDead, theirDead, expected) => {
				render(MessageDetail, {
					message: create(PlayerMessageSchema, {
						type,
						battleNum: 1,
						spec: {
							name: 'Earth',
							battle: {
								numShipsByPlayer: { 1: 2, 2: 6 },
								shipsDestroyedByPlayer: { 1: ourDead, 2: theirDead }
							}
						}
					})
				});
				await expect.element(page.getByText(expected, { exact: false })).toBeInTheDocument();
				if (ourDead === 0 || theirDead === 0) {
					await expect.element(page.getByText('Both', { exact: false })).not.toBeInTheDocument();
				}
			}
		);
	}
});

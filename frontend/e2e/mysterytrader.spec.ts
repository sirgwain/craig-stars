import { MysteryTraderRewardType } from '../src/lib/protogen/craig_stars/v1/mysterytrader_pb';
import { key } from '../src/lib/types/MapObject';
import { expect, submitTurn, test } from './setup';

test('Mystery Trader Test', async ({ testGamePage }) => {
	const { page, universe } = await testGamePage('Mystery Trader');

	// zoom out so we can see the mystery trader
	await page.keyboard.press('-');
	await page.keyboard.press('-');
	await page.keyboard.press('-');
	await page.keyboard.press('-');
	await page.keyboard.press('-');

	for (let i = 0; i < universe.fleets.length; i++) {
		const freighter = universe.fleets[i];
		const mt = universe.mysteryTraders[i];

		// click the homeworld once to cycle to the colonizer
		await page.locator(`[data-id="${key(freighter)}"]`).dblclick({ force: true });
		await page
			.locator(`[data-type="command-tile"][data-id="${freighter.mapObject?.name}"]`)
			.first();

		const updateFleetOrdersResponse = page.waitForResponse(
			(resp) =>
				resp.url().includes('/api/grpc/craig_stars.v1.FleetService/UpdateFleetOrders') &&
				resp.request().method() === 'POST' &&
				resp.status() === 200
		);

		// shift click the planet to set a waypoint at max speed
		await page.locator(`[data-id="${key(mt)}"]`).click({ force: true, modifiers: ['Shift'] });

		// wait for the TransferCargo to complete
		await updateFleetOrdersResponse;
	}

	const { universe: updatedUniverse, player: updatedPlayer } = await submitTurn(page);

	// mystery traders should have moved
	for (let i = 0; i < universe.mysteryTraders.length; i++) {
		const mt = universe?.mysteryTraders[i];
		const updatedMt = updatedUniverse?.mysteryTraders[i];
		// we should colonize the new planet
		expect(updatedMt?.mapObject?.position?.x).not.toBe(mt.mapObject?.position?.x);
	}

	// should have one new fleet, a lifeboat
	// all other freighters should have been absorbed
	expect(updatedUniverse?.fleets.length).toBe(1);
	const rewardFleet = updatedUniverse?.fleets[0];
	const rewardDesign = updatedUniverse?.designs.find((design) => design.num === rewardFleet?.tokens[0].designNum);
	expect(rewardDesign?.mysteryTrader).toBe(true);

	// should have messages for various mystery trader meetups
	expect(
		updatedPlayer?.messages.find(
			(m) => m.spec?.mysteryTrader?.type === MysteryTraderRewardType.RESEARCH
		)
	).toBeDefined();
	expect(
		updatedPlayer?.messages.find(
			(m) => m.spec?.mysteryTrader?.type === MysteryTraderRewardType.SHIP_HULL
		)
	).toBeDefined();
	expect(
		updatedPlayer?.messages.find(
			(m) => m.spec?.mysteryTrader?.type === MysteryTraderRewardType.LIFEBOAT
		)
	).toBeDefined();
	expect(
		updatedPlayer?.messages.find(
			(m) => m.spec?.mysteryTrader?.type === MysteryTraderRewardType.GENESIS
		)
	).toBeDefined();
	expect(
		updatedPlayer?.messages.find(
			(m) => m.spec?.mysteryTrader?.type === MysteryTraderRewardType.JUMP_GATE
		)
	).toBeDefined();
});

import { key } from '../src/lib/types/MapObject';
import { apiErrorsFailTest, expect, loadGamePage, submitTurn, test } from './setup';

test('Kitchen Sink', async ({ authenticatedPage }) => {
	const { page, gameId, universe } = await loadGamePage(authenticatedPage, 'Kitchen Sink');
	apiErrorsFailTest(page, gameId);

	const mapObjectSummary = await page.locator('[data-type="map-object-summary"]').first();

	const mineralsOnHandTile = await page
		.locator('[data-type="command-tile"][data-id="Minerals on Hand"]')
		.first();

	await expect(mineralsOnHandTile).toBeVisible();

	await expect(mineralsOnHandTile).toContainText('Ironium 1000kT');
	await expect(mineralsOnHandTile).toContainText('Boranium 1000kT');
	await expect(mineralsOnHandTile).toContainText('Germanium 1000kT');

	const statusTile = await page.locator('[data-type="command-tile"][data-id="Status"]').first();

	await expect(statusTile).toContainText('Population 250,000');
	await expect(statusTile).toContainText('Defense Type SDI');

	// right click planet
	await page
		.locator(`[data-id="${key(universe.planets[0])}"]`)
		.click({ force: true, button: 'right' });
	const scannerPopup = await page.locator(`[data-id="scanner-context-popup"]`);
	await expect(scannerPopup).toBeVisible();

	// click the minefield from the context menu
	await scannerPopup.getByRole('button', { name: 'Humanoids Standard Minefield #1' }).click();
	await expect(mapObjectSummary).toContainText('Humanoids Standard Minefield #1');
	await expect(mapObjectSummary).toContainText('Location: (0, 0)');
	await expect(mapObjectSummary).toContainText('Field Type: Standard');
	await expect(mapObjectSummary).toContainText('Field Radius: 32 l.y. (1000 mines)');
	await expect(mapObjectSummary).toContainText('Maximum Safe Speed: Warp 4');

	// select mineralPacket
	await page.locator(`[data-id="${key(universe.mineralPackets[0])}"]`).click({ force: true });

	await expect(mapObjectSummary).toContainText('Humanoids Mineral Packet #1');
	await expect(mapObjectSummary).toContainText('Location: (50, 0)');
	await expect(mapObjectSummary).toContainText('Traveling at Warp: 5');
	await expect(mapObjectSummary).toContainText('Destination: Planet 1');
	await expect(mapObjectSummary).toContainText('ETA: 2 years');
	await expect(mapObjectSummary).toContainText('Ironium 50kT');
	await expect(mapObjectSummary).toContainText('Boranium 50kT');
	await expect(mapObjectSummary).toContainText('Germanium 50kT');

	// select other player's fleet
	const otherPlayerFleet = universe.fleets.find((f) => f.mapObject?.playerNum === 2);
	if (!otherPlayerFleet?.mapObject) {
		throw new Error("other player's fleet not found");
	}
	await page.locator(`[data-id="${key(otherPlayerFleet)}"]`).click({ force: true });
	await expect(mapObjectSummary).toContainText(otherPlayerFleet.mapObject.name);
	await expect(mapObjectSummary).toContainText(`Ship Count: 1`);
	await expect(mapObjectSummary).toContainText(
		`Mass: ${otherPlayerFleet.spec?.shipDesignSpec?.mass}`
	);
	await expect(mapObjectSummary).toContainText(`Warp Speed: ${otherPlayerFleet.warpSpeed}`);
	await expect(mapObjectSummary).toContainText(universe.playerIntels[1].racePluralName);

	// select other player's planet
	const otherPlayerPlanet = universe.planets.find((f) => f.mapObject?.playerNum === 2);
	if (!otherPlayerPlanet?.mapObject) {
		throw new Error("other player's planet not found");
	}
	await page.locator(`[data-id="${key(otherPlayerPlanet)}"]`).click({ force: true });
	await expect(mapObjectSummary).toContainText('Report is current');
	await expect(mapObjectSummary).toContainText(universe.playerIntels[1].racePluralName);

	// zoom out so we can see the mystery trader
	await page.keyboard.press('-');
	await page.keyboard.press('-');
	await page.keyboard.press('-');

	// select mystery trader
	const mysteryTrader = universe.mysteryTraders[0];
	if (!mysteryTrader?.mapObject) {
		throw new Error('mysteryTrader not found');
	}
	await page.locator(`[data-id="${key(mysteryTrader)}"]`).click({ force: true });
	await expect(mapObjectSummary).toContainText(
		'The trader requests interested parties to send it a feet with at least 5000kT of minerals on board to be absorbed into the trader. It offers technological assistance in return. Trader is traveling at Warp 7.'
	);

	// select wormhole
	const wormhole1 = universe.wormholes[0];
	if (!wormhole1?.mapObject) {
		throw new Error('wormhole not found');
	}
	await page.locator(`[data-id="${key(wormhole1)}"]`).click({ force: true });
	await expect(mapObjectSummary).toContainText('Location: (10, 10)');
	await expect(mapObjectSummary).toContainText('Stability: Rock Solid');
	await expect(mapObjectSummary).toContainText('Destination: unknown');

	const fleetsInOrbitTile = await page
		.locator('[data-type="command-tile"][data-id="Fleets In Orbit"]')
		.first();

	// goto the scout
	await expect(fleetsInOrbitTile).toBeVisible();
	await fleetsInOrbitTile.getByRole('button', { name: 'Goto' }).click();

	// send the scout to the wormhole
	await page.locator(`[data-id="${key(wormhole1)}"]`).click({ force: true, modifiers: ['Meta'] });

	const { universe: updatedUniverse } = await submitTurn(page);

	// should know about the second wormhole
	expect(updatedUniverse?.wormholes.length).toBe(2);

	// zoom out so we can see the wormhole
	await page.keyboard.press('-');
	await page.keyboard.press('-');
	await page.keyboard.press('-');

	// double click the fleet
	await page
		.locator(`[data-id="${key(updatedUniverse?.fleets[0])}"]`)
		.click({ force: true, clickCount: 2 });

	// fleet should be at wormhole2 location
	const fleetWaypoints = await page
		.locator('[data-type="command-tile"][data-id="Fleet Waypoints"]')
		.first();

	await expect(fleetWaypoints).toContainText('Space: (60, 60)');

	// select wormhole2
	const wormhole2 = updatedUniverse?.wormholes[1];
	if (!wormhole2?.mapObject) {
		throw new Error('wormhole2 not found');
	}
	await page.locator(`[data-id="${key(wormhole2)}"]`).click({ force: true });
	await expect(mapObjectSummary).not.toContainText('Destination: unknown');
});

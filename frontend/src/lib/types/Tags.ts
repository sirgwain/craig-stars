// Tags are short labels a player puts on map objects. Tags that differ only by case are the same tag.
export function sameTag(a: string, b: string): boolean {
	return a.toLowerCase() === b.toLowerCase();
}

// collapse whitespace the way the server does
export function cleanTag(tag: string): string {
	return tag.trim().split(/\s+/).join(' ');
}

// Every tag gets a color from its name, so a tag looks the same everywhere without the player
// having to pick one.
export function tagHue(tag: string): number {
	let hash = 0;
	for (const c of tag.toLowerCase()) {
		hash = (hash * 31 + c.codePointAt(0)!) >>> 0;
	}
	// 12 hues, evenly spaced
	return (hash % 12) * 30;
}

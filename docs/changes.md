# Deviations from original game

`craig-stars`, as a clone of the original 1995 Stars! game, strives to be as close to the original as possible in terms of base game mechanics. However, certain elements either do not translate well to a modern web game or are simply not worth keeping.

This is a non-exhaustive list of intentional differences between `craig-stars` and the original game. Features yet to be implemented are not included.

<!-- TODO: Do we want to make this into a table???-->

## Medium-ish changes

These changes, while none are explicitly "game-breaking", are still good to keep in mind while playing.

- Upgrading a starbase now checks all slots in _both_ designs when calculating refunds, rather than only the corresponding slot in the new design[^1]. Additionally, part refund/transfer checks still occur when swapping hulls (so adding a component before swapping hulls costs the same regardless of order).

- Invalid items inside planet production queues (packets lacking mass drivers, excess planetary installations, etc.) will be automatically removed during turn generation and refund any previously spent minerals/resources upon doing so. The production queue estimator also explicitly highlights such items as being canceled.

- [Tech trades by scrapping](https://wiki.starsautohost.org/wiki/Tech_Trade_by_Scrapping) (both for techs and MT parts) now tally up chances for each _token_ in the fleet, not each individual fleet being scrapped. This has virtually no impact on most forms of gameplay, simply removing the need to split one's trading ships before scrapping.

<!--
- Fleets hitting minefields will reduce mine counts on a per-token basis instead of a per-fleet basis. This (again) has little to no actual bearing on most forms of gameplay, other than making (collision/chaff sweeping)[https://wiki.starsautohost.org/wiki/Collision_sweeping] slightly less annoying to perform.
-->

[^1]: For reference, Stars! charges extra for merely moving a component to a different slot on the same hull.

## Bugs fixed

This is a list of confirmed "bugs" from the base game fixed by `craig-stars`. More will be added as fixes for them are confirmed.

<!-- TODO: Check and fix more bugs -->

- [0.2% min damage bug](https://wiki.starsautohost.org/wiki/Known_Bugs#0.2%_Minimum_Damage)
- [Colonization Module Check](https://wiki.starsautohost.org/wiki/Known_Bugs#Colonization_Module_Check)
- [Cheap Starbase](https://wiki.starsautohost.org/wiki/Known_Bugs#Cheap_Starbase)

## Extremely minor changes

These changes have next to no tangible impact on gameplay, only being mentioned out of a sence of completeness.

- The `DeltaPopulation` variable responsible for tracking pop growths below multiples of 100 would not reset when abandoning or invading a planet. This is fixed in `craig-stars`.

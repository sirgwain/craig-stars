# Deviations from original game

`craig-stars`, as a clone of the original 1995 Stars! game, strives to be as close to the original as possible in terms of base game mechanics. However, certain elements either do not translate well to a modern web game or are simply not worth keeping.

Here is a non-exhaustive list of differences between `craig-stars` and the original game. Features yet to be implemented are not included.

## Major changes

- The AI's ship design algorithm has been effectively redesigned from the ground up, creating much more varied and effective ships closer to those from a human player.
  -# For the record, filling a ship with half beams & torpedoes is not a good design.

## Medium-ish changes

- Upgrading a starbase now checks all slots in _both_ designs when calculating refunds, rather than only the corresponding slot in the new design. Part refund/transfer checks also occur when swapping hulls[^1]. (This has the overall effect of making starbase upgrades slightly cheaper.)

[^1]: For reference, Stars! charges extra for merely moving components to a different slot, and forgoes normal refund/transfer logic for a flat 50% refund when swapping hulls.

<!-- TODO: Check and fix more bugs -->

## Bugs fixed

- [0.2% min damage bug](https://wiki.starsautohost.org/wiki/Known_Bugs#0.2%_Minimum_Damage)
- [Colonization Module Check](https://wiki.starsautohost.org/wiki/Known_Bugs#Colonization_Module_Check)
- [Cheap Starbase](https://wiki.starsautohost.org/wiki/Known_Bugs#Cheap_Starbase)

## Minor changes

- The `DeltaPopulation` variable responsible for tracking pop growths below multiples of 100 did not reset when abandoning or invading a planet. This is fixed in `craig-stars`.

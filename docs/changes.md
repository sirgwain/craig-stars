# Deviations from original game

`craig-stars`, as a clone of the original 1995 Stars! game, strives to be as close to the original as possible in terms of base game mechanics. However, certain elements either do not translate well to a modern web game or are simply not worth keeping.

Here is a non-exhaustive list of differences between `craig-stars` and the original game. Features yet to be implemented are not included.

Small differences in how a number is calculated, such as rounding where the original truncated, are not listed.

## Limits removed

- A player can have any number of ship and starbase designs. Stars! allowed 16 ship designs and 10 starbase designs.
- A player can have any number of fleets. Stars! allowed 512.

## Race design and starting position

- Leftover advantage points that can't be spent on the chosen bonus are spent on surface minerals instead of being lost. This happens when there are too few points for a single mine, factory or defense, or when the race has no use for them (Alternate Reality). The race designer warns when this will happen.
- Starting ships use the same hulls and names as Stars!, but their components are chosen by the ship designer the AI uses, based on the race's starting tech. Stars! used a fixed parts list for each starting design, which sometimes mixed beam weapons and torpedoes on one ship.

## Diplomacy

- A player can share their map with a player they have set as a friend. The friend sees what they see, and learns about the players they have met. Stars! had no map sharing.

## Medium-ish changes

- Once every tech field is at its maximum level, resources that would have gone to research are spent on mineral alchemy instead of being wasted.

- Upgrading a starbase now checks all slots in _both_ designs when calculating refunds, rather than only the corresponding slot in the new design[^1]. Additionally, part refund/transfer checks still occur when swapping hulls (so adding a component and swapping hulls on consecutive turns costs the same regardless of order).
- Players with fleets at a battle's location who don't take part still observe it. They receive the battle record, discover every design in the battle, and see the surviving fleets and starbase, just like the participants. In Stars!, these spectators only received a message that a battle took place.

[^1]: For reference, Stars! charges extra for merely moving components to a different slot.

## Interface

- Commanding a fleet that already has waypoints selects its last waypoint, so new waypoints are added to the end of the route. Stars! selected the fleet's current position.

## Bugs fixed

- [0.2% min damage bug](https://wiki.starsautohost.org/wiki/Known_Bugs#0.2%_Minimum_Damage)
- [Colonization Module Check](https://wiki.starsautohost.org/wiki/Known_Bugs#Colonization_Module_Check)
- [Cheap Starbase](https://wiki.starsautohost.org/wiki/Known_Bugs#Cheap_Starbase)

## Minor changes

- The `DeltaPopulation` variable responsible for tracking pop growths below multiples of 100 did not reset when abandoning or invading a planet. This is fixed in `craig-stars`.

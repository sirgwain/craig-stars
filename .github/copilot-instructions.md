# Custom instructions for Copilot

## Project context

This project is a web clone of the 4x game Stars! released in 1995. It uses a Golang backend and Typescript frontend server.

## Indentation

Use tabs for indentation and spaces for alignment. Gofmt and Prettier typically take care of formatting automatically, but keeping it consistent helps with readability.

## Code Style Guidelines

Use camelCase for variable names and prefer function literals over const declarations.

When writing constant numerals over 10000, use underscores to separate every 3 orders of magnitude.

When writing new functions or moving existing ones, attempt to place them in order of first appearance in the file (or in their parent file for test functions).

Unless explicitly told otherwise, preserve comments when rewriting source material.

## Testing Guidelines

Golang tests use a mixture of assertion-based testing and table-driven tests. Vitest relies near exclusively on asserio

For table-driven tests with more than 2 parameters (other than name and want) have them conventionally wrapped in args and fields structs. Args holds any arguments directly passed to the function being tested, while fields holds values used for instantizing each test case. (Either can be ommitted if they would hold only 1 value.)

In table driven tests, if the object being tested has a Name field, set it to the name of the test case before running to help during debugging.

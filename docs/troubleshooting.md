# Troubleshooting

> When the code hits the fan

Things don't always go as planned (in practice, they hardly ever do). In the event things turn south when setting things up, consult this semi-curated list of past problems and resolutions:

<!-- TODO: Improve intro and sort these in vague order of appearance-->

- "I try to click on the login button on localhost using the admin credentials and it does nothing! Worse, an error pops up in my terminal!"

You might be running the frontend server without the backend. Open a new terminal tab and type `mage launch_backend` to launch the backend to handle all the nitty gritty logic stuff.

- "My computer complains about undefined Sqlite Drivers!"

Make sure you have `GCC` built and in your PATH. If you haven't installed it yet, go do that.

- "While building the server, I get an obscure error about 'executable not found in %PATH%' or 'build target excluding all files in XXX'!"

What's probably happening is you're trying to generate go files or build the server with incorrect GOARCH and GOOS settings. Try running `go env -u GOOS GOARCH` to reset them to their defaults and see if the problems persist.

- "I get an error message mentioning the database on launch!"

In the event MySQL is being unhappy, you have 2 options:

1. Delete the entire `dist` folder (containing the starter admin database) and re-launch the server. This will re-generate the db and _should_ fix 90% of issues related to new installs.
2. In the event you modified exported values of structs saved to the database, you will need to add new .sql files inside `./db/schema` to instruct it to drop the new games. See the relevant section in [architecture.md](architecture.md/#db) for more info.

- "When I boot up the server, all the ships have no icons!"

See the [Assets](development#assets) section for information on how to download art assets.

- "I tried to run `mage images`, but the command failed!

In the event the image download script fails, you'll have to download the [images](https://craig-stars.net/images/images.zip) manually and move the extrated files to `frontend/static/images` yourself.

- "Mage is spitting out weird commands that aren't working!"

In the event mage starts executing warped commands, you can use the `Debug Magefile Target` debug configuration to launch a `dlv` session debugging a particular magefile target.

- "I'm getting some other errors in the command line!"

Consult this ordered checklist of vague general suggestions:

1. Read the error message to try and figure out why it's failing. Errors before running can be a sign of malformed mage commands, while errors during command execution oft lie with the software being run.
2. Try and search online for the error message or similar problems to see if others may have found solutions already.
3. Try updating your packages (either node and/or golang, depending on where the errors occur) to the latest versions. A surprising amount of bugs can be fixed by simply running `npm update` or `go install`.
4. If all else fails, reach out in the #stars-clones or #craig-stars channels in the discord (ideally with images/text of the commands used and/or the resulting error messages - vague comments like "AAA MY BUILD IS BORKING" tend to be hard to troubleshoot).

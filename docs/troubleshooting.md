# Troubleshooting
> When the code hits the fan

Things don't always go as planned (in practice, they hardly ever do).
<!-- TODO: Add a better intro and segue-->

- "I try to click on the login button on localhost using the admin credentials and it does nothing! Worse, an error pops up in my terminal!"

You might be running the frontend server without the backend. Open a new terminal tab and type `mage launch_backend` to launch the backend to handle all the nitty gritty logic stuff.

- "When I run air, my computer complains about undefined Sqlite Drivers!"

You likely haven't installed `go-sqlite3` and `GCC` correctly. Go do that.

- "When I run `mage build`, I get an obscure error about 'executable not found in %PATH%' or 'build target excluding all files in XXX'!"

What's probably happening is you're trying to generate go files or build the server with the incorrect GOARCH and GOOS settings. Try running `go env -u GOOS GOARCH` to reset them to their defaults and see if the problems persist.

- "Running air produces an error message something like `cmd will not recognize XXX file for execution`!"

This is a 100% normal thing and a direct consequence of using `mage` to execute commands.

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

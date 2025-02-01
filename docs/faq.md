## Common Problems
- "I try to click on the login button on localhost using the admin credentials and it does nothing! Worse, an error pops up in my terminal!"

You probably aren't running the backend server. Open a new terminal tab and type `air` to build the backend needed to handle all the nitty gritty logic stuff (including login credentials).

- "When I run air, my computer complains about undefined Sqlite Drivers!"

See the [GCC](#gcc) section for info on how to install `go-sqlite3` and `GCC`.

- "When I run `make build`, I get an obscure error about 'executable not found in %PATH%' or 'build target excluding all files in XXX'!"

What's probably happening is you're trying to generate go files or build the server with the incorrect GOARCH and GOOS settings. Try running `go env -u GOOS GOARCH` to reset them to their defaults and see if the problem persists.

- "When I boot up the server, all the ships have no icons!"

See [Assets](#assets) for information on how to download art assets.

- "I tried to run `make images`, but the command failed!`

In the event `make images` fails, you'll have to download & extract [the images](https://craig-stars.net/images/images.zip) manually and move them to `frontend/static/images` yourself.

- "I'm getting some other errors in the command line!"

Consult this ordered checklist of vague general suggestions:

1. Read the error message to try and figure out why it's failing. Make (as it's being used here) effectively just copy-pastes its commands into the terminal one by one (expanding variables here and there), so errors in the terminal can be a symptom of bad or incorrect launch commands.
2. Try and search online for the error message or similar problems to see if others may have found solutions to the problem for you.
3. If all else fails, reach out in the #stars-clones or #craig-stars channels in the discord (ideally with images/text of the commands used and/or the resulting error messages - vague comments like "AAA MY BUILD IS BORKING" tend to be hard to troubleshoot).

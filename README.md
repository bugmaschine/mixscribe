## mixscribe

mixscribe is like an radio recorder and it checks for looping parts. It works by recording and then checking if it's starts overlaping with a previous part of the already saved audio. If it does it than saves the file and wait's until the a different mix is played.

PS: please don't sue me ard, and yes i know the code is shit.

# Instalation

We recommend installing mixscribe on a VPS with some memory to spare (swap should work), a good uptime and internet connection. Theoreticly you can run Jump Record on a home internet connection, but if you have the telekom as your provider then we don't recommend you try to even think about it.

TODO: Add instalation instructions or add a docker container

# Requirements

TODO: Add requirements

# Developers

Run the program with `go run main.go`

# Current State

The current state of the project is that it's still in development and the part where it checks if the audio is still the same is still very basic and not working very well. The project is open source and any help is welcome. If you have any ideas or want to help, please don't hesitate to open an issue or make a pull request.

# Licensing

See [COPYING](COPYING)
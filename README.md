IF YOU ARE ON GITHUB.COM GO HERE INSTEAD: https://gitea.asdf.cafe/abs3nt/gospt :)

If you open an issue or PR on github I won't see it please use gitea. Register on asdf and open your PRs there

This project is still under heavy development and some things might not work or not work as intended. Don't hesitate to open an issue to let me know.

Project Discord: [Join Here](https://discord.gg/nWEEK6HrUD)
---

to install:

```yay -S gospt```

or to build from source by pulling and building the binary

```git clone https://gitea.asdf.cafe/abs3nt/gospt```

```cd gospt```

```make build```
then

```sudo make install```

to use add your information to ~/.config/gospt/client.yml like this

```
client_id: "idgoeshere"
client_secret: "secretgoeshere"
port: "8888"
```
if you dont want to store your secret in the file in plaintext you can use a command to retreive it:

```
client_secret_cmd: "secret spotify_secret"
```

you should have either client_secret or client_secret_cmd


then run

```gospt```

you will be asked to login, you will only have to do this the first time. After login you will be asked to select your default device, this will also only happen once. To reset your device run ```gospot setdevice```


To use the custom radio feature:

```gospt radio```


or hit ctrl+r on any track in the TUI. This will start an extended radio. To replenish the current radio run ```gospt refillradio``` and all the songs already listened will be removed and that number of new recomendations will be added.

This radio uses slightly different logic than the standard spotify radio to give a longer playlist and more recomendation. With a cronjob you can schedule refill to run to have an infinite and morphing radio station.

To view help:

```gospt --help```

Very open to contributations feel free to open a PR

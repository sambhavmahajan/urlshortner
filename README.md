### UrlShortner

This is a URL shortener project which shortens the URLs. It takes a random number between 100000-999999 and assigns it and maps it to the URL to be shortened. By default it works on the 8080 port, you can modify it in the http.ListenAndServe function call.

#### Clone
```bash
git clone https://github.com/sambhavmahajan/urlshortner.git
```

#### Data Storage
- HashMap (For fast access)
- Database (Postgres, for persistence)
- Mutex lock ( i am not a monster, i know synchronization shit)

#### Collusion Strategy
Linear probing is used, when a collusion occurs, first we look at the right, if the right is full, we look at the left (kinda caught a bug while writing this, in the shortenHandler function, check out issue [#2](https://github.com/sambhavmahajan/urlshortner/issues/2))

There are constants in main.go, the host, port, user, password and dbname, make sure to edit those. Look at line number 60. No i did not make an '.env' file, yes i know its bad practice, no i am not sorry, yes its because it forgot(i even forgot to write this documentation, what do you expect from me)

#### Before running the webapp, run the command
```bash
go mod tidy
```
for the modules, mostly because of database and sql modules, or i guess only database modules(lib/pq and database/sql)

#### Endpoints

```
Post /shorten/?url=https://github.com/sambhavmahajan
```
although i would say, if you are shortening my profile, thats lazy, why would anyone remember a 6 digit number over a name, but hey, i am also hardwired that way lol :)

```
GET /get/?id=100069
```
Get the url

That's It!!!

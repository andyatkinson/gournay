Gournay is a URL shrinker. The name is a horrible pun. I started with "shrink", [referring to psychiatrists](http://timesofindia.indiatimes.com/home/stoi/Why-is-a-psychiatrist-called-a-shrink/articleshow/929514.cms), then found [Kevin Gournay](http://en.wikipedia.org/wiki/Kevin_Gournay), involved in mental health, and who happens to have "Go" as part of his last name. See, I told you it was bad. :)

### building locally

    go build gournay.go

When pushing to heroku, the app is compiled. When adding new dependencies, remember to run `godep save`.

### running in development

    # psql
    create database gournay_development;
    create table entries (
      url varchar(255),
      hash varchar(5)
    );

Local postgres connection string like this, for a database called `gournay_development`:

    export DATABASE_URL="user=andy host=localhost dbname=gournay_development sslmode=disable"

    ./gournay
    # port is 5000 by default

### Heroku

Details on using the [Heroku buildpack](http://stackoverflow.com/a/16655766/126688).

    * Create a heroku instance with the custom buildpack
    * Enable the postgres addon with: `heroku addons:add heroku-postgresql`
    * Heroku binds to a random port, can't hard-code the port. [details](http://stackoverflow.com/a/15693371/126688)
    * connect to heroku psql and set up database: `heroku pg:psql`. Find out database name from connection string, and create the tables.

### resources

  * [go import path](http://blog.golang.org/organizing-go-code)
  * http://golang.org/doc/articles/wiki/
  * http://mmcgrana.github.io/2012/09/getting-started-with-go-on-heroku.html
  * http://blog.zmxv.com/2011/09/go-template-examples.html
  * https://code.google.com/p/go-wiki/wiki/SQLInterface
  * https://github.com/mindreframer/golang-stuff/tree/master/github.com/VividCortex/go-database-sql-tutorial

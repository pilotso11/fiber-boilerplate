# Fiber Boilerplate
A boilerplate for the Fiber web framework

## Configuration
Different to previous versions of this boilerplate, configurations are in a single file called `.env`. You can copy the `.env.example` and change it to your needs.

This `.env` file represents system environment variables on your machine. This change was made with the ease-of-use with Docker in mind.

A full version of all available configurations is located in the `.env.full` file. Various options can be changed depending on your needs such as Database, Fiber and Middleware settings.

Keep in mind if configurations are not set, they default to Fiber's default settings which can be found [here](https://docs.gofiber.io/).

## Routing
Routing examples can be found within the `/routes` directory. Both web and API routes are split, but you can adjust this to your likings.

## Views
Views are located and be edited under the `/resources/views` directory. 

You are able change this behavior using the `.env` file, as well as the ability to modify the Views Engine and other templating configurations using this file or using environment variables.

## Controllers
Example controllers can be found within the `/app/controllers` directory. You can extend or edit these to your preferences.

## Database
We use GORM v2 as an ORM to provide useful features to your models. Please check out their documentation [here](https://gorm.io/index.html).

## Models
Models are located within the `/app/models` directory and are also based on the GORM v2 package.

## Swagger docs
Sample swagger annotations are on the User and Group API's. Docs are exposed with `github.com/swaggo/fiber-swagger and docs generated by swag. 

To install: 
```bash
go get github.com/swaggo/swag/cmd/swag
```
Ensure GOPATH/bin is in your path.

To generate the docs run: 
```bash
swag init # to docs .json and .yaml
swag fmt  # To format the go annotations
````

# Sample authentication
The sample application provides to example authentication models.
* There is a built-in database model exposed via API.  You can use the swagger docs to add users and groups and test it with the login link on the main page. This is for illustration and is not for use in production. 
* There is also OAUTH2 added as a middleware using `goth` and `goth_fiber`.  
A sample configuration for AWS cognito is in the sample .env file.
```
MW_OAUTH2_ENABLED=true
MW_OAUTH2_PROVIDER="cognito"
MW_OAUTH2_KEY="MY_COGNITO_CLIENTKEY"
MW_OAUTH2_ORG_URL="https://CUSTOM-DOMAIN.auth.us-east-1.amazoncognito.com"
MW_OAUTH2_CALLBACK_SERVER="http://localhost:8080"
MW_OAUTH2_AFTER_LOGOUT_REDIRECT_URL="https://CUSTOM-DOMAIN.us-east-1.amazoncognito.com/logout?client_id=MY_COGNITO_CLIENTKEY&logout_uri=http://localhost:8080/"
```
This is a production ready middleware.  The authenticated user details are saved in the session as `user` 
and the middleware copies `userid` `email` and `username` to the Ctx on each page allowing you to easily check for authentication without having to load the session details.
There is a sample secured page that you can only access when authenticated with OAUTH2. 

This can be viewed in the `app/controllers/web/secured.go`
```golang
    userid := c.Locals("userid")
    if userid == nil || userid == "" {
        return c.SendStatus(fiber.StatusUnauthorized)
    }
```

# Secrets
You can load a secrets .env file.  The path is specified in your .env using `SECRETS_PATH`.  Any values set in the secrets file will override anything set in your .env file.

This secret file could be supplied by Docker or your hosting engine.

# Faster JSON 
By default faster JSON is enabled using `github.com/goccy/go-json`.

# Fast logging with zap
In your .env file set `MW_ACCESS_LOGGER_TYPE="console"` to enable console logging for fiber with zap.
To use the zap global logger in in your code `zap.S().Debugf()` and similar zap functions. 

## Compiling assets
This boilerplate uses [Laravel Mix](https://github.com/JeffreyWay/laravel-mix) as an elegant wrapper around [Webpack](https://github.com/webpack/webpack) (a bundler for javascript and friends).

In order to compile your assets, you must first add them in the `webpack.mix.js` file. Examples of the Laravel Mix API can be found [here](https://laravel-mix.com/docs/5.0/mixjs).

Then you must run either `npm install` or `yarn install` to install the packages required to compile your assets.

Next, run one of the following commands to compile your assets with either `npm` or `yarn`:

```bash
# Run all Mix tasks
npm run dev

# Run all Mix tasks and minify output
yarn run production
# Run all Mix tasks and watch for changes (useful when developing)
yarn run watch
# Run all Mix tasks with hot module replacement
yarn run hot
```

## Docker
You can run your own application using the Docker example image.
To build and run the Docker image, you can use the following commands.

Please note, I am using host.docker.internal to point to my Docker host machine. You are free to use Docker's internal networking to point to your desired database host.


```bash
docker build -t fiber-boilerplate .
docker run -it --rm --name fiber-boilerplate -e DB_HOST=host.docker.internal -e DB_USER=fiber -e DB_PASSWORD=secret -e DB_DATABASE=boilerplate -p 8080:8080 fiber-boilerplate
```


## Live Reloading (Air)
Example configuration files for [Air](https://github.com/cosmtrek/air) have also been included.
This allows you to live reload your Go application when you change a model, view or controller which is very useful when developing your application.

To run Air, use the following commands. Also, check out [Air its documentation](https://github.com/cosmtrek/air) about running the `air` command.
```bash
# Windows
air -c .air.windows.conf
# Linux
air -c .air.linux.conf
```

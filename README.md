# envpop

This is a tool for projects that you wish to integrate with CI/CD but may not be containrised and by nature do not source
their configuration from the environment but you wish to integrate with CI and have it build the environment specific
configuration dynamically at build time.

envpop takes a `.env` example or dist file that can be populated with defaults and source environment specific
configuration dynamically. It should correctly detect the type of environment variable value and set them accordingly
(i.e. string, numbers, floats, boolean, null etc).

## Example

Take the following partial Laravel framework `.env.example` file:

```bash
APP_NAME=Laravel
APP_ENV=local
APP_KEY=
APP_DEBUG=true
APP_URL=http://localhost

LOG_CHANNEL=stack
LOG_DEPRECATIONS_CHANNEL=null
LOG_LEVEL=debug

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=laravel
DB_USERNAME=root
DB_PASSWORD=
````

Executing:

```bash
export APP_KEY="SECRET"
export DB_PASSWORD="SUPER_SECRET"

envpop -file ./.env.example
```

Will output the following dynamic configured .env file:

```bash
APP_NAME=Laravel
APP_ENV=local
APP_KEY="SECRET"
APP_DEBUG=true
APP_URL=http://localhost

LOG_CHANNEL=stack
LOG_DEPRECATIONS_CHANNEL=null
LOG_LEVEL=debug

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=laravel
DB_USERNAME=root
DB_PASSWORD="SUPER_SECRET"
```

## License

envpop is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).

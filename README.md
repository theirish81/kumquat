# Kumquat
An application that allows you to expose various types of resources as HTTP endpoints, using simple descriptor files.

## Problem statement
In certain testing and development scenarios, you may need to access resources which are  not directly exposed to the
consumer application. Kumquat allows you to expose them as HTTP endpoints, finely tuning what should be exposed and how.

Resources may include:
* files
* rows generated by SQL queries
* documents generated by MongoDB queries
* output of shell commands
* output of 3rd party programs

## Basic domain objects
Each API action corresponds to a `sequence`.
* `sequence`: is a set of steps Kumquat must perform. Once the sequence is fully executed, the API will return
  the results
* `step`: is a single operation and the source of data
* `filter`: each step can be assigned a number of filters. Their main purpose is transforming the data provided by
  the step

## Variables scope
When a sequence executes, each step has access to a set of variables which is composed as follows:
* `values` file: a YAML file where env specific values can be stored. These variables are accessible through the
  `Values` map, a lot like Helm
* environment: as a default, scripts have no access to the environment variables for security reasons.
  However, in the sequence definition you can list out which environment variables should be passed over
  to the scripts
* the steps output: each step will store their output as a variable in the scope. This data can then be
  accessed by subsequent steps, using the name of the step
* `Params`: disabled by default (but can be enabled at sequence level), will take the JSON request body and
  make it available to the steps as the `Params` variable. Careful though, security issues may arise!

### Accessing the variables
All string fields in the sequence can be parsed as [gowalker](https://github.com/theirish81/gowalker) templates. 
Examples:
* `${Values.key}` returns the `key` variable in the `Values` collection (loaded from the `values.yaml` file)
* `${Params.id}` returns the `id` variable in the `Params` collection (loaded from the input data)
* `${firstProcess}` returns the output of the `firstProcess` step

### Passing parameters to the sequence
You can pass parameters to a sequence in the request body. This functionality is disabled by default for security
reasons as command injection is possible. To enable it, set:
* `accept_params` to `true` in the sequence main definition 
Now you can `POST` a JSON body to the sequence, such as:
```json
{
  "id": 2,
  "name": "foo"
}
```
This will populate the `Params` collection, therefore: `${Params.id}` will translate into `2` and `${Params.name}`
will translate into `foo`.

Additionally, you can also add a `requires` section to the sequence, to list out the parameters which are necessary to
complete the sequence. The expected content of this section is a JSON Schema:
```yaml
requires:
  required:
    - id
    - name
  properties:
    id:
      type: integer
    name:
      type: string
```

## Results
Each successful step execution will produce a result entry, unless the `hide: true` param is set at the step
definition. If a command fails to run, the `errors` array will be populated.


## Steps
The fields common to all steps are:
* `type`: the type of step.
* `name`: the name of the step, it should be one word, because it will be used as variable name
* `hide`: if set to `true`, the step will execute and set the return value in the scope, but will not display in the
  results (default: `false`)
* `config`: a map of values containing the implementation specific parameters

### file
Loads a file from the file system. The config keys are:
* `path`: the path to the file
* `basePath`: the base path for the file
* `raw`: if set to `true` the output will be binary, base64 encoded

### nixShellCommand
Executes *NIX shell command. The config keys are:
* `command`: the *NIX command itself
* `stdin`: if set to `true`, the variable scope will be piped into the standard input of the child process, in the form
  of JSON (default: `false`)
* `timeout`: the timeout for the command in "duration" format (default: `10s`)

### template
Executes a template against the variable scope
* `template`: the [gowalker](https://github.com/theirish81/gowalker) template

### sql
Executes a `select` query against a SQL server
* `driver`: the name of the driver. Currently, the available options are `postgres`, `mysql` and `mssql`
* `URI`: the URI of the service
* `select`: the select query
* `timeout`: the timeout for the query in "duration" format (default: `10s`)

### mongo
Executes a `find` query against a MongoDB server
* `URI`: URI to the MongoDB server
* `db`: the database name
* `collection`: the collection name
* `find`: the query, in the form of a map
* `timeout`: the timeout for the query in "duration" format (default: `10s`)

## Filters
The shared fields are:
* `type`: the type of the filter
* `config`: a map of values containing the implementation specific parameters (optional)

### splitLines
Splits the output string into an array where each item corresponds to one line.

### split
splits the output string into an array, using a provided separator sequence
* `sep`: the separator string

### jsonParse
parses the output string into a JSON data structure

### yamlParse
parses the output string into a YAML data structure

### replace
replaces a substring matching a regular expression, with a provided string
* `regexp`: a regular expression matching what we want to replace
* `replace`: the replacement string

## Authentication
All requests to run a sequence need to be authenticated with an `Authorization: Bearer ...` token.
By using basic authentication against the `authorize` endpoint, Kumquat will deliver an access token that you can
then use as bearer token.
The singing process will need a private/public RSA keypair to sign the JWT token and a file describing the users
(see "Server configuration").

The `users` file should look like the following:
```yaml
username:
    password: $2a$12$MXWDXhdiA02CGrzBSjcfbuAjVMbUYKt7fL0lAkItg61OPjYmGmu3G
    access: restricted
    sequences:
      - http
```
* `username`: the username
* `password`: the password, which must be hashed and salted with `bcrypt`. To generate a password, you can use one
of the many utility websites, `htpasswd` or use the included utility by issuing:
`docker run --rm -ti --entrypoint "bash" kumquat -c '/usr/local/kumquat/pass.sh foobar'` where `foobar` is your
cleartext password.
* `access`: either `all` (access to all sequences), or `restricted` (the user has access only to the sequences listed 
  in `sequences`)
* `sequences`: a list of sequence names the user has access to, if `access` is `restricted`


## API
Please check the [OpenAPI spec](specs/v1_openapi.yaml).

## Server configuration
The server configuration happens via environment variables:
* `VALUES_PATH`: path to the values file (default: `etc/values.yaml`)
* `SEQUENCES_PATH`: path to the directory containing the sequences (default: `etc/sequences`)
* `PRIVATE_KEY_PATH`: path to the private certificate for authentication (default: `etc/keys/private.pem`)
* `PUBLIC_KEY_PATH`: path to the public certificate for authentication (default: `etc/keys/public.pem`)
* `USERS_PATH`: path to the YAML file describing the users (default: `etc/users.yaml`)
* `TOKEN_DURATION`: the TTL for the token, in "duration" format (default: `24h`)

## Examples
Find examples in the [etc folder](etc).

## Running with Docker
The Docker image comes in a fully Debian distribution, and packs Python3 as well. This is important because as you can
shell scripts, the base image determines the availability of commands.

Of course, this image is meant to be extended and customized by the final user, so that more commands can be added on
need.

On *NIX systems, run:
```yaml
docker run -v "$(pwd)/etc:/usr/local/kumquat/etc" -p 5000:5000 theirish81/kumquat
```
where `"$(pwd)/etc` is an `etc` directory in the current working dir.

On any other system, apply the OS-dependant path references and separators.
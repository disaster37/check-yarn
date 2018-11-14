# check-yarn
Monitor Yarn jobs with Nagios from Yarn Ressource Manager API.

Its general purpose is to monitor Yarn jobs using external monitoring tools like Nagios, Centreon, Shinken, Sensu, etc…

The following program must be called by your monitoring tools and it returns a status (using Nagios status normalization) with human-readable messages and optionnaly perfdata.

This program calls the Ressource Manager API to get the state of your Yarn jobs.

You can use it to monitor the following section of your Yarn cluster:
- *Jobs state*: It checks that there is no alert on your host.

## Usage

### Global parameters

You need to set the Ressource Manager API informations for all checks.

```sh
./check-yarn --yarn-url https://yarn.company.com --yarn-login admin --yarn-password admin ... 
```

You need to specify the following parameters:
- **yarn-url**: It's your Ressource Manager URL throught Knox (or directly). Alternatively you can use the environment variable `YARN_URL`.
- **yarn-login**: It's the Knox login to use when calling the API. Alternatively you can use the environment variable `YARN_LOGIN`.
- **yarn-password**: It's the password associated with the login. Alternatively you can use the environment variable `YARN_PASSWORD`.

You can also set this parameters in YAML file(s) (one or many) and use the parameter `--config` with the path of your YAML file.
```yaml
---
yarn-url: https://yarn.company.com
yarn-login: admin
yarn-password: admin
```

### Check the jobs state

You need to lauch the following command:

```sh
./check-yarn --yarn-url https://yarn.company.com --yarn-login admin --yarn-password admin check-jobs --finished-since 24
```

You need to specify the following parameters:
- **--finished-since**: The number of hours since the current datetime.
- **--queue-name** (optionnal): The queue name where to check for jobs.
- **--user-name** (optionnal): The user which runs the jobs that you want to check.

The check follows this logic:
1. `OK` when there is no failed jobs with the given filter
2. `CRITICAL` when there is one or more jobs in the failed state

> All jobs which have failed are displayed on the output.

It returns the following perfdata:
- **nbJobFailed**: the number of failed jobs

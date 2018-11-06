# check-yarn
Monitor Yarn jobs on Nagios from Ressource Manager API.

It's a general purpose to monitore Yarn jobs with external monitoring tools like Nagios, Centreon, Shinken, Sensu, etc.

The following program must be called by your monitoring tools and it return the status (nagios status normalization) with human messages and some times perfdatas.
This program called Ressource Manager API to compute the state of your Yarn jobs.

You can use it to monitore the following section of your Yarn cluster:
- *Jobs state*: It's check that there are no alert in your host.

## Usage

### Global parameters

You need to set the Ressource Manager API informations for all checks.

```sh
./check-yarn --yarn-url https://yarn.company.com --yarn-login admin --yarn-password admin ... 
```

You need to specify the following parameters:
- **yarn-url**: It's your Ressource Manager URL throught Knox (or directly). Alternatively you can use environment variable `YARN_URL`.
- **yarn-login**: It's the Knox login to use when it call the API. Alternatively you can use environment variable `YARN_LOGIN`.
- **yarn-password**: It's the password associated with the login. Alternatively you can use environment variable `YARN_PASSWORD`.

You can set also this parameters on yaml file (one or all) and use the parameters `--config` with the path of your Yaml file.
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
- **--finished-since**: The number of hour since the current datetime.
- **--queue-name** (optionnal): The queue name where you should to check jobs.
- **--user-name** (optionnal): The user that run the jobs where you should to check.

This check follow this logic:
1. `OK` when there are no failed jobs with the given filter
2. `CRITICAL` when there are one or more jobs in failed state

> All jobs that failed is displayed on the outpout.

It's return the following perfdata:
- **nbJobFailed**: the number of jobs failed

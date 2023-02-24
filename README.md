# eye - a little tool to check dependencies by configuring a dependency blacklist.
![](https://komarev.com/ghpvc/?username=itisokey-eye&color=orange&style=flat&label=pv)

Supported project types include:
- [x] Maven
- [ ] Gradle
- [ ] ...  

## Usage
You can use this tool in GitHub Actions, Gitlab CI or local machine.


## GitHub Actions
Add `.dependency.yaml` file to the root directory of your project and add the following.
```yaml
dependency:
  files:
    - pom.xml # If this is a maven project.
  black-list: # Support regular expressions, the priority is groupId > artifactId > version
    - groupId: junit
    - groupId: com.alibaba.*
      artifactId: fastjson
      version:
```
and add the following to GHA `workflows`
```yaml
- name: Dependency Eye
  uses: lv-lifeng/eye@latestTag
  #with:
    #log: debug # optional: set the log level. The default value is `info`.
    #config: .dependency.yaml # optional: set the config file. The default value is `.dependency.yaml`.
    #token: # optional: the token that dependency eye uses when it needs to comment on the pull request. Set to empty ("") to disable commenting on pull request. The default value is ${{ github.token }}
    #mode: # optional: Which mode Dep-Eye should be run in. The default value is `check`.
```
## Gitlab CI
First, dep-eye commands need to be configured in gitlab runner.
```yaml
dep-check-job:       # job name.
  tags: [dep_check] # gitlab runner tag.
  rules:
    # trigger condition.
    - if: $CI_PIPELINE_SOURCE == 'merge_request_event' && $CI_MERGE_REQUEST_TARGET_BRANCH_NAME == 'main'
  script:
    - dep-eye d check
```

## Other
### Download [release](https://github.com/lv-lifeng/eye/releases)
Download binary file `Assets/eye.zip`, and add `.dependency.yaml` file to the root directory of your project or the other specified directory(e.g. `/User/other/dependency.yaml`), execute the following command in specified directory.
```shell
%PATH%/eye/bin/linux/dep-eye dependency(d/dep) -c /User/other/dependency.yaml check
```
or add `%PATH%/eye/bin/linux` to the environment variable and execute the following command everywhere.
```shell
dep-eye dependency(d/dep) -c /User/other/dependency.yaml check
```
if the `-c` parameter is not specified and the current directory does not have `.dependency.yaml` file, then `dependency-default.yaml` will be used.

### Compile from source
```shell
git clone git@github.com:lv-lifeng/eye.git
cd eye
make build 
```
the command same as [download release](#download-releasehttpsgithubcomlv-lifengeyereleases)

## Check result
```shell
dep-eye d check
INFO Loading configuration from file: .dependency.yaml 
INFO Config file .dependency.yaml does not exist, using the default config: eye/dependency-default.yaml 
INFO Start checking dependencies, please wait!    
Black-List           |                                                                                  Path
-------------------- | -------------------------------------------------------------------------------------
com.alibaba:fastjson | org.apache.rocketmq:rocketmq-acl:4.9.2 -> org.apache.rocketmq:rocketmq-remoting:4.9.2

ERROR found 1 dependencies hit the blacklist 
```
`Black-List:` dependence in the blacklist  
`Path:` parent dependency of dependence in the blacklist

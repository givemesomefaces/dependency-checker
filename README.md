# eye
A little tool to check dependencies

## How to use


### download [release](https://github.com/lv-lifeng/eye/releases)
Add `.dependency.yaml` file to the root directory of your project or the other specified directory(e.g. `/User/test.yaml`), and add the following.
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
Execute the following command in specified directory
```shell
%PATH%/eye/bin/linux/dep-eye dependency(d/dep) -c /User/test.yaml check
```
or add `%PATH%/eye/bin/linux` to the environment variable and execute the following command everywhere.
```shell
dep-eye dependency(d/dep) -c /User/test.yaml check
```
if the `-c` parameter is not specified and the current directory does not have `.dependency.yaml` file, then `dependency-default.yaml` will be used,

### compile from source
```shell
git clone git@github.com:lv-lifeng/eye.git
cd eye
make build 
```
the command same as [download release](#download-releasehttpsgithubcomlv-lifengeyereleases)

### check result
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
`Path：` parent dependency of dependence in the blacklist

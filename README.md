# eye
A little tool to check dependencies

## How to use

add `.dependency.yaml` file to the root directory of your project and add the following, if it does not exist, the default file `assets/default-config.yaml` will be used.
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

compile from source
```shell
git clone git@github.com:lv-lifeng/eye.git
cd eye
make build 
```

Execute the following command in specified directory, this directory is the root directory of the project to be checked
```shell
%PATH%/eye/bin/linux/dep-eye dependency(d/dep) check
```
or add `%PATH%/eye/bin/linux` to the environment variable and execute the following command everywhere.
```shell
dep-eye dependency(d/dep) check
```

check result:
```shell
dep-eye d check
INFO Loading configuration from file: .dependency.yaml 
INFO Config file .dependency.yaml does not exist, using the default config 
INFO Start checking dependencies, please wait!    
Black-List           |                                                                                  Path
-------------------- | -------------------------------------------------------------------------------------
com.alibaba:fastjson | org.apache.rocketmq:rocketmq-acl:4.9.2 -> org.apache.rocketmq:rocketmq-remoting:4.9.2

ERROR found 1 dependencies hit the blacklist 
```
`Black-List:` dependence in the blacklist  
`Path：` parent dependency of dependence in the blacklist

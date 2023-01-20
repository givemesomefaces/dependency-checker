# eye
A little tool to check dependencies

## How to use
download by https://github.com/lv-lifeng/eye/releases or make build by source code
```shell
git clone git@github.com:lv-lifeng/eye.git
cd eye
make build 
```

Execute the following command in specified directory, this directory is the root directory of the project to be checked
```shell
%PATH%/eye/bin/linux/dep-eye dependency(d/dep) check
```
or add %PATH%/eye/bin/linux to the environment variable and execute the following command everywhere.
```shell
dep-eye dependency(d/dep) check
```

Add `.dependency.yaml` file to the root directory of your project and add the following, if it does not exist, the default file `default-config.yaml` will be used.
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


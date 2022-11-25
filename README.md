# Requirements

- You must have postgres as an alias to localhost for this to work in /etc/hosts

  #### BEFORE
  ![Before](img/hosts_before.png)

  #### AFTER
  ![Before](img/hosts_after.png)

- You must have postgres port as 5432 in dockerapi/containerctl.sh You have to restart or re-setup your docker after you
  make this change
  #### BEFORE
  ![Before](img/containerctl_before.png)

  #### AFTER
  ![Before](img/containerctl_after.png)

# SETUP

This setups up the package for use

```go 
go mod tidy
```

This is an optional setup. It's only needed when you want to run the application on command line without golang

```shell 
./build.sh
```

# Using Command Line

## sMigration

#### SETUP

This setups up the necessary files needed for migration or seeding. Run one of these commands

```go 
go run main.go -e development -s seed -a setup
```

or

```shell
./build/mac/packages -e development -s seed -a setup
```

```go 
go run main.go -e development -s migration -a setup
```

```shell
./build/mac/packages -e development -s migration -a setup
```

#### RUN

This runs either migration or seeding. Run the necessary command

```go 
go run main.go -e development -s seed -a run
```

or

```shell
./build/mac/packages -e development -s seed -a run
```

```go 
go run main.go -e development -s migration -a run 
```

or

```shell
./build/mac/packages -e development -s migration -a run 
```
   

# Go Exercise
## Configuration
The configuration file is contained in the `/config` folder in YAML format. In order to run the application,
it must contain the following configuration:
- `mappers_number`: an integer representing the number of mappers components
- `reducers_number`: an integer representing the number of reducers components
- `master`: a complex object containing the `address` of the host and the `port`
- `workers`: a list of complex objects, each containing the `address` of the host and the `port`
## Run
Workers must be started before master.
### Note
Before running any component, move into the `/main` folder:
```
cd ./main
```
### Workers
To start a worker, run
```
go run . worker -i {index}
```
where `{index}` is the index of the worker in the configuration file.
### Master
To run the master, run
```
fo run . master
```
# cmsurvivalrecorder
record funded data from the website https://cameroonsurvival.org/fr/dons/

Query the website, parse the data, and save in the local db.
The local data will persist on the machine. 

## Query the data

Obtain last data using endpoint `/records/last`
```
> curl 0.0.0.0:9090/records/last
{"time":"2020-04-12T10:58:30.803482Z","value":522251.99999999994}
```

```
> curl 0.0.0.0:9090/records/last
{"time":"2020-04-12T11:01:35.546632Z","value":522272.00000000006}
```

Obtain all data using endpoint `records/timeseries`

```
> curl 0.0.0.0:9090/records/timeseries
{
    "2020-04-12T10:58:30.803482Z":522251.99999999994,
    "2020-04-12T11:01:35.546632Z":522272.00000000006
}
```

## Required 
- [docker](www.docker.com) for running docker images using `docker-compose up`
- [go 1.13 or higher](https://golang.org/) for development
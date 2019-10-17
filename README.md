# Spring Bridge
Spring Bridge is a tool used to do mass setting of backgrounds to soltice displays

## Create a directory named `img`
Place at least 6 photos in this directory. The script will random select 6 from the directory if you plae more than 6

## Create the MySQL config file

```
{
    "MySQLUser":"User",
    "MySQLPassoword":"Secret",
    "MySQLHost":"Databas host",
    "MySQLDB":"Database Name",
    "MySQLQuery":"Query"
}
```

## Usage

No special build requirements. Run `go build` then execute the ending bianary.

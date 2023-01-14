# API Usage

Our API creates gRPC server which provides 4 different methods which are listed below

- Get Aggregated Scores by Categories Daily/Weekly
- Get Scores by Tickets
- Get Overal Performance Score between period
- Get Score difference between 2 Date Periods

# Get Aggregated Scores by Categories Daily/Weekly

### Input Schema

```
{
    "range_from": "2019-01-17", 
    "range_to": "2019-07-19"
}
```

Our gRPC method `TicketAnalysisService/GetAggregatedCategoryScores` accepts the above structure and returns the following schema 

```
{
    "categoryName": "Randomness",
    "ratingsCount": 27,
    "dateScores": [
    {
        "date": "2019-01-17"
    },
    {
        "date": "2019-05-30",
        "score": 45
    }
    ]
}
```

**NB:** If `score` field does not exist, front-end should consider it as `N/A` value. Ranges are defines based on the input date range. Iteration through dates are possible.

`gRPCurl` command example
```
grpcurl --plaintext -d '{"range_from": "2019-07-17", "range_to": "2019-07-19"}' 127.0.0.1:8080 TicketAnalysisService/GetAggregatedCategoryScores 
```


# Get Scores by Tickets

This function will return the scores aggregated based on `ticket_id`s. 

### Input Schema

```
{
    "range_from": "2019-01-17", 
    "range_to": "2019-07-19"
}
```

Our gRPC method `TicketAnalysisService/GetScoresByTicket` accepts the above structure and returns the following schema 

```
{
    "scores": [
        {
            "ticketId": 123,
            "categoryScores": [
                {
                    "categoryName": "Spelling",
                    "score": 50
                }, 
            ]
        },
    ]
}
```

Empty "score" field means `N/A`

`gRPCurl` Command Example
```
grpcurl --plaintext -d '{"range_from": "2019-07-17", "range_to": "2019-07-19"}' 127.0.0.1:8080 TicketAnalysisService/GetScoresByTicket 
```

# Get Overal Performance Score

### Input Schema

```
{
    "range_from": "2019-01-17", 
    "range_to": "2019-07-19"
}
```

Our gRPC method `TicketAnalysisService/GetScoreOveralForQuality` accepts the above structure and returns the following schema to illustrate the overal quality score during the period

```
{
    "score": 50
}
```

`gRPCurl` command example
```
grpcurl --plaintext -d '{"range_from": "2019-07-01", "range_to": "2019-07-30"}' 127.0.0.1:8080 TicketAnalysisService/GetScoreOveralForQuality    
```

# Get Score difference between 2 Date Periods
### Input Schema

```
{
    "end_period": {
        "range_from": "2019-07-01", 
        "range_to": "2019-07-30"
    },
    "previous_period":{
        "range_from": "2019-08-01", 
        "range_to": "2019-08-30"
    }
    
}
```

Our gRPC method `TicketAnalysisService/GetScoreChangePeriodOverPeriod` accepts the above structure and returns the following schema to illustrate the changes between overal qualities of the periods. 

```
{
    "changeScore": -4
}
```

`gRPCurl` command example
```
grpcurl --plaintext -d '{"end_period": {"range_from": "2019-08-01", "range_to": "2019-08-30"}, "previous_period": {"range_from": "2019-07-01", "range_to": "2019-07-30"}}' 127.0.0.1:8080 TicketAnalysisService/GetScoreChangePeriodOverPeriod
```
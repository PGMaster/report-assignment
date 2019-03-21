## William Test Review

#### 1. The request pattern does not match the Task

```
Provided Task specification:
http://localhost/chapter_versions/27
```
 
Current implementation expects a query string param:

```
http://localhost/chapter_versions?chapter_id=:chapter_id
```

#### 2. Result is wrong

The result provided by the request is as below:

```
[
    {
        "company": "CtrlPrint",
        "project": "Testnot",
        "chapter": "Test Chapter  1",
        "versions": null
    },
    {
        "company": "",
        "project": "",
        "chapter": "",
        "versions": [
            {
                "created_by": "ctrl-romain",
                "chapter_version_id": 262,
                "version_number": 4,
                "created": "2015-08-26T10:49:53.059514Z",
                "appversion": "CC 2015"
            },
            {
                "created_by": "ctrl-jens",
                "chapter_version_id": 261,
                "version_number": 3,
                "created": "2015-08-26T10:48:41.795795Z",
                "appversion": "CC 2015"
            },
            {
                "created_by": "ctrl-jens",
                "chapter_version_id": 211,
                "version_number": 2,
                "created": "2015-07-07T13:12:49.847299Z",
                "appversion": "CC 2015"
            },
            {
                "created_by": "ctrl-jens",
                "chapter_version_id": 210,
                "version_number": 1,
                "created": "2015-07-07T13:12:15.262319Z",
                "appversion": "CC 2015"
            },
            {
                "created_by": "ctrl-jens",
                "chapter_version_id": 209,
                "version_number": 0,
                "created": "2015-07-07T13:12:03.467114Z",
                "appversion": "CC 2015"
            },
            {
                "created_by": "ctrl-an",
                "chapter_version_id": 611,
                "version_number": 7,
                "created": "2019-02-27T16:39:59.037887Z",
                "appversion": "CC 2015"
            },
            {
                "created_by": "ctrl-an",
                "chapter_version_id": 610,
                "version_number": 6,
                "created": "2019-02-27T16:35:24.461144Z",
                "appversion": "CC 2015"
            },
            {
                "created_by": "ctrl-arun",
                "chapter_version_id": 607,
                "version_number": 5,
                "created": "2018-11-27T13:29:17.094967Z",
                "appversion": "CC 2017"
            }
        ]
    }
]
```

The results expands based on the number of requests called.

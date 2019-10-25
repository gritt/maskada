### API Contract

> **Create transaction**
> ```
> curl -X POST {{domain}}/v1/transaction \
>   -H 'Content-Type: application/json' \
>   -d '{
>     "amount": 26,
>     "type": 2,
>     "category": "Food",
>     "name": "Family Flavor"
> }'
> ```
> Response :: 201 Created
> ```
> {
>    "id": 11,
>    "amount": 26,
>    "type": 2,
>    "category": "Food",
>    "date": "2019-10-25T00:26:56.707907Z",
>    "name": "Family Flavor"
> }
> ```

<br>

> **List transactions**
> ```
> curl -X GET {{domain}}/v1/transaction
> ```
> Response :: 200 OK
> ```
> [
>    {
>        "id": 9,
>        "amount": 1300,
>        "type": 1,
>        "category": "Rent",
>        "date": "2019-10-25T00:26:21Z",
>        "name": ""
>    },
>    {
>        "id": 11,
>        "amount": 26,
>        "type": 2,
>        "category": "Food",
>        "date": "2019-10-25T00:26:57Z",
>        "name": "Family Flavor"
>    }
> ]
>```

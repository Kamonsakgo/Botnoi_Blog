# Template

## RUN APP

```bash
# install golang package
$ go mod tidy

# Start APP
$ go run . || go run main.go
```

## Ads Endpoint

### Get All Ads Data (No JWT)

HTTP Method: GET
Endpoint: /api/payment/all_Ads
Function: gateway.GetAllAdsData

### Get Ads Data (No JWT)

HTTP Method: GET
Endpoint: /api/payment/Ads
Function: gateway.GetAdsData

### Set Redis Ads (No JWT)

HTTP Method: GET
Endpoint: /api/payment/set_ads
Function: gateway.SetRedisAds

### Get Ads None Token (No JWT)

HTTP Method: GET
Endpoint: /api/payment/get_ads_none_token
Function: gateway.GetAdsNoneToken

### Get Marketplace Sound (JWT Required)

HTTP Method: GET
Endpoint: /api/payment/get_marketplace_sound
Middleware: middlewares.SetBotnoiJWtHeaderHandler
Function: gateway.GetMarketplaceSound
response:

```json
{
    "message": "get ads success",
    "data": [
        {
            "description": "หม้อไฟฟ้าตัวนี้น่ารักปุ้กปิ้กมาก....",
            "id": "9",
            "language": "th",
            "level": 1,
            "status": true,
            "url": ""
        },
        ...
    ]
}
```

#### Update Ads (JWT Required)

HTTP Method: POST
Endpoint: /api/payment/update_ads
Middleware: middlewares.SetBotnoiJWtHeaderHandler
Function: gateway.CheckCoupon
payload:

```json
{
  "ads_play": [
    {
      "id": "207",
      "play": 1
    }
  ]
}
```

response:

```json
{
  "message": "success update ads"
}
```

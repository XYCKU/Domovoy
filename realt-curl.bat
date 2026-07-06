$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$session.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36"
$session.Cookies.Add((New-Object System.Net.Cookie("consent", "{%22analytics%22:false%2C%22advertising%22:false%2C%22functionality%22:true}", "/", ".realt.by")))
$session.Cookies.Add((New-Object System.Net.Cookie("currency-notice", "1", "/", "realt.by")))
$session.Cookies.Add((New-Object System.Net.Cookie("sl-session", "7r+ZVvgTTWrw4gtZtoTTIw==", "/", "realt.by")))
Invoke-WebRequest -UseBasicParsing -Uri "https://realt.by/bff/graphql" `
-Method "POST" `
-WebSession $session `
-Headers @{
"authority"="realt.by"
  "method"="POST"
  "path"="/bff/graphql"
  "scheme"="https"
  "accept"="*/*"
  "accept-encoding"="gzip, deflate, br, zstd"
  "accept-language"="en-US,en;q=0.9"
  "origin"="https://realt.by"
  "priority"="u=1, i"
  "referer"="https://realt.by/sale/flats/?page=2"
  "sec-ch-ua"="`"Google Chrome`";v=`"149`", `"Chromium`";v=`"149`", `"Not)A;Brand`";v=`"24`""
  "sec-ch-ua-mobile"="?0"
  "sec-ch-ua-platform"="`"Windows`""
  "sec-fetch-dest"="empty"
  "sec-fetch-mode"="cors"
  "sec-fetch-site"="same-origin"
  "x-realt-client"="www@6.14.4"
} `
-ContentType "application/json" `
-Body "[{`"operationName`":`"searchObjects`",`"variables`":{`"data`":{`"where`":{`"addressV2`":[{`"townUuid`":`"4cb07174-7b00-11eb-8943-0cc47adabd66`"}],`"category`":5},`"pagination`":{`"page`":13,`"pageSize`":30},`"sort`":[{`"by`":`"paymentStatus`",`"order`":`"DESC`"},{`"by`":`"priority`",`"order`":`"DESC`"},{`"by`":`"raiseDate`",`"order`":`"DESC`"},{`"by`":`"updatedAt`",`"order`":`"DESC`"}],`"extraFields`":[`"minPriceAggregation`"],`"isReactAdaptiveUA`":false}},`"query`":`"query searchObjects(`$data: GetObjectsByAddressInput!) {\n  searchObjects(data: `$data) {\n    body {\n      results {\n        companyName\n        companyUuid\n        uuid\n        title\n        description\n        headline\n        createdAt\n        updatedAt\n        metroTime\n        metroTimeType\n        price\n        priceCurrency\n        pricePerM2\n        pricePerM2Max\n        pricePerPerson\n        priceMin\n        priceMax\n        priceChangeDirection\n        priceChangeDate\n        storeys\n        storey\n        rooms\n        contactPhones\n        images\n        areaTotal\n        areaLiving\n        areaMax\n        areaMin\n        areaLand\n        objectType\n        code\n        stateRegionName\n        stateDistrictName\n        townType\n        townName\n        streetUuid\n        streetName\n        address\n        contactName\n        contactEmail\n        agencyName\n        metroStationName\n        metroLineId\n        houseNumber\n        buildingNumber\n        paymentStatus\n        comments\n        isFavorite\n        category\n        has3dTour\n        hasVideo\n        stateRegionUuid\n        numberOfBeds\n        directionName\n        townDistance\n        customSorting\n        specialComment\n        userUuid\n        agencyUuid\n        location\n        townUuid\n        buildingYear\n        levels\n        roofMaterial\n        wallMaterial\n        heating\n        infrastructure\n        balconyType\n        houseType\n        furniture\n        areaKitchen\n        appliances\n        objectCategory\n        realEstateDevUuid\n        availableYear\n        availableQuarter\n        availableAlready\n        availableText\n        isSellingCompleted\n        communicationMethod\n        interactiveCatalogToken\n        interactiveCatalogBaseToken\n        isObjectInRealtyDeal\n        repairState\n        __typename\n      }\n      pagination {\n        page\n        pageSize\n        totalCount\n        __typename\n      }\n      rates {\n        from\n        to\n        rate\n        __typename\n      }\n      extraFields {\n        minPriceAggregation\n        __typename\n      }\n      __typename\n    }\n    ...StatusAndErrors\n    __typename\n  }\n}\n\nfragment StatusAndErrors on INullResponse {\n  success\n  errors {\n    code\n    title\n    message\n    field\n    __typename\n  }\n  __typename\n}`"}]"
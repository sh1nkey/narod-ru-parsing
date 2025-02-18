что он щас может:

/api/v1/narod/ GET limit offset - выдавать все записи

/api/v1/narod/save POST {url: httpsUrl, html: htmlContent } - добавить url в БД

/api/v1/narod/set-description PATCH {id: int, description: string} - добавить описание для записи

/api/v1/narod/mark-empty PATCH {id: int, description: string} - пометить запись как пустую

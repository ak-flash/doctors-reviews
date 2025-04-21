## Получение отзывов о врачах с сайтов СберЗдоровье и ПроДокторов

Поддерживаемые сайты с отзывами:
- СберЗдоровье (сменили защиту от парсинга, не работает на данный момент)
- ПроДокторов

***

> POST http://127.0.0.1:8000/api/v1/getReviews  
Accept: application/json  
Content-Type: application/x-www-form-urlencoded

platform=prodoctorov&doctorUrl=https://prodoctorov.ru/ekaterinburg/vrach/ID-Фамилия_Имя

***

*platform* = **sberzdorovie** или **prodoctorov**  
*doctorUrl* = ссылка на профиль врача

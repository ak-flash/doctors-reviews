## Получение отзывов о врачах с сайтов СберЗдоровье и ПроДокторов

Поддерживаемые сайты с отзывами:
- СберЗдоровье
- ПроДокторов

***

> POST http://127.0.0.1:8000/api/v1/getReviews  
Accept: application/json  
Content-Type: application/x-www-form-urlencoded  
platform=sberzdorovie & doctorUrl=https://docdoc.ru/doctor/Фамилия_Имя

***

*platform* = **sberzdorovie** или **prodoctorov**  
*doctorUrl* = ссылка на профиль врача
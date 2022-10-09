# Тест план gRPC

\* - негативные кейсы<br/>
\# - кейсы по заданию 2.2.1

## Тест-план по Device-Api

### Тест 1.1 Регистрация устройства GRPC

Предварительные условия:
1) Система запущена.
2) Проходит Readiness probe
3) Нет записи с `deviceId` = **deviceID**

| Шаг | Запрос | Параметры запроса | Ответ | Параметры ответа | База данных |
|-----|--------|-------------------|-------|------------------|-------------|
| 1.  | Отправьте запрос CreateDeviceV1Request |`user_id`: **userID**, `platform`: **platform** | * Пришел ответ CreateDeviceV1Response | `deviceId` = **deviceID** > 0| В БД появилась 1 запись: `platform`: **platform**, `user_id`: , `entered_at`: **timestamp_1** <= Текущее время, `removed`: false, `created_at` <= entered_at, `updated_at` >= created_at |

                                                                                                                        
### Тест 1.2 Получение информации об устройстве

Предварительные условия:
1) Система запущена.
2) Проходит Readiness probe
3) Есть как минимум одна запись с `device_id` = **existent_id**, нет записи с id = **nonexistentId**

| Шаг | Действие | Ожидаемый результат | Параметры ответа |
|-----|----------|---------------------|------------------|
| 1.  | Отправьте запрос DescribeDeviceV1Request | `device_id`: **existent_id** |  Пришел ответ CreateDeviceV1Response | `deviceId`:  **device_ID_1** > 0; <br/> `id`: **device_ID_1**;<br/> `platform`: **platform_1**;<br/> `userId`: **userId_1**;<br/> `enteredAt`: "**timestamp_1**" |
| 2.  | Отправьте запрос DescribeDeviceV1Request | `device_id`: **nonexistentId** | Пришел ответ "5 NOT_FOUND" | `code`: 5,<br/> `message`: device not found<br/> `details`: [] |              

### Тест 1.3 Редактирование устройства

Предварительные условия:
1) Система запущена.
2) Проходит Readiness probe
3) Есть записи в базе данных не менее одной штуки, записи уникальны по полю `device_id`
4) Пусть первая запись имеет `device_id` = **existent_id**, нет записи с `device_id` = **nonexistentId**
5) Пусть нынешнее время = **timestamp_now**

| Шаг | Запрос | Параметры запроса | Ответ | Параметры ответа |
|-----|--------|-------------------|-------|------------------|
| 1.  | Отправьте запрос UpdateDeviceV1Request | `device_id`: **existent_id**,<br/> `platform`: "Ios", <br/>`userId`: **user_id**<br/> | Пришел ответ UpdateDeviceV1Response | `deviceId`:  **device_ID_1** > 0; <br/> `id`: **device_ID_1**;<br/> `platform`: **platform_1**;<br/> `userId`: **userId_1**;<br/> `enteredAt`: **timestamp_now**<br/> |
| 2.  | Отправьте запрос UpdateDeviceV1Request | `device_id`: **nonexistentId**; `platform`: "Ios",<br/>`userId`: **userID**<br/> | Пришел ответ UpdateDeviceV1Response со статусом 200 | `success`: "false" |

### Тест 1.4 Удаление устройства gRPC

Предварительные условия:
1) Система запущена.
2) Проходит Readiness probe
3) Есть записи в базе данных не менее одной штуки, записи уникальны по полю `device_id`
4) Пусть первая запись имеет `device_id` = **existent_id**, нет записи с `device_id` = **nonexistentId**

| Шаг | Запрос | Параметры запроса | Ответ | Параметры ответа | База данных |
|-----|--------|-------------------|-------|------------------|-------------|
| 1.  | Отправьте запрос RemoveDeviceV1Request | `device_id` = **existent_id**| * Пришел ответ RemoveDeviceV1Response | `found` = "true" типа "bool"|В поле `deleted` у записи с `device_id` = **existent_id** теперь отметка "true"|
| 2.  | Отправьте запрос RemoveDeviceV1Request | `device_id` = **nonexistent_id**| * Пришел ответ RemoveDeviceV1Response | `found` = "false" типа "bool"|Такой записи нет в БД|
|\*3. | Отправьте запрос RemoveDeviceV1Request | `device_id` = **existent_id**| * Пришел ответ RemoveDeviceV1Response | `found` = "false" типа "bool"|поле `deleted` у записи с `device_id` = **existent_id** все еще отметка "true"|
|\*4. | Отправьте запрос RemoveDeviceV1Request | `device_id` = "String"| * Пришел ответ RemoveDeviceV1Response со статусом 400 | `code`: 3,<br/>  `message`: "type mismatch, parameter: device_id, error: strconv.ParseUint: parsing \"\\\"String\\\"\": invalid syntax",<br/>`details`: []|-|


### Тест 1.5 Список устройств gRPC

Предварительные условия:
1) Система запущена.
2) Проходит Readiness probe
3) Есть **n** записей в базе данных, записи уникальны по полю `device_id`
4) Пусть первая запись имеет `device_id` = **existent_id**, нет записи с `device_id` = **nonexistentId**
5) 

| Шаг | Запрос | Параметры запроса | Ответ | Параметры ответа |
|-----|--------|-------------------|-------|------------------|
| 1.  | Отправьте запрос gRPC RemoveDeviceV1Request | `per_page`: **m** >= 1, `page`: **k** <= ceil(**n**, **m**)| Пришел ответ RemoveDeviceV1Response | `items`: type array, содержащий **n** девайсов;<br/>Структура `device` должна быть описана так: `deviceId` > 0 = **deviceID**, `platform`: **platform**, `user_id` > 0 = **user_id**, `entered_at` >0 = ± **timestamp_1**|


##  Тест-план по протоколу act_notification_api

### Тест 2.1 Отправка уведомлений по устройству

Предварительные условия:
1) Система запущена.
2) Проходит Readiness probe
3) Есть устройства с `device_id`: **device_id_1**, **device_id_2**, **device_id_3**, **device_id_4**, нет устройства с **device_id_5**
4) Есть возможность менять время дня `timeOfTheDay`: **t**
5) Запущены четыре подписки по адресам **device_id_1**, **device_id_2**, **device_id_3**, **device_id_4** по gRPC (SubsctibeNotificationV1)

*Здесь я сгенерировал таблицу Pairwise со следующими колонками: `timeOfTheDay`, `OS`, `lang` <br/>

***Morning**: 6:00 < **t** < 10:59; **Afternoon**: 11:00 < **t** < 14:59, **Evening**: 15:00 < **t** < 20:59, **Night**: 21:00 < **t** < 05:59
  
*OS: **Android**: **device_id_2**, **MacOS**: **device_id_4**, **IOS** **device_id_1**, **Windows**: **device_id_3**


| Шаг | Запрос | Параметры запроса | Ответ | Параметры ответа | База данных | Подписка |
|-----|--------|-------------------|-------|------------------|-------------|----------|
| 1.  | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Morning** | `deviceId`: **device_id_1**, <br/>`lang`: **LANG_ENGLISH**,<br/>`message`: **notifText** = Dear **IOS** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Good morning, **notifText**.| В подписке **device_id_1** появилось уведомление **notifID** |
| 2.  | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Afternoon** | `deviceId`: **device_id_2**, <br/>`lang`: **LANG_ENGLISH**,<br/>`message`: **notifText** = Dear **Android** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Good Afternoon, **notifText**.| В подписке **device_id_2** появилось уведомление **notifID** |
| 3.  | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Evening** | `deviceId`: **device_id_3**, <br/>`lang`: **LANG_ENGLISH**,<br/>`message`: **notifText** = Dear **Windows** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Good Evening, **notifText**.| В подписке **device_id_3** появилось уведомление **notifID** |
| 4.  | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Night** | `deviceId`: **device_id_4**, <br/>`lang`: **LANG_ENGLISH**,<br/>`message`: **notifText** = Dear **MacOS** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Good Night, **notifText**.| В подписке **device_id_4** появилось уведомление **notifID** |
| 5.  | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Morning** | `deviceId`: **device_id_3**, <br/>`lang`: **LANG_RUSSIAN**,<br/>`message`: **notifText** = Dear **Windows** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Доброе утро, **notifText**.| В подписке **device_id_3** появилось уведомление **notifID** |
| 6.  | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Afternoon** | `deviceId`: **device_id_4**, <br/>`lang`: **LANG_RUSSIAN**,<br/>`message`: **notifText** = Dear **MacOS** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Добрый День, **notifText**.| В подписке **device_id_4** появилось уведомление **notifID** |
| 7.  | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Evening** | `deviceId`: **device_id_1**, <br/>`lang`: **LANG_RUSSIAN**,<br/>`message`: **notifText** = Dear **IOS** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Добрый Вечер, **notifText**.| В подписке **device_id_1** появилось уведомление **notifID** |
| 8.  | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Night** | `deviceId`: **device_id_2**, <br/>`lang`: **LANG_RUSSIAN**,<br/>`message`: **notifText** = Dear **Android** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Доброй ночи, **notifText**.| В подписке **device_id_2** появилось уведомление **notifID** |
| 9.  | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Morning** | `deviceId`: **device_id_1**, <br/>`lang`: **LANG_ESPANOL**,<br/>`message`: **notifText** = Dear **IOS** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Buenos dias, **notifText**.| В подписке **device_id_1** появилось уведомление **notifID** |
| 10. | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Afternoon** | `deviceId`: **device_id_2**, <br/>`lang`: **LANG_ESPANOL**,<br/>`message`: **notifText** = Dear **Android** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Buenos tardes, **notifText**.| В подписке **device_id_2** появилось уведомление **notifID** |
| 11. | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Evening** | `deviceId`: **device_id_3**, <br/>`lang`: **LANG_ESPANOL**,<br/>`message`: **notifText** = Dear **Windows** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Buenas noches, **notifText**.| В подписке **device_id_3** появилось уведомление **notifID** |
| 12. | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Night** | `deviceId`: **device_id_4**, <br/>`lang`: **LANG_ESPANOL**,<br/>`message`: **notifText** = Dear **MacOS** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Buenas noches, **notifText**.| В подписке **device_id_4** появилось уведомление **notifID** |
| 13. | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Morning** | `deviceId`: **device_id_3**, <br/>`lang`: **LANG_ITALIAN**,<br/>`message`: **notifText** = Dear **Windows** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Buon giorno, **notifText**.| В подписке **device_id_3** появилось уведомление **notifID** |
| 14. | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Afternoon** | `deviceId`: **device_id_4**, <br/>`lang`: **LANG_ITALIAN**,<br/>`message`: **notifText** = Dear **MacOS** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Buon pomeriggio, **notifText**.| В подписке **device_id_4** появилось уведомление **notifID** |
| 15. | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Evening** | `deviceId`: **device_id_1**, <br/>`lang`: **LANG_ITALIAN**,<br/>`message`: **notifText** = Dear **IOS** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Buona serata, **notifText**.| В подписке **device_id_1** появилось уведомление **notifID** |
| 16. | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Night** | `deviceId`: **device_id_2**, <br/>`lang`: **LANG_ITALIAN**,<br/>`message`: **notifText** = Dear **Android** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | `notificationId`: **notifID**| В базе данных появляется аналогичная запись с полем `message`: Buona notte, **notifText**.| В подписке **device_id_2** появилось уведомление **notifID** |
| 17*. | Отправьте запрос gRPC SendNotificationV1<br/>Время дня **Night** | `deviceId`: **device_id_5**, <br/>`lang`: **LANG_ITALIAN**,<br/>`message`: **notifText** = Dear **Android** user,<br/>`notificationId`: "0",<br/>`notificationStatus`: "STATUS_CREATED",<br/>`username`: **username** | Пришел ответ SendNotificationV1Response | 13 INTERNAL (ошибка)| |

### Тест 2.2 Получение уведомлений 


Предварительные условия:
1) Система запущена.
2) Проходит Readiness probe
3) Есть устройство с `device_id`: **device_id_1**, для него есть **m** > 0 уведомлений в базе данных c `notifID`: **1**...**m**


| Шаг | Запрос | Параметры запроса | Ответ | Параметры ответа | База данных |
|-----|--------|-------------------|-------|------------------|-------------|
| 1.  | Отправьте запрос gRPC GetNotificationV1 | `deviceId`: **device_id_1** | Пришел ответ GetNotificationV1Response, с количеством элементов **m**| `notificationId`: **notifID**| В базе данных для `deviceId`: **device_id_1**, где статус STATUS_CREATED, уведомления меняют статус на STATUS_IN_PROGRESS.|
| 2.  | Отправьте запрос gRPC GetNotificationV1 | `notifID`: **n** <= **m** | Пришел ответ GetNotificationV1Response с одним уведомлением| `notificationId`: **notifID** == **n**| В базе данных для **notifID** == **n**, если статус STATUS_CREATED, уведомление меняет статус на STATUS_IN_PROGRESS.|
| 3.  | Отправьте запрос gRPC GetNotificationV1 | `deviceId`: **device_id_1**<br/>, `notifID`: **n** < **m** != 0 | Пришел ответ GetNotificationV1Response с одним уведомлением| `notificationId`: **notifID** == **n**| В базе данных для **notifID** == **n**, если статус STATUS_CREATED, уведомление меняет статус на STATUS_IN_PROGRESS.|
| 4.  | Отправьте запрос gRPC GetNotificationV1 | `deviceId`: **device_id_1**<br/>, `notifID`: **n** <= **m** = 0 | Пришел ответ GetNotificationV1Response, с количеством элементов **m**| `notificationId`: **notifID**| В базе данных для `deviceId`: **device_id_1**, где статус STATUS_CREATED, уведомления меняют статус на STATUS_IN_PROGRESS.|


### 2.3 Подтверждение получения уведомления
Предварительные условия:
1) Система запущена.
2) Проходит Readiness probe
3) Есть **m** > 0 уведомлений в базе данных c `notifID`: **1**...**m**

| Шаг | Запрос | Параметры запроса | Ответ | Параметры ответа | База данных |
|-----|--------|-------------------|-------|------------------|-------------|
| 1.  | Отправьте запрос gRPC AckNotificationV1 | `notifID`: **n** <= **m** | Пришел ответ AckNotificationV1Response| `success`: "true"| В базе данных для `deviceId`: **device_id_1**, где статус STATUS_IN_PROGRESS, уведомления меняют статус на STATUS_DELIVERED.|
|*2.  | Отправьте запрос gRPC AckNotificationV1 | `notifID`: **n** > **m** | Пришел ответ AckNotificationV1Response| `success`: "false"| В базе данных нет этого значения|

### 2.4 Подписка на уведомления
Предварительные условия:
1) Система запущена.
2) Проходит Readiness probe
3) Есть устройство с `device_id`: **device_id_1**, для него есть **m** уведомлений в базе данных c STATUS_CREATED, STATUS_IN_PROGRESS - `notifID`: **1**...**m** 

| Шаг | Запрос | Параметры запроса | Ответ | Параметры ответа | База данных | Подписка |
|-----|--------|-------------------|-------|------------------|-------------|----------|
| 1.  | Отправьте запрос gRPC SubscribeNotification | `device_id`: **device_id_1** |  | в STATUS_CREATED - оно должно сменить статус в БД на STATUS_IN_PROGRESS.| Пришли все уведомления с `notifID`: **1**...**m** |
| 2.  | Отправьте запрос gRPC SendNotificationV1 | Произвольные параметры | Пришел ответ SendNotificationV1Response | `notifID`: **m** + 1 |в STATUS_CREATED - оно должно сменить статус в БД на STATUS_IN_PROGRESS.| Пришли все уведомления с `notifID`: **1**...**m** |

# 4 Интеграция с Kafka
### Тест 4.1 Общий тест интегации с кафкой

Предварительные условия:
1) Система запущена.
2) Проходит Readiness probe
3) Нет устройства с `device_id`: **device_id_1**

| Шаг | Запрос | Параметры запроса | devices_events до ивента | Kafka | devices_events после ивента |
|-----|--------|-------------------|--------------------------|-------|-----------------------------|
| 1.  | Отправьте запрос CreateDeviceV1Request |`user_id`: **userID**,<br/> `platform`: **platform** | `id`: **index**, <br/> `deviceID`: **device_id_1**,<br/>`type`: 1,<br/>`status`:1,<br/>`payload`:**json_content**,<br/>`created_at`: **time_now** = **time_created**,<br/>`updated_at`: **time_now**> | Создался ивент в кафка | Запись `deviceID`: **device_id_1** изменилась на `status`: 2|
| 2.  | Отправьте запрос UpdateDeviceV1Request |`deviceID`: **device_id_1**<br/>`user_id`: **userID**,<br> `platform`: **platform** | `id`: **index**, <br/> `deviceID`: **device_id_1**,<br/>`type`: 2,<br/>`status`:1,<br/>`payload`:**json_content**,<br/>`created_at`: **time_created**,<br/>`updated_at`: **time_now**> | Создался ивент в кафка | Запись `deviceID`: **device_id_1** изменилась на `status`: 2|
| 2.  | Отправьте запрос DeleteDeviceV1Request |`deviceID`: **device_id_1**| `id`: **index**, <br/> `deviceID`: **device_id_1**,<br/>`type`: 3,<br/>`status`:1,<br/>`payload`:**json_content**,<br/>`created_at`: **time_created**,<br/>`updated_at`: **time_now**> | Создался ивент в кафка | Запись `deviceID`: **device_id_1** изменилась на `status`: 2|
